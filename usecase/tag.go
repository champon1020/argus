package usecase

import (
	"github.com/champon1020/argus/domain"
	"github.com/champon1020/argus/domain/repository"
	"gorm.io/gorm"
)

// TagUseCase is usecase interface for tag.
type TagUseCase interface {
	FindPublic(db *gorm.DB) (*[]domain.Tag, error)
	Posts(db *gorm.DB, tags *[]domain.Tag, articleID string) error
	DeleteByArticleID(db *gorm.DB, articleID string) error
}

type tagUseCase struct {
	tR repository.TagRepository
}

// NewTagUseCase creates tagUseCase.
func NewTagUseCase(tR repository.TagRepository) TagUseCase {
	return &tagUseCase{tR: tR}
}

func (tU *tagUseCase) FindPublic(db *gorm.DB) (*[]domain.Tag, error) {
	return tU.tR.Find(db, &domain.Public)
}

func (tU tagUseCase) Posts(db *gorm.DB, tags *[]domain.Tag, articleID string) error {
	return tU.tR.Posts(db, tags, articleID)
}

func (tU tagUseCase) DeleteByArticleID(db *gorm.DB, articleID string) error {
	return tU.tR.DeleteByArticleID(db, articleID)
}
