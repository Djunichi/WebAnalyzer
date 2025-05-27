package service

import (
	"WebAnalyzer/internal/dto"
	"WebAnalyzer/internal/helpers"
	"context"
	"fmt"
	"github.com/dgraph-io/ristretto"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

var cache *ristretto.Cache

func init() {
	var err error
	cache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e4,     // ~10x of expected elements
		MaxCost:     1 << 20, // ~1 MB overall weight
		BufferItems: 64,
	})
	if err != nil {
		log.Fatalf("failed to initialize cache: %v", err)
	}
}

type WebPageService struct {
	WebPageRepo WebPageRepo
}

func NewWebPageService(webPageRepo WebPageRepo) *WebPageService {
	return &WebPageService{
		WebPageRepo: webPageRepo,
	}
}

func (w *WebPageService) AnalyzePage(ctx context.Context, req *dto.AnalyzePageReq) (*dto.AnalyzePageRes, error) {
	parsedURL, err := url.ParseRequestURI(req.Url)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// cache
	if val, found := cache.Get(req.Url); found {
		if entry, ok := val.(*dto.AnalyzePageRes); ok {
			return entry, nil
		}
	}

	result := &dto.AnalyzePageRes{
		Url:      req.Url,
		Headings: make(map[string]int),
	}

	resp, err := http.Get(parsedURL.String())
	if err != nil {
		result.Error = err.Error()
		result.StatusCode = 0

		err = w.saveWebpageRequest(ctx, result)
		if err != nil {
			return nil, err
		}

		return result, nil
	}
	result.StatusCode = resp.StatusCode

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	resp.Body.Close()

	htmlStr := string(bodyBytes)

	reader, err := helpers.GetHTMLReader(bodyBytes, resp.Header.Get("Content-Type"))
	if reader == nil {
		result.Error = resp.Status

		err = w.saveWebpageRequest(ctx, result)
		if err != nil {
			return nil, err
		}
	}

	doc, err := html.Parse(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	result.HTMLVersion = detectHTMLVersion(htmlStr)

	traverse(doc, result, parsedURL)

	err = w.saveWebpageRequest(ctx, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (w *WebPageService) saveWebpageRequest(ctx context.Context, result *dto.AnalyzePageRes) error {

	err := w.WebPageRepo.Add(ctx, result)
	if err != nil {
		return err
	}

	//save to cache
	cache.SetWithTTL(result.Url, result, 1, 24*time.Hour)

	return nil
}
