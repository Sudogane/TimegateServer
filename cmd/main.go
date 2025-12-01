package main

import (
	"context"
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
		return nil, fmt.Errorf("POSTGRES_DATABASE_URL is not set")
	}

	fmt.Println("Connecting to database")
	databaseRepository, err := database.NewRepository(databaseString)

	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to database")
	return databaseRepository, nil
}

func startRedisCache() *cache.RedisClient {
	fmt.Println("Connecting to redis")
	rdb := cache.NewRedisClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USR"),
		Password: os.Getenv("REDIS_PWD"),
		DB:       0,
	})

	fmt.Println("Connected to redis")
	return rdb
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

	redisDb := startRedisCache()
	defer redisDb.Close()

	gameServer := server.NewGameServer(databaseRepository, redisDb)
	router := router.NewRouter(gameServer)

	// Gracefull Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	server := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      nil,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		gameServer.HandleWebsocket(w, r, router)
	})

	go func() {
		fmt.Printf("WebSocket server listening on ws://localhost:%s/ws\n", os.Getenv("PORT"))

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
