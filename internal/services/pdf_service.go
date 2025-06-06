package services

import (
	"bytes"
	"html/template"
	"pdf-service/internal/infrastructure"
	"pdf-service/internal/models"
)


type PDFServiceInterface interface {
	GeneratePDF(req *models.PDFRequest) ([]byte, error)
}
type PDFService struct {
	chromedpClient infrastructure.PDFGenerator
}

func NewPDFService(chromedpClient infrastructure.PDFGenerator) *PDFService {
	return &PDFService{chromedpClient: chromedpClient}
}

func (s *PDFService) GeneratePDF(req *models.PDFRequest) ([]byte, error) {
	if req.HTMLTemplate == "" {
		return nil, ErrEmptyHTMLTemplate
	}
	if req.Data == nil {
		return nil, ErrNilData
	}

	tmpl, err := template.New("dynamic").Parse(req.HTMLTemplate)
	if err != nil {
		return nil, err
	}

	var renderedHTML bytes.Buffer
	if err := tmpl.Execute(&renderedHTML, req.Data); err != nil {
		return nil, err
	}

	return s.chromedpClient.GeneratePDF(renderedHTML.String())
}

var (
	ErrEmptyHTMLTemplate = &AppError{Message: "HTML template cannot be empty"}
	ErrNilData           = &AppError{Message: "Data cannot be nil"}
)

type AppError struct {
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}