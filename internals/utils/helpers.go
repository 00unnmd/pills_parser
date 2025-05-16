package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/00unnmd/pills_parser/models"
	"github.com/xuri/excelize/v2"
)

func PrintStructAsJSON(data []models.ParsedItem) {
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(jsonData))
}

func ParseRawData[T models.ParsedFieldsGetter](rawData []T) []models.ParsedItem {
	var parsed []models.ParsedItem

	for _, item := range rawData {
		fields := item.GetFields()
		parsed = append(parsed, fields)
	}

	return parsed
}

func FilterByProducer[T models.ProducerGetter](filterItems []T, pillValue string) []T {
	var filtered []T
	producerNames := OzonProducerNames

	if pillValue == PillsList["flebofa"] {
		producerNames = append(producerNames, AgroFarmProducerNames...)
	}
	if pillValue == PillsList["naftoderil"] {
		producerNames = append(producerNames, KirovFFProducerNames...)
	}

	for _, item := range filterItems {
		for _, producer := range producerNames {
			itemProducer := item.GetProducer()
			isEqual := strings.EqualFold(itemProducer, producer)

			if isEqual {
				filtered = append(filtered, item)
				break
			}
		}
	}

	return filtered
}

func getFilePath() (string, error) {
	now := time.Now().Format("02.01.2006")
	outputDir := filepath.Join("result/")

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("не удалось создать папку %s: %w", outputDir, err)
	}

	fileName := filepath.Join(outputDir, fmt.Sprintf("parsing-go-%s.xlsx", now))

	return fileName, nil
}

func GenerateXLSX(data map[string][]models.ParsedItem) {
	fmt.Println("Запущен процесс генерации документа...")
	file := excelize.NewFile()

	for key, item := range data {
		sheetName := key
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
		file.SetCellValue(sheetName, "H1", "rating")
		file.SetCellValue(sheetName, "I1", "reviewsCount")

		for i := range item {
			file.SetCellValue(sheetName, "A"+strconv.Itoa(i+2), item[i].Region)
			file.SetCellValue(sheetName, "B"+strconv.Itoa(i+2), item[i].Name)
			file.SetCellValue(sheetName, "C"+strconv.Itoa(i+2), item[i].Price)
			file.SetCellValue(sheetName, "D"+strconv.Itoa(i+2), item[i].Discount)
			file.SetCellValue(sheetName, "E"+strconv.Itoa(i+2), item[i].PriceOld)
			file.SetCellValue(sheetName, "F"+strconv.Itoa(i+2), item[i].MaxQuantity)
			file.SetCellValue(sheetName, "G"+strconv.Itoa(i+2), item[i].Producer)
			file.SetCellValue(sheetName, "H"+strconv.Itoa(i+2), item[i].Rating)
			file.SetCellValue(sheetName, "I"+strconv.Itoa(i+2), item[i].ReviewsCount)
		}
	}

	colWidths := map[string]float64{
		"A": 30,
		"B": 60,
		"G": 40,
	}

	for _, sheet := range file.GetSheetList() {
		for col, width := range colWidths {
			if err := file.SetColWidth(sheet, col, col, width); err != nil {
				fmt.Println(err)
			}
		}
		if sheet == "Sheet1" {
			file.DeleteSheet(sheet)
		}
	}

	filePath, err := getFilePath()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := file.SaveAs(filePath); err != nil {
		fmt.Println("Ошибка сохранения файла. Парсинг не завершен.")
		fmt.Println("err: ", err)
		return
	}

	fmt.Println("Документ успешно сгенерирован (result/). Парсинг завершен.")
}
