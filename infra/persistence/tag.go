package persistence

import (
	"github.com/champon1020/argus/domain"
	"github.com/champon1020/argus/domain/repository"
	"gorm.io/gorm"
)

type tagPersistence struct{}

// NewTagPersistence returns tagPersistence, which implements repository.TagRepository.
func NewTagPersistence() repository.TagRepository {
	return &tagPersistence{}
}

func (tp *tagPersistence) FindByName(db *gorm.DB, name string, isPublic bool) (*domain.Tag, error) {
	return nil, nil
}

func (tp *tagPersistence) FindAll(db *gorm.DB, isPublic bool) (*[]domain.Tag, error) {
	return nil, nil
}

func (tp *tagPersistence) Count(db *gorm.DB, isPublic bool) (int, error) {
	return 0, nil
}

func (tp *tagPersistence) Post(db *gorm.DB) error {
	return nil
}

func (tp *tagPersistence) Delete(db *gorm.DB) error {
	return nil
}
