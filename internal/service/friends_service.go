package service

import (
	"context"

	"github.com/Lucky112/social/internal/models"
)

type FriendsService struct {
	storage FriendsStorage
}

// Хранилище друзей пользователей
type FriendsStorage interface {
	Get(ctx context.Context, id string) (*models.Friends, error)
	Add(ctx context.Context, id, friendId string) error
	Delete(ctx context.Context, id, friendId string) error
}

func NewFriendsService(storage FriendsStorage) FriendsService {
	return FriendsService{
		storage: storage,
	}
}

func (s FriendsService) Get(ctx context.Context, id string) (*models.Friends, error) {
	return s.storage.Get(ctx, id)
}

func (s FriendsService) Add(ctx context.Context, id, friendId string) error {
	return s.storage.Add(ctx, id, friendId)
}

func (s FriendsService) Delete(ctx context.Context, id, friendId string) error {
	return s.storage.Delete(ctx, id, friendId)
}
