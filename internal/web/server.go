package web

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"lexdemo/internal/ai"
	"lexdemo/internal/model"
)

type Server struct {
	ai        ai.Client
	templates *template.Template
	mockMode  bool
}

type PageData struct {
	Title           string
	Active          string
	Actions         []ai.Action
	Prompts         []ai.PromptTemplate
	Examples        []model.Example
	SelectedAction  string
	InputA          string
	InputB          string
	SelectedExample string
	Result          *model.Result
	Error           string
	MockMode        bool
	CurrentYear     int
}

func NewServer(client ai.Client, mockMode bool) (*Server, error) {
	tmpl, err := template.ParseGlob(filepath.Join("templates", "*.html"))
	if err != nil {
		return nil, err
	}
	return &Server{ai: client, templates: tmpl, mockMode: mockMode}, nil
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/", s.home)
	mux.HandleFunc("/about", s.about)
	mux.HandleFunc("/examples", s.examples)
	mux.HandleFunc("/prompts", s.prompts)
	return securityHeaders(mux)
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := s.baseData("LexPilot Demo", "home")
	data.SelectedAction = "contract-analysis"
	if r.Method == http.MethodGet {
		if example, ok := ExampleByID(r.URL.Query().Get("example")); ok {
			data.SelectedExample = example.ID
			data.InputA = example.Content
		}
		if _, ok := ai.PromptTemplateByID(r.URL.Query().Get("prompt")); ok {
			data.SelectedAction = r.URL.Query().Get("prompt")
		}
	}

	if r.Method == http.MethodPost {
		if err := r.ParseMultipartForm(2 << 20); err != nil {
			data.Error = "Formulář se nepodařilo zpracovat."
			s.render(w, "home.html", data)
			return
		}

		data.SelectedAction = r.FormValue("action")
		data.InputA = strings.TrimSpace(r.FormValue("input_a"))
		data.InputB = strings.TrimSpace(r.FormValue("input_b"))
		data.SelectedExample = r.FormValue("example")

		uploadedA, err := readUploadedText(r, "file_a")
		if err != nil {
			data.Error = err.Error()
			s.render(w, "home.html", data)
			return
		}
		if uploadedA != "" {
			data.InputA = uploadedA
		}

		uploadedB, err := readUploadedText(r, "file_b")
		if err != nil {
			data.Error = err.Error()
			s.render(w, "home.html", data)
			return
		}
		if uploadedB != "" {
			data.InputB = uploadedB
		}

		if data.InputA == "" {
			data.Error = "Vložte dokument nebo vyberte ukázkový text."
			s.render(w, "home.html", data)
			return
		}
		if ai.ActionByID(data.SelectedAction).NeedsSecond && data.InputB == "" {
			data.Error = "Pro porovnání dvou verzí vložte také druhý dokument."
			s.render(w, "home.html", data)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 180*time.Second)
		defer cancel()
		result, err := s.ai.Analyze(ctx, data.SelectedAction, data.InputA, data.InputB)
		if err != nil {
			log.Printf("analyze failed: %v", err)
			data.Error = friendlyAnalyzeError(err)
			s.render(w, "home.html", data)
			return
		}
		data.Result = &result
	}

	s.render(w, "home.html", data)
}

func readUploadedText(r *http.Request, field string) (string, error) {
	file, header, err := r.FormFile(field)
	if err != nil {
		if err == http.ErrMissingFile {
			return "", nil
		}
		return "", fmt.Errorf("Soubor se nepodařilo načíst.")
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	switch ext {
	case ".txt", ".md", ".markdown", ".csv", ".rst", ".log":
	default:
		return "", fmt.Errorf("Soubor %s zatím neumím načíst. V tomto MVP nahrajte textový soubor .txt, .md nebo .csv.", header.Filename)
	}

	const maxUploadBytes = 512 * 1024
	content, err := io.ReadAll(io.LimitReader(file, maxUploadBytes+1))
	if err != nil {
		return "", fmt.Errorf("Soubor %s se nepodařilo přečíst.", header.Filename)
	}
	if len(content) > maxUploadBytes {
		return "", fmt.Errorf("Soubor %s je příliš velký. Limit pro demo je 512 KB.", header.Filename)
	}
	return strings.TrimSpace(string(content)), nil
}

func friendlyAnalyzeError(err error) string {
	message := err.Error()
	switch {
	case strings.Contains(message, "401"):
		return "OpenAI API odmítlo klíč. Zkontrolujte, že je v .env platný OPENAI_API_KEY a že nebyl smazán."
	case strings.Contains(message, "404") && strings.Contains(message, "model"):
		return "Zvolený OpenAI model není pro tento API klíč dostupný. Zkuste v .env nastavit OPENAI_MODEL=gpt-4o-mini nebo OPENAI_MODEL=gpt-5-nano."
	case strings.Contains(message, "429"):
		return "OpenAI API hlásí limit nebo nedostupný kredit. Zkontrolujte billing, usage limit nebo chvíli počkejte."
	case strings.Contains(message, "insufficient_quota") || strings.Contains(message, "billing"):
		return "OpenAI API hlásí problém s kreditem nebo billingem. Zkontrolujte nastavení plateb na platform.openai.com."
	case strings.Contains(message, "context_length") || strings.Contains(message, "maximum context"):
		return "Vložený dokument je pro zvolený model příliš dlouhý. Zkuste kratší text nebo silnější model s větším kontextem."
	case strings.Contains(message, "context deadline exceeded") || strings.Contains(message, "Client.Timeout") || strings.Contains(message, "timeout"):
		return "OpenAI API nestihlo odpovědět v časovém limitu. Zkuste kratší dokument, zopakovat požadavek nebo v .env navýšit OPENAI_TIMEOUT_SECONDS."
	case strings.Contains(message, "no such host") || strings.Contains(message, "connection refused"):
		return "Nepodařilo se připojit k OpenAI API. Zkontrolujte internetové připojení a zkuste to znovu."
	default:
		return "Analýza se nepodařila. Detail je v terminálu, kde běží aplikace. Nejčastěji jde o model, API klíč, billing nebo síť."
	}
}

func (s *Server) about(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/about" {
		http.NotFound(w, r)
		return
	}
	s.render(w, "about.html", s.baseData("O projektu", "about"))
}

func (s *Server) examples(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/examples" {
		http.NotFound(w, r)
		return
	}
	s.render(w, "examples.html", s.baseData("Ukázková data", "examples"))
}

func (s *Server) prompts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/prompts" {
		http.NotFound(w, r)
		return
	}
	s.render(w, "prompts.html", s.baseData("Prompt knihovna", "prompts"))
}

func (s *Server) baseData(title, active string) PageData {
	return PageData{
		Title:          title,
		Active:         active,
		Actions:        ai.Actions,
		Prompts:        ai.PromptLibrary,
		Examples:       Examples,
		SelectedAction: "contract-analysis",
		MockMode:       s.mockMode,
		CurrentYear:    time.Now().Year(),
	}
}

func (s *Server) render(w http.ResponseWriter, name string, data PageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.templates.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("template render failed: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}
