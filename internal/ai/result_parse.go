package ai

import (
	"encoding/json"
	"strings"

	"filipaipilot/internal/model"
)

func parseResultContent(content string) (model.Result, bool) {
	content = strings.TrimSpace(content)
	if content == "" {
		return model.Result{}, false
	}

	var result model.Result
	if json.Unmarshal([]byte(content), &result) == nil {
		return result, true
	}

	cleaned := stripMarkdownCodeFence(content)
	if json.Unmarshal([]byte(cleaned), &result) == nil {
		return result, true
	}

	extracted, ok := extractFirstJSONObject(cleaned)
	if !ok {
		return model.Result{}, false
	}
	if json.Unmarshal([]byte(extracted), &result) == nil {
		return result, true
	}

	return model.Result{}, false
}

func stripMarkdownCodeFence(text string) string {
	text = strings.TrimSpace(text)
	if !strings.HasPrefix(text, "```") {
		return text
	}

	lines := strings.Split(text, "\n")
	if len(lines) < 2 {
		return text
	}
	if strings.HasPrefix(strings.TrimSpace(lines[0]), "```") {
		lines = lines[1:]
	}
	if len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "```" {
		lines = lines[:len(lines)-1]
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func extractFirstJSONObject(text string) (string, bool) {
	start := strings.Index(text, "{")
	if start == -1 {
		return "", false
	}

	depth := 0
	for i := start; i < len(text); i++ {
		switch text[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return text[start : i+1], true
			}
		}
	}
	return "", false
}
