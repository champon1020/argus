package persistence

import (
	"github.com/champon1020/argus/domain"
	"github.com/champon1020/argus/domain/filter"
	"github.com/champon1020/argus/domain/repository"
	"github.com/champon1020/argus/infra/dto"
	"gorm.io/gorm"
)

type articlePersistence struct{}

// NewArticlePersistence returns articlePersistence, which implements repository.ArticleRepository.
func NewArticlePersistence() repository.ArticleRepository {
	return &articlePersistence{}
}

func (ap *articlePersistence) FindByID(db *gorm.DB, id string, status *domain.Status) (*domain.Article, error) {
	articleDTO := &dto.ArticleDTO{}

	base := db.Table("articles").Where("id = ?", id)
	if status != nil {
		base.Where("status = ?", *status)
	}
	if err := base.First(articleDTO).Error; err != nil {
		return nil, err
	}

	article := articleDTO.MapToDomain()
	tags, err := (*tagPersistence).FindByArticleID(&tagPersistence{}, db, article.ID)
	if err != nil {
		return nil, err
	}
	article.Tags = *tags

	return article, nil
}

func (ap *articlePersistence) Find(db *gorm.DB, limit int, offset int, filter *filter.ArticleFilter) (*[]domain.Article, error) {
	articleDTOs := []dto.ArticleDTO{}

	base := db.Table("articles").Limit(limit).Offset(offset)
	if filter != nil {
		base = filter.Apply(base)
	}
	if err := base.Find(&articleDTOs).Error; err != nil {
		return nil, err
	}

	articles := make([]domain.Article, len(articleDTOs))
	for i, a := range articleDTOs {
		articles[i] = *a.MapToDomain()
		tags, err := (*tagPersistence).FindByArticleID(&tagPersistence{}, db, articles[i].ID)
		if err != nil {
			return nil, err
		}
		articles[i].Tags = *tags
	}

	return &articles, nil
}

func (ap *articlePersistence) Count(db *gorm.DB, filter *filter.ArticleFilter) (int, error) {
	var cnt int64

	base := db.Table("articles")
	if filter != nil {
		base = filter.Apply(base)
	}
	if err := base.Count(&cnt).Error; err != nil {
		return -1, err
	}

	return int(cnt), nil
}

func (ap *articlePersistence) Post(db *gorm.DB, article *domain.Article) error {
	articleDTO := dto.NewArticleDTO(article)
	if err := db.Table("articles").Create(articleDTO).Error; err != nil {
		return err
	}
	return nil
}

func (ap *articlePersistence) Update(db *gorm.DB, article *domain.Article) error {
	articleDTO := dto.NewArticleDTO(article)
	if err := db.Table("articles").
		Select("title", "updated_at", "content", "image_url", "status").
		Where("id = ?", article.ID).
		Updates(articleDTO).Error; err != nil {
		return err
	}
	return nil
}

func (ap *articlePersistence) Delete(db *gorm.DB, id string) error {
	if err := db.Exec("DELETE FROM articles WHERE id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
