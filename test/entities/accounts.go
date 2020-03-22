package entities

import (
	"github.com/erickoliv/finances-api/domain"
	"github.com/google/uuid"
)

func ValidCompleteAccount() *domain.Account {
	return &domain.Account{
		BaseModel: domain.BaseModel{
			UUID:      uuid.MustParse("2415d0a8-e543-4007-b323-51f19325b74a"),
			CreatedAt: Date20200101,
			UpdatedAt: Date20200101,
			DeletedAt: nil,
		},
		Name:        "a valid name",
		Description: "a valid description, with sóme speci@l s&mbols",
		Owner:       uuid.MustParse("2415d0a8-e543-4007-b323-51f19325b74b"),
	}
}

func ValidAccountWithoutDescription() *domain.Account {
	return &domain.Account{
		BaseModel: domain.BaseModel{
			UUID:      uuid.MustParse("2415d0a8-e543-4007-b323-51f19325b74c"),
			CreatedAt: Date20200101,
			UpdatedAt: Date20200101,
			DeletedAt: nil,
		},
		Name:  "a valid name",
		Owner: uuid.MustParse("2415d0a8-e543-4007-b323-52f19325b74b"),
	}
}

func ValidAccountWithoutName() *domain.Account {
	return &domain.Account{
		BaseModel: domain.BaseModel{
			UUID:      uuid.MustParse("2415d0a8-e543-4007-b323-51f19325b74c"),
			CreatedAt: Date20200101,
			UpdatedAt: Date20200101,
			DeletedAt: nil,
		},
		Description: "a valid description, with sóme speci@l s&mbols",
		Owner:       uuid.MustParse("2415d0a8-e543-4007-b323-52f19325b74b"),
	}
}

func ValidAcccounts() []*domain.Account {
	return []*domain.Account{
		ValidCompleteAccount(), ValidAccountWithoutDescription(), ValidAccountWithoutName(),
	}
}

const ValidAccountPayload = `{
	"name": "a name",
	"description":"a description"
}`
