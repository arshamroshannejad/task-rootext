package domain

import (
	"github.com/arshamroshannejad/task-rootext/internal/entities"
	"github.com/arshamroshannejad/task-rootext/internal/helpers"
	"github.com/arshamroshannejad/task-rootext/internal/model"
)

type PostRepository interface {
	GetAll(filter *helpers.PaginateFilter) (*[]model.Post, helpers.Metadata, error)
	GetByID(postID string) (*model.Post, error)
	GetByTitle(title string) (*model.Post, error)
	Create(post *entities.PostCreateUpdateRequest, userID string) (*model.Post, error)
	Update(post *entities.PostCreateUpdateRequest, postID string) (*model.Post, error)
	Delete(postID string) error
	AddVote(postID, userID, vote string) error
	RemoveVote(postID, userID string) error
}

type PostService interface {
	GetAllPosts(filter *helpers.PaginateFilter) (*[]model.Post, helpers.Metadata, error)
	GetPostByID(postID string) (*model.Post, error)
	GetPostByTitle(title string) (*model.Post, error)
	CreatePost(post *entities.PostCreateUpdateRequest, userID string) (*model.Post, error)
	UpdatePost(post *entities.PostCreateUpdateRequest, postID string) (*model.Post, error)
	DeletePost(postID string) error
	AddPostVote(postID, userID, vote string) error
	RemovePostVote(postID, userID string) error
}
