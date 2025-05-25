package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/00unnmd/pills_parser/internals/utils"
	"github.com/00unnmd/pills_parser/models"
	"net/http"
	"os"
	"path/filepath"
)

type PillsHandler struct {
	db *sql.DB
}

func NewPillsHandler(db *sql.DB) *PillsHandler {
	return &PillsHandler{db: db}
}

func (h *PillsHandler) GetZSPills(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	page, perPage, sortField, sortOrder, filter := utils.GetPillsReqFormValues(r)

	query := "SELECT id, region, name, price, discount, priceOld, maxQuantity, producer, rating, reviewsCount, error FROM zdravcity"
	var whereClause string
	var args []interface{}
	if filter != "" {
		whereClause = " WHERE name ILIKE $1"
		args = append(args, "%"+filter+"%")
	}
	query += whereClause
	query += fmt.Sprintf(" ORDER BY %s %s", sortField, sortOrder)
	offset := (page - 1) * perPage
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", perPage, offset)

	rows, err := h.db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var data []models.ParsedItem

	for rows.Next() {
		var item models.ParsedItem
		err := rows.Scan(&item.Id, &item.Region, &item.Name, &item.Price, &item.Discount, &item.PriceOld, &item.MaxQuantity, &item.Producer, &item.Rating, &item.ReviewsCount, &item.Error)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data = append(data, item)
	}

	var totalCount int
	countQuery := "SELECT COUNT(*) FROM zdravcity" + whereClause
	if err := h.db.QueryRow(countQuery, args...).Scan(&totalCount); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":  data,
		"total": totalCount,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *PillsHandler) GetARPills(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	page, perPage, sortField, sortOrder, filter := utils.GetPillsReqFormValues(r)

	query := "SELECT id, region, name, price, discount, priceOld, maxQuantity, producer, rating, reviewsCount, error FROM aptekaru"
	var whereClause string
	var args []interface{}
	if filter != "" {
		whereClause = " WHERE name ILIKE $1"
		args = append(args, "%"+filter+"%")
	}
	query += whereClause
	query += fmt.Sprintf(" ORDER BY %s %s", sortField, sortOrder)
	offset := (page - 1) * perPage
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", perPage, offset)

	rows, err := h.db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var data []models.ParsedItem

	for rows.Next() {
		var item models.ParsedItem
		err := rows.Scan(&item.Id, &item.Region, &item.Name, &item.Price, &item.Discount, &item.PriceOld, &item.MaxQuantity, &item.Producer, &item.Rating, &item.ReviewsCount, &item.Error)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data = append(data, item)
	}

	var totalCount int
	countQuery := "SELECT COUNT(*) FROM aptekaru" + whereClause
	if err := h.db.QueryRow(countQuery, args...).Scan(&totalCount); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":  data,
		"total": totalCount,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *PillsHandler) GetEAPills(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	page, perPage, sortField, sortOrder, filter := utils.GetPillsReqFormValues(r)

	query := "SELECT id, region, name, price, discount, priceOld, maxQuantity, producer, rating, reviewsCount, error FROM eapteka"
	var whereClause string
	var args []interface{}
	if filter != "" {
		whereClause = " WHERE name ILIKE $1"
		args = append(args, "%"+filter+"%")
	}
	query += whereClause
	query += fmt.Sprintf(" ORDER BY %s %s", sortField, sortOrder)
	offset := (page - 1) * perPage
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", perPage, offset)

	rows, err := h.db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var data []models.ParsedItem

	for rows.Next() {
		var item models.ParsedItem
		err := rows.Scan(&item.Id, &item.Region, &item.Name, &item.Price, &item.Discount, &item.PriceOld, &item.MaxQuantity, &item.Producer, &item.Rating, &item.ReviewsCount, &item.Error)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data = append(data, item)
	}

	var totalCount int
	countQuery := "SELECT COUNT(*) FROM eapteka" + whereClause
	if err := h.db.QueryRow(countQuery, args...).Scan(&totalCount); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":  data,
		"total": totalCount,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *PillsHandler) ExportPillsXLSX(w http.ResponseWriter, r *http.Request) {
	dirPath := filepath.Join("result")

	latestFile, err := utils.FindLatestParsingFile(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Export file not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed finding export file", http.StatusInternalServerError)
		}
		return
	}

	file, err := os.Open(latestFile)
	if err != nil {
		http.Error(w, "Failed to open export file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	fileName := filepath.Base(latestFile)
	if err != nil {
		http.Error(w, "Failed to load file info", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	http.ServeContent(w, r, fileName, fileInfo.ModTime(), file)
}
