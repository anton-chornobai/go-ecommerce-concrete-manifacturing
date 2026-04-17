package dto

type UserContactPostRequest struct {
	Name    string  `json:"name"`
	Email   *string  `json:"email,omitempty"`
	Message string  `json:"message"`
	Number  string `json:"number"`
}
