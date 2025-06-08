package domain

type TableReqParams struct {
	Page        int      `json:"page"`
	PerPage     int      `json:"perPage"`
	Sort        string   `json:"sort"`
	Order       string   `json:"order"`
	CreatedAt   []string `json:"createdAt"`
	Pharmacy    []string `json:"pharmacy"`
	Region      []string `json:"region"`
	Name        string   `json:"name"`
	Mnn         string   `json:"mnn"`
	Producer    string   `json:"producer"`
	SearchValue string   `json:"searchValue"`
}
