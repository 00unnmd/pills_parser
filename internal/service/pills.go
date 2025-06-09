package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/00unnmd/pills_parser/internal/domain"
	"net/http"
	"strconv"
	"time"
)

type PillsHandler struct {
	db *sql.DB
}

func NewPillsHandler(db *sql.DB) *PillsHandler {
	return &PillsHandler{db: db}
}

func (h *PillsHandler) GetOzonPills(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	params, err := parsePillsReqQuery(r)
	if err != nil {
		http.Error(w, "ParsePillsReqQuery err: "+err.Error(), http.StatusInternalServerError)
		return
	}

	queryBuilder := NewTableQueryBuilder("ozon_data", params).
		Select("id, pharmacy, region, name, mnn, price, discount, discountPercent, producer, rating, reviewsCount, searchValue, createdAt, error").
		ApplyPagination().
		ApplySorting().
		ApplyFilters()

	data, err := h.executePillsQuery(queryBuilder.Build(), queryBuilder.Args...)
	if err != nil {
		http.Error(w, "executePillsQuery err: "+err.Error(), http.StatusInternalServerError)
		return
	}

	totalCount, err := h.getTotalCount(queryBuilder.BuildCountQuery(), queryBuilder.Args...)
	if err != nil {
		http.Error(w, "getTotalCount err: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":  data,
		"total": totalCount,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *PillsHandler) GetMNNPills(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	params, err := parsePillsReqQuery(r)
	if err != nil {
		http.Error(w, "ParsePillsReqQuery err: "+err.Error(), http.StatusInternalServerError)
		return
	}

	queryBuilder := NewTableQueryBuilder("mnn_data", params).
		Select("id, pharmacy, region, name, mnn, price, discount, discountPercent, producer, rating, reviewsCount, searchValue, createdAt, error").
		ApplyPagination().
		ApplySorting().
		ApplyFilters()

	data, err := h.executePillsQuery(queryBuilder.Build(), queryBuilder.Args...)
	if err != nil {
		http.Error(w, "executePillsQuery err: "+err.Error(), http.StatusInternalServerError)
		return
	}

	totalCount, err := h.getTotalCount(queryBuilder.BuildCountQuery(), queryBuilder.Args...)
	if err != nil {
		http.Error(w, "getTotalCount err: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":  data,
		"total": totalCount,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *PillsHandler) GetCompetitorsPills(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	params, err := parsePillsReqQuery(r)
	if err != nil {
		http.Error(w, "ParsePillsReqQuery err: "+err.Error(), http.StatusInternalServerError)
		return
	}

	queryBuilder := NewTableQueryBuilder("competitors_data", params).
		Select("id, pharmacy, region, name, mnn, price, discount, discountPercent, producer, rating, reviewsCount, searchValue, createdAt, error").
		ApplyPagination().
		ApplySorting().
		ApplyFilters()

	data, err := h.executePillsQuery(queryBuilder.Build(), queryBuilder.Args...)
	if err != nil {
		http.Error(w, "executePillsQuery err: "+err.Error(), http.StatusInternalServerError)
		return
	}

	totalCount, err := h.getTotalCount(queryBuilder.BuildCountQuery(), queryBuilder.Args...)
	if err != nil {
		http.Error(w, "getTotalCount err: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":  data,
		"total": totalCount,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (h *PillsHandler) executePillsQuery(query string, args ...interface{}) ([]domain.ParsedItem, error) {
	rows, err := h.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data = make([]domain.ParsedItem, 0)
	for rows.Next() {
		var item domain.ParsedItem

		err := rows.Scan(&item.Id, &item.Pharmacy, &item.Region, &item.Name, &item.Mnn, &item.Price, &item.Discount, &item.DiscountPercent, &item.Producer, &item.Rating, &item.ReviewsCount, &item.SearchValue, &item.CreatedAt, &item.Error)
		if err != nil {
			return nil, err
		}

		data = append(data, item)
	}

	return data, nil
}

func (h *PillsHandler) getTotalCount(query string, args ...interface{}) (int, error) {
	var count int
	err := h.db.QueryRow(query, args...).Scan(&count)

	return count, err
}

func parsePillsReqQuery(r *http.Request) (*domain.TableReqParams, error) {
	values := r.URL.Query()
	params := &domain.TableReqParams{}

	if page := values.Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err != nil {
			return nil, fmt.Errorf("error parsing page param to int: %v", err)
		} else {
			params.Page = p
		}
	}
	if perPage := values.Get("perPage"); perPage != "" {
		if pp, err := strconv.Atoi(perPage); err != nil {
			return nil, fmt.Errorf("error parsing perPage param to int: %v", err)
		} else {
			params.PerPage = pp
		}
	}

	params.CreatedAt = values["createdAt"]
	params.Pharmacy = values["pharmacy"]
	params.Region = values["region"]

	params.Sort = values.Get("sort")
	params.Order = values.Get("order")

	params.Name = values.Get("name")
	params.Mnn = values.Get("mnn")
	params.Producer = values.Get("producer")
	params.SearchValue = values.Get("searchValue")

	for _, dateStr := range params.CreatedAt {
		if _, err := time.Parse(time.RFC3339, dateStr); err != nil {
			return nil, fmt.Errorf("invalid date format: %v", dateStr)
		}
	}
	if len(params.Pharmacy) > 3 || len(params.CreatedAt) > 10 || len(params.Region) > 8 {
		return nil, fmt.Errorf("too many values in pharmacy || createdAt || region")
	}

	return params, nil
}
