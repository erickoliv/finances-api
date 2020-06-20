package rest

// PaginatedMessage is a structure which contains standard attributes to be used on paginated services
type PaginatedMessage struct {
	Page  int         `json:"page"`
	Count int         `json:"count"`
	Data  interface{} `json:"data"`
}

// Query is structure which contains standard attributes to parse http parameters for API filter, search and pagination
type Query struct {
	Page    int
	Limit   int
	Sort    string
	Filters map[string]interface{}
}
