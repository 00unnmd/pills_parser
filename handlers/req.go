package handlers

import (
	"context"
	"github.com/00unnmd/pills_parser/handlers/pharmacies"
	"log"
	"time"

	"github.com/00unnmd/pills_parser/internals/utils"
	"github.com/00unnmd/pills_parser/models"
	"github.com/cheggaaa/pb/v3"
)

func getZSAllRegions(pillValue string, bar *pb.ProgressBar) []models.ParsedItem {
	var result []models.ParsedItem

	for key, value := range utils.ZSRegions {
		bar.Set("prefix", value)
		time.Sleep(utils.RequestDelay)

		pillsForRegion, err := pharmacies.GetZSPills(pillValue, key, value)
		if err != nil {
			pillsForRegion = utils.CreatePIWithError(pillValue, value, err)
		}

		result = append(result, pillsForRegion...)
		bar.Increment()
	}

	return result
}

func getZSData(pillsBar *pb.ProgressBar, regionsBar *pb.ProgressBar) []models.ParsedItem {
	var result []models.ParsedItem

	for _, value := range utils.PillsList {
		regionsBar.SetCurrent(0)
		pillsBar.Set("prefix", "ZdravCity: "+value)

		pillsAllRegions := getZSAllRegions(value, regionsBar)
		result = append(result, pillsAllRegions...)
		pillsBar.Increment()
	}

	return result
}

func getARAllPills(regionValue string, bar *pb.ProgressBar) []models.ParsedItem {
	var result []models.ParsedItem

	for _, value := range utils.PillsList {
		bar.Set("prefix", "AptekaRu: "+value)
		time.Sleep(utils.RequestDelay)

		pillsAllRegions, err := pharmacies.GetARPills(value, regionValue)
		if err != nil {
			pillsAllRegions = utils.CreatePIWithError(value, regionValue, err)
		}

		result = append(result, pillsAllRegions...)
		bar.Increment()
	}

	return result
}

func getARData(pillsBar *pb.ProgressBar, regionsBar *pb.ProgressBar) []models.ParsedItem {
	var result []models.ParsedItem

	for id, value := range utils.ARRegions {
		pillsBar.SetCurrent(0)
		regionsBar.Set("prefix", value)

		_, err := pharmacies.ChangeARRegion(id)
		if err != nil {
			log.Println("err: ", err)
			break
		}

		regionAllPills := getARAllPills(value, pillsBar)
		result = append(result, regionAllPills...)
		regionsBar.Increment()
	}

	return result
}

func getEAAllPills(ctx context.Context, bar *pb.ProgressBar, regionKey string, regionValue string) []models.ParsedItem {
	var result []models.ParsedItem

	for _, value := range utils.PillsList {
		bar.Set("prefix", "EApteka: "+value)
		time.Sleep(utils.RequestDelay)

		pillsAllRegions, err := pharmacies.GetEAPills(ctx, value, regionKey, regionValue)
		if err != nil {
			pillsAllRegions = utils.CreatePIWithError(value, regionValue, err)
		}

		result = append(result, pillsAllRegions...)
		bar.Increment()
	}

	return result
}

func getEAAllData(pillsBar *pb.ProgressBar, regionsBar *pb.ProgressBar) []models.ParsedItem {
	var result []models.ParsedItem

	ctx, cancel, err := pharmacies.CreateEAContext()
	if err != nil {
		log.Println("err: ", err)
		return result
	}
	defer cancel()

	for key, value := range utils.EARegions {
		pillsBar.SetCurrent(0)
		regionsBar.Set("prefix", value)

		_, err := pharmacies.ChangeEARegion(ctx, key)
		if err != nil {
			log.Println("err: ", err)
			break
		}

		regionAllPills := getEAAllPills(ctx, pillsBar, key, value)
		result = append(result, regionAllPills...)
		regionsBar.Increment()
	}

	return result
}

func GetAllData() map[string][]models.ParsedItem {
	log.Println("Процесс получения данных...")
	ZSPillsBar := pb.New(len(utils.PillsList)).SetRefreshRate(time.Second * 3)
	ZSRegionsBar := pb.New(len(utils.ZSRegions)).SetRefreshRate(time.Second * 3)
	ARPillsBar := pb.New(len(utils.PillsList)).SetRefreshRate(time.Second * 3)
	ARRegionsBar := pb.New(len(utils.ARRegions)).SetRefreshRate(time.Second * 3)
	EAPillsBar := pb.New(len(utils.PillsList)).SetRefreshRate(time.Second * 3)
	EARegionsBar := pb.New(len(utils.EARegions)).SetRefreshRate(time.Second * 3)

	ZSChan := make(chan []models.ParsedItem)
	ARChan := make(chan []models.ParsedItem)
	EAChan := make(chan []models.ParsedItem)

	pool, _ := pb.StartPool(ZSPillsBar, ZSRegionsBar, ARPillsBar, ARRegionsBar, EAPillsBar, EARegionsBar)

	go func() {
		ZSChan <- getZSData(ZSPillsBar, ZSRegionsBar)
	}()
	go func() {
		ARChan <- getARData(ARPillsBar, ARRegionsBar)
	}()
	go func() {
		EAChan <- getEAAllData(EAPillsBar, EARegionsBar)
	}()

	ZSData := <-ZSChan
	ARData := <-ARChan
	EAData := <-EAChan
	pool.Stop()

	log.Println("Данные успешно получены.")

	return map[string][]models.ParsedItem{
		"zdravcity": ZSData,
		"aptekaRu":  ARData,
		"eapteka":   EAData,
	}
}
