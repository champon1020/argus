package domain

import "time"

// Draft structure.
type Draft struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Tags      string    `json:"tags"`
	UpdatedAt time.Time `json:"updated_at"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url"`
}
