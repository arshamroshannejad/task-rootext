package entities

type UserAuthRequest struct {
	Email    string `json:"email" example:"james@gmail.com" validate:"required,email"`
	Password string `json:"password" example:"1qaz2wsx" validate:"required,min=8"`
}
