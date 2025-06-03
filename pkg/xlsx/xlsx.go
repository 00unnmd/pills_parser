package xlsx

import (
	"fmt"
	"github.com/00unnmd/pills_parser/internal/domain"
	"github.com/xuri/excelize/v2"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func getFilePath() (string, error) {
	now := time.Now().Format("02.01.2006")
	outputDir := filepath.Join("result/")

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("не удалось создать папку %s: %w", outputDir, err)
	}

	fileName := filepath.Join(outputDir, fmt.Sprintf("parsing-go-%s.xlsx", now))

	return fileName, nil
}

func GenerateXLSX(data []domain.ParsedItem, sheetName string) {
	log.Println("Запущен процесс генерации документа...")
	file := excelize.NewFile()

	_, err := file.NewSheet(sheetName)

	if err != nil {
		log.Println(err)
		return
	}

	file.SetCellValue(sheetName, "A1", "Pharmacy")
	file.SetCellValue(sheetName, "B1", "Region")
	file.SetCellValue(sheetName, "C1", "Name")
	file.SetCellValue(sheetName, "D1", "MNN")
	file.SetCellValue(sheetName, "E1", "Price")
	file.SetCellValue(sheetName, "F1", "Discount")
	file.SetCellValue(sheetName, "G1", "DiscountPercent")
	file.SetCellValue(sheetName, "H1", "producer")
	file.SetCellValue(sheetName, "I1", "rating")
	file.SetCellValue(sheetName, "J1", "reviewsCount")
	file.SetCellValue(sheetName, "K1", "SearchValue")
	file.SetCellValue(sheetName, "L1", "Error")

	for i, item := range data {
		file.SetCellValue(sheetName, "A"+strconv.Itoa(i+2), item.Pharmacy)
		file.SetCellValue(sheetName, "B"+strconv.Itoa(i+2), item.Region)
		file.SetCellValue(sheetName, "C"+strconv.Itoa(i+2), item.Name)
		file.SetCellValue(sheetName, "D"+strconv.Itoa(i+2), item.Mnn)
		file.SetCellValue(sheetName, "E"+strconv.Itoa(i+2), item.Price)
		file.SetCellValue(sheetName, "F"+strconv.Itoa(i+2), item.Discount)
		file.SetCellValue(sheetName, "G"+strconv.Itoa(i+2), item.DiscountPercent)
		file.SetCellValue(sheetName, "H"+strconv.Itoa(i+2), item.Producer)
		file.SetCellValue(sheetName, "I"+strconv.Itoa(i+2), item.Rating)
		file.SetCellValue(sheetName, "J"+strconv.Itoa(i+2), item.ReviewsCount)
		file.SetCellValue(sheetName, "K"+strconv.Itoa(i+2), item.SearchValue)
		file.SetCellValue(sheetName, "L"+strconv.Itoa(i+2), item.Error)
	}

	colWidths := map[string]float64{
		"A": 10,
		"B": 30,
		"C": 60,
		"H": 40,
		"K": 30,
		"L": 60,
	}

	for _, sheet := range file.GetSheetList() {
		for col, width := range colWidths {
			if err := file.SetColWidth(sheet, col, col, width); err != nil {
				log.Println(err)
			}
		}
		if sheet == "Sheet1" {
			file.DeleteSheet(sheet)
		}
	}

	filePath, err := getFilePath()
	if err != nil {
		log.Println(err)
		return
	}

	if err := file.SaveAs(filePath); err != nil {
		log.Println("Ошибка сохранения файла. Парсинг не завершен.")
		log.Println("err: ", err)
		return
	}

	log.Println("Документ успешно сгенерирован (result/). Парсинг завершен.")
}
