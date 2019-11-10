package domain

// QueryData is structure which contains standard atributes to parse http parameters for API filter, search and pagination
type QueryData struct {
	Page    int
	Pages   int
	Total   int
	Limit   int
	Sort    string
	Filters map[string]interface{}
}
