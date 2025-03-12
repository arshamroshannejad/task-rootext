package repository

import (
	"context"
	"database/sql"
	"github/arshamroshannejad/task-rootext/internal/domain"
	"github/arshamroshannejad/task-rootext/internal/entities"
	"github/arshamroshannejad/task-rootext/internal/model"
	"time"
)

type userRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

func (u *userRepositoryImpl) GetByID(id string) (*model.User, error) {
	query := "SELECT * FROM users WHERE id = $1"
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := u.db.QueryRowContext(ctx, query, id)
	return collectUserRow(row)
}

func (u *userRepositoryImpl) GetByEmail(email string) (*model.User, error) {
	query := "SELECT * FROM users WHERE email = $1"
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := u.db.QueryRowContext(ctx, query, email)
	return collectUserRow(row)
}

func (u *userRepositoryImpl) Create(user *entities.UserAuthRequest) error {
	query := `INSERT INTO users (email, password) VALUES ($1, $2)`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	args := []interface{}{user.Email, user.Password}
	_, err := u.db.ExecContext(ctx, query, args...)
	return err
}

func collectUserRow(row *sql.Row) (*model.User, error) {
	var user model.User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)
	return &user, err
}
