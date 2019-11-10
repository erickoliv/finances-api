package rest

import (
	"fmt"
	"github.com/ericktm/olivsoft-golang-api/pkg/domain"
	"net/url"
	"strconv"
	"strings"
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

// ExtractFilters can be used to parse query parameters and return a QueryData object, useful to query, filter and paginate requests
func ExtractFilters(f url.Values) domain.QueryData {
	q := domain.QueryData{
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
