package usecase

import (
	"encoding/json"
	"time"

	"github.com/champon1020/argus/domain"
	"github.com/champon1020/argus/domain/repository"
	"github.com/champon1020/argus/infra/filter"
	"github.com/champon1020/argus/usecase/pagenation"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ArticleUseCase is usecase interface for article.
type ArticleUseCase interface {
	FindPublicByID(db *gorm.DB, id string) (*domain.Article, error)
	FindPublic(db *gorm.DB, p pagenation.Pagenation) (*[]domain.Article, error)
	FindPublicByTitle(db *gorm.DB, p pagenation.Pagenation, title string) (*[]domain.Article, error)
	FindPublicByTag(db *gorm.DB, p pagenation.Pagenation, tag string) (*[]domain.Article, error)
	FindByID(db *gorm.DB, id string) (*domain.Article, error)
	Find(db *gorm.DB, p pagenation.Pagenation) (*[]domain.Article, error)
	FindDraftArticles(db *gorm.DB, p pagenation.Pagenation) (*[]domain.Article, error)
	CountPublic(db *gorm.DB) (int, error)
	CountPublicByTitle(db *gorm.DB, title string) (int, error)
	CountPublicByTag(db *gorm.DB, tag string) (int, error)
	Count(db *gorm.DB) (int, error)
	CountDraftArticles(db *gorm.DB) (int, error)
	Post(db *gorm.DB, jsonBody []byte) (string, error)
	Update(db *gorm.DB, jsonBody []byte) (string, error)
	UpdateStatus(db *gorm.DB, jsonBody []byte) (string, error)
	Delete(db *gorm.DB, jsonBody []byte) error
}

type articleUseCase struct {
	aR repository.ArticleRepository
	tR repository.TagRepository
}

// NewArticleUseCase creates articleUseCase.
func NewArticleUseCase(aR repository.ArticleRepository, tR repository.TagRepository) ArticleUseCase {
	return &articleUseCase{aR: aR, tR: tR}
}

// FindPublicByID fetches the public article by id.
func (aU articleUseCase) FindPublicByID(db *gorm.DB, id string) (*domain.Article, error) {
	return aU.aR.FindByID(db, id, &domain.Public)
}

// FindPublic fetches the public articles.
func (aU articleUseCase) FindPublic(db *gorm.DB, p pagenation.Pagenation) (*[]domain.Article, error) {
	filter := &filter.ArticleFilter{
		Status: &domain.Public,
		Order:  "created_at desc",
	}
	return aU.aR.Find(db, p.Limit, (p.Page-1)*p.Limit, filter)
}

// FindPublicByTitle fetches the public articles by title.
func (aU articleUseCase) FindPublicByTitle(db *gorm.DB, p pagenation.Pagenation, title string) (*[]domain.Article, error) {
	filter := &filter.ArticleFilter{
		Status: &domain.Public,
		Title:  &title,
		Order:  "created_at desc",
	}
	return aU.aR.Find(db, p.Limit, (p.Page-1)*p.Limit, filter)
}

// FindPublicByTag fetches the public articles by tag.
func (aU articleUseCase) FindPublicByTag(db *gorm.DB, p pagenation.Pagenation, tag string) (*[]domain.Article, error) {
	filter := &filter.ArticleFilter{
		Status: &domain.Public,
		Tags:   []string{tag},
		Order:  "created_at desc",
	}
	return aU.aR.Find(db, p.Limit, (p.Page-1)*p.Limit, filter)
}

// CountPublic counts the number of public articles.
func (aU articleUseCase) CountPublic(db *gorm.DB) (int, error) {
	filter := &filter.ArticleFilter{
		Status: &domain.Public,
		Order:  "created_at desc",
	}
	return aU.aR.Count(db, filter)
}

// CountPublicByTitle counts the public number of articles by title.
func (aU articleUseCase) CountPublicByTitle(db *gorm.DB, title string) (int, error) {
	filter := &filter.ArticleFilter{
		Status: &domain.Public,
		Title:  &title,
	}
	return aU.aR.Count(db, filter)
}

// CountPublicByTag counts the number of public articles by tag.
func (aU articleUseCase) CountPublicByTag(db *gorm.DB, tag string) (int, error) {
	filter := &filter.ArticleFilter{
		Status: &domain.Public,
		Tags:   []string{tag},
	}
	return aU.aR.Count(db, filter)
}

// FindByID fetches the article by id.
func (aU articleUseCase) FindByID(db *gorm.DB, id string) (*domain.Article, error) {
	return aU.aR.FindByID(db, id, nil)
}

// Find fetches the articles.
func (aU articleUseCase) Find(db *gorm.DB, p pagenation.Pagenation) (*[]domain.Article, error) {
	filter := &filter.ArticleFilter{
		Order: "created_at desc",
	}
	return aU.aR.Find(db, p.Limit, (p.Page-1)*p.Limit, filter)
}

// FindDraftArticles fetche the draft articles.
func (aU articleUseCase) FindDraftArticles(db *gorm.DB, p pagenation.Pagenation) (*[]domain.Article, error) {
	filter := &filter.ArticleFilter{
		Status: &domain.Draft,
		Order:  "created_at desc",
	}
	return aU.aR.Find(db, p.Limit, (p.Page-1)*p.Limit, filter)
}

// Count counts the number of articles.
func (aU articleUseCase) Count(db *gorm.DB) (int, error) {
	return aU.aR.Count(db, nil)
}

// CountDraftArticles counts the number of draft articles.
func (aU articleUseCase) CountDraftArticles(db *gorm.DB) (int, error) {
	filter := &filter.ArticleFilter{
		Status: &domain.Draft,
	}
	return aU.aR.Count(db, filter)
}

// Post registers a new article.
func (aU articleUseCase) Post(db *gorm.DB, jsonBody []byte) (string, error) {
	article := &domain.Article{}
	if err := json.Unmarshal(jsonBody, article); err != nil {
		return "", err
	}

	uuid, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	article.ID = uuid.String()
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := aU.aR.Post(tx, article); err != nil {
			return err
		}

		if err := aU.tR.Posts(tx, &article.Tags, article.ID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return "", err
	}

	return article.ID, nil
}

// Update updates the article.
func (aU articleUseCase) Update(db *gorm.DB, jsonBody []byte) (string, error) {
	article := &domain.Article{}
	if err := json.Unmarshal(jsonBody, article); err != nil {
		return "", err
	}

	article.UpdatedAt = time.Now()

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := aU.aR.Update(tx, article); err != nil {
			return err
		}

		if err := aU.tR.DeleteByArticleID(tx, article.ID); err != nil {
			return err
		}

		if err := aU.tR.Posts(tx, &article.Tags, article.ID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return "", err
	}

	return article.ID, nil
}

// UpdateStatus updates the article status.
func (aU articleUseCase) UpdateStatus(db *gorm.DB, jsonBody []byte) (string, error) {
	article := &domain.Article{}
	if err := json.Unmarshal(jsonBody, article); err != nil {
		return "", err
	}

	article.UpdatedAt = time.Now()

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := aU.aR.UpdateStatus(tx, article); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return "", err
	}

	return article.ID, nil
}

// Delete removes the article.
func (aU articleUseCase) Delete(db *gorm.DB, jsonBody []byte) error {
	article := &domain.Article{}
	if err := json.Unmarshal(jsonBody, article); err != nil {
		return err
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := aU.aR.Delete(tx, article.ID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
