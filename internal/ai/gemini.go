package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"filipaipilot/internal/model"
)

type GeminiClient struct {
	apiKey string
	model  string
	http   *http.Client
}

func NewGeminiClient(apiKey, modelName string) *GeminiClient {
	return NewGeminiClientWithTimeout(apiKey, modelName, 180*time.Second)
}

func NewGeminiClientWithTimeout(apiKey, modelName string, timeout time.Duration) *GeminiClient {
	if modelName == "" {
		modelName = "gemini-2.5-flash"
	}
	if timeout <= 0 {
		timeout = 180 * time.Second
	}
	return &GeminiClient{
		apiKey: apiKey,
		model:  modelName,
		http:   &http.Client{Timeout: timeout},
	}
}

func (c *GeminiClient) Analyze(ctx context.Context, action string, inputA string, inputB string, options Options) (model.Result, error) {
	if strings.TrimSpace(c.apiKey) == "" {
		return model.Result{}, errors.New("missing Gemini API key")
	}

	payload := geminiRequest{
		Contents: []geminiContent{
			{
				Role: "user",
				Parts: []geminiPart{
					{Text: SystemPrompt + "\n\n" + buildPrompt(action, inputA, inputB, options)},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return model.Result{}, err
	}

	endpoint := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", c.model, url.QueryEscape(c.apiKey))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return model.Result{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return model.Result{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return model.Result{}, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return model.Result{}, fmt.Errorf("gemini status %d: %s", resp.StatusCode, sanitizeGeminiError(respBody))
	}

	var parsed geminiResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return model.Result{}, err
	}
	if len(parsed.Candidates) == 0 {
		return model.Result{}, errors.New("gemini returned no candidates")
	}
	if len(parsed.Candidates[0].Content.Parts) == 0 {
		return model.Result{}, errors.New("gemini returned empty content")
	}

	content := parsed.Candidates[0].Content.Parts[0].Text
	result, ok := parseResultContent(content)
	if !ok {
		return model.Result{
			Title:   PromptTemplateByIDOrDefault(action).Label,
			Summary: "AI vrátila odpověď mimo očekávaný JSON formát. Níže je surový výstup.",
			Raw:     content,
			Warnings: []string{
				"Výstup je pracovní podklad pro právníka, nikoli právní stanovisko.",
			},
		}, nil
	}
	if result.Title == "" {
		result.Title = PromptTemplateByIDOrDefault(action).Label
	}
	if len(result.Warnings) == 0 {
		result.Warnings = []string{"Výstup je pracovní podklad pro právníka, nikoli právní stanovisko."}
	}
	return result, nil
}

func sanitizeGeminiError(body []byte) string {
	var parsed struct {
		Error struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
			Status  string `json:"status"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &parsed); err == nil && parsed.Error.Message != "" {
		parts := []string{parsed.Error.Message}
		if parsed.Error.Status != "" {
			parts = append(parts, "status="+parsed.Error.Status)
		}
		if parsed.Error.Code != 0 {
			parts = append(parts, fmt.Sprintf("code=%d", parsed.Error.Code))
		}
		return strings.Join(parts, " | ")
	}
	text := strings.TrimSpace(string(body))
	if len(text) > 500 {
		text = text[:500] + "..."
	}
	return text
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []struct {
		Content geminiContent `json:"content"`
	} `json:"candidates"`
}
