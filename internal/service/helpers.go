package service

import (
	"WebAnalyzer/internal/dto"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
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

	// Asynchronous checking of collected links
	if len(hrefs) > 0 {
		result.InaccessibleLinks += checkLinksConcurrently(hrefs, baseURL, 10)
	}
}

func checkLinksConcurrently(links []string, baseURL *url.URL, limit int) int {
	var inaccessibleCount int
	var wg sync.WaitGroup
	sem := make(chan struct{}, limit)
	results := make(chan bool, len(links))

	for _, href := range links {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			if !isLinkAccessible(link, baseURL) {
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

// isLinkAccessible checks is link accessible
func isLinkAccessible(href string, baseURL *url.URL) bool {
	linkURL, err := url.Parse(href)
	if err != nil {
		return false
	}

	if linkURL.Host == "" {
		linkURL.Scheme = baseURL.Scheme
		linkURL.Host = baseURL.Host
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Head(linkURL.String())
	if err != nil {
		return false
	}
	defer io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 400
}
