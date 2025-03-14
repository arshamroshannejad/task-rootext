package handler

import (
	"database/sql"
	"errors"
	"github.com/arshamroshannejad/task-rootext/internal/domain"
	"github.com/arshamroshannejad/task-rootext/internal/entities"
	"github.com/arshamroshannejad/task-rootext/internal/helpers"
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

// RegisterHandler godoc
//
//	@Summary		Register
//	@Description	Register a new user
//	@Accept			json
//	@Produce		json
//	@Tags			Auth
//	@Param			registerRequest	body		entities.UserAuthRequest	true	"make sure send a valid email and password must be grater than 8 character"
//	@Success		201				{object}	helpers.UserCreated
//	@Failure		400				{object}	helpers.BadRequest
//	@Failure		409				{object}	helpers.UserExists
//	@Failure		500				{object}	helpers.InternalServerError
//	@Router			/auth/register [post]
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

// LoginHandler godoc
//
//	@Summary		Login
//	@Description	Login with register credential
//	@Accept			json
//	@Produce		json
//	@Tags			Auth
//	@Param			loginRequest	body		entities.UserAuthRequest	true	"make sure send a valid email and password must be grater than 8 character"
//	@Success		200				{object}	helpers.LoginOk
//	@Failure		400				{object}	helpers.BadRequest
//	@Failure		404				{object}	helpers.UserNotFound
//	@Failure		500				{object}	helpers.InternalServerError
//	@Router			/auth/login [post]
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

// LogoutHandler godoc
//
//	@Summary		Logout
//	@Description	Logout and expire jwt token. authenticate required!
//	@Accept			json
//	@Produce		json
//	@Tags			Auth
//	@Security		BearerAuth
//	@Success		200	{object}	helpers.LogoutOk
//	@Failure		500	{object}	helpers.InternalServerError
//	@Router			/auth/logout [post]
func (u *UserHandlerImpl) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	exp := r.Context().Value("exp").(float64)
	if err := u.UserService.BlockJwtToken(tokenString, exp); err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.M{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	helpers.WriteJson(w, http.StatusOK, helpers.M{"response": "logged out"})
}
