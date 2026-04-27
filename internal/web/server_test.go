package web

import (
	"archive/zip"
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExtractDOCXText(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create("word/document.xml")
	if err != nil {
		t.Fatal(err)
	}
	_, err = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
    <w:p><w:r><w:t>Prvni odstavec.</w:t></w:r></w:p>
    <w:p><w:r><w:t>Druhy odstavec.</w:t></w:r></w:p>
  </w:body>
</w:document>`))
	if err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}

	text, err := extractDOCXText(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(text, "Prvni odstavec.") || !strings.Contains(text, "Druhy odstavec.") {
		t.Fatalf("expected DOCX text, got %q", text)
	}
}

func TestEnsureExtractedTextRejectsEmptyText(t *testing.T) {
	_, err := ensureExtractedText("empty.pdf", "   \n\t")
	if err == nil {
		t.Fatal("expected empty extracted text error")
	}
}

func TestUploadTextEndpointReturnsExtractedText(t *testing.T) {
	server := &Server{}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "contract.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte("Text smlouvy")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/upload-text", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	server.uploadText(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Text smlouvy") {
		t.Fatalf("expected extracted text in response, got %s", rec.Body.String())
	}
}
