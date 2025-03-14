package service

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github/arshamroshannejad/task-rootext/internal/domain"
	"github/arshamroshannejad/task-rootext/internal/entities"
	"github/arshamroshannejad/task-rootext/internal/helpers"
	"github/arshamroshannejad/task-rootext/internal/model"
	"go.uber.org/zap"
	"time"
)

type postServiceImpl struct {
	postRepository domain.PostRepository
	redisDB        *redis.Client
	zapLogger      *zap.Logger
}

func NewPostService(postRepository domain.PostRepository, redisDB *redis.Client, zapLogger *zap.Logger) domain.PostService {
	return &postServiceImpl{
		postRepository: postRepository,
		redisDB:        redisDB,
		zapLogger:      zapLogger,
	}
}

func (p *postServiceImpl) GetAllPosts(filter *helpers.PaginateFilter) (*[]model.Post, helpers.Metadata, error) {
	if filter.Page == 1 && filter.Sort == "-vote_count" {
		cachedData, err := p.redisDB.Get(context.Background(), "top_5_posts").Result()
		if err == nil {
			var cachedResponse struct {
				Posts    []model.Post     `json:"posts"`
				Metadata helpers.Metadata `json:"metadata"`
			}
			if err := json.Unmarshal([]byte(cachedData), &cachedResponse); err == nil {
				p.zapLogger.Info("Serving top 5 posts from Redis cache")
				return &cachedResponse.Posts, cachedResponse.Metadata, nil
			}
		}
	}
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
	go p.refreshTopVotedCache()
	return nil
}

func (p *postServiceImpl) RemovePostVote(postID, userID string) error {
	if err := p.postRepository.RemoveVote(postID, userID); err != nil {
		p.zapLogger.Error("Failed to remove vote on post", zap.Error(err))
		return err
	}
	go p.refreshTopVotedCache()
	return nil
}

func (p *postServiceImpl) refreshTopVotedCache() {
	filter := &helpers.PaginateFilter{
		Page:         1,
		PageSize:     5,
		Sort:         "-vote_count",
		SortSafeList: []string{"vote_count", "-vote_count"},
	}
	posts, metadata, err := p.postRepository.GetAll(filter)
	if err != nil {
		p.zapLogger.Error("Failed to refresh top voted cache", zap.Error(err))
		return
	}
	cacheData := struct {
		Posts    []model.Post     `json:"posts"`
		Metadata helpers.Metadata `json:"metadata"`
	}{
		Posts:    *posts,
		Metadata: metadata,
	}
	jsonData, err := json.Marshal(cacheData)
	if err != nil {
		p.zapLogger.Error("Failed to marshal cache data", zap.Error(err))
		return
	}
	err = p.redisDB.Set(context.Background(), "top_5_posts", jsonData, time.Hour).Err()
	if err != nil {
		p.zapLogger.Error("Failed to store top voted posts in Redis", zap.Error(err))
	}
}
