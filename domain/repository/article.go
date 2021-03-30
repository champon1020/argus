package repository

import (
	"github.com/champon1020/argus/domain"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	FindByID(db *gorm.DB, id string) (*domain.Article, error)
}
