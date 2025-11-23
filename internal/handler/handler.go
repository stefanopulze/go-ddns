package handler

import (
	"go-ddns/internal/api/middleware"
	"go-ddns/internal/provider"
	"log/slog"
	"net/http"
)

type Handler struct {
	provider       provider.Client
	authMiddleware *middleware.Authorization
}

func NewHandler(pc provider.Client) *Handler {
	return &Handler{
		provider: pc,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UpdateDNS(writer http.ResponseWriter, request *http.Request) {
	slog.Info("Request to update dns")
	ip := request.URL.Query().Get("ip")
	hostname := request.URL.Query().Get("hostname")

	if err := h.provider.UpdateIp(hostname, ip); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
