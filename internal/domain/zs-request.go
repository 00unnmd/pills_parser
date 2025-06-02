package domain

type ZSSimple struct {
	Id  string `json:"id"`
	Val string `json:"val"`
}

type ZSSearchFilters struct {
	Simple [1]ZSSimple `json:"simple"`
}

type ZSSearchPaginator struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ZSSearch struct {
	AdvertisementKey string            `json:"advertisementKey"`
	Filters          ZSSearchFilters   `json:"filters"`
	Paginator        ZSSearchPaginator `json:"paginator"`
	RegionID         string            `json:"regionID"`
	Query            string            `json:"query"`
}

type ZSRequestBodyVariables struct {
	Search            ZSSearch `json:"search"`
	RegionID          string   `json:"regionID"`
	AdvertisementType string   `json:"advertisementType"`
	Query             string   `json:"query"`
	SkipFeaturing     bool     `json:"skipFeaturing"`
}

type ZSRequestBody struct {
	OperationName string                 `json:"operationName"`
	Query         string                 `json:"query"`
	Variables     ZSRequestBodyVariables `json:"variables"`
}
