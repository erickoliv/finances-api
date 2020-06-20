package rest

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ExtractFilters can be used to parse query parameters and return a Query object, useful to query, filter and paginate requests
func ExtractFilters(c *gin.Context, needUser bool) (*Query, error) {
	user, err := ExtractUser(c)
	if err != nil && needUser {
		return nil, err
	}
	filters := map[string]interface{}{}
	q := Query{
		Page:  1,
		Limit: 100,
	}
	if user != uuid.Nil {
		filters["owner = ?"] = user
	}

	f := c.Request.URL.Query()

	if limit, err := strconv.Atoi(f.Get("limit")); err == nil {
		q.Limit = limit
	}

	if page, err := strconv.Atoi(f.Get("page")); err == nil {
		q.Page = page
	}

	q.Sort = f.Get("sort")
	for key := range f {
		if !strings.HasPrefix(key, "q_") {
			continue
		}

		switch {
		case strings.HasSuffix(key, "__eq"):
			field := fmt.Sprintf("%s = ?", key[2:len(key)-4])
			filters[field] = f.Get(key)

		case strings.HasSuffix(key, "__like"):
			field := fmt.Sprintf("%s LIKE ?", key[2:len(key)-6])
			filters[field] = fmt.Sprintf("%%%s%%", f.Get(key))

		case strings.HasSuffix(key, "__gte"):
			field := fmt.Sprintf("%s >= ?", key[2:len(key)-5])
			filters[field] = f.Get(key)

		case strings.HasSuffix(key, "__lte"):
			field := fmt.Sprintf("%s <= ?", key[2:len(key)-5])
			filters[field] = f.Get(key)
		}
	}

	q.Filters = filters

	return &q, nil
}
