package main

import (
	"fmt"
	"go-ddns/internal/api/middleware"
	"go-ddns/internal/config"
	"go-ddns/internal/handler"
	"go-ddns/internal/provider/cloudflare"
	"log/slog"
	"net/http"
)

var version = "dev"

func main() {
	//ctx, cancel := context.WithCancel(context.Background())
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	slog.Info(fmt.Sprintf("Starting go-ddns %s", version))

	cc, err := cloudflare.NewCloudflareClient(cfg.Providers.Cloudflare)
	if err != nil {
		panic(err)
	}
	am := middleware.NewAuthorizationMiddleware(cfg.Authorization.Username, cfg.Authorization.Password)

	mux := http.NewServeMux()
	h := handler.NewHandler(cc)

	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/update", am.Secure(http.HandlerFunc(h.UpdateDNS)))

	addr := fmt.Sprintf(":%d", cfg.Port)
	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	slog.Info(fmt.Sprintf("Starting server at %s", addr))
	if err = server.ListenAndServe(); err != nil {
		panic(err)
	}
}
