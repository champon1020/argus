package persistence

import (
	"github.com/champon1020/argus/domain"
	"github.com/champon1020/argus/domain/repository"
	"github.com/champon1020/argus/infra/dto"
	"gorm.io/gorm"
)

type tagPersistence struct{}

// NewTagPersistence returns tagPersistence, which implements repository.TagRepository.
func NewTagPersistence() repository.TagRepository {
	return &tagPersistence{}
}

func (tP *tagPersistence) FindByArticleID(db *gorm.DB, articleID string) (*[]domain.Tag, error) {
	tagDTOs := []dto.TagDTO{}

	if err := db.Table("tags").Where("article_id = ?", articleID).Find(&tagDTOs).Error; err != nil {
		return nil, err
	}

	tags := make([]domain.Tag, len(tagDTOs))
	for i, t := range tagDTOs {
		tags[i] = *t.MapToDomain()
	}

	return &tags, nil
}

func (tP *tagPersistence) Find(db *gorm.DB, articleStatus *domain.Status) (*[]domain.Tag, error) {
	tagDTOs := []dto.TagDTO{}

	base := db.Table("tags").Select("name", "COUNT(name) AS n_articles").Group("name")
	if articleStatus != nil {
		base = base.Where("article_id IN (?)",
			db.Table("articles").
				Where("status = ?", articleStatus).
				Select("id"),
		)
	}
	if err := base.Find(&tagDTOs).Error; err != nil {
		return nil, err
	}

	tags := make([]domain.Tag, len(tagDTOs))
	for i, t := range tagDTOs {
		tags[i] = *t.MapToDomain()
	}

	return &tags, nil
}

func (tP *tagPersistence) Posts(db *gorm.DB, tags *[]domain.Tag, articleID string) error {
	if len(*tags) == 0 {
		return nil
	}

	tagDTOs := dto.NewTagDTOs(tags, articleID)
	if err := db.Table("tags").Select("name", "article_id").Create(tagDTOs).Error; err != nil {
		return err
	}

	return nil
}

func (tP *tagPersistence) DeleteByArticleID(db *gorm.DB, articleID string) error {
	if err := db.Exec("DELETE FROM tags WHERE article_id = ?", articleID).Error; err != nil {
		return err
	}
	return nil
}
