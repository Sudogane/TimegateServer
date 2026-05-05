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

func startDatabase() (*database.Repository, error) {
	databaseString := os.Getenv("POSTGRES_DATABASE_URL")
	if databaseString == "" {
		return nil, errors.New("POSTGRES_DATABASE_URL is not set")
	}

	fmt.Println("Connecting to database")
	databaseRepository, err := database.NewRepository(databaseString)

	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to database")
	return databaseRepository, nil
}

func startRedisCache() (*cache.RedisClient, error) {
	fmt.Println("Connecting to redis")
	redisAddr := os.Getenv("REDIS_ADDR")
	redisUsr := os.Getenv("REDIS_USR")
	redisPwd := os.Getenv("REDIS_PWD")

	if redisAddr == "" {
		redisAddr = "localhost:6379"
		fmt.Println("REDIS_ADDR not defined, running on local redis.")
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

	fmt.Println("Connected to redis")
	return rdb, nil
}

func main() {
	fmt.Println("Starting Server")

	err := loadEnv()
	if err != nil {
		fmt.Println("Error loading .env: ", err)
		return
	}

	databaseRepository, err := startDatabase()
	if err != nil {
		fmt.Println("Error starting database: ", err)
		return
	}
	defer databaseRepository.Close()

	redisDb, err := startRedisCache()
	if err != nil {
		fmt.Println("Error starting redis: ", err)
		return
	}
	defer redisDb.Close()

	gameServer := server.NewGameServer(databaseRepository, redisDb)
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
		fmt.Printf("WebSocket server listening on ws://localhost:%s/ws\n", port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Error starting server: ", err)
			stop <- os.Interrupt
		}
	}()

	// Wait for shutdown
	<-stop
	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Error shutting down server: ", err)
	}

	fmt.Println("Server shut down")
}
