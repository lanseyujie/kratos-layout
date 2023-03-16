package data

import (
	"context"
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

func (repo *postRepo) bizToData(bo *biz.Post, do *ent.Post) *ent.Post {
	if do == nil {
		do = &ent.Post{}
	}

	do.ID = bo.ID
	do.Title = bo.Title
	do.Content = bo.Content
	do.CreatedAt = bo.CreatedAt
	do.UpdatedAt = bo.UpdatedAt

	return do
}

func (repo *postRepo) dataToBiz(do *ent.Post, bo *biz.Post) *biz.Post {
	if bo == nil {
		bo = &biz.Post{}
	}

	bo.ID = do.ID
	bo.Title = do.Title
	bo.Content = do.Content
	bo.CreatedAt = do.CreatedAt
	bo.UpdatedAt = do.UpdatedAt

	return bo
}

// ------------

func (repo *postRepo) CreatePost(ctx context.Context, bizPost *biz.Post) (id string, err error) {
	var entPost *ent.Post
	entPost, err = repo.data.PostClient(ctx).Create().
		SetTitle(bizPost.Title).SetContent(bizPost.Content).
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

	id = entPost.ID

	return
}

func (repo *postRepo) GetPost(ctx context.Context, id string) (bizPost *biz.Post, err error) {
	var entPost *ent.Post
	entPost, err = repo.data.PostClient(ctx).Query().
		Where(post.ID(id), post.DeletedAtEQ(0)).
		First(ctx)
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

	bizPost = repo.dataToBiz(entPost, nil)

	return
}

func (repo *postRepo) list(ctx context.Context, ids, keywords []string) *ent.PostQuery {
	q := repo.data.PostClient(ctx).Query()
	if len(ids) > 0 {
		q.Where(post.IDIn(ids...))
	}
	q.Where(post.DeletedAt(0))

	conds := make([]predicate.Post, 0, len(keywords))
	for _, keyword := range keywords {
		conds = append(conds, post.TitleContainsFold(keyword))
	}
	if len(conds) > 0 {
		q.Where(post.Or(conds...))
	}

	return q
}

func (repo *postRepo) ListPosts(ctx context.Context, ids, keywords []string) (bizPosts []*biz.Post, err error) {
	var entPosts []*ent.Post
	entPosts, err = repo.list(ctx, ids, keywords).
		Order(ent.Desc(post.FieldUpdatedAt), ent.Asc(post.FieldTitle), ent.Desc(post.FieldCreatedAt)).
		Order(ent.Desc(post.FieldUpdatedAt)).
		// Offset(0).Limit(100).
		All(ctx)
	if err != nil {
		err = errors.InternalServer(postv1.ErrorReason_ERROR_REASON_INTERNAL_ERROR.String(), "invalid query").
			WithCause(err)

		return
	}

	bizPosts = make([]*biz.Post, len(entPosts))
	for i, entPost := range entPosts {
		bizPosts[i] = repo.dataToBiz(entPost, nil)
	}

	return
}

func (repo *postRepo) CountAccount(ctx context.Context, ids, keywords []string) (count int, err error) {
	count, err = repo.list(ctx, ids, keywords).Count(ctx)

	return
}

func (repo *postRepo) ExistAccount(ctx context.Context, title string) (bool, error) {
	count, err := repo.data.PostClient(ctx).Query().
		Where(post.Title(title)).
		Where(post.DeletedAt(0)).
		Count(ctx)

	return count > 0, err
}

func (repo *postRepo) UpdatePost(ctx context.Context, id string, bizPost *biz.Post) (err error) {
	if len(id) == 0 {
		return
	}

	_, err = repo.data.PostClient(ctx).UpdateOneID(id).
		Where(post.DeletedAt(0)).
		SetTitle(bizPost.Title).
		SetContent(bizPost.Content).
		SetUpdatedAt(bizPost.UpdatedAt).
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
		Where(post.DeletedAt(0)).
		SetDeletedAt(int(time.Now().Unix())).
		Save(ctx)
	if err != nil {
		err = errors.InternalServer(postv1.ErrorReason_ERROR_REASON_INTERNAL_ERROR.String(), "invalid query").
			WithCause(err)

		return
	}

	return
}
