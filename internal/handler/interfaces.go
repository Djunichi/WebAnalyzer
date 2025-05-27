package handler

import (
	"WebAnalyzer/internal/dto"
	"context"
	"github.com/google/uuid"
)

type WebPageSvc interface {
	AnalyzePage(context.Context, *dto.AnalyzePageReq) (*dto.AnalyzePageRes, error)
}

type AnalysisSvc interface {
	GetAllAnalyses(context.Context) (*dto.GetAllAnalysesRes, error)
	GetAnalysisById(context.Context, uuid.UUID) (*dto.AnalyzePageRes, error)
}
