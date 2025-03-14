package helpers

import "github.com/arshamroshannejad/task-rootext/internal/model"

type UserCreated struct {
	Response string `json:"response" example:"user created"`
}

type BadRequest struct {
	Error string `json:"error" example:"bad request"`
}

type UserExists struct {
	Error string `json:"error" example:"user already exists"`
}

type InternalServerError struct {
	Error string `json:"error" example:"Internal Server Error"`
}

type UserNotFound struct {
	Error string `json:"error" example:"user not found"`
}

type LoginOk struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type LogoutOk struct {
	Response string `json:"response" example:"logged out"`
}

type PostQueryParams struct {
	Page     *int    `json:"page"        example:"1" default:"1"`
	PageSize *int    `json:"page_size"   example:"3" default:"5"`
	Sort     *string `json:"sort"        example:"created_at -created_at vote_count -vote_count" default:"-vote_count"`
}

type AllPosts struct {
	Posts    []model.Post `json:"posts"`
	Metadata Metadata     `json:"metadata"`
}

type PostNotFound struct {
	Error string `json:"error" example:"post not found"`
}

type Post model.Post

type Forbidden struct {
	Error string `json:"error" example:"Forbidden"`
}

type VoteSuccessful struct {
	Response string `json:"response" example:"successful"`
}
