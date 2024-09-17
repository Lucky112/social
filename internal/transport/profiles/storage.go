package profiles

import (
	"context"

	"github.com/Lucky112/social/internal/models"
)

// Хранилище профилей пользователей
type ProfilesStorage interface {
	GetAll(ctx context.Context) ([]*models.Profile, error)
	Get(ctx context.Context, id string) (*models.Profile, error)
	Add(ctx context.Context, profile *models.Profile) (string, error)
}
