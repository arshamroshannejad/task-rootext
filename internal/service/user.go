package service

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github/arshamroshannejad/task-rootext/config"
	"github/arshamroshannejad/task-rootext/internal/domain"
	"github/arshamroshannejad/task-rootext/internal/entities"
	"github/arshamroshannejad/task-rootext/internal/model"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type userServiceImpl struct {
	userRepository domain.UserRepository
	redisDB        *redis.Client
	zapLogger      *zap.Logger
	cfg            *config.Config
}

func NewUserService(userRepository domain.UserRepository, redisDB *redis.Client, zapLogger *zap.Logger, cfg *config.Config) domain.UserService {
	return &userServiceImpl{
		userRepository: userRepository,
		redisDB:        redisDB,
		zapLogger:      zapLogger,
		cfg:            cfg,
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

func (u *userServiceImpl) VerifyPassword(hashPass, plainPass string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(plainPass)); err != nil {
		u.zapLogger.Error("Failed to verify user password", zap.Error(err))
		return err
	}
	return nil
}

func (u *userServiceImpl) CreateAccessToken(userID, email string) (string, error) {
	exp := time.Now().Add(u.cfg.App.AccessHourTTL).Unix()
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     exp,
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(u.cfg.App.Secret))
	if err != nil {
		u.zapLogger.Error("Failed to create user access token", zap.Error(err))
	}
	return token, nil
}

func (u *userServiceImpl) BlockJwtToken(token string, exp float64) error {
	remaining := time.Until(time.Unix(int64(exp), 0))
	result := u.redisDB.Set(context.Background(), token, "blocked", remaining)
	if result.Err() != nil {
		u.zapLogger.Error("Failed to add token in blacklist", zap.Error(result.Err()))
		return result.Err()
	}
	return nil
}
