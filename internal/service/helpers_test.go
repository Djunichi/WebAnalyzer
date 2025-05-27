package service

import (
	"WebAnalyzer/internal/dto"
	"golang.org/x/net/html"
	"net/url"
	"strings"
	"testing"
)

func TestIsCheckableLink(t *testing.T) {
	cases := []struct {
		link     string
		expected bool
	}{
		{"mailto:someone@example.com", false},
		{"javascript:void(0)", false},
		{"tel:+123456789", false},
		{"#section", false},
		{"data:image/png;base64,...", false},
		{"https://example.com", true},
		{"/relative/path", true},
	}

	for _, c := range cases {
		got := isCheckableLink(c.link)
		if got != c.expected {
			t.Errorf("isCheckableLink(%q) = %v; want %v", c.link, got, c.expected)
		}
	}
}

func TestDetectHTMLVersion(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{`<!DOCTYPE html><html></html>`, "HTML5"},
		{`<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">`, "HTML 4.01"},
		{`<!DOCTYPE XHTML 1.0 Strict>`, "XHTML"},
		{`<html><body>No doctype</body></html>`, ""},
	}

	for _, c := range cases {
		got := detectHTMLVersion(c.input)
		if got != c.expected {
			t.Errorf("detectHTMLVersion() = %q; want %q", got, c.expected)
		}
	}
}

func TestExtractHeadingsAndLoginForm(t *testing.T) {
	htmlContent := `
    <!DOCTYPE html>
    <html>
    <head><title>Test</title></head>
    <body>
        <h1>Main Heading</h1>
        <h2>Sub Heading</h2>
        <form>
            <input type="text" name="user">
            <input type="password" name="pass">
        </form>
    </body>
    </html>`

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("failed to parse html: %v", err)
	}

	parsedURL, _ := url.Parse("https://example.com")
	result := &dto.AnalyzePageRes{
		Headings: make(map[string]int),
	}

	traverse(doc, result, parsedURL)

	if result.Headings["h1"] != 1 || result.Headings["h2"] != 1 {
		t.Errorf("Headings parsed incorrectly: %v", result.Headings)
	}

	if !result.HasLoginForm {
		t.Error("Expected login form to be detected")
	}
}

func TestIsLinkAccessibleCached_SkipsInvalid(t *testing.T) {
	base, _ := url.Parse("https://example.com")
	if isLinkAccessibleCached("mailto:someone@example.com", base) {
		t.Log("Correctly skipped mailto:")
	} else {
		t.Error("Expected mailto to be considered accessible (skipped)")
	}
}
