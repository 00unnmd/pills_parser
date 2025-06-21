package domain

import (
	"math"
	"strings"
)

type arPhoto struct {
	Original     string
	Medium       string
	Small        string
	Preview      string
	Png          string
	OriginalWebp string
	MediumWebp   string
	SmallWebp    string
	PreviewWebp  string
}

type arUniqueItemInfo struct {
	Id         string
	GoodNaming struct {
		TradeName           string
		Dosage              float64
		DosageStr           string
		DosageAdd           int
		Packing             float64
		PrimaryPackingShort string
		PrimaryPackingFull  string
		FormReleaseShort    string
		FormReleaseFull     string
		FormReleaseAdd      struct{}
		Features            []struct{}
		CustomFields        []struct {
			Key   string
			Value string
		}
	}
}

type arItemVariantInfo struct {
	Id           string
	Name         string
	HumanableUrl string
}

type ARResultItem struct {
	GroupRating            float64
	MaxGroupReviews        int
	NotGeneric             bool
	PrescriptionDrug       bool
	RecipeInPh             bool
	Photos                 []arPhoto
	HumanableUrl           string
	Id                     string
	InterNames             []string
	ItemsCount             int
	UniqueItemInfo         arUniqueItemInfo
	Disclaimer             string
	ItemVariantsInfo       []arItemVariantInfo
	TradeName              string
	Vendor                 string
	MaxPromoVits           int
	Brand                  struct{}
	Line                   struct{}
	FileInsts              struct{}
	ItemsInCart            struct{}
	Fund                   bool
	EDrug                  bool
	IsCourse               bool
	BuyRating              int
	ReviewRating           int
	ViewRating             int
	ThermoLab              bool
	IsInFavorites          bool
	AvailableInSets        bool
	VideoUrl               string
	VideoPreview           string
	Discount               bool
	IndexType              string
	MinPrice               float64
	LastMinPrice           float64
	Profit                 float64
	NoDiscPrice            float64
	VitaminsToBeCredited   int
	DiscountPercent        int
	HasProgressiveDiscount bool
	NoApplyDiscSelf        bool
	Gift                   bool
	ItemGroupType          string
	Incoming               string
	IncomingString         string
	Type                   string
	TypeKey                string
}

func (p ARResultItem) GetProducer() string {
	return p.Vendor
}

func (p ARResultItem) GetFields(pharmacy string, region string, pill string) ParsedItem {
	mnn := strings.Join(p.InterNames, ", ")

	return ParsedItem{
		Pharmacy:        pharmacy,
		Region:          region,
		Name:            p.TradeName,
		Mnn:             mnn,
		Price:           int(math.Round(p.MinPrice)),
		Discount:        int(math.Round(p.Profit)),
		DiscountPercent: p.DiscountPercent,
		Producer:        p.Vendor,
		Rating:          p.GroupRating,
		ReviewsCount:    p.ReviewRating,
		SearchValue:     pill,
	}
}

type arMultiValue struct {
	Url       string
	Name      string
	HumanName string
	Used      bool
	HasGoods  bool
}

type arAttribute struct {
	Name         string
	AttributeUrl string
	Type         string
	NumValues    struct{}
	MultiValues  []arMultiValue
}

type arSubCategory struct {
	Uid         string
	Name        string
	Url         string
	SubCategory struct{}
	ItemsCount  int
	Photo       arPhoto
}

type arCategory struct {
	Uid         string
	Name        string
	Url         string
	ItemsCount  int
	SubCategory arSubCategory
	Photo       arPhoto
}

type ARSearchBody struct {
	Took          int
	EsTook        int
	Page          int
	PageSize      int
	CurrentCount  int
	TotalCount    int
	Result        []ARResultItem
	Attributes    []arAttribute
	GroupInfo     struct{}
	Categories    []arCategory
	Suggestion    string
	MinGroupPrice float32
	AnalogGoods   []struct{}
}
