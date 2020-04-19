package rest

// Query is structure which contains standard attributes to parse http parameters for API filter, search and pagination
type Query struct {
	Page    int
	Limit   int
	Sort    string
	Filters map[string]interface{}
}
