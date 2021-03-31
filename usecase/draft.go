package usecase

import "github.com/champon1020/argus/domain/repository"

// DraftUseCase is usecase interface for draft.
type DraftUseCase interface{}

type draftUseCase struct {
	dr repository.DraftRepository
}

// NewDraftUseCase creates draftUseCase.
func NewDraftUseCase(dr repository.DraftRepository) DraftUseCase {
	return &draftUseCase{dr: dr}
}
