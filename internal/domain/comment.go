package domain

import "time"

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	ReplyToID int       `json:"reply_to_id,omitempty"`
	Author_id int       `json:"author_id"`
	Content   string    `json:"content"`
	AvatarURL string    `json:"avatar_url"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	ImageLink string    `json:"image_link"`
}
