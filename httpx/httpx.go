package httpx

import (
	"encoding/json"
	"net/http"
)

type Responder interface {
	JSON(w http.ResponseWriter, status int, body any)
	Error(w http.ResponseWriter, status int, code string, description string)
}

type DefaultResponder struct{}

func (r *DefaultResponder) JSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func (r *DefaultResponder) Error(w http.ResponseWriter, status int, code string, description string) {
	r.JSON(w, status, map[string]any{
		"type":        "error",
		"code":        code,
		"description": description,
	})
}
