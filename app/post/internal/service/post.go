package service

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

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

// ------------

func (p *PostService) RegisterServiceServer(srv *grpc.Server) {
	postv1.RegisterPostServiceServer(srv, p)
}

func (p *PostService) RegisterServiceHTTPServer(srv *http.Server) {
	postv1.RegisterPostServiceHTTPServer(srv, p)
}

// ------------

func (p *PostService) serviceToBiz(so *postv1.Post, bo *biz.Post) *biz.Post {
	if bo == nil {
		bo = &biz.Post{}
	}

	bo.ID = so.Id
	bo.Title = so.Title
	bo.Content = so.Content
	bo.CreatedAt = time.Unix(so.CreatedAt, 0)
	bo.UpdatedAt = time.Unix(so.UpdatedAt, 0)

	return bo
}

func (p *PostService) bizToService(bo *biz.Post, so *postv1.Post) *postv1.Post {
	if so == nil {
		so = &postv1.Post{}
	}

	so.Id = bo.ID
	so.Title = bo.Title
	so.Content = bo.Content
	so.CreatedAt = bo.CreatedAt.Unix()
	so.UpdatedAt = bo.UpdatedAt.Unix()

	return so
}

// ------------

func (p *PostService) CreatePost(ctx context.Context, req *postv1.CreatePostRequest) (*postv1.CreatePostResponse, error) {
	id, err := p.postUseCase.Create(ctx, p.serviceToBiz(req.Post, nil))
	if err != nil {
		return nil, err
	}

	resp := &postv1.CreatePostResponse{Id: id}

	return resp, nil
}

func (p *PostService) GetPost(ctx context.Context, req *postv1.GetPostRequest) (*postv1.GetPostResponse, error) {
	bizPost, err := p.postUseCase.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	resp := &postv1.GetPostResponse{Post: p.bizToService(bizPost, nil)}

	return resp, nil
}

func (p *PostService) ListPosts(ctx context.Context, req *postv1.ListPostsRequest) (*postv1.ListPostsResponse, error) {
	bizPosts, err := p.postUseCase.List(ctx, req.Ids, req.Keyword)
	if err != nil {
		return nil, err
	}

	resp := &postv1.ListPostsResponse{
		Posts: make([]*postv1.Post, len(bizPosts)),
	}

	for i, bizPost := range bizPosts {
		resp.Posts[i] = p.bizToService(bizPost, nil)
	}

	return resp, nil
}

func (p *PostService) UpdatePost(ctx context.Context, req *postv1.UpdatePostRequest) (*postv1.UpdatePostResponse, error) {
	err := p.postUseCase.Update(ctx, req.Post.Id, p.serviceToBiz(req.Post, nil))
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
