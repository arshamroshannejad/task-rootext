package handler

import (
	"database/sql"
	"errors"
	"github/arshamroshannejad/task-rootext/internal/domain"
	"github/arshamroshannejad/task-rootext/internal/entities"
	"github/arshamroshannejad/task-rootext/internal/helpers"
	"net/http"
)

type UserHandlerImpl struct {
	UserService domain.UserService
}

func NewUserHandler(userService domain.UserService) *UserHandlerImpl {
	return &UserHandlerImpl{
		UserService: userService,
	}
}

func (u *UserHandlerImpl) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	reqBody := new(entities.UserAuthRequest)
	if err := helpers.ReadJson(r, reqBody); err != nil {
		helpers.WriteJson(w, http.StatusBadRequest, helpers.M{"error": err.Error()})
		return
	}
	if _, err := u.UserService.GetUserByEmail(reqBody.Email); !errors.Is(err, sql.ErrNoRows) {
		helpers.WriteJson(w, http.StatusConflict, helpers.M{"error": "user already exists"})
		return
	}
	hashPassword, err := u.UserService.EncryptPassword(reqBody.Password)
	if err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.M{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	reqBody.Password = hashPassword
	if err := u.UserService.CreateUser(reqBody); err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.M{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	helpers.WriteJson(w, http.StatusCreated, helpers.M{"response": "user created"})
}

func (u *UserHandlerImpl) LoginHandler(w http.ResponseWriter, r *http.Request) {
	reqBody := new(entities.UserAuthRequest)
	if err := helpers.ReadJson(r, reqBody); err != nil {
		helpers.WriteJson(w, http.StatusBadRequest, helpers.M{"error": err.Error()})
		return
	}
	user, err := u.UserService.GetUserByEmail(reqBody.Email)
	if errors.Is(err, sql.ErrNoRows) {
		helpers.WriteJson(w, http.StatusNotFound, helpers.M{"error": "user not found"})
		return
	}
	if err := u.UserService.VerifyPassword(user.Password, reqBody.Password); err != nil {
		helpers.WriteJson(w, http.StatusBadRequest, helpers.M{"error": "wrong password"})
		return
	}
	accessToken, err := u.UserService.CreateAccessToken(user.ID, user.Email)
	if err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.M{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	helpers.WriteJson(w, http.StatusOK, helpers.M{"access_token": accessToken})
}

func (u *UserHandlerImpl) LogoutHandler(w http.ResponseWriter, r *http.Request) {

}
