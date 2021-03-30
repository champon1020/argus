package persistence

import (
	"github.com/champon1020/argus/domain"
	"github.com/champon1020/argus/domain/repository"
	"gorm.io/gorm"
)

type articlePersistence struct{}

func NewArticlePersistence() repository.ArticleRepository {
	return &articlePersistence{}
}

func (ap *articlePersistence) FindByID(db *gorm.DB, id string) (*domain.Article, error) {
	return nil, nil
}
