package filter

import "gorm.io/gorm"

// ArticleFilter contains search filters.
type ArticleFilter struct {
	Title    string
	Tags     []string
	IsPublic bool
}

// Apply applies filter.
func (af *ArticleFilter) Apply(db *gorm.DB) *gorm.DB {
	if af.Title != "" {
		db.Where("title LIKE %s", "%%"+af.Title+"%%")
	}
	if len(af.Tags) > 0 {
		db.Where("tag IN (?)", db.Table("article_tag").
			Where("tag_name IN (?)", af.Tags).
			Select("article_id"),
		)
	}
	if af.IsPublic {
		db.Where("is_public = ?", 1)
	}
	return db
}
