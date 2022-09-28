package models

// HealthCheckResponse returns the health check response duh
type HealthCheckResponse struct {
	Alive bool `json:"alive"`
}

// UserResponse is a general response structure with a status, message and optional json data
type UserResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// ErrorMessageResponse returns the error message response struct
type ErrorMessageResponse struct {
	Response MessageError
}

// MessageError contains the inner details for the error message response
type MessageError struct {
	Message string
	Error   string
}
