package service

import (
	"fmt"
	"github.com/00unnmd/pills_parser/internal/domain"
	"strings"
	"time"
)

type TableQueryBuilder struct {
	Table      string
	SelectCols string
	Where      []string
	Args       []interface{}
	Sort       string
	Order      string
	Limit      int
	Offset     int
	Params     *domain.TableReqParams
}

func NewTableQueryBuilder(table string, params *domain.TableReqParams) *TableQueryBuilder {
	return &TableQueryBuilder{
		Table:  table,
		Params: params,
	}
}

func (qb *TableQueryBuilder) Select(cols string) *TableQueryBuilder {
	qb.SelectCols = cols
	return qb
}

func (qb *TableQueryBuilder) ApplyPagination() *TableQueryBuilder {
	qb.Limit = qb.Params.PerPage
	qb.Offset = (qb.Params.Page - 1) * qb.Params.PerPage

	return qb
}

func (qb *TableQueryBuilder) ApplySorting() *TableQueryBuilder {
	if qb.Params.Sort != "" {
		qb.Sort = qb.Params.Sort
		qb.Order = qb.Params.Order
	}

	return qb
}

func (qb *TableQueryBuilder) ApplyFilters() *TableQueryBuilder {
	params := qb.Params
	if len(params.CreatedAt) > 0 {
		parsedDates := make([]string, 0, len(params.CreatedAt))
		for _, ca := range params.CreatedAt {
			t, _ := time.Parse(time.RFC3339, ca)
			dateStr := t.Format("2006-01-02")
			parsedDates = append(parsedDates, dateStr)
		}

		qb.Where = append(qb.Where, "DATE(createdAt) IN (?"+strings.Repeat(",?", len(parsedDates)-1)+")")
		for _, pd := range parsedDates {
			qb.Args = append(qb.Args, pd)
		}
	}
	if len(params.Pharmacy) > 0 {
		qb.Where = append(qb.Where, "pharmacy IN (?"+strings.Repeat(",?", len(params.Pharmacy)-1)+")")
		for _, ph := range params.Pharmacy {
			qb.Args = append(qb.Args, ph)
		}
	}
	if len(params.Region) > 0 {
		qb.Where = append(qb.Where, "region IN (?"+strings.Repeat(",?", len(params.Region)-1)+")")
		for _, rg := range params.Region {
			qb.Args = append(qb.Args, rg)
		}
	}

	if params.Name != "" {
		qb.Where = append(qb.Where, "name ILIKE ?")
		qb.Args = append(qb.Args, "%"+params.Name+"%")
	}
	if params.Mnn != "" {
		qb.Where = append(qb.Where, "mnn ILIKE ?")
		qb.Args = append(qb.Args, "%"+params.Mnn+"%")
	}
	if params.Producer != "" {
		qb.Where = append(qb.Where, "producer ILIKE ?")
		qb.Args = append(qb.Args, "%"+params.Producer+"%")
	}
	if params.SearchValue != "" {
		qb.Where = append(qb.Where, "searchValue ILIKE ?")
		qb.Args = append(qb.Args, "%"+params.SearchValue+"%")
	}

	return qb
}

func (qb *TableQueryBuilder) Build() string {
	query := fmt.Sprintf("SELECT %s FROM %s", qb.SelectCols, qb.Table)

	if len(qb.Where) > 0 {
		query += " WHERE " + strings.Join(qb.Where, " AND ")
	}
	if qb.Sort != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", qb.Sort, qb.Order)
	}
	if qb.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", qb.Limit, qb.Offset)
	}

	for i := 1; i <= len(qb.Args); i++ {
		query = strings.Replace(query, "?", fmt.Sprintf("$%d", i), 1)
	}

	return query
}

func (qb *TableQueryBuilder) BuildCountQuery() string {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", qb.Table)

	if len(qb.Where) > 0 {
		query += " WHERE " + strings.Join(qb.Where, " AND ")
	}

	for i := 1; i <= len(qb.Args); i++ {
		query = strings.Replace(query, "?", fmt.Sprintf("$%d", i), 1)
	}

	return query
}
