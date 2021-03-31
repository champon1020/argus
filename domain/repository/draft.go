package repository

import (
	"github.com/champon1020/argus/domain"
	"gorm.io/gorm"
)

// DraftRepository is repository interface for draft.
type DraftRepository interface {
	FindByID(db *gorm.DB, id string) (*domain.Draft, error)
	Find(db *gorm.DB, num int) (*[]domain.Draft, error)
	Count(db *gorm.DB) (int, error)
	Post(db *gorm.DB) error
	Update(db *gorm.DB) error
	Delete(db *gorm.DB) error
}
