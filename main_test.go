package main

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"pdf-service/internal/handlers"
	"pdf-service/internal/models"
	"testing"
	"io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPDFGenerator struct {
	mock.Mock
}

func (m *MockPDFGenerator) GeneratePDF(htmlContent string) ([]byte, error) {
	args := m.Called(htmlContent)
	return args.Get(0).([]byte), args.Error(1)
}

type MockPDFService struct {
	mock.Mock
}

func (m *MockPDFService) GeneratePDF(req *models.PDFRequest) ([]byte, error) {
	args := m.Called(req)
	return args.Get(0).([]byte), args.Error(1)
}

func TestMainHandler(t *testing.T) {
	pdfService := &MockPDFService{}
	pdfHandler := handlers.NewPDFHandler(pdfService)

	expectedPDF := []byte("%PDF-1.4 mock")
	pdfService.On("GeneratePDF", mock.AnythingOfType("*models.PDFRequest")).Return(expectedPDF, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("/generate-pdf", pdfHandler.GeneratePDFHandler)
	server := httptest.NewServer(mux)
	defer server.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("template_file", "template.html")
	part.Write([]byte("<html><body>{{.Name}}</body></html>"))
	data := map[string]interface{}{"Name": "John Doe"}
	dataBytes, _ := json.Marshal(data)
	writer.WriteField("data", string(dataBytes))
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+"/generate-pdf", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/pdf", resp.Header.Get("Content-Type"))
	assert.Equal(t, `attachment; filename=dynamic_document.pdf`, resp.Header.Get("Content-Disposition"))

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, expectedPDF, bodyBytes)

	pdfService.AssertExpectations(t)
}

func TestMainHandler_InvalidMethod(t *testing.T) {
	pdfService := &MockPDFService{}
	pdfHandler := handlers.NewPDFHandler(pdfService)

	mux := http.NewServeMux()
	mux.HandleFunc("/generate-pdf", pdfHandler.GeneratePDFHandler)
	server := httptest.NewServer(mux)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL+"/generate-pdf", nil)
	assert.NoError(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Method not allowed\n", string(bodyBytes))
}