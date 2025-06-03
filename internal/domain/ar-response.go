package domain

import (
	"strings"
	"time"
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

func (p ARItemInfo) GetFields() ParsedItem {
	return ParsedItem{
		Name:            p.Name,
		Mnn:             strings.Join(p.InterNames, ", "),
		Price:           p.Price,
		Discount:        p.Profit,
		DiscountPercent: p.DiscountPercent,
		Producer:        p.Vendor,
		Rating:          p.Rating,
		ReviewsCount:    p.ReviewsCount,
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
