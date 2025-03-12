package model

import "time"

type User struct {
	ID        string    `json:"id" example:"23"`
	Email     string    `json:"email" example:"james@gmail.com"`
	Password  string    `json:"password" example:"1qaz2wsx"`
	CreatedAt time.Time `json:"created_at" example:"2023-10-05T14:30:45Z"`
}
