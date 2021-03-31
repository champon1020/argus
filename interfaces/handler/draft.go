package handler

import (
	"github.com/champon1020/argus/config"
	"github.com/champon1020/argus/usecase"
	"github.com/labstack/echo/v4"
)

// DraftHandler is handler interface for draft.
type DraftHandler interface {
	DraftByID(c echo.Context) error
	Drafts(c echo.Context) error
	CountDrafts(c echo.Context) error
	PostDraft(c echo.Context) error
	UpdateDraft(c echo.Context) error
	DeleteDraft(c echo.Context) error
}

type draftHandler struct {
	config *config.Config
	du     usecase.DraftUseCase
}

// NewDraftHandler creates draftHandler.
func NewDraftHandler(du usecase.DraftUseCase, config *config.Config) DraftHandler {
	return &draftHandler{config: config, du: du}
}

func (h *draftHandler) DraftByID(c echo.Context) error {
	return nil
}

func (h *draftHandler) Drafts(c echo.Context) error {
	return nil
}

func (h *draftHandler) CountDrafts(c echo.Context) error {
	return nil
}

func (h *draftHandler) PostDraft(c echo.Context) error {
	return nil
}

func (h *draftHandler) UpdateDraft(c echo.Context) error {
	return nil
}

func (h *draftHandler) DeleteDraft(c echo.Context) error {
	return nil
}
