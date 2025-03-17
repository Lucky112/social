package friends

import (
	"context"

	"github.com/Lucky112/social/internal/models"
)

type FriendsService interface {
	Get(ctx context.Context, id string) (*models.Friends, error)
	Add(ctx context.Context, id, friendId string) error
	Delete(ctx context.Context, id, friendId string) error
}

type ProfileService interface {
	Get(ctx context.Context, id string) (*models.Profile, error)
}
