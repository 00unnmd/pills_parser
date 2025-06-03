package domain

type ParsedItem struct {
	Id              string  `json:"id"`
	Pharmacy        string  `json:"pharmacy"`
	Region          string  `json:"region"`
	Name            string  `json:"name"`
	Mnn             string  `json:"mnn"`
	Price           float64 `json:"price"`
	PriceOld        float64 `json:"priceOld"`
	Discount        float64 `json:"discount"`
	DiscountPercent int     `json:"discountPercent"`
	Producer        string  `json:"producer"`
	Rating          float64 `json:"rating"`
	ReviewsCount    int     `json:"reviewsCount"`
	SearchValue     string  `json:"searchValue"`
	CreatedAt       string  `json:"createdAt"`
	Error           string  `json:"error"`
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
