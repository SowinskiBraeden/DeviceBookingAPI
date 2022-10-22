package models

// This is io.go (input/output) for json queries and responses

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
	Response MessageError `json:"response"`
}

// MessageError contains the inner details for the error message response
type MessageError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

// Query for search queries by name
type Query struct {
	Name string `json:"name"`
}

// NewDeviceToCow is to add device object id to cow
type NewDeviceToCow struct {
	ID string `json:"_id"`
}
