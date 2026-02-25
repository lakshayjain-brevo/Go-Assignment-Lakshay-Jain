package model

type HashRequest struct {
	Input string `json:"input"`
}

type HashResponse struct {
	Input string `json:"input"`
	Hash  string `json:"hash"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
