package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"pdf-service/internal/models"
	"pdf-service/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPDFService struct {
	mock.Mock
}

func (m *MockPDFService) GeneratePDF(req *models.PDFRequest) ([]byte, error) {
	args := m.Called(req)
	return args.Get(0).([]byte), args.Error(1)
}

func TestNewPDFHandler(t *testing.T) {
	pdfService := &MockPDFService{}
	handler := NewPDFHandler(pdfService)
	assert.NotNil(t, handler)
	assert.Equal(t, pdfService, handler.pdfService)
}

func TestGeneratePDFHandler_MethodNotAllowed(t *testing.T) {
	pdfService := &MockPDFService{}
	handler := NewPDFHandler(pdfService)

	req := httptest.NewRequest(http.MethodGet, "/generate-pdf", nil)
	rr := httptest.NewRecorder()

	handler.GeneratePDFHandler(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	assert.Equal(t, "Method not allowed\n", rr.Body.String())
}

func TestGeneratePDFHandler_Success(t *testing.T) {
	pdfService := &MockPDFService{}
	handler := NewPDFHandler(pdfService)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormFile("template_file", "template.html")
	part.Write([]byte("<html><body>{{.Name}}</body></html>"))

	data := map[string]interface{}{"Name": "John Doe"}
	dataBytes, _ := json.Marshal(data)
	writer.WriteField("data", string(dataBytes))

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/generate-pdf", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	expectedPDF := []byte("%PDF-1.4 mock")
	pdfService.On("GeneratePDF", mock.AnythingOfType("*models.PDFRequest")).Return(expectedPDF, nil)

	handler.GeneratePDFHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/pdf", rr.Header().Get("Content-Type"))
	assert.Equal(t, `attachment; filename=dynamic_document.pdf`, rr.Header().Get("Content-Disposition"))
	assert.Equal(t, expectedPDF, rr.Body.Bytes())
	pdfService.AssertExpectations(t)
}

func TestGeneratePDFHandler_MissingTemplateFile(t *testing.T) {
	pdfService := &MockPDFService{}
	handler := NewPDFHandler(pdfService)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("data", `{"Name":"John Doe"}`)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/generate-pdf", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	handler.GeneratePDFHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to get template file")
}

func TestGeneratePDFHandler_MissingData(t *testing.T) {
	pdfService := &MockPDFService{}
	handler := NewPDFHandler(pdfService)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("template_file", "template.html")
	part.Write([]byte("<html><body>{{.Name}}</body></html>"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/generate-pdf", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	handler.GeneratePDFHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "Data field is required\n", rr.Body.String())
}

func TestGeneratePDFHandler_InvalidJSONData(t *testing.T) {
	pdfService := &MockPDFService{}
	handler := NewPDFHandler(pdfService)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("template_file", "template.html")
	part.Write([]byte("<html><body>{{.Name}}</body></html>"))
	writer.WriteField("data", "invalid_json")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/generate-pdf", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	handler.GeneratePDFHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid JSON data")
}

func TestGeneratePDFHandler_PDFServiceAppError(t *testing.T) {
	pdfService := &MockPDFService{}
	handler := NewPDFHandler(pdfService)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("template_file", "template.html")
	part.Write([]byte("<html><body>{{.Name}}</body></html>"))
	writer.WriteField("data", `{"Name":"John Doe"}`)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/generate-pdf", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	pdfService.On("GeneratePDF", mock.AnythingOfType("*models.PDFRequest")).Return([]byte(nil), &services.AppError{Message: "Invalid template"})

	handler.GeneratePDFHandler(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "Invalid template\n", rr.Body.String())
	pdfService.AssertExpectations(t)
}

func TestGeneratePDFHandler_PDFServiceGenericError(t *testing.T) {
	pdfService := &MockPDFService{}
	handler := NewPDFHandler(pdfService)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("template_file", "template.html")
	part.Write([]byte("<html><body>{{.Name}}</body></html>"))
	writer.WriteField("data", `{"Name":"John Doe"}`)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/generate-pdf", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	pdfService.On("GeneratePDF", mock.AnythingOfType("*models.PDFRequest")).Return([]byte(nil), errors.New("internal error"))

	handler.GeneratePDFHandler(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to generate PDF: internal error")
	pdfService.AssertExpectations(t)
}