package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"pdf-service/internal/models"
	"pdf-service/internal/services"
)

type PDFHandler struct {
	pdfService *services.PDFService
}

func NewPDFHandler(pdfService *services.PDFService) *PDFHandler {
	return &PDFHandler{pdfService: pdfService}
}

func (h *PDFHandler) GeneratePDFHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form (set a reasonable max memory limit, e.g., 10MB)
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "Failed to parse multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("template_file")
	if err != nil {
		http.Error(w, "Failed to get template file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	htmlBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read template file: "+err.Error(), http.StatusBadRequest)
		return
	}
	htmlTemplate := string(htmlBytes)

	dataStr := r.FormValue("data")
	if dataStr == "" {
		http.Error(w, "Data field is required", http.StatusBadRequest)
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		http.Error(w, "Invalid JSON data: "+err.Error(), http.StatusBadRequest)
		return
	}

	req := &models.PDFRequest{
		HTMLTemplate: htmlTemplate,
		Data:         data,
	}

	pdfBuffer, err := h.pdfService.GeneratePDF(req)
	if err != nil {
		if appErr, ok := err.(*services.AppError); ok {
			http.Error(w, appErr.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to generate PDF: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=dynamic_document.pdf")

	if _, err := w.Write(pdfBuffer); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}