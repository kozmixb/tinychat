package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = io.WriteString(w, indexHTML)
}

func handleConfig(cfg appConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, cfg)
	}
}

func proxyRequest(cfg appConfig, method, path string) http.HandlerFunc {
	client := &http.Client{Timeout: 0}
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		if r.Method != method {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		baseURL, err := requestBaseURL(cfg, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		target := baseURL + path
		body, err := readBody(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req, err := http.NewRequestWithContext(r.Context(), method, target, bytes.NewReader(body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "text/event-stream, application/json")
		log.Printf("proxy request method=%s target=%s content_type=%q accept=%q bytes=%d", method, target, req.Header.Get("Content-Type"), req.Header.Get("Accept"), len(body))

		resp, err := client.Do(req)
		if err != nil {
			status := http.StatusBadGateway
			if errors.Is(err, context.Canceled) {
				status = http.StatusRequestTimeout
			}
			log.Printf("proxy response method=%s target=%s status=%d error=%q duration=%s", method, target, status, err.Error(), time.Since(start).Round(time.Millisecond))
			http.Error(w, err.Error(), status)
			return
		}
		defer resp.Body.Close()
		log.Printf("proxy response method=%s target=%s status=%d content_type=%q duration=%s", method, target, resp.StatusCode, resp.Header.Get("Content-Type"), time.Since(start).Round(time.Millisecond))

		copyHeader(w.Header(), resp.Header)
		w.WriteHeader(resp.StatusCode)
		_, _ = copyResponse(w, resp.Body)
	}
}

func copyResponse(w http.ResponseWriter, r io.Reader) (int64, error) {
	flusher, _ := w.(http.Flusher)
	buf := make([]byte, 32*1024)
	var written int64
	for {
		nr, er := r.Read(buf)
		if nr > 0 {
			nw, ew := w.Write(buf[:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if flusher != nil {
				flusher.Flush()
			}
			if ew != nil {
				return written, ew
			}
			if nr != nw {
				return written, io.ErrShortWrite
			}
		}
		if er != nil {
			if errors.Is(er, io.EOF) {
				return written, nil
			}
			return written, er
		}
	}
}

func requestBaseURL(cfg appConfig, r *http.Request) (string, error) {
	if cfg.OpenAIChatURL != "" {
		return cfg.OpenAIChatURL, nil
	}

	raw := strings.TrimSpace(r.Header.Get("X-OpenAI-Chat-Host"))
	if raw == "" {
		return "", errors.New("OPENAI_CHAT_HOST is not configured; provide a runtime host such as https://api.example.com/v1")
	}

	baseURL := normalizeOpenAIBaseURL(raw)
	parsed, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid OpenAI chat host: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", errors.New("OpenAI chat host must use http or https")
	}
	if parsed.Host == "" {
		return "", errors.New("OpenAI chat host must include a host")
	}

	return baseURL, nil
}

func readBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}
	defer r.Body.Close()
	return io.ReadAll(io.LimitReader(r.Body, 8<<20))
}

func copyHeader(dst, src http.Header) {
	for key, values := range src {
		lower := strings.ToLower(key)
		if lower == "content-length" || lower == "connection" {
			continue
		}
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
