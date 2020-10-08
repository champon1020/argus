package model

import (
	"errors"
	"time"

	"github.com/champon1020/argus/v2"
)

var (
	errDraftDbNil       = errors.New("model.draft: model.Database.DB is nil")
	errDraftTxNil       = errors.New("model.draft: model.Database.TX is nil")
	errDraftQueryFailed = errors.New("model.draft: Failed to execute query")
	errDraftNoResult    = errors.New("model.draft: Query result is nothing")
)

// Draft is the struct including draft information.
type Draft struct {
	// unique id (primary key)
	ID string `mgorm:"id" json:"id"`

	// id for sorting drafts
	SortedID int `mgorm:"sorted_id" json:"sortedId"`

	// draft title
	Title string `mgorm:"title" json:"title"`

	// categories of draft
	Categories string `mgorm:"categories" json:"categories"`

	// date draft is updated
	UpdateDate time.Time `mgorm:"update_date" json:"updateDate"`

	// content of draft
	Content string `mgorm:"content" json:"content"`

	// image file name
	ImageHash string `mgorm:"image_hash" json:"imageHash"`
}

// FindDraftByID searches for draft
// whose id is the specified id string.
func (db *Database) FindDraftByID(d *Draft, id string) error {
	if db.DB == nil {
		return argus.NewError(errDraftDbNil, nil)
	}

	var _d []Draft
	ctx := db.DB.Select(&_d, "drafts").
		Where("id = ?", id)

	if err := ctx.Do(); err != nil {
		return argus.NewError(errDraftQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	if len(_d) == 0 {
		return argus.NewError(errDraftNoResult, nil)
	}
	*d = _d[0]

	return nil
}

// FindDrafts searches for drafts.
func (db *Database) FindDrafts(d *[]Draft, op *QueryOptions) error {
	if db.DB == nil {
		return argus.NewError(errDraftDbNil, nil)
	}

	ctx := db.DB.Select(d, "drafts")
	op.apply(ctx)

	if err := ctx.Do(); err != nil {
		return argus.NewError(errDraftQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	return nil
}

// InsertDraft inserts new draft.
func (db *Database) InsertDraft(d *Draft) error {
	if db.TX == nil {
		return argus.NewError(errDraftTxNil, nil)
	}

	ctx := db.TX.InsertWithModel(d, "drafts")

	if err := ctx.Do(); err != nil {
		return argus.NewError(errDraftQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	return nil
}

// UpdateDraft updates the draft contents.
func (db *Database) UpdateDraft(d *Draft) error {
	if db.TX == nil {
		return argus.NewError(errDraftTxNil, nil)
	}

	ctx := db.TX.UpdateWithModel(d, "drafts").
		Where("id = ?", d.ID)

	if err := ctx.Do(); err != nil {
		return argus.NewError(errDraftQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	return nil
}

// DeleteDraft deletes the draft.
func (db *Database) DeleteDraft(id int) error {
	if db.TX == nil {
		return argus.NewError(errDraftTxNil, nil)
	}

	ctx := db.TX.Delete("drafts").
		Where("id = ?", id)

	if err := ctx.Do(); err != nil {
		return argus.NewError(errDraftQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	return nil
}
