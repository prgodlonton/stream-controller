package internal

import (
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// NewRouter creates a new router with HTTP handlers
func NewRouter(logger *zap.SugaredLogger, store Store) http.Handler {
	router := chi.NewRouter()
	router.Route("/v1/users/{userID}", func(r chi.Router) {
		r.Route("/streams/{streamID}", func(r chi.Router) {
			r.Delete("/", deleteStream(logger, store))
			r.Put("/", createStream(logger, store))
		})
		r.Get("/", listStreams(logger, store))
	})
	return router
}

func createStream(logger *zap.SugaredLogger, store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, streamID := getURLParams(r)
		if err := store.AddStream(userID, streamID); err != nil {
			if err == exceededStreamsQuota {
				logger.Debugw(
					"user exceeded streaming quota",
					"userID", userID,
					"streamID", streamID,
				)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			logger.Errorw(
				"cannot create stream",
				"userID", userID,
				"streamID", streamID,
				"error", err,
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func listStreams(logger *zap.SugaredLogger, store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")
		streamIDs, err := store.GetStreams(userID)
		if err != nil {
			logger.Debugw(
				"cannot list streams",
				"userID", userID,
				"error", err,
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if _, err = w.Write([]byte(strings.Join(streamIDs, ","))); err != nil {
			logger.Errorw(
				"cannot write to http response",
				"userID", userID,
				"error", err,
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func deleteStream(logger *zap.SugaredLogger, store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, streamID := getURLParams(r)
		if err := store.RemoveStream(userID, streamID); err != nil {
			logger.Errorw(
				"cannot remove stream",
				"userID", userID,
				"streamID", streamID,
				"error", err,
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func getURLParams(r *http.Request) (string, string) {
	return chi.URLParam(r, "userID"), chi.URLParam(r, "streamID")
}
