package domain

import "time"

type User struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ImageURL   string `json:"image"`
	Token_ID   string
	TokenDate  time.Time
	Expires_at time.Time
}
