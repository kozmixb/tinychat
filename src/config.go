package main

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
)

type appConfig struct {
	Addr          string `json:"-"`
	OpenAIChatURL string `json:"openai_chat_host"`
}

func loadConfig() (appConfig, error) {
	host := env("APP_HOST", "127.0.0.1")
	port := env("APP_PORT", "8080")
	openAIHost := strings.TrimSpace(os.Getenv("OPENAI_CHAT_HOST"))

	if openAIHost != "" {
		openAIHost = normalizeOpenAIBaseURL(openAIHost)
		if _, err := url.ParseRequestURI(openAIHost); err != nil {
			return appConfig{}, fmt.Errorf("invalid OPENAI_CHAT_HOST: %w", err)
		}
	}

	return appConfig{
		Addr:          net.JoinHostPort(host, port),
		OpenAIChatURL: strings.TrimRight(openAIHost, "/"),
	}, nil
}

func env(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func normalizeBaseURL(raw string) string {
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return raw
	}
	return "http://" + raw
}

func normalizeOpenAIBaseURL(raw string) string {
	baseURL := strings.TrimRight(normalizeBaseURL(strings.TrimSpace(raw)), "/")
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}
	if parsed.Path == "" || parsed.Path == "/" {
		parsed.Path = "/v1"
		return strings.TrimRight(parsed.String(), "/")
	}
	return baseURL
}
