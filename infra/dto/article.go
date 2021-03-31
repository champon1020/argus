package dto

import (
	"time"

	"github.com/champon1020/argus/domain"
)

// ArticleDTO is article data transfer object.
type ArticleDTO struct {
	ID        string    `gorm:"id"`
	Title     string    `gorm:"title"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
	Content   string    `gorm:"content"`
	ImageURL  string    `gorm:"image_url"`
	Status    int       `gorm:"status"`
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
