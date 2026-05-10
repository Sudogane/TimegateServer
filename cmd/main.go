package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sudogane/project_timegate/internal/database"
	"github.com/sudogane/project_timegate/internal/database/cache"
	"github.com/sudogane/project_timegate/internal/logger"
	"github.com/sudogane/project_timegate/internal/router"
	"github.com/sudogane/project_timegate/internal/server"
)

func loadEnv() error {
	err := godotenv.Load("../.env")
	if err != nil {
		return err
	}
	return nil
}

func startDatabase(logger *logger.Logger) (*database.Repository, error) {
	databaseString := os.Getenv("POSTGRES_DATABASE_URL")
	if databaseString == "" {
		return nil, errors.New("POSTGRES_DATABASE_URL is not set")
	}

	logger.Log("INFO", "Connecting to the Database...")
	databaseRepository, err := database.NewRepository(databaseString)

	if err != nil {
		return nil, err
	}

	logger.Log("INFO", "Sucessfully connected to the Database")
	return databaseRepository, nil
}

func startRedisCache(logger *logger.Logger) (*cache.RedisClient, error) {
	logger.Log("INFO", "Connecting to redis...")
	redisAddr := os.Getenv("REDIS_ADDR")
	redisUsr := os.Getenv("REDIS_USR")
	redisPwd := os.Getenv("REDIS_PWD")

	if redisAddr == "" {
		redisAddr = "localhost:6379"
		logger.Log("INFO", "RedisAddr not Defined on Env, Running on local.")
	}

	rdb := cache.NewRedisClient(&redis.Options{
		Addr:     redisAddr,
		Username: redisUsr,
		Password: redisPwd,
		DB:       0,
	})

	if rdb == nil {
		return nil, errors.New("Failed to connect to redis")
	}

	logger.Log("INFO", "Connected to Redis")
	return rdb, nil
}

func main() {
	logger, err := logger.NewLogger()
	if err != nil {
		fmt.Println("Error loading custom logger: ", err)
		return
	}

	err = loadEnv()
	if err != nil {
		logger.Log("ERROR", "Env Loading Error: "+err.Error())
		return
	}

	databaseRepository, err := startDatabase(logger)
	if err != nil {
		logger.Log("ERROR", "Database Error: "+err.Error())
		return
	}
	defer databaseRepository.Close()

	redisDb, err := startRedisCache(logger)
	if err != nil {
		logger.Log("ERROR", "Redis Error: "+err.Error())
		return
	}
	defer redisDb.Close()

	gameServer := server.NewGameServer(databaseRepository, redisDb, logger)
	router := router.NewRouter(gameServer)

	// Gracefull Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      nil,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		gameServer.HandleWebsocket(w, r, router)
	})

	go func() {
		logger.Log("INFO", "WebSocket server listening on ws://localhost:"+port+"/ws")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log("ERROR", "Server Startup Error: "+err.Error())
			stop <- os.Interrupt
		}
	}()

	// Wait for shutdown
	<-stop
	logger.Log("INFO", "Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Log("ERROR", "Error while shutting down server: "+err.Error())
	}

	logger.Log("INFO", "Server down.")
}
