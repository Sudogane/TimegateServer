package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/sudogane/project_timegate/internal/database"
	"github.com/sudogane/project_timegate/internal/router"
	"github.com/sudogane/project_timegate/internal/server"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env: ", err)
		return
	}

	databaseString := os.Getenv("POSTGRES_DATABASE_URL")
	databaseRepository, err := database.NewRepository(databaseString)

	if err != nil {
		fmt.Println("database err: ", err)
		return
	}
	defer databaseRepository.Close()

	gameServer := server.NewGameServer(databaseRepository)
	router := router.NewRouter(gameServer)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		gameServer.HandleWebsocket(w, r, router)
	})

	err = http.ListenAndServe(":8888", nil)
	if err != nil {
		fmt.Println(err)
	}
}
