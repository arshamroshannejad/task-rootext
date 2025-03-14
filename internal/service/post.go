package service

import (
	"github/arshamroshannejad/task-rootext/internal/domain"
	"github/arshamroshannejad/task-rootext/internal/entities"
	"github/arshamroshannejad/task-rootext/internal/helpers"
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

func (p *postServiceImpl) GetAllPosts(filter *helpers.PaginateFilter) (*[]model.Post, helpers.Metadata, error) {
	posts, metaData, err := p.postRepository.GetAll(filter)
	if err != nil {
		p.zapLogger.Error("Failed to get all posts", zap.Error(err))
		return nil, helpers.Metadata{}, err
	}
	return posts, metaData, err
}

func (p *postServiceImpl) GetPostByID(postID string) (*model.Post, error) {
	post, err := p.postRepository.GetByID(postID)
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

func (p *postServiceImpl) UpdatePost(post *entities.PostCreateUpdateRequest, postID string) (*model.Post, error) {
	updatedPost, err := p.postRepository.Update(post, postID)
	if err != nil {
		p.zapLogger.Error("Failed to update post", zap.Error(err))
		return nil, err
	}
	return updatedPost, nil
}

func (p *postServiceImpl) DeletePost(postID string) error {
	err := p.postRepository.Delete(postID)
	if err != nil {
		p.zapLogger.Error("Failed to delete post")
		return err
	}
	return nil
}

func (p *postServiceImpl) AddPostVote(postID, userID, vote string) error {
	if err := p.postRepository.AddVote(postID, userID, vote); err != nil {
		p.zapLogger.Error("Failed to add vote on post", zap.Error(err))
		return err
	}
	return nil
}

func (p *postServiceImpl) RemovePostVote(postID, userID string) error {
	if err := p.postRepository.RemoveVote(postID, userID); err != nil {
		p.zapLogger.Error("Failed to remove vote on post", zap.Error(err))
		return err
	}
	return nil
}
