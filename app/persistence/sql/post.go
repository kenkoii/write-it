package sql

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rbo13/write-it/app"
)

var (
	errEmpty       = errors.New("error: Post is required")
	errNotInserted = errors.New("error: Not inserted")
	errNoID        = errors.New("error: ID is required")
	errPostDelete  = errors.New("error: Post deletion")
	errPostUpdate  = errors.New("error: Post update")
)

// PostService implements the app.UserService
type PostService interface {
	app.PostService
}

// Post implements the PostService interface
type Post struct {
	DB       *sqlx.DB
	PostSrvc *app.Post
}

// NewPostSQLService returns the interface that implements the app.PostService
func NewPostSQLService(db *sqlx.DB) PostService {
	return &Post{
		DB:       db,
		PostSrvc: new(app.Post),
	}
}

// CreatePost ...
func (p *Post) CreatePost(post *app.Post) error {
	if post == nil {
		return errEmpty
	}

	tx := p.DB.MustBegin()

	post.CreatedAt = time.Now().Unix()

	res, err := tx.NamedExec("INSERT INTO posts (creator_id, post_title, post_body, created_at, deleted_at, updated_at) VALUES(:creator_id, :post_title, :post_body, :created_at, :deleted_at, :updated_at)", &post)

	if err != nil && res == nil {
		tx.Rollback()
		return errNotInserted
	}
	tx.Commit()
	return nil
}

// Post ...
func (p *Post) Post(id int64) (*app.Post, error) {
	if id <= 0 {
		return nil, errNoID
	}

	post := new(app.Post)

	err := p.DB.Get(post, "SELECT * FROM posts WHERE id = ? LIMIT 1;", id)

	if err != nil {
		return nil, err
	}

	return post, nil
}

// Posts ...
func (p *Post) Posts() ([]*app.Post, error) {
	posts := []*app.Post{}

	err := p.DB.Select(&posts, "SELECT * FROM posts ORDER BY id DESC;")

	if err != nil {
		return nil, err
	}
	return posts, nil
}

// UpdatePost ...
func (p *Post) UpdatePost(post *app.Post) error {
	post.UpdatedAt = time.Now().Unix()

	tx := p.DB.MustBegin()
	res := tx.MustExec("UPDATE posts SET post_title = ?, post_body = ?, created_at = ?, updated_at = ? WHERE id = ? AND creator_id = ? LIMIT 1;", post.PostTitle, post.PostBody, post.CreatedAt, post.UpdatedAt, post.ID, post.CreatorID)

	if res == nil {
		tx.Rollback()
		return errPostUpdate
	}

	tx.Commit()
	return nil
}

// DeletePost ...
func (p *Post) DeletePost(id int64) error {
	tx := p.DB.MustBegin()

	res := tx.MustExec("DELETE FROM posts WHERE id = ?;", id)

	if res == nil {
		tx.Rollback()
		return errPostDelete
	}

	tx.Commit()
	return nil
}
