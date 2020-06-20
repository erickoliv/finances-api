package mocks

import (
	"time"

	"github.com/erickoliv/finances-api/entries"
	"github.com/google/uuid"
)

var date20200101 = time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC)

func ValidCompleteEntry() *entries.Entry {
	return &entries.Entry{
		UUID:        uuid.MustParse("2415d0a8-e543-4007-b323-51f19325b74a"),
		CreatedAt:   date20200101,
		UpdatedAt:   date20200101,
		DeletedAt:   nil,
		Name:        "a valid entry name",
		Description: "a valid description, with sóme speci@l s&mbols",
		Owner:       uuid.MustParse("2415d0a8-e543-4007-b323-51f19325b74b"),
	}
}

func ValidEntryWithoutDescription() *entries.Entry {
	return &entries.Entry{
		UUID:      uuid.MustParse("2415d0a8-e543-4007-b323-51f19325b74c"),
		Type:      true,
		Value:     float64(1.12),
		Category:  uuid.MustParse("d8babef3-110e-4259-9bcc-e68847b57e77"),
		Account:   uuid.MustParse("580a77ce-9ea0-445d-be97-508451d87328"),
		Date:      date20200101,
		CreatedAt: date20200101,
		UpdatedAt: date20200101,
		DeletedAt: nil,
		Name:      "a valid name",
		Owner:     uuid.MustParse("2415d0a8-e543-4007-b323-52f19325b74b"),
	}
}

func InvalidValidEntryWithoutName() *entries.Entry {
	return &entries.Entry{
		UUID:        uuid.MustParse("2415d0a8-e543-4007-b323-51f19325b74c"),
		CreatedAt:   date20200101,
		UpdatedAt:   date20200101,
		Date:        date20200101,
		Type:        true,
		Category:    uuid.MustParse("24f00f4b-0173-4f73-89ee-77bc303f7fbd"),
		Account:     uuid.MustParse("3fcc4248-9569-47db-a265-96114dcc60dc"),
		Value:       float64(125.00),
		DeletedAt:   nil,
		Description: "a valid description, with sóme speci@l s&mbols",
		Owner:       uuid.MustParse("2415d0a8-e543-4007-b323-52f19325b74b"),
	}
}

func ValidEntries() []entries.Entry {
	return []entries.Entry{
		*ValidCompleteEntry(),
		*ValidEntryWithoutDescription(),
		*InvalidValidEntryWithoutName(),
	}
}

const ValidEntryPayload = `{
	"name": "a name",
	"description":"a description",
	"type": true,
	"date": "2020-06-20T16:42:09.235571-03:00",
	"value": 25.26,
	"category": "3fcc4248-9569-47db-a265-96114dcc60dc",
	"account": "e4a3f901-ee6d-4716-b080-7c9a09f0d28b"
}`
