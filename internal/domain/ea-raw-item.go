package domain

type EARawItem struct {
	Name            string
	Mnn             string
	Price           int
	Discount        int
	DiscountPercent int
	Producer        string
	Rating          float64
	ReviewsCount    int
	Error           string
}

func (p EARawItem) GetProducer() string {
	return p.Producer
}

func (p EARawItem) GetFields(pharmacy string, region string, pill string) ParsedItem {
	return ParsedItem{
		Pharmacy:        pharmacy,
		Region:          region,
		Name:            p.Name,
		Mnn:             p.Mnn,
		Price:           p.Price,
		Discount:        p.Discount,
		DiscountPercent: p.DiscountPercent,
		Producer:        p.Producer,
		Rating:          p.Rating,
		ReviewsCount:    p.ReviewsCount,
		SearchValue:     pill,
	}
}
