package pagenation

import "github.com/champon1020/argus/domain"

// Pagenation is the pagenation utility model.
type Pagenation struct {
	Page  int
	Total int
	Limit int
}

// NewPagenation creates Pagenation instance.
func NewPagenation(page, total, limit int) *Pagenation {
	return &Pagenation{Page: page, Total: total, Limit: limit}
}

// HasNext returns whether or not there is next page.
func (p Pagenation) HasNext() bool {
	return p.Page*p.Limit < p.Total
}

// HasPrev returns whether or not there is previous page.
func (p Pagenation) HasPrev() bool {
	return p.Page > 1
}

// MapToDomain maps to domain model.
func (p Pagenation) MapToDomain() *domain.Pagenation {
	return &domain.Pagenation{
		Next:    p.HasNext(),
		Current: p.Page,
		Prev:    p.HasPrev(),
	}
}
