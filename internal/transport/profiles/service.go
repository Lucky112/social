package profiles

import (
	"context"

	"github.com/Lucky112/social/internal/models"
)

// Сервис профилей пользователей
type ProfilesService interface {
	GetAll(ctx context.Context) ([]*models.Profile, error)
	Search(ctx context.Context, params *models.SearchParams) ([]*models.Profile, error)
	Get(ctx context.Context, id string) (*models.Profile, error)
	Add(ctx context.Context, profile *models.Profile) (string, error)
}
