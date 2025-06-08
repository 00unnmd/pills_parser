package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/00unnmd/pills_parser/internal/domain"
	"github.com/00unnmd/pills_parser/pkg/utils"
	"net/http"
	"time"
)

type OptionsHandler struct {
	db *sql.DB
}

func NewOptionsHandler(db *sql.DB) *OptionsHandler {
	return &OptionsHandler{db: db}
}

type GetOptionsResponseItem struct {
	Regions    []domain.OptionItem `json:"regions"`
	Dates      []domain.OptionItem `json:"dates"`
	Pharmacies []domain.OptionItem `json:"pharmacies"`
}

func getRegions(regionsList []utils.Region) []domain.OptionItem {
	result := make([]domain.OptionItem, 0)

	for _, regionItem := range regionsList {
		result = append(result, domain.OptionItem{
			Id:   regionItem.Value,
			Name: regionItem.Value,
		})
	}

	return result
}

func getPharmacies() []domain.OptionItem {
	result := make([]domain.OptionItem, 0)

	for key, value := range utils.PharmaciesList {
		result = append(result, domain.OptionItem{
			Id:   key,
			Name: value,
		})
	}

	return result
}

func (h *OptionsHandler) getUniqueDates() (map[string][]domain.OptionItem, error) {
	query := `
        SELECT 'ozon' AS table_name, DATE(createdAt) AS date FROM ozon_data
        UNION SELECT 'mnn', DATE(createdAt) FROM mnn_data
        UNION SELECT 'competitors', DATE(createdAt) FROM competitors_data
        ORDER BY table_name, date DESC
    `

	rows, err := h.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("getUniqueDates query error %w", err)
	}
	defer rows.Close()

	datesMap := make(map[string][]domain.OptionItem)
	for rows.Next() {
		var tableName, dateStr string
		if err := rows.Scan(&tableName, &dateStr); err != nil {
			return nil, fmt.Errorf("getUniqueDates row scanning error %w", err)
		}

		t, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return nil, fmt.Errorf("getUniqueDates parsing date err: %w", err)
		}
		datesMap[tableName] = append(datesMap[tableName], domain.OptionItem{
			Id:   dateStr,
			Name: t.Format("02.01.2006"),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getUniqueDates row iteration err: %w", err)
	}

	return datesMap, nil
}

func (h *OptionsHandler) GetOptions(w http.ResponseWriter, r *http.Request) {
	ozonRegions := getRegions(utils.RegionsList)
	competitorsRegions := getRegions([]utils.Region{utils.RegionsList[0]})
	pharmacies := getPharmacies()

	uniqueDates, err := h.getUniqueDates()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get dates: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]GetOptionsResponseItem{
		"ozon": {
			Regions:    ozonRegions,
			Dates:      uniqueDates["ozon"],
			Pharmacies: pharmacies,
		},
		"mnn": {
			Regions:    competitorsRegions,
			Dates:      uniqueDates["mnn"],
			Pharmacies: pharmacies,
		},
		"competitors": {
			Regions:    competitorsRegions,
			Dates:      uniqueDates["competitors"],
			Pharmacies: pharmacies,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
