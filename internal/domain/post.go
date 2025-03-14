package domain

import (
	"github/arshamroshannejad/task-rootext/internal/entities"
	"github/arshamroshannejad/task-rootext/internal/model"
)

type PostRepository interface {
	GetAll() (*[]model.Post, error)
	GetByID(id string) (*model.Post, error)
	GetByTitle(title string) (*model.Post, error)
	Create(post *entities.PostCreateUpdateRequest, userID string) (*model.Post, error)
	Update(post *entities.PostCreateUpdateRequest, userID string) (*model.Post, error)
	Delete(id string) error
}

type PostService interface {
	GetAllPosts() (*[]model.Post, error)
	GetPostByID(id string) (*model.Post, error)
	GetPostByTitle(title string) (*model.Post, error)
	CreatePost(post *entities.PostCreateUpdateRequest, userID string) (*model.Post, error)
	UpdatePost(post *entities.PostCreateUpdateRequest, userID string) (*model.Post, error)
	DeletePost(id string) error
}
