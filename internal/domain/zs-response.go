package domain

import (
	"strings"
)

type MNN struct {
	Title string `json:"title"`
}

type ZSProductItem struct {
	Region       string
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

func (p ZSProductItem) GetFields() ParsedItem {
	titles := make([]string, len(p.Mnns))
	for i, item := range p.Mnns {
		titles[i] = item.Title
	}

	return ParsedItem{
		Name:         p.Name,
		Mnn:          strings.Join(titles, ", "),
		Price:        p.Price,
		Discount:     p.Discount,
		PriceOld:     p.PriceOld,
		Producer:     p.Producer,
		Rating:       p.Rating,
		ReviewsCount: p.ReviewsCount,
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
