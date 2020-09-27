package mocks

import (
	"time"

	accounts "github.com/erickoliv/finances-api/accounts"
	"github.com/google/uuid"
)

var date20200101 = time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC)

func ValidCompleteAccount() *accounts.Account {
	return &accounts.Account{
		UUID:        uuid.MustParse("2415d0a8-e543-4007-b323-51f19325b74a"),
		CreatedAt:   date20200101,
		UpdatedAt:   date20200101,
		Name:        "a valid name",
		Description: "a valid description, with sóme speci@l s&mbols",
		Owner:       uuid.MustParse("2415d0a8-e543-4007-b323-51f19325b74b"),
	}
}

func ValidAccountWithoutDescription() *accounts.Account {
	return &accounts.Account{
		UUID:      uuid.MustParse("2415d0a8-e543-4007-b323-51f19325b74c"),
		CreatedAt: date20200101,
		UpdatedAt: date20200101,
		Name:      "a valid name",
		Owner:     uuid.MustParse("2415d0a8-e543-4007-b323-52f19325b74b"),
	}
}

func ValidAccountWithoutName() *accounts.Account {
	return &accounts.Account{
		UUID:        uuid.MustParse("2415d0a8-e543-4007-b323-51f19325b74c"),
		CreatedAt:   date20200101,
		UpdatedAt:   date20200101,
		Description: "a valid description, with sóme speci@l s&mbols",
		Owner:       uuid.MustParse("2415d0a8-e543-4007-b323-52f19325b74b"),
	}
}

func ValidAcccounts() []*accounts.Account {
	return []*accounts.Account{
		ValidCompleteAccount(), ValidAccountWithoutDescription(), ValidAccountWithoutName(),
	}
}

const ValidAccountPayload = `{
	"name": "a name",
	"description":"a description"
}`
