package main

import (
	"context"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
)

type RequestData struct {
	CustomerName            string `json:"customer_name"`
	CustomerNumber          string `json:"customer_number"`
	CustomerBankerName      string `json:"customer_banker_name"`
	CustomerHasUbankContract string `json:"customer_has_ubank_contract"`
	ServiceRequestType      string `json:"service_request_type"`
	ServiceRequestTitle     string `json:"service_request_title"`
	ServiceRequestNumber    string `json:"service_request_number"`
	ServiceRequestDate      string `json:"service_request_date"`
	ServiceRequestTime      string `json:"service_request_time"`
	ServiceRequestStatus    string `json:"service_request_status"`
	ServiceRequestDetails   string `json:"service_request_details"`
}

func main() {

	r := gin.Default()

	r.POST("/generate-pdf", func(c *gin.Context) {

		var data RequestData
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
			return
		}


		tmpl, err := template.ParseFiles("templates/service_request.html")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load template: " + err.Error()})
			return
		}


		var renderedHTML strings.Builder
		if err := tmpl.Execute(&renderedHTML, data); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render template: " + err.Error()})
			return
		}


		pdfBuffer, err := generatePDF(renderedHTML.String())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF: " + err.Error()})
			return
		}


		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", "attachment; filename=service_request.pdf")


		c.Writer.Write(pdfBuffer)
	})


	r.Run(":8080")
}


func generatePDF(htmlContent string) ([]byte, error) {

	chromePath := "/usr/bin/chromium-browser"
	if os.Getenv("CHROME_PATH") != "" {
		chromePath = os.Getenv("CHROME_PATH")
	} else if _, err := os.Stat("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"); err == nil {
		chromePath = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	}


	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(chromePath),
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true), // Required in Docker, safe for local testing
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true), // Avoid shared memory issues in Docker
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()


	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()


	var pdfBuffer []byte


	err := chromedp.Run(ctx,

		chromedp.Navigate("about:blank"),

		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			return page.SetDocumentContent(frameTree.Frame.ID, htmlContent).Do(ctx)
		}),

		chromedp.WaitVisible("body", chromedp.ByQuery),

		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuffer, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(8.27).  // A4 width in inches
				WithPaperHeight(11.69). // A4 height in inches
				Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, err
	}

	return pdfBuffer, nil
}