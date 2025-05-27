package service

import (
	"WebAnalyzer/internal/dto"
	"context"
	"github.com/google/uuid"
)

type AnalysisService struct {
	WebPageRepo WebPageRepo
}

func NewAnalysisService(webPageRepo WebPageRepo) *AnalysisService {
	return &AnalysisService{
		WebPageRepo: webPageRepo,
	}
}

func (a *AnalysisService) GetAllAnalyses(ctx context.Context) (*dto.GetAllAnalysesRes, error) {
	list, err := a.WebPageRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	result := &dto.GetAllAnalysesRes{
		Analyses: list,
	}

	return result, nil
}

func (a *AnalysisService) GetAnalysisById(ctx context.Context, id uuid.UUID) (*dto.AnalyzePageRes, error) {
	res, err := a.WebPageRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	result := dto.AnalyzePageRes{}.FromModel(res)

	return result, nil
}
