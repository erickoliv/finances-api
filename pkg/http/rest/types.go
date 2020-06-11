package rest

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/erickoliv/finances-api/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PaginatedMessage is a structure which contains standard attributes to be used on paginated services
type PaginatedMessage struct {
	Page  int         `json:"page"`
	Count int         `json:"count"`
	Data  interface{} `json:"data"`
}

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

	filters["owner = ?"] = user

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

func ExtractUser(c *gin.Context) (uuid.UUID, error) {
	identifier := c.GetString(domain.LoggedUser)
	if len(identifier) == 0 {
		return uuid.Nil, errors.New("user not present in context")
	}

	return uuid.Parse(identifier)
}

func ExtractUUID(c *gin.Context) (uuid.UUID, error) {

	pk, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return uuid.Nil, errors.New("uuid parameter is invalid")
	}

	return pk, nil
}
