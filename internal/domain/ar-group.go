package domain

import (
	"math"
	"strings"
	"time"
)

type ARItemInfo struct {
	ID               string      `json:"id"`
	SourceUID        string      `json:"sourceUid"`
	Name             string      `json:"name"`
	Vendor           string      `json:"vendor"`
	VendorBespoke    interface{} `json:"vendorBespoke"`
	Brand            interface{} `json:"brand"`
	Line             interface{} `json:"line"`
	Indic            string      `json:"indic"`
	ItemRankValue    float64     `json:"itemRankValue"`
	Fund             bool        `json:"fund"`
	NotGeneric       bool        `json:"notGeneric"`
	PromoVits        int         `json:"promoVits"`
	PrescriptionDrug bool        `json:"prescriptionDrug"`
	RecipeInPh       bool        `json:"recipeInPh"`
	SaleLimit        bool        `json:"saleLimit"`
	PhotoPack        struct {
		Photos []struct {
			Original     string      `json:"original"`
			Medium       interface{} `json:"medium"`
			Small        interface{} `json:"small"`
			Preview      interface{} `json:"preview"`
			Png          interface{} `json:"png"`
			OriginalWebp interface{} `json:"originalWebp"`
			MediumWebp   interface{} `json:"mediumWebp"`
			SmallWebp    interface{} `json:"smallWebp"`
			PreviewWebp  interface{} `json:"previewWebp"`
		} `json:"photos"`
	} `json:"photoPack"`
	Category struct {
		Name        string      `json:"name"`
		URL         string      `json:"url"`
		Ccgs        interface{} `json:"ccgs"`
		SubCategory struct {
			Name        string      `json:"name"`
			URL         string      `json:"url"`
			Ccgs        string      `json:"ccgs"`
			SubCategory interface{} `json:"subCategory"`
			SeoName     interface{} `json:"seoName"`
		} `json:"subCategory"`
		SeoName interface{} `json:"seoName"`
	} `json:"category"`
	GoodGroupInfo struct {
		UID      string `json:"uid"`
		Name     string `json:"name"`
		URL      string `json:"url"`
		SubGroup struct {
			UID      string      `json:"uid"`
			Name     string      `json:"name"`
			URL      string      `json:"url"`
			SubGroup interface{} `json:"subGroup"`
		} `json:"subGroup"`
	} `json:"goodGroupInfo"`
	FileInst []struct {
		FilePath interface{} `json:"filePath"`
		Photos   struct {
			Original     string      `json:"original"`
			Medium       string      `json:"medium"`
			Small        string      `json:"small"`
			Preview      string      `json:"preview"`
			Png          interface{} `json:"png"`
			OriginalWebp string      `json:"originalWebp"`
			MediumWebp   string      `json:"mediumWebp"`
			SmallWebp    string      `json:"smallWebp"`
			PreviewWebp  string      `json:"previewWebp"`
		} `json:"photos"`
		AdInfo      interface{} `json:"adInfo"`
		OrdObjectID interface{} `json:"ordObjectId"`
	} `json:"fileInst"`
	VariantValues struct {
		NAMING_FAILED string `json:"В упаковке"`
	} `json:"variantValues"`
	InterNames  []string `json:"interNames"`
	TypedPrices []struct {
		Type     string      `json:"type"`
		TypeKey  interface{} `json:"typeKey"`
		PackInfo struct {
			PackCalcPrice float64 `json:"packCalcPrice"`
			Desc          string  `json:"desc"`
		} `json:"packInfo"`
		Incoming        interface{} `json:"incoming"`
		DbPrice         float64     `json:"dbPrice"`
		Rest            int         `json:"rest"`
		DiscountPercent int         `json:"discountPercent"`
		NoDiscPrice     float64     `json:"noDiscPrice"`
		DbProfit        float64     `json:"dbProfit"`
		AmountInCart    int         `json:"amountInCart"`
		LifeTime        time.Time   `json:"lifeTime"`
	} `json:"typedPrices"`
	Discount bool `json:"discount"`
	PackInfo struct {
		PackCalcPrice float64 `json:"packCalcPrice"`
		Desc          string  `json:"desc"`
	} `json:"packInfo"`
	Type           string      `json:"type"`
	TypeKey        interface{} `json:"typeKey"`
	Price          float64     `json:"price"`
	LastPrice      interface{} `json:"lastPrice"`
	NoDiscPrice    float64     `json:"noDiscPrice"`
	Profit         float64     `json:"profit"`
	Incoming       interface{} `json:"incoming"`
	IncomingString interface{} `json:"incomingString"`
	LifeTimeInfo   struct {
		Type        string `json:"type"`
		Description string `json:"description"`
		Hint        string `json:"hint"`
	} `json:"lifeTimeInfo"`
	AmountInCart           int         `json:"amountInCart"`
	AmountInGoodSet        int         `json:"amountInGoodSet"`
	HumanableURL           string      `json:"humanableUrl"`
	IsInFavorites          bool        `json:"isInFavorites"`
	NotifyAppearance       bool        `json:"notifyAppearance"`
	Default                bool        `json:"default"`
	EDrug                  bool        `json:"eDrug"`
	OutOfStock             bool        `json:"outOfStock"`
	IPhGoodSetsIds         interface{} `json:"iPhGoodSetsIds"`
	IPhGoodSetsDetails     interface{} `json:"iPhGoodSetsDetails"`
	IsCourse               bool        `json:"isCourse"`
	NoApplyDiscSelf        bool        `json:"noApplyDiscSelf"`
	HasProgressiveDiscount bool        `json:"hasProgressiveDiscount"`
	VitaminsToBeCredited   int         `json:"vitaminsToBeCredited"`
	DiscountPercent        int         `json:"discountPercent"`
	Rating                 float64     `json:"rating"`
	ReviewsCount           int         `json:"reviewsCount"`
	VideoDesc              interface{} `json:"videoDesc"`
}

func (p ARItemInfo) GetFields(pharmacy string, region string, pill string) ParsedItem {
	mnn := strings.Join(p.InterNames, ", ")

	return ParsedItem{
		Pharmacy:        pharmacy,
		Region:          region,
		Name:            p.Name,
		Mnn:             mnn,
		Price:           int(math.Round(p.Price)),
		Discount:        int(math.Round(p.Profit)),
		DiscountPercent: p.DiscountPercent,
		Producer:        p.Vendor,
		Rating:          p.Rating,
		ReviewsCount:    p.ReviewsCount,
		SearchValue:     pill,
	}
}

type arGroupItem struct {
	LevelDescription string       `json:"levelDescription"`
	LevelName        string       `json:"levelName"`
	LevelNameShort   string       `json:"levelNameShort"`
	LevelType        string       `json:"levelType"`
	ItemInfos        []ARItemInfo `json:"itemInfos"`
	Default          bool         `json:"default"`
	CommonProperties interface{}  `json:"commonProperties"`
}

type arCommonProperty struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	NameShort   string `json:"nameShort"`
	Type        string `json:"type"`
}

type ARGroupInfoBody struct {
	GroupType               string             `json:"groupType"`
	GroupItems              []arGroupItem      `json:"groupItems"`
	CommonProperties        []arCommonProperty `json:"commonProperties"`
	HumanableURL            string             `json:"humanableUrl"`
	Incoming                interface{}        `json:"incoming"`
	GoodSetItems            interface{}        `json:"goodSetItems"`
	GoodSetInfo             interface{}        `json:"goodSetInfo"`
	PreparationInfo         interface{}        `json:"preparationInfo"`
	ItemsTotalCount         int                `json:"itemsTotalCount"`
	HasProgressiveDiscount  bool               `json:"hasProgressiveDiscount"`
	ItemGroupName           string             `json:"itemGroupName"`
	NeedGroupForLargeScreen bool               `json:"needGroupForLargeScreen"`
	NeedGroupForSmallScreen bool               `json:"needGroupForSmallScreen"`
}
