package repository

import (
	"github.com/champon1020/argus/domain"
	"gorm.io/gorm"
)

// TagRepository is repository interface for tag.
type TagRepository interface {
	FindByArticleID(db *gorm.DB, articleID string) (*[]domain.Tag, error)
	FindByName(db *gorm.DB, name string, articleStatus *domain.Status) (*[]domain.Tag, error)
	Find(db *gorm.DB, articleStatus *domain.Status) (*[]domain.Tag, error)
	Post(db *gorm.DB, tag *domain.Tag) error
	Delete(db *gorm.DB, articleID string, name string) error
}
