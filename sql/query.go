package sql

import (
	"math"

	"github.com/erickoliv/finances-api/pkg/http/rest"
	"github.com/jinzhu/gorm"
)

// Build is the function where the Query attributes are translated into a gorm.DB instance. Can be used to generic filter, order and pagination
func BuildQuery(base *gorm.DB, q rest.Query) *gorm.DB {
	for k, v := range q.Filters {
		base = base.Where(k, v)
	}

	base.Count(&q.Total)
	base = base.Offset(q.Limit * (q.Page - 1)).Limit(q.Limit).Order(q.Sort)
	q.Pages = int(math.Ceil(float64(q.Total) / float64(q.Limit)))

	return base
}
