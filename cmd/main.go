package main

import (
	"api_v2/database"
	"api_v2/internal/routers"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.InitDatabaseConnection()
}

func main() {

	router := routers.SetUpRouter()
	url := os.Getenv("MAIN_HOST")
	err := router.Listen(url)
	if err != nil {
		log.Fatal("Error setup router!", err)
	}
}
