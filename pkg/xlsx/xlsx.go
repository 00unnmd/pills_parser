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

func GenerateXLSX(data []domain.ParsedItem, sheetName string) (*excelize.File, error) {
	file := excelize.NewFile()

	_, err := file.NewSheet(sheetName)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("GenerateXLSX err: %w", err)
	}

	file.SetCellValue(sheetName, "A1", "CreatedAt")
	file.SetCellValue(sheetName, "B1", "Pharmacy")
	file.SetCellValue(sheetName, "C1", "Region")
	file.SetCellValue(sheetName, "D1", "Name")
	file.SetCellValue(sheetName, "E1", "MNN")
	file.SetCellValue(sheetName, "F1", "Price")
	file.SetCellValue(sheetName, "G1", "Discount")
	file.SetCellValue(sheetName, "H1", "DiscountPercent")
	file.SetCellValue(sheetName, "I1", "producer")
	file.SetCellValue(sheetName, "J1", "rating")
	file.SetCellValue(sheetName, "K1", "reviewsCount")
	file.SetCellValue(sheetName, "L1", "SearchValue")
	file.SetCellValue(sheetName, "M1", "Error")

	for i, item := range data {
		file.SetCellValue(sheetName, "A"+strconv.Itoa(i+2), item.CreatedAt)
		file.SetCellValue(sheetName, "B"+strconv.Itoa(i+2), item.Pharmacy)
		file.SetCellValue(sheetName, "C"+strconv.Itoa(i+2), item.Region)
		file.SetCellValue(sheetName, "D"+strconv.Itoa(i+2), item.Name)
		file.SetCellValue(sheetName, "E"+strconv.Itoa(i+2), item.Mnn)
		file.SetCellValue(sheetName, "F"+strconv.Itoa(i+2), item.Price)
		file.SetCellValue(sheetName, "G"+strconv.Itoa(i+2), item.Discount)
		file.SetCellValue(sheetName, "H"+strconv.Itoa(i+2), item.DiscountPercent)
		file.SetCellValue(sheetName, "I"+strconv.Itoa(i+2), item.Producer)
		file.SetCellValue(sheetName, "J"+strconv.Itoa(i+2), item.Rating)
		file.SetCellValue(sheetName, "K"+strconv.Itoa(i+2), item.ReviewsCount)
		file.SetCellValue(sheetName, "L"+strconv.Itoa(i+2), item.SearchValue)
		file.SetCellValue(sheetName, "M"+strconv.Itoa(i+2), item.Error)
	}

	colWidths := map[string]float64{
		"B": 10,
		"C": 30,
		"D": 60,
		"I": 40,
		"L": 30,
		"M": 60,
	}

	for _, sheet := range file.GetSheetList() {
		for col, width := range colWidths {
			if err := file.SetColWidth(sheet, col, col, width); err != nil {
				log.Println(err)
				return nil, fmt.Errorf("GenerateXLSX err: %w", err)
			}
		}
		if sheet == "Sheet1" {
			file.DeleteSheet(sheet)
		}
	}

	return file, nil
}

func getFilePath(filePrefix string) (string, error) {
	now := time.Now().Format("02.01.2006.1504")
	outputDir := filepath.Join("result/")

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("не удалось создать папку %s: %w", outputDir, err)
	}

	fileName := filepath.Join(outputDir, fmt.Sprintf("%s-%s.xlsx", filePrefix, now))
	return fileName, nil
}

func SaveFileToResult(file *excelize.File, filePrefix string) error {
	filePath, err := getFilePath(filePrefix)
	if err != nil {
		return err
	}

	if err := file.SaveAs(filePath); err != nil {
		return fmt.Errorf("ошибка сохранения файла. Парсинг не завершен. %w", err)
	}

	log.Println("Документ успешно сгенерирован (result/). Парсинг завершен.")
	return nil
}
