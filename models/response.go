package models

type PillsItem struct {
	Region       string
	Name         string `json:"name"`
	Price        int    `json:"price"`
	Discount     int    `json:"discount"`
	PriceOld     int    `json:"priceOld"`
	MaxQuantity  int    `json:"maxQuantity"`
	Producer     string `json:"producer"`
	IsBundle     bool   `json:"isBundle"`
	Rating       int    `json:"rating"`
	ReviewsCount int    `json:"reviewsCount"`
}

type ResponseBody struct {
	Data struct {
		Products struct {
			Items []PillsItem `json:"items"`
			Total int         `json:"total"`
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
