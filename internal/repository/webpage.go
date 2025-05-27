package repository

import (
	"WebAnalyzer/internal/dto"
	"WebAnalyzer/internal/model"
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WebPageRepository struct {
	db *gorm.DB
}

func NewWebPageRepository(db *gorm.DB) *WebPageRepository {
	return &WebPageRepository{db: db}
}

func (w *WebPageRepository) Add(ctx context.Context, req *dto.AnalyzePageRes, url string) error {
	e := req.ToModel(url)

	err := w.db.WithContext(ctx).Create(e).Error
	if err != nil {
		return fmt.Errorf("[WebPageRepository] failed to create WebpageRequest: %w", err)
	}

	return nil
}

func (w *WebPageRepository) Remove(ctx context.Context, id uuid.UUID) error {
	err := w.db.WithContext(ctx).Delete(&model.WebpageRequest{}, "request_id = ?", id)
	if err != nil {
		return fmt.Errorf("[WebPageRepository] failed to delete WebpageRequest: %w", err)
	}

	return nil
}

func (w *WebPageRepository) GetAll(ctx context.Context) ([]model.Analysis, error) {
	var result []model.Analysis

	err := w.db.WithContext(ctx).Table("webpage_requests").Select("request_id, url, title, created_at").Find(&result).Error
	if err != nil {
		return nil, fmt.Errorf("[WebPageRepository] failed to get all WebPageRequests: %w", err)
	}

	return result, nil
}

func (w *WebPageRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.WebpageRequest, error) {
	var result model.WebpageRequest
	err := w.db.WithContext(ctx).Table("webpage_requests").First(&result, "request_id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("[WebPageRepository] failed to get WebpageRequest: %w", err)
	}
	return &result, nil
}
