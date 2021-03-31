package usecase

import (
	"github.com/champon1020/argus/domain"
	"github.com/champon1020/argus/domain/filter"
	"github.com/champon1020/argus/domain/repository"
	"gorm.io/gorm"
)

// ArticleUseCase is usecase interface for article.
type ArticleUseCase interface {
	FindByID(db *gorm.DB, id string) (*domain.Article, error)
}

type articleUseCase struct {
	ar repository.ArticleRepository
}

// NewArticleUseCase creates articleUseCase.
func NewArticleUseCase(ar repository.ArticleRepository) ArticleUseCase {
	return &articleUseCase{ar: ar}
}

func (au articleUseCase) FindPublicByID(db *gorm.DB, id string) (*domain.Article, error) {
	return au.ar.FindByID(db, id, true)
}

func (au articleUseCase) FindPublic(db *gorm.DB, num int, page int) (*[]domain.Article, error) {
	return au.ar.Find(db, num, page, &filter.ArticleFilter{IsPublic: true})
}

func (au articleUseCase) FindPublicByTitle(db *gorm.DB, num int, page int, title string) (*[]domain.Article, error) {
	return au.ar.Find(db, num, page, &filter.ArticleFilter{IsPublic: true, Title: title})
}

func (au articleUseCase) FindPublicByTags(db *gorm.DB, num int, page int, tags []string) (*[]domain.Article, error) {
	return au.ar.Find(db, num, page, &filter.ArticleFilter{IsPublic: true, Tags: tags})
}

func (au articleUseCase) CountPublic(db *gorm.DB) (int, error) {
	return au.ar.Count(db, &filter.ArticleFilter{IsPublic: true})
}

func (au articleUseCase) CountPublicByTitle(db *gorm.DB, title string) (int, error) {
	return au.ar.Count(db, &filter.ArticleFilter{IsPublic: true, Title: title})
}

func (au articleUseCase) CountPublicByTags(db *gorm.DB, tags []string) (int, error) {
	return au.ar.Count(db, &filter.ArticleFilter{IsPublic: true, Tags: tags})
}

func (au articleUseCase) FindByID(db *gorm.DB, id string) (*domain.Article, error) {
	return au.ar.FindByID(db, id, false)
}

func (au articleUseCase) Find(db *gorm.DB, num int, page int) (*[]domain.Article, error) {
	return au.ar.Find(db, num, page, nil)
}

func (au articleUseCase) Count(db *gorm.DB) (int, error) {
	return au.ar.Count(db, nil)
}

func (au articleUseCase) Post(db *gorm.DB, article *domain.Article) error {
	return nil
}

func (au articleUseCase) Update(db *gorm.DB, article *domain.Article) error {
	return nil
}

func (au articleUseCase) Delete(db *gorm.DB, id string) error {
	return nil
}
