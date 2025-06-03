package parser

import (
	"github.com/00unnmd/pills_parser/pkg/database/postgres"
	"github.com/00unnmd/pills_parser/pkg/xlsx"
	"github.com/joho/godotenv"
	"log"
)

func parseOZON() {
	OZONData := GetOzonAllData()
	postgres.SaveToDB(OZONData, "ozon_data")
	xlsx.GenerateXLSX(OZONData, "ozon_data")
}

func parseMNN() {
	MNNData := GetMNNAllData()
	postgres.SaveToDB(MNNData, "mnn_data")
	xlsx.GenerateXLSX(MNNData, "mnn_data")
}

func parseCompetitors() {
	competitorsData := GetCompetitorsAllData()
	postgres.SaveToDB(competitorsData, "competitors_data")
	xlsx.GenerateXLSX(competitorsData, "competitors_data")
}

func ParseNow() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file: %w", err)
	}

	parseOZON()
	//parseMNN()
	//parseCompetitors()
}
