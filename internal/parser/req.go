package parser

import (
	"github.com/00unnmd/pills_parser/internal/domain"
	"github.com/00unnmd/pills_parser/internal/parser/calls"
	"github.com/00unnmd/pills_parser/pkg/utils"
	"log"
	"time"

	"github.com/cheggaaa/pb/v3"
)

func getZSData(pillsList map[int]string, pillsBar *pb.ProgressBar, regionsBar *pb.ProgressBar, ozonFilter bool) []domain.ParsedItem {
	var result []domain.ParsedItem

	for _, regionItem := range utils.RegionsList {
		pillsBar.SetCurrent(0)
		regionsBar.Set("prefix", regionItem.Value)

		for _, pillValue := range pillsList {
			pillsBar.Set("prefix", "Zdravcity: "+pillValue)
			utils.RunRandomReqDelay(3, 4)

			pillsForRegion, err := calls.GetZSPills(pillValue, regionItem.ZSKey, regionItem.Value, ozonFilter)
			if err != nil {
				pillsForRegion = utils.CreatePIWithError(pillValue, regionItem.Value, err, "zdravcity")
			}

			result = append(result, pillsForRegion...)
			pillsBar.Increment()
		}

		regionsBar.Increment()
	}

	return result
}

func getARData(pillsList map[int]string, pillsBar *pb.ProgressBar, regionsBar *pb.ProgressBar, ozonFilter bool) []domain.ParsedItem {
	var result []domain.ParsedItem

	for _, regionItem := range utils.RegionsList {
		pillsBar.SetCurrent(0)
		regionsBar.Set("prefix", regionItem.Value)

		_, err := calls.ChangeARRegion(regionItem.ARKey)
		if err != nil {
			log.Println("err: ", err)
			break
		}

		for _, value := range pillsList {
			pillsBar.Set("prefix", "AptekaRu: "+value)
			utils.RunRandomReqDelay(3, 4)

			pillsAllRegions, err := calls.GetARPills(value, regionItem.Value, ozonFilter)
			if err != nil {
				pillsAllRegions = utils.CreatePIWithError(value, regionItem.Value, err, "aptekaru")
			}

			result = append(result, pillsAllRegions...)
			pillsBar.Increment()
		}

		regionsBar.Increment()
	}

	return result
}

func getEAData(pillsList map[int]string, pillsBar *pb.ProgressBar, regionsBar *pb.ProgressBar, ozonFilter bool) []domain.ParsedItem {
	var result []domain.ParsedItem

	ctx, cancel, err := calls.CreateEAContext()
	if err != nil {
		log.Println("err: ", err)
		return result
	}
	defer cancel()

	for _, regionItem := range utils.RegionsList {
		pillsBar.SetCurrent(0)
		regionsBar.Set("prefix", regionItem.Value)

		_, err := calls.ChangeEARegion(ctx, regionItem.EAKey)
		if err != nil {
			log.Println("err: ", err)
			break
		}

		for _, value := range pillsList {
			pillsBar.Set("prefix", "EApteka: "+value)
			utils.RunRandomReqDelay(3, 4)

			pillsAllRegions, err := calls.GetEAPills(ctx, value, regionItem.EAKey, regionItem.Value, ozonFilter)
			if err != nil {
				pillsAllRegions = utils.CreatePIWithError(value, regionItem.Value, err, "eapteka")
			}

			result = append(result, pillsAllRegions...)
			pillsBar.Increment()
		}

		regionsBar.Increment()
	}

	return result
}

func GetOzonAllData() []domain.ParsedItem {
	log.Println("Процесс получения данных...")
	ZSPillsBar := pb.New(len(utils.PillsList.OZON)).SetRefreshRate(time.Second * 3)
	ZSRegionsBar := pb.New(len(utils.RegionsList)).SetRefreshRate(time.Second * 3)
	ARPillsBar := pb.New(len(utils.PillsList.OZON)).SetRefreshRate(time.Second * 3)
	ARRegionsBar := pb.New(len(utils.RegionsList)).SetRefreshRate(time.Second * 3)
	EAPillsBar := pb.New(len(utils.PillsList.OZON)).SetRefreshRate(time.Second * 3)
	EARegionsBar := pb.New(len(utils.RegionsList)).SetRefreshRate(time.Second * 3)

	ZSChan := make(chan []domain.ParsedItem)
	ARChan := make(chan []domain.ParsedItem)
	EAChan := make(chan []domain.ParsedItem)

	pool, _ := pb.StartPool(ZSPillsBar, ZSRegionsBar, ARPillsBar, ARRegionsBar, EAPillsBar, EARegionsBar)

	go func() {
		ZSChan <- getZSData(utils.PillsList.OZON, ZSPillsBar, ZSRegionsBar, true)
	}()
	go func() {
		ARChan <- getARData(utils.PillsList.OZON, ARPillsBar, ARRegionsBar, true)
	}()
	go func() {
		EAChan <- getEAData(utils.PillsList.OZON, EAPillsBar, EARegionsBar, true)
	}()

	ZSData := <-ZSChan
	ARData := <-ARChan
	EAData := <-EAChan
	pool.Stop()

	log.Println("Данные успешно получены.")
	return append(append(ZSData, ARData...), EAData...)
}

func GetMNNAllData() []domain.ParsedItem {
	log.Println("Процесс получения данных...")
	ZSPillsBar := pb.New(len(utils.PillsList.MNN)).SetRefreshRate(time.Second * 3)
	ZSRegionsBar := pb.New(len(utils.RegionsList)).SetRefreshRate(time.Second * 3)
	ARPillsBar := pb.New(len(utils.PillsList.MNN)).SetRefreshRate(time.Second * 3)
	ARRegionsBar := pb.New(len(utils.RegionsList)).SetRefreshRate(time.Second * 3)
	EAPillsBar := pb.New(len(utils.PillsList.MNN)).SetRefreshRate(time.Second * 3)
	EARegionsBar := pb.New(len(utils.RegionsList)).SetRefreshRate(time.Second * 3)

	ZSChan := make(chan []domain.ParsedItem)
	ARChan := make(chan []domain.ParsedItem)
	EAChan := make(chan []domain.ParsedItem)

	pool, _ := pb.StartPool(ZSPillsBar, ZSRegionsBar, ARPillsBar, ARRegionsBar, EAPillsBar, EARegionsBar)

	go func() {
		ZSChan <- getZSData(utils.PillsList.MNN, ZSPillsBar, ZSRegionsBar, false)
	}()
	go func() {
		ARChan <- getARData(utils.PillsList.MNN, ARPillsBar, ARRegionsBar, false)
	}()
	go func() {
		EAChan <- getEAData(utils.PillsList.MNN, EAPillsBar, EARegionsBar, false)
	}()

	ZSData := <-ZSChan
	ARData := <-ARChan
	EAData := <-EAChan
	pool.Stop()

	log.Println("Данные успешно получены.")
	return append(append(ZSData, ARData...), EAData...)
}

func GetCompetitorsAllData() []domain.ParsedItem {
	log.Println("Процесс получения данных...")
	ZSPillsBar := pb.New(len(utils.PillsList.Competitors)).SetRefreshRate(time.Second * 3)
	ZSRegionsBar := pb.New(len(utils.RegionsList)).SetRefreshRate(time.Second * 3)
	ARPillsBar := pb.New(len(utils.PillsList.Competitors)).SetRefreshRate(time.Second * 3)
	ARRegionsBar := pb.New(len(utils.RegionsList)).SetRefreshRate(time.Second * 3)
	EAPillsBar := pb.New(len(utils.PillsList.Competitors)).SetRefreshRate(time.Second * 3)
	EARegionsBar := pb.New(len(utils.RegionsList)).SetRefreshRate(time.Second * 3)

	ZSChan := make(chan []domain.ParsedItem)
	ARChan := make(chan []domain.ParsedItem)
	EAChan := make(chan []domain.ParsedItem)

	pool, _ := pb.StartPool(ZSPillsBar, ZSRegionsBar, ARPillsBar, ARRegionsBar, EAPillsBar, EARegionsBar)

	go func() {
		ZSChan <- getZSData(utils.PillsList.Competitors, ZSPillsBar, ZSRegionsBar, false)
	}()
	go func() {
		ARChan <- getARData(utils.PillsList.Competitors, ARPillsBar, ARRegionsBar, false)
	}()
	go func() {
		EAChan <- getEAData(utils.PillsList.Competitors, EAPillsBar, EARegionsBar, false)
	}()

	ZSData := <-ZSChan
	ARData := <-ARChan
	EAData := <-EAChan
	pool.Stop()

	log.Println("Данные успешно получены.")
	return append(append(ZSData, ARData...), EAData...)
}
