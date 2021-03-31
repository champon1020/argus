package domain

import "time"

// Article structure.
type Article struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Tags      []Tag     `json:"tag"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url"`
	IsPublic  bool      `json:"is_public"`
}
