package main

import (
	"encoding/json"
	"net/http"
)

type ErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseStruct struct {
	Success bool         `json:"success"`
	Data    interface{}  `json:"data,omitempty"`
	Message string       `json:"message,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

func Send(w http.ResponseWriter, code int, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	var resp ResponseStruct
	if code >= 200 && code < 300 {
		// Success
		resp = ResponseStruct{
			Success: true,
			Data:    data,
			Message: message,
		}
	} else {
		// Failure
		resp = ResponseStruct{
			Success: false,
			Error: &ErrorDetail{
				Code:    code,
				Message: message,
			},
		}
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, `{"success":false,"error":{"code":500,"message":"Failed to encode response"}}`, http.StatusInternalServerError)
	}
}
