package persistence

import (
	"github.com/champon1020/argus/domain"
	"github.com/champon1020/argus/domain/repository"
	"gorm.io/gorm"
)

type draftPersistence struct{}

// NewDraftPersistence returns draftPersistence, which implements repository.DraftRepository.
func NewDraftPersistence() repository.DraftRepository {
	return &draftPersistence{}
}

func (dp *draftPersistence) FindByID(db *gorm.DB, id string) (*domain.Draft, error) {
	return nil, nil
}

func (dp *draftPersistence) Find(db *gorm.DB, num int) (*[]domain.Draft, error) {
	return nil, nil
}

func (dp *draftPersistence) Count(db *gorm.DB) (int, error) {
	return 0, nil
}

func (dp *draftPersistence) Post(db *gorm.DB) error {
	return nil
}

func (dp *draftPersistence) Update(db *gorm.DB) error {
	return nil
}

func (dp *draftPersistence) Delete(db *gorm.DB) error {
	return nil
}
