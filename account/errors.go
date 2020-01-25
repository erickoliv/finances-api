package account

import "github.com/erickoliv/finances-api/pkg/http/rest"

var (
	accountNotFound = rest.ErrorMessage{Message: "account not found"}
)
