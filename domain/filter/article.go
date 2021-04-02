package filter

import (
	"github.com/champon1020/argus/domain"
	"gorm.io/gorm"
)

// ArticleFilter contains search filters.
type ArticleFilter struct {
	Title  *string
	Tags   []string
	Status *domain.Status
}

// Apply applies filter.
func (af *ArticleFilter) Apply(base *gorm.DB) *gorm.DB {
	if af.Title != nil {
		base = base.Where("title LIKE ?", "%"+*af.Title+"%")
	}
	if len(af.Tags) > 0 {
		base = base.Where("id IN (SELECT article_id FROM tags WHERE name IN (?))", af.Tags)
	}
	if af.Status != nil {
		base = base.Where("status = ?", *af.Status)
	}
	return base
}
