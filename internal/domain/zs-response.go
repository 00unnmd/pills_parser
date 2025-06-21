package domain

import (
	"math"
	"strings"
)

type MNN struct {
	Title string `json:"title"`
}

type ZSProductItem struct {
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Mnns         []MNN   `json:"mnns"`
	Price        float64 `json:"price"`
	Discount     float64 `json:"discount"`
	PriceOld     float64 `json:"priceOld"`
	Producer     string  `json:"producer"`
	Rating       float64 `json:"rating"`
	ReviewsCount int     `json:"reviewsCount"`
}

func (p ZSProductItem) GetProducer() string {
	return p.Producer
}

func (p ZSProductItem) GetFields(pharmacy string, region string, pill string) ParsedItem {
	titles := make([]string, len(p.Mnns))
	for i, item := range p.Mnns {
		titles[i] = item.Title
	}
	mnn := strings.Join(titles, ", ")

	discountPercent := 0
	if p.PriceOld != 0 {
		discountPercent = int(math.Round(p.Discount / (p.PriceOld / 100)))
	}

	return ParsedItem{
		Pharmacy:        pharmacy,
		Region:          region,
		Name:            p.Name,
		Mnn:             mnn,
		Price:           int(math.Round(p.Price)),
		Discount:        int(math.Round(p.Discount)),
		DiscountPercent: discountPercent,
		Producer:        p.Producer,
		Rating:          p.Rating,
		ReviewsCount:    p.ReviewsCount,
		SearchValue:     pill,
	}
}

type ZSSearchBody struct {
	Data struct {
		Products struct {
			Items []ZSProductItem `json:"items"`
			Total int             `json:"total"`
		} `json:"products"`
		Featuring []struct {
			Title      string `json:"title"`
			Banner     string `json:"banner"`
			Advertiser string `json:"advertiser"`
			Counters   struct {
				Imps     interface{} `json:"imps"`
				Track    interface{} `json:"track"`
				Creative string      `json:"creative"`
				ID       string      `json:"id"`
				Name     string      `json:"name"`
				Position string      `json:"position"`
			} `json:"counters"`
			Products []interface{} `json:"products"`
		} `json:"featuring"`
	} `json:"data"`
}
