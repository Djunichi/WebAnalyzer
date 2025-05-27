package repository

import "gorm.io/gorm"

type Container struct {
	WebPages *WebPageRepository
}

func NewRepositoryContainer(db *gorm.DB) *Container {
	return &Container{
		WebPages: NewWebPageRepository(db),
	}
}
