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
	"os"
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
		AllowedOrigins:   []string{"http://localhost:5173", os.Getenv("PROD_ORIGIN")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Content-Disposition"},
		AllowCredentials: true,
	})
	handler := c.Handler(r)

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("Сервер запущен на :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
