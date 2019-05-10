package handlers

import (
	"fmt"
	"net/http"
)

type Handler struct {
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("server is responding")); err != nil {
		fmt.Printf("cannot write response due to %v\n", err.Error())
	}
}
