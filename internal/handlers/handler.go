package handlers

import (
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	logger *zap.SugaredLogger
}

func NewHandler(logger *zap.SugaredLogger) http.Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("server is responding")); err != nil {
		h.logger.Errorw("writing response body", "error", err)
	}
}
