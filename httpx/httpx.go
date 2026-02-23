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

func (r *DefaultResponder) JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (r *DefaultResponder) Error(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	r.JSON(w, status, payload)
}
