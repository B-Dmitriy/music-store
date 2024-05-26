package web

import (
	"encoding/json"
	"net/http"
)

type WebError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func WriteJSON(w http.ResponseWriter, data any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// TODO: implement 401,403

func WriteNotFound(w http.ResponseWriter, e error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	err := json.NewEncoder(w).Encode(&WebError{
		Code:    http.StatusNotFound,
		Message: e.Error(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func WriteBadRequest(w http.ResponseWriter, e error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	err := json.NewEncoder(w).Encode(&WebError{
		Code:    http.StatusBadRequest,
		Message: e.Error(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func WriteServerError(w http.ResponseWriter, e error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	err := json.NewEncoder(w).Encode(&WebError{
		Code:    http.StatusInternalServerError,
		Message: e.Error(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
