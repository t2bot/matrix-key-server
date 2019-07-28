package common

import (
	"net/http"
)

type EmptyResponse struct{}

type ErrorResponse struct {
	Code       string `json:"errcode"`
	Message    string `json:"error"`
	HttpStatus int    `json:"http_status"`
}

func InternalServerError(message string) *ErrorResponse {
	return &ErrorResponse{"M_UNKNOWN", message, http.StatusInternalServerError}
}

func MethodNotAllowed() *ErrorResponse {
	return &ErrorResponse{"M_UNKNOWN", "Method Not Allowed", http.StatusMethodNotAllowed}
}

func NotFoundError() *ErrorResponse {
	return &ErrorResponse{"M_NOT_FOUND", "Resource Not Found", http.StatusNotFound}
}

func UnauthorizedError() *ErrorResponse {
	return &ErrorResponse{"M_UNAUTHORIZED", "Authentication Failed", http.StatusUnauthorized}
}

func BadRequest(message string) *ErrorResponse {
	return &ErrorResponse{"M_UNKNOWN", message, http.StatusBadRequest}
}
