package service

import "WebAnalyzer/internal/repository"

type Container struct {
	WebPageSvc  *WebPageService
	AnalysisSvc *AnalysisService
}

func NewServiceContainer(repos *repository.Container) *Container {
	return &Container{
		WebPageSvc:  NewWebPageService(repos.WebPages),
		AnalysisSvc: NewAnalysisService(repos.WebPages),
	}
}
