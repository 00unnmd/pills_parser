package main

import (
	"github.com/00unnmd/pills_parser/handlers"
	"github.com/00unnmd/pills_parser/internals/database"
	"github.com/00unnmd/pills_parser/internals/router"
	"github.com/00unnmd/pills_parser/internals/utils"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func parseNow() {
	data := handlers.GetAllData()
	database.SaveToDB(data)
	utils.GenerateXLSX(data)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file: %w", err)
	}

	db, err := database.ConnectToPostgres()
	if err != nil {
		log.Fatal("ConnectToPostgres err: %w", err)
	}
	defer db.Close()

	r, err := router.SetupRouter(db)
	if err != nil {
		log.Printf("setupRouter err: %s", err)
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://192.168.31.106:5173", "http://192.168.31.106:5175"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Content-Disposition"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
