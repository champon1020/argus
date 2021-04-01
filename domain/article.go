package domain

import "time"

// Article domain model of article.
type Article struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Tags      []Tag     `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url"`
	Status    int       `json:"status"`
}

// Status is status of article.
type Status int

// Status types of article.
var (
	Private Status = 0
	Public  Status = 1
	Draft   Status = 2
)
