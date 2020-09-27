package database

// Category is the struct including category information.
type Category struct {
	// unique id (primary key)
	ID string `json:"id"`

	// category name
	Name string `json:"name"`
}

// FindCategories searches for article categories.
func (db *Database) FindCategories(c *[]Category, op *QueryOptions) error {
	ctx := db.DB.Select(c, "categories")
	op.apply(ctx)
	return ctx.Do()
}
