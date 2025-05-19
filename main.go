package main

import (
	"log"

	"github.com/00unnmd/pills_parser/handlers"
	"github.com/00unnmd/pills_parser/internals/database"
	"github.com/00unnmd/pills_parser/internals/utils"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func parseNow() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	data := handlers.GetAllData()
	database.SaveToDB(data)
	utils.GenerateXLSX(data)
}

func main() {
	parseNow()
}
