package domain

import "time"

type Post struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	AuthorID     int       `json:"author_id"`
	Changed_name string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	ImageLink    string    `json:"image_url"`
	File         []byte    `json:"file_content"`
	Archived     bool      `json:"archived"`
}
