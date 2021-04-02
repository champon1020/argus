package domain

// Pagenation is domain model for pagenation.
type Pagenation struct {
	Next    bool `json:"next"`
	Current int  `json:"current"`
	Prev    bool `json:"prev"`
}
