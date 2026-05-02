package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"filipaipilot/internal/ai"
	"filipaipilot/internal/web"
)

func main() {
	loadDotEnv(".env")

	addr := env("ADDR", ":8080")
	provider := strings.ToLower(env("AI_PROVIDER", "openai"))
	openAIKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	openAIModel := env("OPENAI_MODEL", "gpt-4o-mini")
	geminiKey := strings.TrimSpace(os.Getenv("GEMINI_API_KEY"))
	geminiModel := env("GEMINI_MODEL", "gemini-2.5-flash")
	aiTimeout := envDurationSeconds("AI_TIMEOUT_SECONDS", 180*time.Second)

	var client ai.Client
	mockMode := false
	switch provider {
	case "gemini":
		if geminiKey == "" {
			mockMode = true
			client = ai.NewMockAIClient()
			log.Println("GEMINI_API_KEY is not set; running in mock demo mode")
		} else {
			client = ai.NewGeminiClientWithTimeout(geminiKey, geminiModel, aiTimeout)
			log.Printf("running with provider gemini, model %s and timeout %s", geminiModel, aiTimeout)
		}
	case "openai":
		fallthrough
	default:
		if openAIKey == "" {
			mockMode = true
			client = ai.NewMockAIClient()
			log.Println("OPENAI_API_KEY is not set; running in mock demo mode")
		} else {
			client = ai.NewOpenAIClientWithTimeout(openAIKey, openAIModel, aiTimeout)
			log.Printf("running with provider openai, model %s and timeout %s", openAIModel, aiTimeout)
		}
	}

	if client == nil {
		client = ai.NewMockAIClient()
		mockMode = true
		log.Println("AI client was not initialized; running in mock demo mode")
	}

	server, err := web.NewServer(client, mockMode)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("FilipAiPilot listening on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, server.Routes()))
}

func envDurationSeconds(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	seconds, err := strconv.Atoi(value)
	if err != nil || seconds <= 0 {
		return fallback
	}
	return time.Duration(seconds) * time.Second
}

func env(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func loadDotEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if key != "" && os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}
}
