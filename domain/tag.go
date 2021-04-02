package domain

// Tag is domain model of tag.
type Tag struct {
	Name        string `json:"name"`
	NumArticles int    `json:"n_articles"`
}
