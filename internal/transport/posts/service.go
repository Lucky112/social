package posts

import (
	"context"

	"github.com/Lucky112/social/internal/models"
)

// Сервис постов пользователей
type PostsService interface {
	GetAll(ctx context.Context, userID string) ([]*models.Post, error)
	Get(ctx context.Context, userID, postID string) (*models.Post, error)
	Add(ctx context.Context, profile *models.Post) (string, error)
}
