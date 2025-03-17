package service

import (
	"context"

	"github.com/Lucky112/social/internal/models"
)

type PostsService struct {
	storage PostsStorage
}

// Хранилище постов пользователей
type PostsStorage interface {
	GetAll(ctx context.Context, userID string) ([]*models.Post, error)
	Get(ctx context.Context, userID, postID string) (*models.Post, error)
	Add(ctx context.Context, post *models.Post) (string, error)
}

func NewPostsService(storage PostsStorage) PostsService {
	return PostsService{
		storage: storage,
	}
}

func (s PostsService) GetAll(ctx context.Context, userID string) ([]*models.Post, error) {
	return s.storage.GetAll(ctx, userID)
}

func (s PostsService) Get(ctx context.Context, userID, postID string) (*models.Post, error) {
	return s.storage.Get(ctx, userID, postID)
}

func (s PostsService) Add(ctx context.Context, post *models.Post) (string, error) {
	return s.storage.Add(ctx, post)
}
