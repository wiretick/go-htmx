package core

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
)

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status) // Needs to be after the other Header() stuff, since reasons of ResponseWriter
	return json.NewEncoder(w).Encode(v)
}

func APIHandler(h APIFunc) http.HandlerFunc {
	// The Make function is an adapter between the default http.HandlerFunc
	// and the custom handler APIFunc which returns an error. Because I want
	// to be able to centralize the error handling
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			// TODO: consider trying errors.As() or .Is() later
			if apiErr, ok := err.(APIError); ok {
				//w.WriteHeader(apiErr.Status)
				//w.Write([]byte("Error: " + apiErr.Msg))
				WriteJSON(w, apiErr.Status, apiErr)
				slog.Error("API", "detail", apiErr.Msg)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				slog.Error("Internal server error", "err", err.Error())
			}
		}
	}
}
