package service

import (
	"WebAnalyzer/internal/dto"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

func traverse(n *html.Node, result *dto.AnalyzePageRes, baseURL *url.URL) {
	if n.Type == html.ErrorNode {
		result.Error = n.Data
		return
	}

	if n.Type == html.DocumentNode {
		if n.FirstChild != nil && n.FirstChild.Type == html.DoctypeNode {
			if strings.EqualFold(n.FirstChild.Data, "html") {
				result.HTMLVersion = "HTML5"
			} else {
				result.HTMLVersion = n.FirstChild.Data
			}
		}
	}

	var hrefs []string

	if n.Type == html.ElementNode {
		switch n.Data {
		case "title":
			if n.FirstChild != nil {
				result.Title = n.FirstChild.Data
			}
		case "h1", "h2", "h3", "h4", "h5", "h6":
			result.Headings[n.Data]++
		case "a":
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link := strings.TrimSpace(attr.Val)
					if link == "" || strings.HasPrefix(link, "#") {
						break
					}

					if isInternalLink(link, baseURL) {
						result.InternalLinks++
					} else {
						result.ExternalLinks++
					}
					hrefs = append(hrefs, link)
				}
			}
		case "input":
			for _, attr := range n.Attr {
				if attr.Key == "type" && strings.ToLower(attr.Val) == "password" {
					result.HasLoginForm = true
				}
			}
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		traverse(child, result, baseURL)
	}

	if len(hrefs) > 0 {
		result.InaccessibleLinks += checkLinksConcurrently(hrefs, baseURL)
	}
}

func detectHTMLVersion(rawHTML string) string {
	doctypeRe := regexp.MustCompile(`(?i)<!DOCTYPE\s+([^>]+)>`)
	matches := doctypeRe.FindStringSubmatch(rawHTML)
	if len(matches) < 2 {
		return ""
	}

	dt := strings.ToLower(matches[1])

	switch {
	case strings.Contains(dt, "html 4.01"):
		return "HTML 4.01"
	case strings.Contains(dt, "xhtml"):
		return "XHTML"
	case strings.Contains(dt, "html"):
		return "HTML5"
	default:
		return dt
	}
}

// isInternalLink checks is link is internal
func isInternalLink(href string, baseURL *url.URL) bool {
	linkURL, err := url.Parse(href)
	if err != nil {
		return false
	}

	if linkURL.Host == "" && strings.HasPrefix(href, "/") {
		return true
	}

	if linkURL.Host == baseURL.Host {
		return true
	}

	return false
}

var (
	linkCache  = sync.Map{}
	httpClient = &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       100,
			IdleConnTimeout:    90 * time.Second,
			DisableCompression: false,
		},
	}
)

// isLinkAccessibleCached checks if a link is accessible with caching
func isLinkAccessibleCached(href string, baseURL *url.URL) bool {
	cacheKey := href
	if val, ok := linkCache.Load(cacheKey); ok {
		return val.(bool)
	}

	if !isCheckableLink(href) {
		return true // skip non-checkable links
	}

	linkURL, err := url.Parse(href)
	if err != nil {
		linkCache.Store(cacheKey, false)
		return false
	}
	if linkURL.Host == "" {
		linkURL.Scheme = baseURL.Scheme
		linkURL.Host = baseURL.Host
	}

	req, err := http.NewRequest("GET", linkURL.String(), nil)
	if err != nil {
		linkCache.Store(cacheKey, false)
		return false
	}
	req.Header.Set("Range", "bytes=0-0")

	resp, err := httpClient.Do(req)
	if err != nil {
		linkCache.Store(cacheKey, false)
		return false
	}
	defer io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	accessible := resp.StatusCode >= 200 && resp.StatusCode < 400
	linkCache.Store(cacheKey, accessible)
	return accessible
}

func isCheckableLink(link string) bool {
	l := strings.ToLower(strings.TrimSpace(link))
	return l != "" &&
		!strings.HasPrefix(l, "mailto:") &&
		!strings.HasPrefix(l, "javascript:") &&
		!strings.HasPrefix(l, "tel:") &&
		!strings.HasPrefix(l, "#") &&
		!strings.HasPrefix(l, "data:") &&
		!strings.HasPrefix(l, "ftp:") &&
		!strings.HasPrefix(l, "ws:") &&
		!strings.HasPrefix(l, "wss:")
}

func checkLinksConcurrently(links []string, baseURL *url.URL) int {
	var inaccessibleCount int
	var wg sync.WaitGroup
	limit := runtime.NumCPU() * 16
	sem := make(chan struct{}, limit)
	results := make(chan bool, len(links))

	for _, href := range links {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			if !isLinkAccessibleCached(link, baseURL) {
				results <- true
			} else {
				results <- false
			}
		}(href)
	}

	wg.Wait()
	close(results)

	for r := range results {
		if r {
			inaccessibleCount++
		}
	}

	return inaccessibleCount
}
