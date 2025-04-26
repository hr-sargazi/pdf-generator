package models

type PDFRequest struct {
	HTMLTemplate string                 `json:"html_template"`
	Data         map[string]interface{} `json:"data"`
}