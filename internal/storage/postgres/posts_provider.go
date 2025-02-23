package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Lucky112/social/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type PostsProvider struct {
	querier pgxscan.Querier
}

func NewPostsProvider(querier pgxscan.Querier) PostsProvider {
	return PostsProvider{querier}
}

// TODO : add pagination
func (p PostsProvider) GetAll(ctx context.Context, userID string) ([]*models.Post, error) {
	var res []*models.Post

	postsInfo, err := p.getAllPostInfo(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("getting all posts info: %v", err)
	}

	for _, postInfo := range postsInfo {
		post, err := postInfo.toModel()
		if err != nil {
			return nil, fmt.Errorf("converting post info of '%d': %v", postInfo.Id, err)
		}

		res = append(res, post)
	}

	return res, nil
}

func (p PostsProvider) Get(ctx context.Context, userID, postID string) (*models.Post, error) {
	pid, err := strconv.ParseInt(postID, 10, 0)
	if err != nil {
		return nil, fmt.Errorf("illegal id '%s': %v : int64 expected", postID, err)
	}

	postInfo, err := p.getPostInfo(ctx, userID, pid)
	if err != nil {
		return nil, fmt.Errorf("getting post info of '%d': %w", pid, err)
	}

	post, err := postInfo.toModel()
	if err != nil {
		return nil, fmt.Errorf("converting post info of '%d': %v", pid, err)
	}

	return post, nil
}

func (p PostsProvider) Add(ctx context.Context, post *models.Post) (string, error) {
	query := `
		insert into scl.posts(user_id, content)
		values (@user, @content)
		returning id
	`

	args := pgx.NamedArgs{
		"user":    post.UserId,
		"content": post.Content,
	}

	rows, err := p.querier.Query(ctx, query, args)
	if err != nil {
		return "", fmt.Errorf("inserting into db: %v", err)
	}

	id, err := pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (int64, error) {
		var id int64
		err := row.Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("scanning post id: %v", err)
		}

		return id, nil
	})
	if err != nil {
		return "", fmt.Errorf("collecting new post id: %v", err)
	}

	return fmt.Sprintf("%d", id), nil
}

func (p PostsProvider) getPostInfo(ctx context.Context, userID string, postID int64) (*post, error) {
	var posts []post

	query := `
		select
			ps.id,
			ps.user_id,
			content,
			created_at
		from scl.posts as ps
		where id = $1 and user_id = $2
		order by created_at desc
	`

	err := pgxscan.Select(ctx, p.querier, &posts, query, postID, userID)
	if err != nil {
		return nil, fmt.Errorf("executing query `%s`: %v", query, err)
	}

	if len(posts) == 0 {
		return nil, fmt.Errorf("querying db: %w", models.PostNotFound)
	}

	return &posts[0], nil
}

func (p PostsProvider) getAllPostInfo(ctx context.Context, userID string) ([]post, error) {
	var posts []post

	query := `
		select
			ps.id,
			ps.user_id,
			content,
			created_at
		from scl.posts as ps
		where user_id = $1
		order by created_at desc
	`

	err := pgxscan.Select(ctx, p.querier, &posts, query, userID)
	if err != nil {
		return nil, fmt.Errorf("executing query `%s`: %w", query, err)
	}
	return posts, nil
}
