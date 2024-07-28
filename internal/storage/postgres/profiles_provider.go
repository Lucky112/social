package postgres

import (
	"context"
	"fmt"

	"github.com/Lucky112/social/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
)

type ProfilesProvider struct {
	querier pgxscan.Querier
}

func (p ProfilesProvider) GetAll(ctx context.Context) ([]Profile, error) {
	var profiles []Profile

	query := `
		select
			id,
			name,
			surname,
			age,
			sex
		from scl.profiles
	`

	err := pgxscan.Select(ctx, p.querier, &profiles, query)
	if err != nil {
		return nil, fmt.Errorf("executing query `%s`: %v", query, err)
	}

	return profiles, nil
}

func (p ProfilesProvider) Get(ctx context.Context, profileID int64) (*Profile, error) {
	var profiles []Profile

	query := `
		select
			id,
			name,
			surname,
			age,
			sex
		from scl.profiles
		where id = $1
	`

	err := pgxscan.Select(ctx, p.querier, &profiles, query, profileID)
	if err != nil {
		return nil, fmt.Errorf("executing query `%s`: %v", query, err)
	}

	if len(profiles) == 0 {
		return nil, fmt.Errorf("querying db: %v", models.ProfileNotFound)
	}

	return &profiles[0], nil
}
