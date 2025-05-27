package service

import (
	"WebAnalyzer/internal/dto"
	"WebAnalyzer/internal/model"
	"context"
	"github.com/google/uuid"
)

type WebPageRepo interface {
	Add(context.Context, *dto.AnalyzePageRes) error
	Remove(context.Context, uuid.UUID) error
	GetAll(context.Context) ([]model.Analysis, error)
	GetByID(context.Context, uuid.UUID) (*model.WebpageRequest, error)
}
