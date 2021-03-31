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

func (tp *tagPersistence) FindByArticleID(db *gorm.DB, articleID string) (*[]domain.Tag, error) {
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

func (tp *tagPersistence) FindByName(db *gorm.DB, name string, articleStatus *domain.Status) (*[]domain.Tag, error) {
	tagDTOs := []dto.TagDTO{}

	base := db.Table("tags").Where("name LIKE ?", "%%"+name+"%%")
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

	tags := []domain.Tag{}
	for i, t := range tagDTOs {
		tags[i] = *t.MapToDomain()
	}

	return &tags, nil
}

func (tp *tagPersistence) Find(db *gorm.DB, articleStatus *domain.Status) (*[]domain.Tag, error) {
	tagDTOs := []dto.TagDTO{}

	base := db.Table("tags")
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

	tags := []domain.Tag{}
	for i, t := range tagDTOs {
		tags[i] = *t.MapToDomain()
	}

	return &tags, nil
}

func (tp *tagPersistence) Post(db *gorm.DB, tag *domain.Tag) error {
	tagDTO := dto.NewTagDTO(tag)
	if err := db.Table("tags").Create(tagDTO).Error; err != nil {
		return err
	}
	return nil
}

func (tp *tagPersistence) Delete(db *gorm.DB, articleID string, name string) error {
	if err := db.Exec("DELETE FROM tags WHERE article_id = ? AND name = ?", articleID, name).Error; err != nil {
		return err
	}
	return nil
}
