package infrastructure

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChromedpClient_DefaultPath(t *testing.T) {
	os.Unsetenv("CHROME_PATH")

	stat := func(path string) (os.FileInfo, error) {
		return nil, os.ErrNotExist
	}

	client := NewChromedpClientWithStat(stat)
	assert.Equal(t, "/usr/bin/chromium-browser", client.chromePath)
}

func TestNewChromedpClient_EnvPath(t *testing.T) {
	os.Setenv("CHROME_PATH", "/custom/path/to/chrome")
	defer os.Unsetenv("CHROME_PATH")

	// Mock stat function
	stat := func(path string) (os.FileInfo, error) {
		return nil, os.ErrNotExist
	}

	client := NewChromedpClientWithStat(stat)
	assert.Equal(t, "/custom/path/to/chrome", client.chromePath)
}

func TestNewChromedpClient_MacPath(t *testing.T) {
	os.Unsetenv("CHROME_PATH")

	stat := func(path string) (os.FileInfo, error) {
		if path == "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" {
			return nil, nil
		}
		return nil, os.ErrNotExist
	}

	client := NewChromedpClientWithStat(stat)
	assert.Equal(t, "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", client.chromePath)
}

func TestGeneratePDF_Integration(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test; set RUN_INTEGRATION_TESTS=true to run")
	}

	client := NewChromedpClient()
	htmlContent := `<html><body><h1>Hello, World!</h1></body></html>`

	pdf, err := client.GeneratePDF(htmlContent)
	assert.NoError(t, err)
	assert.NotEmpty(t, pdf)

	assert.True(t, len(pdf) > 4 && string(pdf[:4]) == "%PDF")
}