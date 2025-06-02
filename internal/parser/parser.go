package parser

import (
	"github.com/00unnmd/pills_parser/pkg/database/postgres"
	"github.com/00unnmd/pills_parser/pkg/xlsx"
	"github.com/joho/godotenv"
	"log"
)

func ParseNow() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file: %w", err)
	}

	data := GetMNNAllData()
	postgres.SaveToDB(data)
	xlsx.GenerateXLSX(data)
}
