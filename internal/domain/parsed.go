package domain

type ParsedItem struct {
	Id              string  `json:"id"`
	Pharmacy        string  `json:"pharmacy"`
	Region          string  `json:"region"`
	Name            string  `json:"name"`
	Mnn             string  `json:"mnn"`
	Price           int     `json:"price"`
	Discount        int     `json:"discount"`
	DiscountPercent int     `json:"discountPercent"`
	Producer        string  `json:"producer"`
	Rating          float64 `json:"rating"`
	ReviewsCount    int     `json:"reviewsCount"`
	SearchValue     string  `json:"searchValue"`
	CreatedAt       string  `json:"createdAt"`
	Error           string  `json:"error"`
}

type ParsedFieldsGetter interface {
	GetFields(pharmacy string, region string, pill string) ParsedItem
}

type ProducerGetter interface {
	GetProducer() string
}
