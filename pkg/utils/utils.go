package utils

import (
	"github.com/00unnmd/pills_parser/internal/domain"
	"strings"
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
