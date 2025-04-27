package services

import (
	"errors"
	"pdf-service/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockChromedpClient struct {
	mock.Mock
}

func (m *MockChromedpClient) GeneratePDF(htmlContent string) ([]byte, error) {
	args := m.Called(htmlContent)
	return args.Get(0).([]byte), args.Error(1)
}

func TestNewPDFService(t *testing.T) {
	chromedpClient := &MockChromedpClient{}
	service := NewPDFService(chromedpClient)
	assert.NotNil(t, service)
	assert.Equal(t, chromedpClient, service.chromedpClient)
}

func TestGeneratePDF_Success(t *testing.T) {
	chromedpClient := &MockChromedpClient{}
	service := NewPDFService(chromedpClient)

	req := &models.PDFRequest{
		HTMLTemplate: "<html><body>{{.Name}}</body></html>",
		Data:         map[string]interface{}{"Name": "John Doe"},
	}

	expectedPDF := []byte("mocked_pdf_content")
	chromedpClient.On("GeneratePDF", mock.AnythingOfType("string")).Return(expectedPDF, nil)

	pdf, err := service.GeneratePDF(req)
	assert.NoError(t, err)
	assert.Equal(t, expectedPDF, pdf)
	chromedpClient.AssertExpectations(t)
}

func TestGeneratePDF_EmptyTemplate(t *testing.T) {
	chromedpClient := &MockChromedpClient{}
	service := NewPDFService(chromedpClient)

	req := &models.PDFRequest{
		HTMLTemplate: "",
		Data:         map[string]interface{}{"Name": "John Doe"},
	}

	pdf, err := service.GeneratePDF(req)
	assert.Error(t, err)
	assert.Equal(t, ErrEmptyHTMLTemplate, err)
	assert.Nil(t, pdf)
}

func TestGeneratePDF_NilData(t *testing.T) {
	chromedpClient := &MockChromedpClient{}
	service := NewPDFService(chromedpClient)

	req := &models.PDFRequest{
		HTMLTemplate: "<html><body>{{.Name}}</body></html>",
		Data:         nil,
	}

	pdf, err := service.GeneratePDF(req)
	assert.Error(t, err)
	assert.Equal(t, ErrNilData, err)
	assert.Nil(t, pdf)
}

func TestGeneratePDF_InvalidTemplate(t *testing.T) {
	chromedpClient := &MockChromedpClient{}
	service := NewPDFService(chromedpClient)

	req := &models.PDFRequest{
		HTMLTemplate: "<html>{{.InvalidSyntax", // Malformed template
		Data:         map[string]interface{}{"Name": "John Doe"},
	}

	pdf, err := service.GeneratePDF(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "template: dynamic")
	assert.Nil(t, pdf)
}

func TestGeneratePDF_ChromedpError(t *testing.T) {
	chromedpClient := &MockChromedpClient{}
	service := NewPDFService(chromedpClient)

	req := &models.PDFRequest{
		HTMLTemplate: "<html><body>{{.Name}}</body></html>",
		Data:         map[string]interface{}{"Name": "John Doe"},
	}

	chromedpClient.On("GeneratePDF", mock.AnythingOfType("string")).Return([]byte(nil), errors.New("chromedp error"))

	pdf, err := service.GeneratePDF(req)
	assert.Error(t, err)
	assert.Equal(t, "chromedp error", err.Error())
	assert.Nil(t, pdf)
	chromedpClient.AssertExpectations(t)
}