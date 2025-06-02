package parser

import (
	"github.com/00unnmd/pills_parser/internal/domain"
	"github.com/00unnmd/pills_parser/internal/parser/calls"
	"github.com/00unnmd/pills_parser/pkg/utils"
	"log"
	"time"

	"github.com/cheggaaa/pb/v3"
)

func getZSData(pillsList map[int]string, regionsList map[string]string, pillsBar *pb.ProgressBar, regionsBar *pb.ProgressBar, ozonFilter bool) []domain.ParsedItem {
	var result []domain.ParsedItem

	for _, pillValue := range pillsList {
		regionsBar.SetCurrent(0)
		pillsBar.Set("prefix", "ZdravCity: "+pillValue)

		for key, value := range regionsList {
			regionsBar.Set("prefix", value)
			time.Sleep(utils.RequestDelay)

			pillsForRegion, err := calls.GetZSPills(pillValue, key, value, ozonFilter)
			if err != nil {
				pillsForRegion = utils.CreatePIWithError(pillValue, value, err)
			}

			result = append(result, pillsForRegion...)
			regionsBar.Increment()
		}

		pillsBar.Increment()
	}

	return result
}

func getARData(pillsList map[int]string, regionsList map[string]string, pillsBar *pb.ProgressBar, regionsBar *pb.ProgressBar, ozonFilter bool) []domain.ParsedItem {
	var result []domain.ParsedItem

	for regionId, regionValue := range regionsList {
		pillsBar.SetCurrent(0)
		regionsBar.Set("prefix", regionValue)

		_, err := calls.ChangeARRegion(regionId)
		if err != nil {
			log.Println("err: ", err)
			break
		}

		for _, value := range pillsList {
			pillsBar.Set("prefix", "AptekaRu: "+value)
			time.Sleep(utils.RequestDelay)

			pillsAllRegions, err := calls.GetARPills(value, regionValue, ozonFilter)
			if err != nil {
				pillsAllRegions = utils.CreatePIWithError(value, regionValue, err)
			}

			result = append(result, pillsAllRegions...)
			pillsBar.Increment()
		}

		regionsBar.Increment()
	}

	return result
}

func getEAData(pillsList map[int]string, regionsList map[string]string, pillsBar *pb.ProgressBar, regionsBar *pb.ProgressBar, ozonFilter bool) []domain.ParsedItem {
	var result []domain.ParsedItem

	ctx, cancel, err := calls.CreateEAContext()
	if err != nil {
		log.Println("err: ", err)
		return result
	}
	defer cancel()

	for regionKey, regionValue := range regionsList {
		pillsBar.SetCurrent(0)
		regionsBar.Set("prefix", regionValue)

		_, err := calls.ChangeEARegion(ctx, regionKey)
		if err != nil {
			log.Println("err: ", err)
			break
		}

		for _, value := range pillsList {
			pillsBar.Set("prefix", "EApteka: "+value)
			time.Sleep(utils.RequestDelay)

			pillsAllRegions, err := calls.GetEAPills(ctx, value, regionKey, regionValue, ozonFilter)
			if err != nil {
				pillsAllRegions = utils.CreatePIWithError(value, regionValue, err)
			}

			result = append(result, pillsAllRegions...)
			pillsBar.Increment()
		}

		regionsBar.Increment()
	}

	return result
}

func GetOzonAllData() map[string][]domain.ParsedItem {
	log.Println("Процесс получения данных...")
	ZSPillsBar := pb.New(len(utils.OzonPillsList)).SetRefreshRate(time.Second * 3)
	ZSRegionsBar := pb.New(len(utils.ZSRegions)).SetRefreshRate(time.Second * 3)
	ARPillsBar := pb.New(len(utils.OzonPillsList)).SetRefreshRate(time.Second * 3)
	ARRegionsBar := pb.New(len(utils.ARRegions)).SetRefreshRate(time.Second * 3)
	EAPillsBar := pb.New(len(utils.OzonPillsList)).SetRefreshRate(time.Second * 3)
	EARegionsBar := pb.New(len(utils.EARegions)).SetRefreshRate(time.Second * 3)

	ZSChan := make(chan []domain.ParsedItem)
	ARChan := make(chan []domain.ParsedItem)
	EAChan := make(chan []domain.ParsedItem)

	pool, _ := pb.StartPool(ZSPillsBar, ZSRegionsBar, ARPillsBar, ARRegionsBar, EAPillsBar, EARegionsBar)

	go func() {
		ZSChan <- getZSData(utils.OzonPillsList, utils.ZSRegions, ZSPillsBar, ZSRegionsBar, true)
	}()
	go func() {
		ARChan <- getARData(utils.OzonPillsList, utils.ARRegions, ARPillsBar, ARRegionsBar, true)
	}()
	go func() {
		EAChan <- getEAData(utils.OzonPillsList, utils.EARegions, EAPillsBar, EARegionsBar, true)
	}()

	ZSData := <-ZSChan
	ARData := <-ARChan
	EAData := <-EAChan
	pool.Stop()

	log.Println("Данные успешно получены.")

	return map[string][]domain.ParsedItem{
		"zdravcity": ZSData,
		"aptekaRu":  ARData,
		"eapteka":   EAData,
	}
}

func GetMNNAllData() map[string][]domain.ParsedItem {
	log.Println("Процесс получения данных...")
	ZSPillsBar := pb.New(len(utils.CompetitorsMNNList)).SetRefreshRate(time.Second * 3)
	ZSRegionsBar := pb.New(len(utils.ZSRegions)).SetRefreshRate(time.Second * 3)
	ARPillsBar := pb.New(len(utils.CompetitorsMNNList)).SetRefreshRate(time.Second * 3)
	ARRegionsBar := pb.New(len(utils.ARRegions)).SetRefreshRate(time.Second * 3)
	EAPillsBar := pb.New(len(utils.CompetitorsMNNList)).SetRefreshRate(time.Second * 3)
	EARegionsBar := pb.New(len(utils.EARegions)).SetRefreshRate(time.Second * 3)

	ZSChan := make(chan []domain.ParsedItem)
	ARChan := make(chan []domain.ParsedItem)
	EAChan := make(chan []domain.ParsedItem)

	pool, _ := pb.StartPool(ZSPillsBar, ZSRegionsBar, ARPillsBar, ARRegionsBar, EAPillsBar, EARegionsBar)

	go func() {
		ZSChan <- getZSData(utils.CompetitorsMNNList, utils.ZSRegions, ZSPillsBar, ZSRegionsBar, false)
	}()
	go func() {
		ARChan <- getARData(utils.CompetitorsMNNList, utils.ARRegions, ARPillsBar, ARRegionsBar, false)
	}()
	go func() {
		EAChan <- getEAData(utils.CompetitorsMNNList, utils.EARegions, EAPillsBar, EARegionsBar, false)
	}()

	ZSData := <-ZSChan
	ARData := <-ARChan
	EAData := <-EAChan
	pool.Stop()

	log.Println("Данные успешно получены.")

	return map[string][]domain.ParsedItem{
		"zdravcity": ZSData,
		"aptekaRu":  ARData,
		"eapteka":   EAData,
	}
}

func GetCompetitorsAllData() map[string][]domain.ParsedItem {
	log.Println("Процесс получения данных...")
	ZSPillsBar := pb.New(len(utils.CompetitorsPillsList)).SetRefreshRate(time.Second * 3)
	ZSRegionsBar := pb.New(len(utils.ZSRegions)).SetRefreshRate(time.Second * 3)
	ARPillsBar := pb.New(len(utils.CompetitorsPillsList)).SetRefreshRate(time.Second * 3)
	ARRegionsBar := pb.New(len(utils.ARRegions)).SetRefreshRate(time.Second * 3)
	EAPillsBar := pb.New(len(utils.CompetitorsPillsList)).SetRefreshRate(time.Second * 3)
	EARegionsBar := pb.New(len(utils.EARegions)).SetRefreshRate(time.Second * 3)

	ZSChan := make(chan []domain.ParsedItem)
	ARChan := make(chan []domain.ParsedItem)
	EAChan := make(chan []domain.ParsedItem)

	pool, _ := pb.StartPool(ZSPillsBar, ZSRegionsBar, ARPillsBar, ARRegionsBar, EAPillsBar, EARegionsBar)

	go func() {
		ZSChan <- getZSData(utils.CompetitorsPillsList, utils.ZSRegions, ZSPillsBar, ZSRegionsBar, false)
	}()
	go func() {
		ARChan <- getARData(utils.CompetitorsPillsList, utils.ARRegions, ARPillsBar, ARRegionsBar, false)
	}()
	go func() {
		EAChan <- getEAData(utils.CompetitorsPillsList, utils.EARegions, EAPillsBar, EARegionsBar, false)
	}()

	ZSData := <-ZSChan
	ARData := <-ARChan
	EAData := <-EAChan
	pool.Stop()

	log.Println("Данные успешно получены.")

	return map[string][]domain.ParsedItem{
		"zdravcity": ZSData,
		"aptekaRu":  ARData,
		"eapteka":   EAData,
	}
}
