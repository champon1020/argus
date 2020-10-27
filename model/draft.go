package model

import (
	"errors"
	"time"

	"github.com/champon1020/argus"
	"github.com/champon1020/minigorm"
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

	// draft title
	Title string `mgorm:"title" json:"title"`

	// categories of draft
	Categories string `mgorm:"categories" json:"categories"`

	// date draft is updated
	UpdatedDate time.Time `mgorm:"updated_date" json:"updatedDate"`

	// content of draft
	Content string `mgorm:"content" json:"content"`

	// image file name
	ImageHash string `mgorm:"image_name" json:"imageName"`
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

// RegisterDraft registers new draft.
func (db *Database) RegisterDraft(d *Draft) error {
	// Create transaction instance.
	tx, err := db.DB.NewTX()
	if err != nil {
		return argus.NewError(errFailedBeginTx, err)
	}

	err = tx.Transact(func(tx *minigorm.TX) error {
		if err := insertDraft(tx, d); err != nil {
			return err
		}

		return nil
	})

	return err
}

// InsertDraft inserts new draft.
func insertDraft(tx *minigorm.TX, d *Draft) error {
	ctx := tx.Insert("drafts").
		AddColumn("id", d.ID).
		AddColumn("title", d.Title).
		AddColumn("categories", d.Categories).
		AddColumn("updated_date", time.Now()).
		AddColumn("content", d.Content).
		AddColumn("image_name", d.ImageHash)

	if err := ctx.DoTx(); err != nil {
		return argus.NewError(errDraftQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	return nil
}

// UpdateDraft updates existed draft.
func (db *Database) UpdateDraft(d *Draft) error {
	// Create transaction instance.
	tx, err := db.DB.NewTX()
	if err != nil {
		return argus.NewError(errFailedBeginTx, err)
	}

	err = tx.Transact(func(tx *minigorm.TX) error {
		if err := updateDraft(tx, d); err != nil {
			return err
		}

		return nil
	})

	return err
}

// updateDraft updates the draft contents.
func updateDraft(tx *minigorm.TX, d *Draft) error {
	ctx := tx.Update("drafts").
		AddColumn("title", d.Title).
		AddColumn("categories", d.Categories).
		AddColumn("updated_date", time.Now()).
		AddColumn("content", d.Content).
		AddColumn("image_name", d.ImageHash).
		Where("id = ?", d.ID)

	if err := ctx.DoTx(); err != nil {
		return argus.NewError(errDraftQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	return nil
}

// DeleteDraft deletes the draft.
func (db *Database) DeleteDraft(draftID string) error {
	// Create transaction instance.
	tx, err := db.DB.NewTX()
	if err != nil {
		return argus.NewError(errFailedBeginTx, err)
	}

	err = tx.Transact(func(tx *minigorm.TX) error {
		return deleteDraft(tx, draftID)
	})

	return err
}

func deleteDraft(tx *minigorm.TX, draftID string) error {
	ctx := tx.Delete("drafts").
		Where("id = ?", draftID)

	if err := ctx.DoTx(); err != nil {
		return argus.NewError(errDraftQueryFailed, err).
			AppendValue("query", ctx.ToSQLString())
	}

	return nil
}
