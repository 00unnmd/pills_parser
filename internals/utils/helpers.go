package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/00unnmd/pills_parser/models"
	"github.com/xuri/excelize/v2"
)

func FilterByProducer(filterItems []models.PillsItem, producersArr []string) []models.PillsItem {
	var filtered []models.PillsItem

	for _, item := range filterItems {
		for _, producer := range producersArr {
			if item.Producer == producer {
				filtered = append(filtered, item)
				break
			}
		}
	}

	return filtered
}

func getFilePath() string {
	now := time.Now().Format("02.01.2006") // TODO проверить как работает формат с часами и минутами
	outputDir := filepath.Join("result/")

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Errorf("Не удалось создать папку %s: %v", outputDir, err)
	}

	fileName := filepath.Join(outputDir, fmt.Sprintf("parsing-go-%s.xlsx", now))

	return fileName
}

func GenerateXLSX(data []models.PillsItem) {
	fmt.Println("Данные получены. Генерация документа...")
	file := excelize.NewFile()
	sheetName := "Sheet1"
	_, err := file.NewSheet(sheetName)

	if err != nil {
		fmt.Println(err)
		return
	}

	file.SetCellValue(sheetName, "A1", "Region")
	file.SetCellValue(sheetName, "B1", "Name")
	file.SetCellValue(sheetName, "C1", "Price")
	file.SetCellValue(sheetName, "D1", "Discount")
	file.SetCellValue(sheetName, "E1", "priceOld")
	file.SetCellValue(sheetName, "F1", "maxQuantity")
	file.SetCellValue(sheetName, "G1", "producer")
	file.SetCellValue(sheetName, "H1", "isBundle")
	file.SetCellValue(sheetName, "I1", "rating")
	file.SetCellValue(sheetName, "J1", "reviewsCount")

	for i := 0; i < len(data); i++ {
		file.SetCellValue(sheetName, "A"+strconv.Itoa(i+2), data[i].Region)
		file.SetCellValue(sheetName, "B"+strconv.Itoa(i+2), data[i].Name)
		file.SetCellValue(sheetName, "C"+strconv.Itoa(i+2), data[i].Price)
		file.SetCellValue(sheetName, "D"+strconv.Itoa(i+2), data[i].Discount)
		file.SetCellValue(sheetName, "E"+strconv.Itoa(i+2), data[i].PriceOld)
		file.SetCellValue(sheetName, "F"+strconv.Itoa(i+2), data[i].MaxQuantity)
		file.SetCellValue(sheetName, "G"+strconv.Itoa(i+2), data[i].Producer)
		file.SetCellValue(sheetName, "H"+strconv.Itoa(i+2), data[i].IsBundle)
		file.SetCellValue(sheetName, "I"+strconv.Itoa(i+2), data[i].Rating)
		file.SetCellValue(sheetName, "J"+strconv.Itoa(i+2), data[i].ReviewsCount)
	}

	filePath := getFilePath()

	if err := file.SaveAs(filePath); err != nil {
		fmt.Println("Ошибка сохранения файла. Парсинг не завершен.")
		fmt.Println("err: ", err)
		return
	}

	fmt.Println("Документ успешно сгенерирован (result/). Парсинг завершен.")
}
