package server

import (
	"github.com/00unnmd/pills_parser/internal/transport"
	"github.com/00unnmd/pills_parser/pkg/database/postgres"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

func Start() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file: %w", err)
	}

	db, err := postgres.ConnectToPostgres()
	if err != nil {
		log.Fatal("ConnectToPostgres err: %w", err)
	}
	defer db.Close()

	r, err := transport.SetupRouter(db)
	if err != nil {
		log.Printf("setupRouter err: %s", err)
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:4173", os.Getenv("PROD_ORIGIN")},
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
