package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/errors"

	v1 "sns/api/sns/post/v1"
)

// DDD: DO, Domain Object

type Post struct {
	ID        string    `json:"id"`         // ID
	Title     string    `json:"title"`      // 标题
	Content   string    `json:"content"`    // 正文
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// Repository Interface

type PostRepo interface {
	CreatePost(ctx context.Context, post *Post) (id string, err error)
	GetPost(ctx context.Context, id string) (post *Post, err error)
	ListPosts(ctx context.Context, keyword string) (list []*Post, err error)
	UpdatePost(ctx context.Context, id string, post *Post) error
	DeletePost(ctx context.Context, id string) error
}

// ------------

// DDD: Service

type PostUseCase struct {
	postRepo PostRepo
}

func NewPostUseCase(postRepo PostRepo) *PostUseCase {
	return &PostUseCase{
		postRepo: postRepo,
	}
}

func (uc *PostUseCase) Create(ctx context.Context, post *Post) (string, error) {
	if post.Title == "" {
		return "", errors.New(500, v1.ErrorReason_ERROR_REASON_INVALID_PARAMS.String(), "invalid title")
	}

	post.UpdatedAt = time.Now()

	return uc.postRepo.CreatePost(ctx, post)
}

func (uc *PostUseCase) Get(ctx context.Context, id string) (*Post, error) {
	return uc.postRepo.GetPost(ctx, id)
}

func (uc *PostUseCase) List(ctx context.Context, keyword string) ([]*Post, error) {
	return uc.postRepo.ListPosts(ctx, keyword)
}

func (uc *PostUseCase) Update(ctx context.Context, id string, post *Post) error {
	if post.ID == "" {
		return errors.New(500, v1.ErrorReason_ERROR_REASON_INVALID_PARAMS.String(), "invalid id")
	}
	if post.Title == "" {
		return errors.New(500, v1.ErrorReason_ERROR_REASON_INVALID_PARAMS.String(), "invalid title")
	}

	post.UpdatedAt = time.Now()

	return uc.postRepo.UpdatePost(ctx, id, post)
}

func (uc *PostUseCase) Delete(ctx context.Context, id string) error {
	return uc.postRepo.DeletePost(ctx, id)
}
