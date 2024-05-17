package core

import (
	"log/slog"
	"net/http"
)

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func APIHandler(h APIFunc) http.HandlerFunc {
	// The Make function is an adapter between the default http.HandlerFunc
	// and the custom handler APIFunc which returns an error. Because I want
	// to be able to centralize the error handling
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			// TODO: consider trying errors.As() later
			if apiErr, ok := err.(APIError); ok {
				w.WriteHeader(apiErr.Status)
				w.Write([]byte("Error: " + apiErr.Msg))
				slog.Error("API", "detail", apiErr.Msg)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				slog.Error("Internal server error", "err", err.Error())
			}
		}
	}
}
