package repository

import (
	"github.com/champon1020/argus/domain"
	"github.com/champon1020/argus/domain/filter"
	"gorm.io/gorm"
)

// ArticleRepository is repository interface for article.
type ArticleRepository interface {
	FindByID(db *gorm.DB, id string, isPublic bool) (*domain.Article, error)
	Find(db *gorm.DB, limit int, offset int, filter *filter.ArticleFilter) (*[]domain.Article, error)
	Count(db *gorm.DB, filter *filter.ArticleFilter) (int, error)
	Post(db *gorm.DB, article *domain.Article) error
	Update(db *gorm.DB, article *domain.Article) error
	Delete(db *gorm.DB, id string) error
}
