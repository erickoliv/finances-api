package querybuilder

import (
	"errors"

	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/jinzhu/gorm"
)

var errInvalidQuery = errors.New("invalid query")

// Build is the function where the Query attributes are translated into a gorm.DB instance. Can be used to generic filter, order and pagination
func Build(base *gorm.DB, q *rest.Query) (*gorm.DB, error) {
	if q == nil {
		return nil, errInvalidQuery
	}
	for k, v := range q.Filters {
		base = base.Where(k, v)
	}

	base = base.Offset(q.Limit * (q.Page - 1)).Limit(q.Limit).Order(q.Sort)

	return base, nil
}
