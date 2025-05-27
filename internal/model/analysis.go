package model

import (
	"github.com/google/uuid"
	"time"
)

type Analysis struct {
	Id            uuid.UUID `json:"id" gorm:"column:request_id"`
	Url           string    `json:"url" gorm:"column:url"`
	Title         string    `json:"title" gorm:"column:title"`
	TimeRequested time.Time `json:"timeRequested" gorm:"column:created_at"`
}
