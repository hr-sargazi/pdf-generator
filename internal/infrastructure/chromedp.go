package infrastructure

import (
	"context"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type PDFGenerator interface {
	GeneratePDF(htmlContent string) ([]byte, error)
}

type ChromedpClient struct {
	chromePath string
}

type StatFunc func(string) (os.FileInfo, error)

func NewChromedpClientWithStat(stat StatFunc) *ChromedpClient {
	chromePath := "/usr/bin/chromium-browser"
	if os.Getenv("CHROME_PATH") != "" {
		chromePath = os.Getenv("CHROME_PATH")
	} else if _, err := stat("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"); err == nil {
		chromePath = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	}
	return &ChromedpClient{chromePath: chromePath}
}

func NewChromedpClient() *ChromedpClient {
	return NewChromedpClientWithStat(os.Stat)
}

func (c *ChromedpClient) GeneratePDF(htmlContent string) ([]byte, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(c.chromePath),
		chromedp.Flag("headless", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("disable-dev-shm-usage", true),
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
				WithPaperWidth(8.27).
				WithPaperHeight(11.69).
				Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, err
	}

	return pdfBuffer, nil
}