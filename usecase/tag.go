package usecase

import (
	"github.com/champon1020/argus/domain"
	"github.com/champon1020/argus/domain/repository"
	"gorm.io/gorm"
)

// TagUseCase is usecase interface for tag.
type TagUseCase interface {
	FindByName(db *gorm.DB, name string) (*[]domain.Tag, error)
}

type tagUseCase struct {
	tr repository.TagRepository
}

// NewTagUseCase creates tagUseCase.
func NewTagUseCase(tr repository.TagRepository) TagUseCase {
	return &tagUseCase{tr: tr}
}

func (tu *tagUseCase) FindByName(db *gorm.DB, name string) (*[]domain.Tag, error) {
	return nil, nil
}

func (tu tagUseCase) Post(db *gorm.DB, tag *domain.Tag) error {
	return nil
}

func (tu tagUseCase) Delete(db *gorm.DB, id string) error {
	return nil
}
