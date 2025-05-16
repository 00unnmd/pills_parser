package models

type zsProductItem struct {
	Region       string
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Discount     float64 `json:"discount"`
	PriceOld     float64 `json:"priceOld"`
	MaxQuantity  int     `json:"maxQuantity"`
	Producer     string  `json:"producer"`
	IsBundle     bool    `json:"isBundle"`
	Rating       float64 `json:"rating"`
	ReviewsCount int     `json:"reviewsCount"`
}

func (p zsProductItem) GetProducer() string {
	return p.Producer
}

func (p zsProductItem) GetFields() ParsedItem {
	return ParsedItem{
		// Region
		Id:           p.Id,
		Name:         p.Name,
		Price:        p.Price,
		Discount:     p.Discount,
		PriceOld:     p.PriceOld,
		MaxQuantity:  p.MaxQuantity,
		Producer:     p.Producer,
		Rating:       p.Rating,
		ReviewsCount: p.ReviewsCount,
	}
}

type ZSSearchBody struct {
	Data struct {
		Products struct {
			Items []zsProductItem `json:"items"`
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
