package api

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
)

// PaginatedMessage is a structure which contains standard attributes to be used on paginated services
type PaginatedMessage struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Pages int `json:"pages"`
	//Next string `json:"next"`
	//Previous string `json:"previous"`
	Limit int         `json:"limit"`
	Count int         `json:"count"`
	Data  interface{} `json:"data"`
}

// QueryData is structure which contains standard atributes to parse http parameters for API filter, search and pagination
type QueryData struct {
	Page    int
	Pages   int
	Total   int
	Limit   int
	Sort    string
	Filters map[string]interface{}
}

// ExtractFilters can be used to parse query parameters and return a QueryData object, useful to query, filter and paginate requests
func ExtractFilters(f url.Values) QueryData {
	q := QueryData{
		Page:  1,
		Limit: 100,
	}

	if limit, err := strconv.Atoi(f.Get("limit")); err == nil {
		q.Limit = limit
	}

	if page, err := strconv.Atoi(f.Get("page")); err == nil {
		q.Page = page
	}

	q.Sort = f.Get("sort")

	// TODO: Create Generic Midleware to put filters inside context
	filters := map[string]interface{}{}
	for key := range f {
		if strings.HasPrefix(key, "q_") {
			if strings.HasSuffix(key, "__like") {
				field := fmt.Sprintf("%s LIKE ?", key[2:len(key)-6])
				filters[field] = f.Get(key)
				continue
			}
			if strings.HasSuffix(key, "__eq") {
				field := fmt.Sprintf("%s = ?", key[2:len(key)-4])
				filters[field] = f.Get(key)
				continue
			}
			if strings.HasSuffix(key, "__gte") {
				field := fmt.Sprintf("%s >= ?", key[2:len(key)-5])
				filters[field] = f.Get(key)
				continue
			}
			if strings.HasSuffix(key, "__lte") {
				field := fmt.Sprintf("%s <= ?", key[2:len(key)-5])
				filters[field] = f.Get(key)
				continue
			}
		}
	}

	q.Filters = filters

	return q
}

// Build is the function where the QueryData attributes are translated into a gorm.DB instance. Can be used to generic filter, order and pagination
func (q *QueryData) Build(db *gorm.DB) *gorm.DB {
	base := db

	for k, v := range q.Filters {
		base = base.Where(k, v)
	}

	base.Count(&q.Total)
	base = base.Offset(q.Limit * (q.Page - 1)).Limit(q.Limit).Order(q.Sort)
	q.Pages = int(math.Ceil(float64(q.Total) / float64(q.Limit)))

	return base
}
