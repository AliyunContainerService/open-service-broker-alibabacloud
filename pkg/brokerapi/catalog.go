package brokerapi

// Catalog is a JSON-compatible type to be used to decode the result from a /v2/catalog call
// to an open service broker compatible API
type Catalog struct {
	Services []*Service `json:"services"`
}
