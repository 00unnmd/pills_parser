package service

import (
	"bytes"
	"fmt"
	"github.com/00unnmd/pills_parser/internal/domain"
	"github.com/00unnmd/pills_parser/pkg/xlsx"
	"net/http"
	"time"
)

func (h *PillsHandler) ExportPillsXLSX(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid query parameters", http.StatusBadRequest)
		return
	}

	params, tableName, err := parseExportReqQuery(r)
	if err != nil {
		http.Error(w, "parseExportReqQuery err: "+err.Error(), http.StatusInternalServerError)
		return
	}

	queryBuilder := NewTableQueryBuilder(tableName, params).
		Select("createdAt, pharmacy, region, name, mnn, price, discount, discountPercent, producer, rating, reviewsCount, searchValue, error").
		// not calling pagination
		ApplySorting().
		ApplyFilters()

	data, err := h.executeExportQuery(queryBuilder.Build(), queryBuilder.Args...)
	if err != nil {
		http.Error(w, "executeExportQuery err: "+err.Error(), http.StatusInternalServerError)
		return
	}

	xlsxFile, err := xlsx.GenerateXLSX(data, tableName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf, err := xlsxFile.WriteToBuffer()
	if err != nil {
		http.Error(w, "failed writing xlsx to buffer", http.StatusInternalServerError)
		return
	}

	modTime := time.Now()
	fileName := fmt.Sprintf("%s-%s.xlsx", tableName, modTime.Format("02.01.2006.1504"))

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))

	http.ServeContent(w, r, fileName, modTime, bytes.NewReader(buf.Bytes()))
}

func parseExportReqQuery(r *http.Request) (*domain.TableReqParams, string, error) {
	values := r.URL.Query()
	params := &domain.TableReqParams{}

	// pagination not needed in export because service download all data by filters
	params.Page = 1
	params.PerPage = 10

	params.CreatedAt = values["createdAt"]
	params.Pharmacy = values["pharmacy"]
	params.Region = values["region"]

	params.Sort = values.Get("sort")
	params.Order = values.Get("order")

	params.Name = values.Get("name")
	params.Mnn = values.Get("mnn")
	params.Producer = values.Get("producer")
	params.SearchValue = values.Get("searchValue")

	tableName := values.Get("tableName")
	if tableName != "ozon" && tableName != "mnn" && tableName != "competitors" {
		return nil, "", fmt.Errorf("invalid tableName(resource)")
	} else {
		tableName = tableName + "_data"
	}

	for _, dateStr := range params.CreatedAt {
		if _, err := time.Parse(time.RFC3339, dateStr); err != nil {
			return nil, "", fmt.Errorf("invalid date format: %v", dateStr)
		}
	}
	if len(params.Pharmacy) > 3 || len(params.CreatedAt) > 10 || len(params.Region) > 8 {
		return nil, "", fmt.Errorf("too many values in pharmacy || createdAt || region")
	}

	return params, tableName, nil
}

func (h *PillsHandler) executeExportQuery(query string, args ...interface{}) ([]domain.ParsedItem, error) {
	rows, err := h.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data = make([]domain.ParsedItem, 0)
	for rows.Next() {
		var item domain.ParsedItem

		err := rows.Scan(&item.CreatedAt, &item.Pharmacy, &item.Region, &item.Name, &item.Mnn, &item.Price, &item.Discount, &item.DiscountPercent, &item.Producer, &item.Rating, &item.ReviewsCount, &item.SearchValue, &item.Error)
		if err != nil {
			return nil, err
		}

		data = append(data, item)
	}

	return data, nil
}
