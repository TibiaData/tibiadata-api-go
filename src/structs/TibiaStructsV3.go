package structs

// Information stores some API data
type Information struct {
	APIVersion int    `json:"api_version"`
	Timestamp  string `json:"timestamp"`
}
