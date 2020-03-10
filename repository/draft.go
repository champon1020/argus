package repository

import (
	"database/sql"
	"time"
)

// Id: primary key
// Title: article title
// Categories: categories of article
// UpdateDate: last updated date
// ContentHash: content file name (html file)
// ImageHash: image file name
type Draft struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Categories  string    `json:"categories"`
	UpdateDate  time.Time `json:"updateDate"`
	ContentHash string    `json:"contentHash"`
	ImageHash   string    `json:"imageHash"`
}

func (draft *Draft) InsertDraft(tx *sql.Tx) (err error) {
	cmd := "INSERT INTO drafts " +
		"(id, title, categories, update_date, content_hash, image_hash)" +
		"VALUES (?, ?, ?, ?, ?, ?)"

	if _, err := tx.Exec(cmd,
		draft.Id,
		draft.Title,
		draft.Categories,
		draft.UpdateDate,
		draft.ContentHash,
		draft.ImageHash,
	); err != nil {
		CmdError.SetErr(err).AppendTo(Errors)
	}
	return
}

func (draft *Draft) UpdateDraft(tx *sql.Tx) (err error) {
	cmd := "UPDATE drafts " +
		"SET title=?, categories=?, update_date=?, content_hash=?, image_hash=? " +
		"WHERE id=?"

	if _, err := tx.Exec(cmd,
		draft.Title,
		draft.Categories,
		draft.UpdateDate,
		draft.ContentHash,
		draft.ImageHash,
		draft.Id,
	); err != nil {
		CmdError.SetErr(err).AppendTo(Errors)
	}
	return
}

func (draft *Draft) DeleteDraft(tx *sql.Tx) (err error) {
	cmd := "DELETE FROM drafts WHERE id=?"
	if _, err := tx.Exec(cmd, draft.Id); err != nil {
		CmdError.SetErr(err).AppendTo(Errors)
	}
	return
}

func (draft *Draft) FindDrafts(db *sql.DB, argsFlg uint32) (drafts []Draft, err error) {
	args := GenArgsSlice(argsFlg, draft)
	whereQuery, limitQuery := GenArgsQuery(argsFlg, draft)
	query := "SELECT * FROM drafts " + whereQuery + "ORDER BY id DESC " + limitQuery

	var rows *sql.Rows
	defer RowsClose(rows)
	if rows, err = db.Query(query, args...); err != nil || rows == nil {
		QueryError.
			SetErr(err).
			SetValues("query", query).
			SetValues("args", args).
			AppendTo(Errors)
		return
	}

	var d Draft
	for rows.Next() {
		if err := rows.Scan(
			&d.Id,
			&d.Title,
			&d.Categories,
			&d.UpdateDate,
			&d.ContentHash,
			&d.ImageHash,
		); err != nil {
			ScanError.SetErr(err).AppendTo(Errors)
			break
		}
		drafts = append(drafts, d)
	}
	return
}