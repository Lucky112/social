package profiles

import "github.com/Lucky112/social/internal/models"

// Хранилище профилей пользователей
type ProfilesStorage interface {
	GetAll() ([]*models.Profile, error)
	Get(id string) (*models.Profile, error)
	Add(id string, profile *models.Profile) error
}
