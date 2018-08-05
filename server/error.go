package server

type ApiError struct {
	Message string `json:"message"`
	Key     string `json:"key"`
}

func NewApiError(message string, key string) *ApiError {
	return &ApiError{
		Message: message,
		Key:     key,
	}
}
