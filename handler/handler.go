package handler

import (
	"github.com/champon1020/argus/repository"
)

type Handler struct {
	Repo *repository.Repository
}

// NewHandler creates Handler instance.
func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{Repo: repo}
}
