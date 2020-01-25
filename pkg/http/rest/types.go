package rest

import (
	"errors"
	"fmt"
	"github.com/erickoliv/finances-api/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// ExtractFilters can be used to parse query parameters and return a Query object, useful to query, filter and paginate requests
func ExtractFilters(f url.Values) domain.Query {
	q := domain.Query{
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

func extractUser(c *gin.Context) (uuid.UUID, error) {
	user, found := c.Get(domain.LoggedUser)
	if !found {
		return uuid.New(), errors.New("user not present in context")
	}

	pk, ok := user.(uuid.UUID)
	if !ok {
		return uuid.New(), errors.New("user in context is invalid")
	}

	return pk, nil
}

func extractUUID(c *gin.Context) (uuid.UUID, error) {

	pk, err := uuid.Parse(c.Param("uuid"))
	if err != nil {
		return uuid.New(), errors.New("uuid parameter is invalid")
	}

	return pk, nil
}
