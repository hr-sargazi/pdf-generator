package main

import (
	"log"
	"net/http"
	"pdf-service/internal/handlers"
	"pdf-service/internal/infrastructure"
	"pdf-service/internal/services"
)

func main() {
	chromedpClient := infrastructure.NewChromedpClient()

	pdfService := services.NewPDFService(chromedpClient)

	pdfHandler := handlers.NewPDFHandler(pdfService)

	http.HandleFunc("/generate-pdf", pdfHandler.GeneratePDFHandler)

	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}