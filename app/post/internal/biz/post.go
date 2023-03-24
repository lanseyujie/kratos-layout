package biz

import (
	"context"
	"strings"
	"time"

	postv1 "sns/api/sns/post/v1"
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
	CreatePost(ctx context.Context, bizPost *Post) (id string, err error)
	GetPost(ctx context.Context, id string) (bizPost *Post, err error)
	ListPosts(ctx context.Context, ids, keywords []string) (bizPosts []*Post, err error)
	CountAccount(ctx context.Context, ids, keywords []string) (count int, err error)
	ExistAccount(ctx context.Context, title string) (bool, error)
	UpdatePost(ctx context.Context, id string, bizPost *Post) error
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

func (uc *PostUseCase) Create(ctx context.Context, bizPost *Post) (string, error) {
	if bizPost.Title == "" {
		return "", postv1.ErrorErrorReasonInvalidParams("invalid title")
	}

	bizPost.UpdatedAt = time.Now()

	return uc.postRepo.CreatePost(ctx, bizPost)
}

func (uc *PostUseCase) Get(ctx context.Context, id string) (*Post, error) {
	return uc.postRepo.GetPost(ctx, id)
}

func (uc *PostUseCase) List(ctx context.Context, ids []string, keyword string) ([]*Post, error) {
	keywords := strings.Fields(keyword)

	return uc.postRepo.ListPosts(ctx, ids, keywords)
}

func (uc *PostUseCase) Update(ctx context.Context, id string, bizPost *Post) error {
	if bizPost.ID == "" {
		return postv1.ErrorErrorReasonInvalidParams("invalid id")
	}
	if bizPost.Title == "" {
		return postv1.ErrorErrorReasonInvalidParams("invalid title")
	}

	bizPost.UpdatedAt = time.Now()

	return uc.postRepo.UpdatePost(ctx, id, bizPost)
}

func (uc *PostUseCase) Delete(ctx context.Context, id string) error {
	return uc.postRepo.DeletePost(ctx, id)
}
