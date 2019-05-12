package internal

import (
	"github.com/go-chi/chi"
	"net/http"
	"strings"
)

// NewRouter creates a new router with HTTP handlers
func NewRouter(store Store) http.Handler {
	router := chi.NewRouter()
	router.Route("/v1/users/{userID}", func(r chi.Router) {
		r.Route("/streams/{streamID}", func(r chi.Router) {
			r.Delete("/", deleteStream(store))
			r.Put("/", createStream(store))
		})
		r.Get("/", listStreams(store))
	})
	return router
}

func createStream(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, streamID := getURLParams(r)
		if err := store.Add(userID, streamID); err != nil {
			if err == exceededStreamsQuota {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func listStreams(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		streamIDs, err := store.Get(userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, err = w.Write([]byte(strings.Join(streamIDs, ","))); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
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
