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
	apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	modelName := env("OPENAI_MODEL", "gpt-4o-mini")
	openAITimeout := envDurationSeconds("OPENAI_TIMEOUT_SECONDS", 180*time.Second)

	var client ai.Client
	mockMode := apiKey == ""
	if mockMode {
		client = ai.NewMockAIClient()
		log.Println("OPENAI_API_KEY is not set; running in mock demo mode")
	} else {
		client = ai.NewOpenAIClientWithTimeout(apiKey, modelName, openAITimeout)
		log.Printf("running with OpenAI model %s and timeout %s", modelName, openAITimeout)
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
