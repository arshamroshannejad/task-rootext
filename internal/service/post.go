package service

import (
	"github/arshamroshannejad/task-rootext/internal/domain"
	"github/arshamroshannejad/task-rootext/internal/entities"
	"github/arshamroshannejad/task-rootext/internal/model"
	"go.uber.org/zap"
)

type postServiceImpl struct {
	postRepository domain.PostRepository
	zapLogger      *zap.Logger
}

func NewPostService(postRepository domain.PostRepository, zapLogger *zap.Logger) domain.PostService {
	return &postServiceImpl{
		postRepository: postRepository,
		zapLogger:      zapLogger,
	}
}

func (p *postServiceImpl) GetAllPosts() (*[]model.Post, error) {
	posts, err := p.postRepository.GetAll()
	if err != nil {
		p.zapLogger.Error("Failed to get all posts", zap.Error(err))
		return nil, err
	}
	return posts, err
}

func (p *postServiceImpl) GetPostByID(id string) (*model.Post, error) {
	post, err := p.postRepository.GetByID(id)
	if err != nil {
		p.zapLogger.Error("Failed to get post with id", zap.Error(err))
		return nil, err
	}
	return post, err
}
func (p *postServiceImpl) GetPostByTitle(title string) (*model.Post, error) {
	post, err := p.postRepository.GetByTitle(title)
	if err != nil {
		p.zapLogger.Error("Failed to get post with title", zap.Error(err))
		return nil, err
	}
	return post, err
}

func (p *postServiceImpl) CreatePost(post *entities.PostCreateUpdateRequest, userID string) (*model.Post, error) {
	createdPost, err := p.postRepository.Create(post, userID)
	if err != nil {
		p.zapLogger.Error("Failed to create post", zap.Error(err))
		return nil, err
	}
	return createdPost, nil
}

func (p *postServiceImpl) UpdatePost(post *entities.PostCreateUpdateRequest, userID string) (*model.Post, error) {
	updatedPost, err := p.postRepository.Update(post, userID)
	if err != nil {
		p.zapLogger.Error("Failed to update post", zap.Error(err))
		return nil, err
	}
	return updatedPost, nil
}

func (p *postServiceImpl) DeletePost(id string) error {
	err := p.postRepository.Delete(id)
	if err != nil {
		p.zapLogger.Error("Failed to delete post")
		return err
	}
	return nil
}
