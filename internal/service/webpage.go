package service

import (
	"WebAnalyzer/internal/dto"
	"WebAnalyzer/internal/helpers"
	"context"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
)

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

	resp, err := http.Get(parsedURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch: %w", err)
	}
	defer resp.Body.Close()

	reader, err := helpers.GetHTMLReader(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode body: %w", err)
	}

	doc, err := html.Parse(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	result := &dto.AnalyzePageRes{
		StatusCode: resp.StatusCode,
		Headings:   make(map[string]int),
	}

	traverse(doc, result, parsedURL)

	err = w.WebPageRepo.Add(ctx, result, req.Url)
	if err != nil {
		return nil, err
	}

	return result, nil
}
