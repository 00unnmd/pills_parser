package models

type FiltersSimple struct {
	Id  string `json:"id"`
	Val string `json:"val"`
}

type SearchFilters struct {
	Simple [1]FiltersSimple `json:"simple"`
}

type SearchPaginator struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type VariablesSearch struct {
	AdvertisementKey string          `json:"advertisementKey"`
	Filters          SearchFilters   `json:"filters"`
	Paginator        SearchPaginator `json:"paginator"`
	RegionID         string          `json:"regionID"`
	Query            string          `json:"query"`
}

type RequestBodyVariables struct {
	Search            VariablesSearch `json:"search"`
	RegionID          string          `json:"regionID"`
	AdvertisementType string          `json:"advertisementType"`
	Query             string          `json:"query"`
	SkipFeaturing     bool            `json:"skipFeaturing"`
}

type RequestBody struct {
	OperationName string               `json:"operationName"`
	Query         string               `json:"query"`
	Variables     RequestBodyVariables `json:"variables"`
}
