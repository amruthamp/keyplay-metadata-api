package util

import (
	"context"
	"encoding/json"
	"net/http"
)

type Error struct {
	Error ErrorDetails `json:"error"`
}

type ErrorDetails struct {
	Message string `json:"message"`
}

func JsonError(ctx context.Context, w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)

	data := Error{Error: ErrorDetails{
		Message: err.Error(),
	}}
	json.NewEncoder(w).Encode(data)

}

func Json(ctx context.Context, w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
