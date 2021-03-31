package dto

import "github.com/champon1020/argus/domain"

// TagDTO is data transfer object of tag.
type TagDTO struct {
	ArticleID string `gorm:"article_id"`
	Name      string `gorm:"name"`
}

// MapToDomain maps dto to the domain model.
func (t *TagDTO) MapToDomain() *domain.Tag {
	return &domain.Tag{Name: t.Name}
}

// NewTagDTO creates TagDTO from domain.Tag.
func NewTagDTO(t *domain.Tag) *TagDTO {
	return &TagDTO{Name: t.Name}
}
