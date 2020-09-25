package database

import mgorm "github.com/champon1020/minigorm"

// Category is the struct including category information.
type Category struct {
	// unique id (primary key)
	ID string `json:"id"`

	// category name
	Name string `json:"name"`
}

// FindCategories searches for article categories.
func FindCategories(db *mgorm.DB, c *[]Category) error {
	ctx := db.Select(c, "categories")
	return ctx.Do()
}
