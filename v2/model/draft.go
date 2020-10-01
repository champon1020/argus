package model

import (
	"errors"
	"time"

	"github.com/champon1020/argus/v2"
)

var (
	errDraftDbNil = errors.New("model.draft: model.Database.DB is nil")
	errDraftTxNil = errors.New("model.draft: model.Database.TX is nil")
)

// Draft is the struct including draft information.
type Draft struct {
	// unique id (primary key)
	ID string `json:"id"`

	// id for sorting drafts
	SortedID int `json:"sortedId"`

	// draft title
	Title string `json:"title"`

	// categories of draft
	Categories string `json:"categories"`

	// date draft is updated
	UpdateDate time.Time `json:"updateDate"`

	// content of draft
	Content string `json:"content"`

	// image file name
	ImageHash string `json:"imageHash"`
}

// FindDrafts searches for drafts.
func (db *Database) FindDrafts(d *[]Draft, op *QueryOptions) error {
	if db.DB == nil {
		return argus.NewError(errDraftDbNil, nil)
	}

	ctx := db.DB.Select(d, "drafts")
	op.apply(ctx)
	return ctx.Do()
}

// FindDraftByID searches for draft
// whose id is the specified id string.
func (db *Database) FindDraftByID(d *[]Draft, id string) error {
	if db.DB == nil {
		return argus.NewError(errDraftDbNil, nil)
	}

	ctx := db.DB.Select(d, "drafts").
		Where("id = ?", id)

	return ctx.Do()
}

// InsertDraft inserts new draft.
func (db *Database) InsertDraft(d *Draft) error {
	if db.TX == nil {
		return argus.NewError(errDraftTxNil, nil)
	}

	ctx := db.TX.InsertWithModel(d, "drafts")
	return ctx.Do()
}

// UpdateDraft updates the draft contents.
func (db *Database) UpdateDraft(d *Draft) error {
	if db.TX == nil {
		return argus.NewError(errDraftTxNil, nil)
	}

	ctx := db.TX.UpdateWithModel(d, "drafts").
		Where("id = ?", d.ID)
	return ctx.Do()
}

// DeleteDraft deletes the draft.
func (db *Database) DeleteDraft(id int) error {
	if db.TX == nil {
		return argus.NewError(errDraftTxNil, nil)
	}

	ctx := db.TX.Delete("drafts").
		Where("id = ?", id)

	return ctx.Do()
}
