package usecase

import (
	"github.com/champon1020/argus/domain"
	"github.com/champon1020/argus/domain/repository"
	"gorm.io/gorm"
)

type ArticleUseCase interface {
	FindByID(db *gorm.DB, id string) (*domain.Article, error)
}

type articleUseCase struct {
	articleRepository repository.ArticleRepository
}

func NewArticleUseCase(ar repository.ArticleRepository) ArticleUseCase {
	return &articleUseCase{
		articleRepository: ar,
	}
}

func (au articleUseCase) FindByID(db *gorm.DB, id string) (*domain.Article, error) {
	article, err := au.articleRepository.FindByID(db, id)
	if err != nil {
		return nil, err
	}
	return article, nil
}
