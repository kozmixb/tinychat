package main

import (
	"errors"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", serveIndex)
	mux.Handle("/static/", http.FileServer(http.FS(staticFiles)))
	mux.HandleFunc("/api/config", handleConfig(cfg))
	mux.HandleFunc("/api/health", proxyRequest(cfg, http.MethodGet, "/health"))
	mux.HandleFunc("/api/models", proxyRequest(cfg, http.MethodGet, "/models"))
	mux.HandleFunc("/api/chat", proxyRequest(cfg, http.MethodPost, "/chat/completions"))

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           logRequests(mux),
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("listening on http://%s", cfg.Addr)
	log.Printf("proxying OpenAI-compatible chat API at %s", cfg.OpenAIChatURL)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
