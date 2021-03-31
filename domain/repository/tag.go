package repository

import (
	"github.com/champon1020/argus/domain"
	"gorm.io/gorm"
)

// TagRepository is repository interface for tag.
type TagRepository interface {
	FindByName(db *gorm.DB, name string, isPublic bool) (*domain.Tag, error)
	FindAll(db *gorm.DB, isPublic bool) (*[]domain.Tag, error)
	Count(db *gorm.DB, isPublic bool) (int, error)
	Post(db *gorm.DB) error
	Delete(db *gorm.DB) error
}
