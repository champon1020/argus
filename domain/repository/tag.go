package repository

import (
	"github.com/champon1020/argus/domain"
	"gorm.io/gorm"
)

// TagRepository is repository interface for tag.
type TagRepository interface {
	FindByArticleID(db *gorm.DB, articleID string) (*[]domain.Tag, error)
	Find(db *gorm.DB, articleStatus *domain.Status) (*[]domain.Tag, error)
	Posts(db *gorm.DB, tags *[]domain.Tag, articleID string) error
	DeleteByArticleID(db *gorm.DB, articleID string) error
}
