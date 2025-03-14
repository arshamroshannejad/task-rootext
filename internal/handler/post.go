package handler

import (
	"database/sql"
	"errors"
	"github.com/go-chi/chi/v5"
	"github/arshamroshannejad/task-rootext/internal/domain"
	"github/arshamroshannejad/task-rootext/internal/entities"
	"github/arshamroshannejad/task-rootext/internal/helpers"
	"net/http"
)

type PostHandlerImpl struct {
	PostService domain.PostService
}

func NewPostHandler(postService domain.PostService) *PostHandlerImpl {
	return &PostHandlerImpl{
		PostService: postService,
	}
}

func (p *PostHandlerImpl) GetAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := p.PostService.GetAllPosts()
	if err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.M{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	helpers.WriteJson(w, http.StatusOK, posts)
}

func (p *PostHandlerImpl) GetPostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")
	post, err := p.PostService.GetPostByID(postID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			helpers.WriteJson(w, http.StatusNotFound, helpers.M{"error": "post not found"})
		default:
			helpers.WriteJson(w, http.StatusInternalServerError, helpers.M{"error": http.StatusText(http.StatusInternalServerError)})
		}
		return
	}
	helpers.WriteJson(w, http.StatusOK, post)
}

func (p *PostHandlerImpl) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	userID := helpers.GetUserID(r)
	reqBody := new(entities.PostCreateUpdateRequest)
	if err := helpers.ReadJson(r, reqBody); err != nil {
		helpers.WriteJson(w, http.StatusBadRequest, helpers.M{"error": err.Error()})
		return
	}
	createdPost, err := p.PostService.CreatePost(reqBody, userID)
	if err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.M{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	helpers.WriteJson(w, http.StatusCreated, createdPost)
}

func (p *PostHandlerImpl) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	userID := helpers.GetUserID(r)
	postID := chi.URLParam(r, "id")
	reqBody := new(entities.PostCreateUpdateRequest)
	if err := helpers.ReadJson(r, reqBody); err != nil {
		helpers.WriteJson(w, http.StatusBadRequest, helpers.M{"error": err.Error()})
		return
	}
	if _, err := p.PostService.GetPostByID(postID); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			helpers.WriteJson(w, http.StatusNotFound, helpers.M{"error": "post not found"})
		default:
			helpers.WriteJson(w, http.StatusInternalServerError, helpers.M{"error": http.StatusText(http.StatusInternalServerError)})
		}
		return
	}
	updatedPost, err := p.PostService.UpdatePost(reqBody, userID)
	if err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.M{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	helpers.WriteJson(w, http.StatusCreated, updatedPost)
}

func (p *PostHandlerImpl) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")
	if _, err := p.PostService.GetPostByID(postID); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			helpers.WriteJson(w, http.StatusNotFound, helpers.M{"error": "post not found"})
		default:
			helpers.WriteJson(w, http.StatusInternalServerError, helpers.M{"error": http.StatusText(http.StatusInternalServerError)})
		}
		return
	}
	if err := p.PostService.DeletePost(postID); err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.M{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	helpers.WriteJson(w, http.StatusNoContent, nil)
}
