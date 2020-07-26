package mocks

import (
	"time"

	"github.com/erickoliv/finances-api/tags"
	"github.com/google/uuid"
)

var date20200101 = time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC)

func ValidCompleteTag() *tags.Tag {
	return &tags.Tag{
		UUID:        uuid.MustParse("2415d0a8-e543-4007-b323-51f19325b74a"),
		CreatedAt:   date20200101,
		UpdatedAt:   date20200101,
		DeletedAt:   nil,
		Name:        "a valid name",
		Description: "a valid description, with sóme speci@l s&mbols",
		Owner:       uuid.MustParse("2415d0a8-e543-4007-b323-51f19325b74b"),
	}
}

func ValidTagWithoutDescription() *tags.Tag {
	return &tags.Tag{
		UUID:      uuid.MustParse("2415d0a8-e543-4007-b323-51f19325b74c"),
		CreatedAt: date20200101,
		UpdatedAt: date20200101,
		DeletedAt: nil,
		Name:      "a valid name",
		Owner:     uuid.MustParse("2415d0a8-e543-4007-b323-52f19325b74b"),
	}
}

func InvalidValidTagWithoutName() *tags.Tag {
	return &tags.Tag{
		UUID:        uuid.MustParse("2415d0a8-e543-4007-b323-51f19325b74c"),
		CreatedAt:   date20200101,
		UpdatedAt:   date20200101,
		DeletedAt:   nil,
		Description: "a valid description, with sóme speci@l s&mbols",
		Owner:       uuid.MustParse("2415d0a8-e543-4007-b323-52f19325b74b"),
	}
}

func ValidTags() []tags.Tag {
	return []tags.Tag{
		*ValidCompleteTag(), *ValidTagWithoutDescription(), *InvalidValidTagWithoutName(),
	}
}

const ValidTagPayload = `{
	"name": "a name",
	"description":"a description"
}`
