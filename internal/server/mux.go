package server

import (
	"net/http"
)

func NewServeMux(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/ports/import", h.ImportPorts)
	return mux
}
