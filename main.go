package main

import (
	"log"

	"github.com/00unnmd/pills_parser/handlers"
	"github.com/00unnmd/pills_parser/internals/utils"
	"github.com/joho/godotenv"
)

func parseNow() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	data := handlers.GetAllData()
	utils.GenerateXLSX(data)
}

func main() {
	parseNow()
}
