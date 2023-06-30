package entities

// APIResponseMessage represents the response for generic happy cases for the RESTFUL apis
type APIResponseMessage struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Body       string `json:"body"`
}
