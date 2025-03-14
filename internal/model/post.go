package model

import "time"

type Post struct {
	ID        int       `json:"id" example:"1"`
	Title     string    `json:"title" example:"My First Post"`
	Text      string    `json:"text" example:"This is the content of my first post."`
	CreatedAt time.Time `json:"created_at" example:"2023-10-27T10:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-10-27T10:30:00Z"`
	UserID    int       `json:"user_id" example:"123"`
}
