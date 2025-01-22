package responses

// Generic API response
// swagger:response Response
type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}
