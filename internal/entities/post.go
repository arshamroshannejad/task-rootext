package entities

type PostCreateUpdateRequest struct {
	Title string `json:"title" example:"My New Post" validate:"required"`
	Text  string `json:"text" example:"Content of my new post." validate:"required"`
}
