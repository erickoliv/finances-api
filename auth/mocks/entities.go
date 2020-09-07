package mocks

import (
	"time"

	"github.com/erickoliv/finances-api/auth"
	"github.com/google/uuid"
)

var date20200101 = time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC)

func ValidUser() *auth.User {
	return &auth.User{
		UUID:      uuid.MustParse("2415d0a8-e543-4007-b323-51f19325b74a"),
		CreatedAt: date20200101,
		UpdatedAt: date20200101,
		Username:  "johndoe",
		Email:     "validmail@mail.com",
		FirstName: "John",
		LastName:  "Doe",
		Active:    true,
	}
}
