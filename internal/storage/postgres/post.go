package postgres

import (
	"time"

	"github.com/Lucky112/social/internal/models"
)

type post struct {
	Id        int64     `db:"id"`
	UserId    string    `db:"user_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}

func (p *post) toModel() (*models.Post, error) {
	return &models.Post{
		UserId:    p.UserId,
		Content:   p.Content,
		CreatedAt: p.CreatedAt,
	}, nil
}
