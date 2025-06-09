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

	file, err := xlsx.GenerateXLSX(OZONData, "ozon_data")
	if err != nil {
		log.Println(err)
		return
	}

	err = xlsx.SaveFileToResult(file, "ozon_data")
	if err != nil {
		log.Println(err)
	}
}

func parseMNN() {
	MNNData := GetMNNAllData()
	postgres.SaveToDB(MNNData, "mnn_data")

	file, err := xlsx.GenerateXLSX(MNNData, "mnn_data")
	if err != nil {
		log.Println(err)
		return
	}

	err = xlsx.SaveFileToResult(file, "mnn_data")
	if err != nil {
		log.Println(err)
	}
}

func parseCompetitors() {
	competitorsData := GetCompetitorsAllData()
	postgres.SaveToDB(competitorsData, "competitors_data")

	file, err := xlsx.GenerateXLSX(competitorsData, "competitors_data")
	if err != nil {
		log.Println(err)
		return
	}

	err = xlsx.SaveFileToResult(file, "competitors_data")
	if err != nil {
		log.Println(err)
	}
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
