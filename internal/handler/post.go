package handler

import (
	"database/sql"
	"errors"
	"github.com/arshamroshannejad/task-rootext/internal/domain"
	"github.com/arshamroshannejad/task-rootext/internal/entities"
	"github.com/arshamroshannejad/task-rootext/internal/helpers"
	"github.com/go-chi/chi/v5"
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

// GetAllPostsHandler godoc
//
//	@Summary		Get All Posts
//	@Description	this endpoint provide all posts. also (pagination, sort, order) is available.
//	@Accept			json
//	@Produce		json
//	@Tags			Posts
//	@Param			_	query		helpers.PostQueryParams	false	"Query Params"
//	@Success		200	{object}	helpers.AllPosts
//	@Success		400	{object}	helpers.BadRequest
//	@Failure		500	{object}	helpers.InternalServerError
//	@router			/post [get]
func (p *PostHandlerImpl) GetAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	var filter helpers.PaginateFilter
	v := helpers.NewValidator()
	qs := r.URL.Query()
	filter.Page = v.ReadQsInt(qs, "page", 1)
	filter.PageSize = v.ReadQsInt(qs, "page_size", 5)
	filter.Sort = v.ReadQsString(qs, "sort", "-vote_count")
	filter.SortSafeList = []string{"created_at", "-created_at", "vote_count", "-vote_count"}
	if filter.Validate(v); !v.IsValid() {
		helpers.WriteJson(w, http.StatusBadRequest, helpers.M{"error": v.Errors})
		return
	}
	posts, metaData, err := p.PostService.GetAllPosts(&filter)
	if err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.M{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	helpers.WriteJson(w, http.StatusOK, helpers.M{"metadata": metaData, "posts": posts})
}

// GetPostHandler  godoc
//
//	@Summary		Get a single post
//	@Description	Get a sing post with id
//	@Accept			json
//	@Produce		json
//	@Tags			Posts
//	@Param			id	path		int	true	"Post ID"
//	@Success		200	{object}	helpers.Post
//	@Failure		404	{object}	helpers.PostNotFound
//	@Failure		500	{object}	helpers.InternalServerError
//	@router			/post/{id} [get]
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

// CreatePostHandler godoc
//
//	@Summary		Create a new post
//	@Description	Create a new post with the provided data
//	@Accept			json
//	@Produce		json
//	@Tags			Posts
//	@Security		BearerAuth
//	@Param			postBody	body		entities.PostCreateUpdateRequest	true	"just send title and text. authenticated required!"
//	@Success		201			{object}	helpers.Post
//	@Failure		400			{object}	helpers.BadRequest
//	@Failure		500			{object}	helpers.InternalServerError
//	@Router			/post [post]
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

// UpdatePostHandler godoc
//
//	@Summary		Update an existing post
//	@Description	Update a post with the provided ID and data
//	@Accept			json
//	@Produce		json
//	@Tags			Posts
//	@Security		BearerAuth
//	@Param			id			path		int									true	"Post ID"
//	@Param			postBody	body		entities.PostCreateUpdateRequest	true	"just send title and text. authenticated required!"
//	@Success		200			{object}	helpers.Post
//	@Failure		400			{object}	helpers.BadRequest
//	@Failure		404			{object}	helpers.PostNotFound
//	@Failure		403			{object}	helpers.Forbidden
//	@Failure		500			{object}	helpers.InternalServerError
//	@Router			/post/{id} [put]
func (p *PostHandlerImpl) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	userID := helpers.GetUserID(r)
	postID := chi.URLParam(r, "id")
	reqBody := new(entities.PostCreateUpdateRequest)
	if err := helpers.ReadJson(r, reqBody); err != nil {
		helpers.WriteJson(w, http.StatusBadRequest, helpers.M{"error": err.Error()})
		return
	}
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
	if post.UserID != userID {
		helpers.WriteJson(w, http.StatusForbidden, helpers.M{"error": http.StatusText(http.StatusForbidden)})
		return
	}
	updatedPost, err := p.PostService.UpdatePost(reqBody, postID)
	if err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.M{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	helpers.WriteJson(w, http.StatusOK, updatedPost)
}

// DeletePostHandler godoc
//
//	@Summary		Delete a post
//	@Description	Delete a post by ID. authenticated required!
//	@Accept			json
//	@Produce		json
//	@Tags			Posts
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Post ID"
//	@Success		204	{object}	nil
//	@Failure		404	{object}	helpers.PostNotFound
//	@Failure		403	{object}	helpers.Forbidden
//	@Failure		500	{object}	helpers.InternalServerError
//	@Router			/post/{id} [delete]
func (p *PostHandlerImpl) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	userID := helpers.GetUserID(r)
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
	if post.UserID != userID {
		helpers.WriteJson(w, http.StatusForbidden, helpers.M{"error": http.StatusText(http.StatusForbidden)})
		return
	}
	if err := p.PostService.DeletePost(postID); err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.M{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	helpers.WriteJson(w, http.StatusNoContent, nil)
}

// AddPostVoteHandler godoc
//
//	@Summary		Add a vote to a post
//	@Description	Add a vote for a post by ID. authenticated required!
//	@Accept			json
//	@Produce		json
//	@Tags			Posts
//	@Security		BearerAuth
//	@Param			id			path		int						true	"Post ID"
//	@Param			voteBody	body		entities.VoteRequest	true	"value must 1 or -1"
//	@Success		200			{object}	helpers.VoteSuccessful
//	@Failure		400			{object}	helpers.BadRequest
//	@Failure		404			{object}	helpers.PostNotFound
//	@Failure		403			{object}	helpers.Forbidden
//	@Failure		500			{object}	helpers.InternalServerError
//	@Router			/post/{id}/vote [post]
func (p *PostHandlerImpl) AddPostVoteHandler(w http.ResponseWriter, r *http.Request) {
	userID := helpers.GetUserID(r)
	postID := chi.URLParam(r, "id")
	reqBody := new(entities.VoteRequest)
	if err := helpers.ReadJson(r, reqBody); err != nil {
		helpers.WriteJson(w, http.StatusBadRequest, helpers.M{"error": err.Error()})
		return
	}
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
	if post.UserID == userID {
		helpers.WriteJson(w, http.StatusForbidden, helpers.M{"error": http.StatusText(http.StatusForbidden)})
		return
	}
	if err := p.PostService.AddPostVote(postID, userID, reqBody.Value); err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.M{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	helpers.WriteJson(w, http.StatusOK, helpers.M{"response": "successful"})
}

// RemovePostVoteHandler godoc
//
//	@Summary		Remove a vote from a post
//	@Description	Remove a vote for a post by ID. authenticated required!
//	@Accept			json
//	@Produce		json
//	@Tags			Posts
//	@Security		BearerAuth
//	@Param			id	path		int	true	"Post ID"
//	@Success		204	{object}	nil
//	@Failure		404	{object}	helpers.PostNotFound
//	@Failure		500	{object}	helpers.InternalServerError
//	@Router			/post/{id}/vote [delete]
func (p *PostHandlerImpl) RemovePostVoteHandler(w http.ResponseWriter, r *http.Request) {
	userID := helpers.GetUserID(r)
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
	if err := p.PostService.RemovePostVote(postID, userID); err != nil {
		helpers.WriteJson(w, http.StatusInternalServerError, helpers.M{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	helpers.WriteJson(w, http.StatusNoContent, nil)
}
