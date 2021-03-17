package repository

import "github.com/champon1020/mgorm"

type Repository struct {
	DB *mgorm.DB
}

// NewRepository creates Repository instance.
func NewRepository(db *mgorm.DB) *Repository {
	return &Repository{DB: db}
}
