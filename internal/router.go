package internal

import (
	"github.com/go-chi/chi"
	"net/http"
)

// NewRouter creates a new router with HTTP handlers
func NewRouter(store Store) http.Handler {
	router := chi.NewRouter()
	router.Route("/v1/users/{userID}", func(r chi.Router) {
		r.Route("/streams/{streamID}", func(r chi.Router) {
			r.Put("/", createStream(store))
			r.Delete("/", deleteStream(store))
		})
	})
	return router
}

func createStream(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, streamID := getURLParams(r)
		if err := store.Add(userID, streamID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func deleteStream(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, streamID := getURLParams(r)
		if err := store.Remove(userID, streamID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func getURLParams(r *http.Request) (string, string) {
	return chi.URLParam(r, "userID"), chi.URLParam(r, "streamID")
}
