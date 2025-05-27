package model

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"time"
)

type WebpageRequest struct {
	RequestID               uuid.UUID         `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:request_id"`
	URL                     string            `gorm:"not null;type:text;column:url"`
	StatusCode              int               `gorm:"not null;column:status_code"`
	HTMLVersion             *string           `gorm:"type:text;column:html_version"`
	Title                   *string           `gorm:"type:text;column:title"`
	Headings                datatypes.JSONMap `gorm:"type:jsonb;column:headings"` // map[string]int as JSONB
	InternalLinksNumber     *int              `gorm:"column:internal_links_number"`
	ExternalLinksNumber     *int              `gorm:"column:external_links_number"`
	InaccessibleLinksNumber *int              `gorm:"column:inaccessible_links_number"`
	ContainsLoginForm       bool              `gorm:"not null;column:contains_login_form"`
	ErrorDescription        *string           `gorm:"type:text;column:error_description"`
	CreatedAt               time.Time         `gorm:"not null;default:now();column:created_at"`
}
