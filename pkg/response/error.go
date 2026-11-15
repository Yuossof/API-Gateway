package response

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func JSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
		Code:  status,
	})
}

func NotFound(w http.ResponseWriter, message string) {
	JSONError(w, message, http.StatusNotFound) // 404
}

func BadRequest(w http.ResponseWriter, message string) {
	JSONError(w, message, http.StatusBadRequest) // 400
}

func Unauthorized(w http.ResponseWriter, message string) {
	JSONError(w, message, http.StatusUnauthorized) // 401
}

func TooManyRequests(w http.ResponseWriter, message string) {
	JSONError(w, message, http.StatusTooManyRequests) // 429
}

func InternalError(w http.ResponseWriter, message string) {
	JSONError(w, message, http.StatusInternalServerError) // 500
}
