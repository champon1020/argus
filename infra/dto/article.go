package dto

import (
	"time"

	"github.com/champon1020/argus/domain"
)

// ArticleDTO is article data transfer object.
type ArticleDTO struct {
	ID        string    `gorm:"column:id"`
	Title     string    `gorm:"column:title"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	Content   string    `gorm:"column:content"`
	ImageURL  string    `gorm:"column:image_url"`
	Status    int       `gorm:"column:status"`
}

// MapToDomain maps dto to the domain model.
func (a *ArticleDTO) MapToDomain() *domain.Article {
	return &domain.Article{
		ID:        a.ID,
		Title:     a.Title,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		Content:   a.Content,
		ImageURL:  a.ImageURL,
		Status:    a.Status,
	}
}

// NewArticleDTO creates ArticleDTO from domain.Article.
func NewArticleDTO(a *domain.Article) *ArticleDTO {
	return &ArticleDTO{
		ID:        a.ID,
		Title:     a.Title,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		Content:   a.Content,
		ImageURL:  a.ImageURL,
		Status:    a.Status,
	}
}
