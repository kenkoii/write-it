package usecase

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"

	"github.com/rbo13/write-it/app"
	"github.com/rbo13/write-it/app/response"
)

type postUsecase struct {
	postService app.PostService
}

type postResponse struct {
	StatusCode uint        `json:"status_code"`
	Message    string      `json:"message"`
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
}

// NewPost ...
func NewPost(postService app.PostService) app.Handler {
	return &postUsecase{
		postService,
	}
}

func (p *postUsecase) Create(w http.ResponseWriter, r *http.Request) {
	var post app.Post

	_, claims, err := jwtauth.FromContext(r.Context())

	if err != nil {
		config := response.Configure(err.Error(), http.StatusForbidden, nil)
		response.JSONError(w, r, config)
		return
	}

	post.CreatorID = int64(claims["user_id"].(float64))

	err = json.NewDecoder(r.Body).Decode(&post)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusBadRequest, nil)
		response.JSONError(w, r, config)
		return
	}

	err = p.postService.CreatePost(&post)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusBadRequest, nil)
		response.JSONError(w, r, config)
		return
	}

	config := response.Configure("Post created successfully", http.StatusOK, post)
	response.JSONOK(w, r, config)
	return
}

func (p *postUsecase) Get(w http.ResponseWriter, r *http.Request) {
	posts, err := p.postService.Posts()

	if err != nil {
		config := response.Configure(err.Error(), http.StatusInternalServerError, nil)
		response.JSONError(w, r, config)
		return
	}

	config := response.Configure("Posts successfully retrieved", http.StatusOK, posts)
	response.JSONOK(w, r, config)
}

func (p *postUsecase) GetByID(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusBadRequest, nil)
		response.JSONError(w, r, config)
		return
	}

	post, err := p.postService.Post(postID)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusNotFound, nil)
		response.JSONError(w, r, config)
		return
	}

	config := response.Configure("Post successfully retrieved", http.StatusOK, post)
	response.JSONOK(w, r, config)
}

func (p *postUsecase) Update(w http.ResponseWriter, r *http.Request) {
	var post app.Post
	postID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	check(err, w, r)

	_, claims, err := jwtauth.FromContext(r.Context())

	check(err, w, r)

	// find a user by the given id
	postFetchRes, err := p.postService.Post(postID)

	check(err, w, r)

	post.ID = postFetchRes.ID
	post.CreatorID = int64(claims["user_id"].(float64))
	post.CreatedAt = postFetchRes.CreatedAt

	err = json.NewDecoder(r.Body).Decode(&post)

	check(err, w, r)

	err = p.postService.UpdatePost(&post)

	check(err, w, r)

	config := response.Configure("Post Successfully Updated", http.StatusOK, post)
	response.JSONOK(w, r, config)
}

func (p *postUsecase) Delete(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusBadRequest, nil)
		response.JSONError(w, r, config)
		return
	}

	err = p.postService.DeletePost(postID)

	if err != nil {
		config := response.Configure(err.Error(), http.StatusBadRequest, nil)
		response.JSONError(w, r, config)
		return
	}

	config := response.Configure("Post Successfully Deleted", http.StatusOK, nil)
	response.JSONOK(w, r, config)
}

func check(err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		config := response.Configure(err.Error(), http.StatusBadRequest, nil)
		response.JSONError(w, r, config)
		return
	}
	return
}
