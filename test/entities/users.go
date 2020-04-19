package entities

import (
	"github.com/erickoliv/finances-api/domain"
	"github.com/google/uuid"
)

func ValidUser() *domain.User {
	return &domain.User{
		BaseModel: domain.BaseModel{
			UUID:      uuid.MustParse("2415d0a8-e543-4007-b323-51f19325b74a"),
			CreatedAt: Date20200101,
			UpdatedAt: Date20200101,
			DeletedAt: nil,
		},
		Username:  "johndoe",
		Email:     "validmail@mail.com",
		FirstName: "John",
		LastName:  "Doe",
		Active:    true,
	}
}
