package dto

import (
	"github.com/champon1020/argus/domain"
)

// TagDTO is data transfer object of tag.
type TagDTO struct {
	ArticleID   string `gorm:"column:article_id"`
	Name        string `gorm:"column:name"`
	NumArticles int    `gorm:"column:n_articles"`
}

// MapToDomain maps dto to the domain model.
func (t *TagDTO) MapToDomain() *domain.Tag {
	return &domain.Tag{Name: t.Name, NumArticles: t.NumArticles}
}

// NewTagDTO creates TagDTO from domain.Tag.
func NewTagDTO(t *domain.Tag) *TagDTO {
	return &TagDTO{Name: t.Name}
}

func NewTagDTOs(tags *[]domain.Tag, articleID string) *[]TagDTO {
	tagDTOs := make([]TagDTO, len(*tags))
	for i, t := range *tags {
		tagDTOs[i] = TagDTO{Name: t.Name, ArticleID: articleID}
	}
	return &tagDTOs
}
