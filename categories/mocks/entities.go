package mocks

import (
	"time"

	"github.com/erickoliv/finances-api/categories"
	"github.com/google/uuid"
)

var date20200101 = time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC)

func ValidCompleteCategory() *categories.Category {
	return &categories.Category{
		UUID:        uuid.MustParse("2415d0a8-e543-4007-b323-51f19325b74a"),
		CreatedAt:   date20200101,
		UpdatedAt:   date20200101,
		Name:        "a valid category name",
		Description: "a valid description, with sóme speci@l s&mbols",
		Owner:       uuid.MustParse("2415d0a8-e543-4007-b323-51f19325b74b"),
	}
}

func ValidCategoryWithoutDescription() *categories.Category {
	return &categories.Category{
		UUID:      uuid.MustParse("2415d0a8-e543-4007-b323-51f19325b74c"),
		CreatedAt: date20200101,
		UpdatedAt: date20200101,
		Name:      "a valid name",
		Owner:     uuid.MustParse("2415d0a8-e543-4007-b323-52f19325b74b"),
	}
}

func InvalidValidCategoryWithoutName() *categories.Category {
	return &categories.Category{
		UUID:        uuid.MustParse("2415d0a8-e543-4007-b323-51f19325b74c"),
		CreatedAt:   date20200101,
		UpdatedAt:   date20200101,
		Description: "a valid description, with sóme speci@l s&mbols",
		Owner:       uuid.MustParse("2415d0a8-e543-4007-b323-52f19325b74b"),
	}
}

func ValidCategories() []categories.Category {
	return []categories.Category{
		*ValidCompleteCategory(),
		*ValidCategoryWithoutDescription(),
		*InvalidValidCategoryWithoutName(),
	}
}

const ValidCategoryPayload = `{
	"name": "a name",
	"description":"a description"
}`
