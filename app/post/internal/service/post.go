package service

import (
	"context"

	postv1 "sns/api/sns/post/v1"
	"sns/app/post/internal/biz"
)

type PostService struct {
	postv1.UnimplementedPostServiceServer

	postUseCase *biz.PostUseCase
}

// NewPostService new an org service.
func NewPostService(postUseCase *biz.PostUseCase) *PostService {
	return &PostService{
		postUseCase: postUseCase,
	}
}

func (p *PostService) CreatePost(ctx context.Context, req *postv1.CreatePostRequest) (*postv1.CreatePostResponse, error) {
	id, err := p.postUseCase.Create(ctx, &biz.Post{
		ID:      req.Post.Id,
		Title:   req.Post.Title,
		Content: req.Post.Content,
	})
	if err != nil {
		return nil, err
	}

	resp := &postv1.CreatePostResponse{Id: id}

	return resp, nil
}

func (p *PostService) GetPost(ctx context.Context, req *postv1.GetPostRequest) (*postv1.GetPostResponse, error) {
	post, err := p.postUseCase.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	resp := &postv1.GetPostResponse{
		Post: &postv1.Post{
			Id:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			CreatedAt: post.CreatedAt.Unix(),
			UpdatedAt: post.UpdatedAt.Unix(),
		},
	}

	return resp, nil
}

func (p *PostService) ListPosts(ctx context.Context, req *postv1.ListPostsRequest) (*postv1.ListPostsResponse, error) {
	posts, err := p.postUseCase.List(ctx, req.Keyword)
	if err != nil {
		return nil, err
	}

	resp := &postv1.ListPostsResponse{Posts: make([]*postv1.Post, 0, len(posts))}
	for _, post := range posts {
		resp.Posts = append(resp.Posts, &postv1.Post{
			Id:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			CreatedAt: post.CreatedAt.Unix(),
			UpdatedAt: post.UpdatedAt.Unix(),
		})
	}

	return resp, nil
}

func (p *PostService) UpdatePost(ctx context.Context, req *postv1.UpdatePostRequest) (*postv1.UpdatePostResponse, error) {
	err := p.postUseCase.Update(ctx, req.Post.Id, &biz.Post{
		ID:      req.Post.Id,
		Title:   req.Post.Title,
		Content: req.Post.Content,
	})
	if err != nil {
		return nil, err
	}

	return &postv1.UpdatePostResponse{}, nil
}

func (p *PostService) DeletePost(ctx context.Context, req *postv1.DeletePostRequest) (*postv1.DeletePostResponse, error) {
	err := p.postUseCase.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &postv1.DeletePostResponse{}, nil
}
