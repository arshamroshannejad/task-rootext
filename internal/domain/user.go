package domain

import (
	"github/arshamroshannejad/task-rootext/internal/entities"
	"github/arshamroshannejad/task-rootext/internal/model"
)

type UserRepository interface {
	GetByID(id string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Create(user *entities.UserAuthRequest) error
}

type UserService interface {
	GetUserByID(id string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	CreateUser(user *entities.UserAuthRequest) error
	EncryptPassword(plainPass string) (string, error)
}
