package utils

import (
	"github.com/00unnmd/pills_parser/internal/domain"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func CreatePIWithError(pillValue string, regionValue string, err error, pharmacy string) []domain.ParsedItem {
	return []domain.ParsedItem{
		{
			Pharmacy:        pharmacy,
			Region:          regionValue,
			Name:            pillValue,
			Mnn:             "",
			Price:           0,
			Discount:        0,
			DiscountPercent: 0,
			Producer:        "",
			Rating:          0,
			ReviewsCount:    0,
			SearchValue:     "",
			Error:           err.Error(),
		},
	}
}

func ParseRawData[T domain.ParsedFieldsGetter](pharmacy string, region string, pill string, rawData []T) []domain.ParsedItem {
	var parsed []domain.ParsedItem

	for _, item := range rawData {
		fields := item.GetFields(pharmacy, region, pill)
		parsed = append(parsed, fields)
	}

	return parsed
}

func FilterByProducer[T domain.ProducerGetter](filterItems []T, pillValue string) []T {
	var filtered []T
	producerNames := OzonProducerNames

	if pillValue == OzonPillsList[13] {
		producerNames = append(producerNames, AgroFarmProducerNames...)
	}
	if pillValue == OzonPillsList[26] {
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

	var validFiles []domain.FileInfo

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

		validFiles = append(validFiles, domain.FileInfo{
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
