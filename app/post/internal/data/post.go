package data

import (
	"context"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"

	postv1 "sns/api/sns/post/v1"
	"sns/app/post/internal/biz"
	"sns/app/post/internal/data/ent"
	"sns/app/post/internal/data/ent/post"
	"sns/app/post/internal/data/ent/predicate"
)

var _ biz.PostRepo = (*postRepo)(nil)

type postRepo struct {
	data *EntData
	log  *log.Helper
}

func NewPostRepo(data *EntData, logger log.Logger) biz.PostRepo {
	return &postRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// ------------

func (repo *postRepo) CreatePost(ctx context.Context, bp *biz.Post) (id string, err error) {
	var ep *ent.Post
	ep, err = repo.data.PostClient(ctx).Create().
		SetTitle(bp.Title).SetContent(bp.Content).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			err = errors.Conflict(postv1.ErrorReason_ERROR_REASON_ALREADY_EXISTS.String(), "post already exists").
				WithCause(err)
		} else {
			err = errors.InternalServer(postv1.ErrorReason_ERROR_REASON_INTERNAL_ERROR.String(), "invalid query").
				WithCause(err)
		}

		return
	}

	id = ep.ID

	return
}

func (repo *postRepo) GetPost(ctx context.Context, id string) (bp *biz.Post, err error) {
	var ep *ent.Post
	ep, err = repo.data.PostClient(ctx).Query().
		Where(post.ID(id), post.DeletedAtEQ(0)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = errors.Conflict(postv1.ErrorReason_ERROR_REASON_NOT_FOUND.String(), "post not found").
				WithCause(err)
		} else {
			err = errors.InternalServer(postv1.ErrorReason_ERROR_REASON_INTERNAL_ERROR.String(), "invalid query").
				WithCause(err)
		}

		return
	}

	bp = &biz.Post{
		ID:        ep.ID,
		Title:     ep.Title,
		Content:   ep.Content,
		CreatedAt: ep.CreatedAt,
		UpdatedAt: ep.UpdatedAt,
	}

	return
}

func (repo *postRepo) ListPosts(ctx context.Context, keyword string) (bps []*biz.Post, err error) {
	// 空格分隔关键词
	keywords := strings.Split(strings.ReplaceAll(strings.TrimSpace(keyword), "  ", " "), " ")
	wheres := make([]predicate.Post, 0, len(keywords))
	for _, word := range keywords {
		if len(word) == 0 {
			continue
		}

		wheres = append(wheres, post.TitleContains(word))
	}

	var eps []*ent.Post
	eps, err = repo.data.PostClient(ctx).Query().
		Where(post.Or(wheres...), post.DeletedAtEQ(0)).
		Order(ent.Desc(post.FieldUpdatedAt)).
		// Offset(0).Limit(100).
		All(ctx)
	if err != nil {
		err = errors.InternalServer(postv1.ErrorReason_ERROR_REASON_INTERNAL_ERROR.String(), "invalid query").
			WithCause(err)

		return
	}

	bps = make([]*biz.Post, 0, len(eps))
	for _, ep := range eps {
		bps = append(bps, &biz.Post{
			ID:        ep.ID,
			Title:     ep.Title,
			Content:   ep.Content,
			CreatedAt: ep.CreatedAt,
			UpdatedAt: ep.UpdatedAt,
		})
	}

	return
}

func (repo *postRepo) UpdatePost(ctx context.Context, id string, bp *biz.Post) (err error) {
	if len(id) == 0 {
		return
	}

	_, err = repo.data.PostClient(ctx).UpdateOneID(id).
		SetTitle(bp.Title).SetContent(bp.Content).SetUpdatedAt(bp.UpdatedAt).
		Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = errors.Conflict(postv1.ErrorReason_ERROR_REASON_NOT_FOUND.String(), "post not found").
				WithCause(err)
		} else if ent.IsConstraintError(err) {
			err = errors.Conflict(postv1.ErrorReason_ERROR_REASON_ALREADY_EXISTS.String(), "post already exists").
				WithCause(err)
		} else {
			err = errors.InternalServer(postv1.ErrorReason_ERROR_REASON_INTERNAL_ERROR.String(), "invalid query").
				WithCause(err)
		}

		return
	}

	return
}

func (repo *postRepo) DeletePost(ctx context.Context, id string) (err error) {
	if len(id) == 0 {
		return
	}

	_, err = repo.data.PostClient(ctx).UpdateOneID(id).
		SetDeletedAt(int(time.Now().Unix())).
		Save(ctx)
	if err != nil {
		err = errors.InternalServer(postv1.ErrorReason_ERROR_REASON_INTERNAL_ERROR.String(), "invalid query").
			WithCause(err)

		return
	}

	return
}
