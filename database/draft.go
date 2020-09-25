package database

import (
	"time"

	mgorm "github.com/champon1020/minigorm"
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
func FindDrafts(db *mgorm.DB, d *[]Draft) error {
	ctx := db.Select(d, "drafts")
	return ctx.Do()
}

// InsertDraft inserts new draft.
func InsertDraft(db *mgorm.DB, d *Draft) error {
	ctx := db.InsertWithModel(d, "drafts")
	return ctx.Do()
}

// UpdateDraft updates the draft contents.
func UpdateDraft(db *mgorm.DB, d *Draft) error {
	ctx := db.UpdateWithModel(d, "drafts").
		Where("id = ?", d.ID)
	return ctx.Do()
}

// DeleteDraft deletes the draft.
func DeleteDraft(db *mgorm.DB, id int) error {
	ctx := db.Delete("drafts").
		Where("id = ?", id)

	return ctx.Do()
}
