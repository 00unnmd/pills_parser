package handlers

import (
	"fmt"
	"time"

	"github.com/00unnmd/pills_parser/internals/utils"
	"github.com/00unnmd/pills_parser/models"
	"github.com/cheggaaa/pb/v3"
)

func reqAllRegions(pillValue string, bar *pb.ProgressBar) []models.PillsItem {
	var result []models.PillsItem
	delay := 2 * time.Second // delay between requests

	for key, value := range utils.Regions {
		bar.Set("prefix", value)
		bar.Increment()

		time.Sleep(delay)
		pillsForRegion := RequestZdravsitiData(pillValue, key, value)
		result = append(result, pillsForRegion...)
	}

	return result
}

func ReqAllPills() []models.PillsItem {
	fmt.Println("Получение данных...")
	var result []models.PillsItem
	pillsB := pb.New(len(utils.PillsList))
	regionsB := pb.New(len(utils.Regions))
	pool, _ := pb.StartPool(pillsB, regionsB)

	for _, value := range utils.PillsList {
		regionsB.SetCurrent(0)
		pillsB.Set("prefix", value)
		pillsB.Increment()

		pillsAllRegions := reqAllRegions(value, regionsB)
		result = append(result, pillsAllRegions...)
	}

	pool.Stop()
	return result
}
