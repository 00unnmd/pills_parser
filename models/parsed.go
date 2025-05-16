package models

type ParsedItem struct {
	Id           string  `json:"id"`
	Region       string  `json:"region"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Discount     float64 `json:"discount"`
	PriceOld     float64 `json:"priceOld"`
	MaxQuantity  int     `json:"maxQuantity"`
	Producer     string  `json:"producer"`
	Rating       float64 `json:"rating"`
	ReviewsCount int     `json:"reviewsCount"`
}

func (p ParsedItem) GetProducer() string {
	return p.Producer
}

type ParsedFieldsGetter interface {
	GetFields() ParsedItem
}

type ProducerGetter interface {
	GetProducer() string
}
