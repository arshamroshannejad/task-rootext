package service

import (
	"github/arshamroshannejad/task-rootext/internal/domain"
	"github/arshamroshannejad/task-rootext/internal/entities"
	"github/arshamroshannejad/task-rootext/internal/model"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type userServiceImpl struct {
	userRepository domain.UserRepository
	zapLogger      *zap.Logger
}

func NewUserService(userRepository domain.UserRepository, zapLogger *zap.Logger) domain.UserService {
	return &userServiceImpl{
		userRepository: userRepository,
		zapLogger:      zapLogger,
	}
}

func (u *userServiceImpl) GetUserByID(id string) (*model.User, error) {
	user, err := u.userRepository.GetByID(id)
	if err != nil {
		u.zapLogger.Error("Failed to get user with id", zap.Error(err))
		return nil, err
	}
	return user, nil
}

func (u *userServiceImpl) GetUserByEmail(email string) (*model.User, error) {
	user, err := u.userRepository.GetByEmail(email)
	if err != nil {
		u.zapLogger.Error("Failed to get user with email", zap.Error(err))
		return nil, err
	}
	return user, nil
}

func (u *userServiceImpl) CreateUser(user *entities.UserAuthRequest) error {
	if err := u.userRepository.Create(user); err != nil {
		u.zapLogger.Error("Failed to create user", zap.Error(err))
		return err
	}
	return nil
}

func (u *userServiceImpl) EncryptPassword(plainPass string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(plainPass), bcrypt.DefaultCost)
	if err != nil {
		u.zapLogger.Error("Failed to create hash password", zap.Error(err))
		return "", err
	}
	return string(hashedPass), nil
}
