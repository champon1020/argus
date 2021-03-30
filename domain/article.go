package domain

import "time"

type Article struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	//Categories []Category `json:"categories"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url"`
	IsPublic  bool      `json:"is_public"`
}
