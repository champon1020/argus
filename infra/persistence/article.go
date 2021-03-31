package persistence

import (
	"github.com/champon1020/argus/domain"
	"github.com/champon1020/argus/domain/filter"
	"github.com/champon1020/argus/domain/repository"
	"gorm.io/gorm"
)

type articlePersistence struct{}

// NewArticlePersistence returns articlePersistence, which implements repository.ArticleRepository.
func NewArticlePersistence() repository.ArticleRepository {
	return &articlePersistence{}
}

func (ap *articlePersistence) FindByID(db *gorm.DB, id string, isPublic bool) (*domain.Article, error) {
	article := &domain.Article{}

	if err := db.Table("articles").Where("id = ?", id).First(article).Error; err != nil {
		return nil, err
	}

	return article, nil
}

func (ap *articlePersistence) Find(db *gorm.DB, limit int, offset int, filter *filter.ArticleFilter) (*[]domain.Article, error) {
	articles := []domain.Article{}

	base := db.Table("articles").Limit(limit).Offset(offset)
	if filter != nil {
		base = filter.Apply(base)
	}
	if err := base.Find(&articles).Error; err != nil {
		return nil, err
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
	if err := db.Create(article).Error; err != nil {
		return err
	}
	return nil
}

func (ap *articlePersistence) Update(db *gorm.DB, article *domain.Article) error {
	if err := db.Model(article).Where("id = ?", article.ID).Updates(article).Error; err != nil {
		return err
	}
	return nil
}

func (ap *articlePersistence) Delete(db *gorm.DB, id string) error {
	if err := db.Where("id = ?", id).Delete(&domain.Article{}).Error; err != nil {
		return err
	}
	return nil
}
