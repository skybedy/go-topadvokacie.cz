package web

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
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
	Prompts         []ai.PromptTemplate
	Examples        []model.Example
	SelectedAction  string
	DetailLevel     string
	Perspective     string
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
	mux.HandleFunc("/upload-text", s.uploadText)
	return securityHeaders(mux)
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := s.baseData("LexPilot Demo", "home")
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
		if err := r.ParseMultipartForm(8 << 20); err != nil {
			data.Error = "Formulář se nepodařilo zpracovat."
			s.render(w, "home.html", data)
			return
		}

		data.SelectedAction = r.FormValue("action")
		data.DetailLevel = normalizeDetailLevel(r.FormValue("detail_level"))
		data.Perspective = normalizePerspective(r.FormValue("perspective"))
		data.InputA = strings.TrimSpace(r.FormValue("input_a"))
		data.InputB = strings.TrimSpace(r.FormValue("input_b"))
		data.SelectedExample = r.FormValue("example")
		data.SelectedAction = ai.PromptTemplateByIDOrDefault(data.SelectedAction).ID

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
		if ai.PromptTemplateByIDOrDefault(data.SelectedAction).NeedsSecond && data.InputB == "" {
			data.Error = "Pro porovnání dvou verzí vložte také druhý dokument."
			s.render(w, "home.html", data)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 180*time.Second)
		defer cancel()
		result, err := s.ai.Analyze(ctx, data.SelectedAction, data.InputA, data.InputB, ai.Options{
			DetailLevel: data.DetailLevel,
			Perspective: data.Perspective,
		})
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

func (s *Server) uploadText(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/upload-text" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(8 << 20); err != nil {
		writeUploadJSON(w, http.StatusBadRequest, "", "Soubor se nepodařilo zpracovat.")
		return
	}
	text, err := readUploadedText(r, "file")
	if err != nil {
		writeUploadJSON(w, http.StatusBadRequest, "", err.Error())
		return
	}
	if strings.TrimSpace(text) == "" {
		writeUploadJSON(w, http.StatusBadRequest, "", "Soubor neobsahuje žádný čitelný text.")
		return
	}
	writeUploadJSON(w, http.StatusOK, text, "")
}

func writeUploadJSON(w http.ResponseWriter, status int, text string, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(struct {
		Text  string `json:"text,omitempty"`
		Error string `json:"error,omitempty"`
	}{
		Text:  text,
		Error: message,
	}); err != nil {
		log.Printf("upload json encode failed: %v", err)
	}
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

	const maxUploadBytes = 5 * 1024 * 1024
	content, err := io.ReadAll(io.LimitReader(file, maxUploadBytes+1))
	if err != nil {
		return "", fmt.Errorf("Soubor %s se nepodařilo přečíst.", header.Filename)
	}
	if len(content) > maxUploadBytes {
		return "", fmt.Errorf("Soubor %s je příliš velký. Limit pro demo je 5 MB.", header.Filename)
	}

	switch ext {
	case ".txt", ".md", ".markdown", ".csv", ".rst", ".log":
		return strings.TrimSpace(string(content)), nil
	case ".docx":
		text, err := extractDOCXText(content)
		if err != nil {
			return "", fmt.Errorf("Soubor %s se nepodařilo převést z DOCX na text.", header.Filename)
		}
		return ensureExtractedText(header.Filename, text)
	case ".pdf":
		text, err := extractPDFText(content)
		if err != nil {
			return "", err
		}
		return ensureExtractedText(header.Filename, text)
	case ".doc":
		return "", fmt.Errorf("Starší .doc zatím nepodporuji. Uložte dokument jako .docx nebo PDF a nahrajte jej znovu.")
	default:
		return "", fmt.Errorf("Soubor %s zatím neumím načíst. Nahrajte .txt, .md, .csv, .docx nebo .pdf.", header.Filename)
	}
}

func ensureExtractedText(filename string, text string) (string, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return "", fmt.Errorf("Ze souboru %s se nepodařilo vytěžit žádný text.", filename)
	}
	const maxExtractedTextBytes = 512 * 1024
	if len(text) > maxExtractedTextBytes {
		text = text[:maxExtractedTextBytes]
	}
	return text, nil
}

func extractDOCXText(content []byte) (string, error) {
	reader, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return "", err
	}

	var parts []string
	for _, file := range reader.File {
		name := strings.ToLower(file.Name)
		if name != "word/document.xml" && !strings.HasPrefix(name, "word/header") && !strings.HasPrefix(name, "word/footer") {
			continue
		}
		text, err := extractWordXMLText(file)
		if err != nil {
			return "", err
		}
		if strings.TrimSpace(text) != "" {
			parts = append(parts, text)
		}
	}
	return strings.Join(parts, "\n\n"), nil
}

func extractWordXMLText(file *zip.File) (string, error) {
	rc, err := file.Open()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	decoder := xml.NewDecoder(rc)
	var b strings.Builder
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		switch t := token.(type) {
		case xml.CharData:
			text := string(t)
			if strings.TrimSpace(text) != "" {
				b.WriteString(text)
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "p", "br", "tab":
				if b.Len() > 0 && !strings.HasSuffix(b.String(), "\n") {
					b.WriteString("\n")
				}
			}
		}
	}
	return b.String(), nil
}

func extractPDFText(content []byte) (string, error) {
	if _, err := exec.LookPath("pdftotext"); err != nil {
		return "", fmt.Errorf("PDF parser není v systému dostupný. Nainstalujte balíček poppler-utils nebo nahrajte textový soubor.")
	}

	tmp, err := os.CreateTemp("", "lexpilot-*.pdf")
	if err != nil {
		return "", fmt.Errorf("PDF se nepodařilo připravit ke čtení.")
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(content); err != nil {
		tmp.Close()
		return "", fmt.Errorf("PDF se nepodařilo uložit pro parser.")
	}
	if err := tmp.Close(); err != nil {
		return "", fmt.Errorf("PDF se nepodařilo uzavřít pro parser.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "pdftotext", "-layout", "-enc", "UTF-8", tmp.Name(), "-")
	output, err := cmd.Output()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("PDF parser nestihl soubor zpracovat v časovém limitu.")
	}
	if err != nil {
		return "", fmt.Errorf("PDF se nepodařilo převést na text. Zkuste jiné PDF nebo dokument uložte jako .txt/.docx.")
	}
	return string(output), nil
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

func normalizeDetailLevel(value string) string {
	switch value {
	case "brief", "standard", "detailed":
		return value
	default:
		return "standard"
	}
}

func normalizePerspective(value string) string {
	switch value {
	case "lawyer", "client", "negotiation":
		return value
	default:
		return "lawyer"
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
		Prompts:        ai.PromptLibrary,
		Examples:       Examples,
		SelectedAction: ai.DefaultPromptID,
		DetailLevel:    "standard",
		Perspective:    "lawyer",
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
