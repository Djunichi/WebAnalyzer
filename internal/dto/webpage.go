package dto

import (
	"WebAnalyzer/internal/helpers"
	"WebAnalyzer/internal/model"
	"net/http"
)

type AnalyzePageReq struct {
	Url string
}

type AnalyzePageRes struct {
	StatusCode        int            `json:"statusCode"`
	Url               string         `json:"url"`
	HTMLVersion       string         `json:"HTMLVersion"`
	Title             string         `json:"title"`
	Headings          map[string]int `json:"headings"`
	InternalLinks     int            `json:"internalLinks"`
	ExternalLinks     int            `json:"externalLinks"`
	InaccessibleLinks int            `json:"inaccessible_links"`
	HasLoginForm      bool           `json:"hasLoginForm"`
	Error             string         `json:"error,omitempty"`
}

func (a AnalyzePageRes) ToModel() *model.WebpageRequest {
	return &model.WebpageRequest{
		URL:                     a.Url,
		StatusCode:              a.StatusCode,
		HTMLVersion:             &a.HTMLVersion,
		Title:                   &a.Title,
		Headings:                helpers.ToJSONMap(a.Headings),
		InternalLinksNumber:     &a.InternalLinks,
		ExternalLinksNumber:     &a.ExternalLinks,
		InaccessibleLinksNumber: &a.InaccessibleLinks,
		ContainsLoginForm:       a.HasLoginForm,
		ErrorDescription:        &a.Error,
	}
}

func (a AnalyzePageRes) FromModel(model *model.WebpageRequest) (*AnalyzePageRes, error) {

	res := &AnalyzePageRes{
		StatusCode: model.StatusCode,
	}

	if res.StatusCode != http.StatusOK {
		res.Error = *model.ErrorDescription
	} else {
		res.HTMLVersion = *model.HTMLVersion
		res.Url = model.URL
		res.Title = *model.Title
		res.InternalLinks = *model.InternalLinksNumber
		res.ExternalLinks = *model.ExternalLinksNumber
		res.InaccessibleLinks = *model.InaccessibleLinksNumber
		res.HasLoginForm = model.ContainsLoginForm
	}

	h, err := helpers.FromJSONMap(model.Headings)
	if err != nil {
		return nil, err
	}

	res.Headings = h

	return res, nil
}

type GetAllAnalysesRes struct {
	Analyses []model.Analysis `json:"analyses"`
}
