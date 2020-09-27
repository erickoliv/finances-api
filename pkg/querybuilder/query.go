package querybuilder

import (
	"errors"

	"github.com/erickoliv/finances-api/pkg/http/rest"
	"gorm.io/gorm"
)

var errInvalidQuery = errors.New("invalid query")

// Build is the function where the Query attributes are translated into a gorm.DB instance. Can be used to generic filter, order and pagination
func Build(db *gorm.DB, q *rest.Query) (*gorm.DB, error) {
	if q == nil {
		return nil, errInvalidQuery
	}

	query := db.Offset(q.Limit * (q.Page - 1)).Limit(q.Limit)

	if q.Sort != "" {
		query = query.Order(q.Sort)
	}

	for k, v := range q.Filters {
		query = query.Where(k, v)
	}

	return query, nil
}
