package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/00unnmd/pills_parser/models"
	"github.com/xuri/excelize/v2"
)

func CreatePIWithError(pillValue string, regionValue string, err error) []models.ParsedItem {
	log.Println(err.Error())
	return []models.ParsedItem{
		{
			Region:       regionValue,
			Name:         pillValue,
			Price:        0,
			Discount:     0,
			PriceOld:     0,
			MaxQuantity:  0,
			Producer:     "",
			Rating:       0,
			ReviewsCount: 0,
			Error:        err.Error(),
		},
	}
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

func GetPillsReqFormValues(r *http.Request) (int, int, string, string, string) {
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(r.FormValue("perPage"))
	if err != nil || perPage < 1 {
		perPage = 10
	}

	sortField := r.FormValue("sort")
	if sortField == "" {
		sortField = "name"
	}

	sortOrder := r.FormValue("order")
	if sortOrder == "" {
		sortOrder = "ASC"
	}

	filter := r.FormValue("filter")

	return page, perPage, sortField, sortOrder, filter
}

func FindLatestParsingFile(dir string) (string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	patternReg := regexp.MustCompile(`^parsing-go-(\d{2})\.(\d{2})\.(\d{4})\.xlsx$`)

	var validFiles []models.FileInfo

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		matches := patternReg.FindStringSubmatch(file.Name())
		if len(matches) != 4 {
			continue
		}

		fileDate, err := time.Parse("02.01.2006", matches[1]+"."+matches[2]+"."+matches[3])
		if err != nil {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		validFiles = append(validFiles, models.FileInfo{
			Path:    filepath.Join(dir, file.Name()),
			ModTime: info.ModTime(),
			Date:    fileDate,
		})
	}

	if len(validFiles) == 0 {
		return "", os.ErrNotExist
	}

	sort.Slice(validFiles, func(i, j int) bool {
		return validFiles[i].Date.After(validFiles[j].Date)
	})

	return validFiles[0].Path, nil
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
	log.Println("Запущен процесс генерации документа...")
	file := excelize.NewFile()

	for key, item := range data {
		sheetName := key
		_, err := file.NewSheet(sheetName)

		if err != nil {
			log.Println(err)
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
		file.SetCellValue(sheetName, "J1", "Error")

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
			file.SetCellValue(sheetName, "J"+strconv.Itoa(i+2), item[i].Error)
		}
	}

	colWidths := map[string]float64{
		"A": 30,
		"B": 60,
		"G": 40,
		"J": 60,
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
