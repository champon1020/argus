package usecase

import (
	"github.com/champon1020/argus/domain"
	"github.com/champon1020/argus/domain/filter"
	"github.com/champon1020/argus/domain/repository"
	"github.com/champon1020/argus/usecase/pagenation"
	"gorm.io/gorm"
)

// ArticleUseCase is usecase interface for article.
type ArticleUseCase interface {
	FindPublicByID(db *gorm.DB, id string) (*domain.Article, error)
	FindPublic(db *gorm.DB, p pagenation.Pagenation) (*[]domain.Article, error)
	FindPublicByTitle(db *gorm.DB, p pagenation.Pagenation, title string) (*[]domain.Article, error)
	FindPublicByTag(db *gorm.DB, p pagenation.Pagenation, tag string) (*[]domain.Article, error)
	CountPublic(db *gorm.DB) (int, error)
	CountPublicByTitle(db *gorm.DB, title string) (int, error)
	CountPublicByTag(db *gorm.DB, tag string) (int, error)
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
	return au.ar.FindByID(db, id, &domain.Public)
}

func (au articleUseCase) FindPublic(db *gorm.DB, p pagenation.Pagenation) (*[]domain.Article, error) {
	filter := &filter.ArticleFilter{
		Status: &domain.Public,
	}
	return au.ar.Find(db, p.Limit, (p.Page-1)*p.Limit, filter)
}

func (au articleUseCase) FindPublicByTitle(db *gorm.DB, p pagenation.Pagenation, title string) (*[]domain.Article, error) {
	filter := &filter.ArticleFilter{
		Status: &domain.Public,
		Title:  &title,
	}
	return au.ar.Find(db, p.Limit, (p.Page-1)*p.Limit, filter)
}

func (au articleUseCase) FindPublicByTag(db *gorm.DB, p pagenation.Pagenation, tag string) (*[]domain.Article, error) {
	filter := &filter.ArticleFilter{
		Status: &domain.Public,
		Tags:   []string{tag},
	}
	return au.ar.Find(db, p.Limit, (p.Page-1)*p.Limit, filter)
}

func (au articleUseCase) CountPublic(db *gorm.DB) (int, error) {
	filter := &filter.ArticleFilter{
		Status: &domain.Public,
	}
	return au.ar.Count(db, filter)
}

func (au articleUseCase) CountPublicByTitle(db *gorm.DB, title string) (int, error) {
	filter := &filter.ArticleFilter{
		Status: &domain.Public,
		Title:  &title,
	}
	return au.ar.Count(db, filter)
}

func (au articleUseCase) CountPublicByTag(db *gorm.DB, tag string) (int, error) {
	filter := &filter.ArticleFilter{
		Status: &domain.Public,
		Tags:   []string{tag},
	}
	return au.ar.Count(db, filter)
}

func (au articleUseCase) FindByID(db *gorm.DB, id string) (*domain.Article, error) {
	return au.ar.FindByID(db, id, nil)
}

func (au articleUseCase) Find(db *gorm.DB, p pagenation.Pagenation) (*[]domain.Article, error) {
	return au.ar.Find(db, p.Limit, (p.Page-1)*p.Limit, nil)
}

func (au articleUseCase) Count(db *gorm.DB) (int, error) {
	return au.ar.Count(db, nil)
}

func (au articleUseCase) Post(db *gorm.DB, article *domain.Article) error {
	return au.ar.Post(db, article)
}

func (au articleUseCase) Update(db *gorm.DB, article *domain.Article) error {
	return au.ar.Update(db, article)
}

func (au articleUseCase) Delete(db *gorm.DB, id string) error {
	return au.ar.Delete(db, id)
}
