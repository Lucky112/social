package inmemory

import (
	"context"
	"fmt"

	"github.com/Lucky112/social/internal/models"
)

type ProfileStorage map[string]*models.Profile

func NewProfileStorage() ProfileStorage {
	return make(ProfileStorage)
}

func (ps ProfileStorage) GetAll(ctx context.Context) ([]*models.Profile, error) {
	res := make([]*models.Profile, 0, len(ps))

	for _, p := range ps {
		res = append(res, p)
	}

	return res, nil
}

func (ps ProfileStorage) Get(ctx context.Context, id string) (*models.Profile, error) {
	p, exists := ps[id]
	if !exists {
		return nil, fmt.Errorf("looking '%s' up: %w", id, models.ProfileNotFound)
	}

	return p, nil
}

func (ps ProfileStorage) Add(ctx context.Context, profile *models.Profile) (string, error) {
	if profile == nil {
		return "", fmt.Errorf("attempt to store nil profile")
	}

	id := generateId()
	ps[id] = profile

	return id, nil
}
