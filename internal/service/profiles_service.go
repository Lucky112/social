package service

import (
	"context"

	"github.com/Lucky112/social/internal/models"
)

type ProfilesService struct {
	storage ProfilesStorage
}

// Хранилище зарегистрированных пользователей
type ProfilesStorage interface {
	GetAll(ctx context.Context) ([]*models.Profile, error)
	Get(ctx context.Context, id string) (*models.Profile, error)
	Add(ctx context.Context, profile *models.Profile) (string, error)
}

func NewProfilesService(storage ProfilesStorage) ProfilesService {
	return ProfilesService{
		storage: storage,
	}
}

func (s ProfilesService) GetAll(ctx context.Context) ([]*models.Profile, error) {
	return s.storage.GetAll(ctx)
}

func (s ProfilesService) Get(ctx context.Context, id string) (*models.Profile, error) {
	return s.storage.Get(ctx, id)
}

func (s ProfilesService) Add(ctx context.Context, profile *models.Profile) (string, error) {
	return s.storage.Add(ctx, profile)
}
