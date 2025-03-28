package domain

import (
	"time"
)

type UserSession struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url"`
	Expires   time.Time `json:"expires"`
}
