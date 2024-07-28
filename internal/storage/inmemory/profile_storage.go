package inmemory

import (
	"fmt"

	"github.com/Lucky112/social/internal/models"
)

type ProfileStorage map[string]*models.Profile

func NewProfileStorage() ProfileStorage {
	return make(ProfileStorage)
}

func (ps ProfileStorage) GetAll() ([]*models.Profile, error) {
	res := make([]*models.Profile, 0, len(ps))

	for _, p := range ps {
		res = append(res, p)
	}

	return res, nil
}

func (ps ProfileStorage) Get(id string) (*models.Profile, error) {
	p, exists := ps[id]
	if !exists {
		return nil, fmt.Errorf("looking '%s' up: %w", id, models.ProfileNotFound)
	}

	return p, nil
}

func (ps ProfileStorage) Add(id string, profile *models.Profile) error {
	if profile == nil {
		return fmt.Errorf("attempt to store nil profile for id '%s'", id)
	}

	ps[id] = profile
	return nil
}
