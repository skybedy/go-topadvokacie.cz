package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"lexdemo/internal/model"
)

type OpenAIClient struct {
	apiKey string
	model  string
	http   *http.Client
}

func NewOpenAIClient(apiKey, modelName string) *OpenAIClient {
	return NewOpenAIClientWithTimeout(apiKey, modelName, 180*time.Second)
}

func NewOpenAIClientWithTimeout(apiKey, modelName string, timeout time.Duration) *OpenAIClient {
	if modelName == "" {
		modelName = "gpt-4o-mini"
	}
	if timeout <= 0 {
		timeout = 180 * time.Second
	}
	return &OpenAIClient{
		apiKey: apiKey,
		model:  modelName,
		http:   &http.Client{Timeout: timeout},
	}
}

func (c *OpenAIClient) Analyze(ctx context.Context, action string, inputA string, inputB string, options Options) (model.Result, error) {
	if strings.TrimSpace(c.apiKey) == "" {
		return model.Result{}, errors.New("missing OpenAI API key")
	}

	payload := chatRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "system", Content: SystemPrompt},
			{Role: "user", Content: buildPrompt(action, inputA, inputB, options)},
		},
		ResponseFormat: map[string]string{"type": "json_object"},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return model.Result{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return model.Result{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
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
		return model.Result{}, fmt.Errorf("openai status %d: %s", resp.StatusCode, sanitizeOpenAIError(respBody))
	}

	var chat chatResponse
	if err := json.Unmarshal(respBody, &chat); err != nil {
		return model.Result{}, err
	}
	if len(chat.Choices) == 0 {
		return model.Result{}, errors.New("openai returned no choices")
	}

	content := chat.Choices[0].Message.Content
	var result model.Result
	if err := json.Unmarshal([]byte(content), &result); err != nil {
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

func sanitizeOpenAIError(body []byte) string {
	var parsed struct {
		Error struct {
			Message string `json:"message"`
			Type    string `json:"type"`
			Code    string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &parsed); err == nil && parsed.Error.Message != "" {
		parts := []string{parsed.Error.Message}
		if parsed.Error.Type != "" {
			parts = append(parts, "type="+parsed.Error.Type)
		}
		if parsed.Error.Code != "" {
			parts = append(parts, "code="+parsed.Error.Code)
		}
		return strings.Join(parts, " | ")
	}
	text := strings.TrimSpace(string(body))
	if len(text) > 500 {
		text = text[:500] + "..."
	}
	return text
}

func buildPrompt(action string, inputA string, inputB string, options Options) string {
	selected := PromptTemplateByIDOrDefault(action)
	var b strings.Builder
	fmt.Fprintf(&b, "Uložený prompt: %s %s\n", selected.Label, selected.Version)
	fmt.Fprintf(&b, "Kategorie: %s\n", selected.Category)
	fmt.Fprintf(&b, "Instrukce promptu: %s\n", selected.Instruction)
	fmt.Fprintf(&b, "Preferovaná délka výstupu: %s\n", outputDetailLabel(options.DetailLevel))
	fmt.Fprintf(&b, "Perspektiva výstupu: %s\n", outputPerspectiveLabel(options.Perspective))
	b.WriteString("\n")
	b.WriteString("Vrať výhradně validní JSON ve tvaru:\n")
	b.WriteString(`{"title":"...","summary":"...","sections":[{"title":"...","items":["..."]}],"warnings":["..."],"raw":""}`)
	b.WriteString("\n\nPřizpůsob rozsah a tón zvolené délce a perspektivě. Pracuj pouze s informacemi ze vstupu. Když něco chybí, napiš to jako nejasnost nebo otázku.\n\n")
	b.WriteString("Dokument A:\n")
	b.WriteString(inputA)
	if selected.NeedsSecond || strings.TrimSpace(inputB) != "" {
		b.WriteString("\n\nDokument B:\n")
		b.WriteString(inputB)
	}
	return b.String()
}

func outputDetailLabel(value string) string {
	switch value {
	case "brief":
		return "stručně"
	case "detailed":
		return "detailně"
	default:
		return "standardně"
	}
}

func outputPerspectiveLabel(value string) string {
	switch value {
	case "client":
		return "pro klienta, s důrazem na praktické dopady"
	case "negotiation":
		return "pro vyjednávání, s důrazem na pozice a argumenty"
	default:
		return "pro právníka, s důrazem na pracovní právní analýzu"
	}
}

type chatRequest struct {
	Model          string            `json:"model"`
	Messages       []chatMessage     `json:"messages"`
	ResponseFormat map[string]string `json:"response_format,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}
