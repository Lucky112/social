package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Lucky112/social/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type ProfilesProvider struct {
	querier pgxscan.Querier
}

func NewProfilesProvider(querier pgxscan.Querier) ProfilesProvider {
	return ProfilesProvider{querier}
}

// TODO : add pagination
func (p ProfilesProvider) GetAll(ctx context.Context) ([]*models.Profile, error) {
	var res []*models.Profile

	profilesInfo, err := p.getAllProfileInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting all profiles info: %v", err)
	}

	for _, profileInfo := range profilesInfo {
		profile, err := profileInfo.toModel()
		if err != nil {
			return nil, fmt.Errorf("converting profile info of '%d': %v", profileInfo.Id, err)
		}

		res = append(res, profile)
	}

	return res, nil
}

func (p ProfilesProvider) Search(ctx context.Context, params *models.SearchParams) ([]*models.Profile, error) {
	var res []*models.Profile

	profilesInfo, err := p.getProfilesInfoByParams(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("getting all profiles info: %v", err)
	}

	for _, profileInfo := range profilesInfo {
		profile, err := profileInfo.toModel()
		if err != nil {
			return nil, fmt.Errorf("converting profile info of '%d': %v", profileInfo.Id, err)
		}

		res = append(res, profile)
	}

	return res, nil
}

func (p ProfilesProvider) Get(ctx context.Context, profileID string) (*models.Profile, error) {
	id, err := strconv.ParseInt(profileID, 10, 0)
	if err != nil {
		return nil, fmt.Errorf("illegal id '%s': %v : int64 expected", profileID, err)
	}

	profileInfo, err := p.getProfileInfo(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting profile info of '%d': %v", id, err)
	}

	profile, err := profileInfo.toModel()
	if err != nil {
		return nil, fmt.Errorf("converting profile info of '%d': %v", id, err)
	}

	return profile, nil
}

func (p ProfilesProvider) Add(ctx context.Context, profile *models.Profile) (string, error) {
	query := `
		insert into scl.profiles(user_id, name, surname, birthdate, sex, address, hobbies)
		values (@user, @name, @surname, @birthdate, @sex, @address, @hobbies)
		returning id
	`

	args := pgx.NamedArgs{
		"user":      profile.UserId,
		"name":      profile.Name,
		"surname":   profile.Surname,
		"birthdate": profile.Birthdate,
		"sex":       profile.Sex.String(),
		"address":   profile.Address,
		"hobbies":   profile.Hobbies,
	}

	rows, err := p.querier.Query(ctx, query, args)
	if err != nil {
		return "", fmt.Errorf("inserting into db: %v", err)
	}

	id, err := pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (int64, error) {
		var id int64
		err := row.Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("scanning profile id: %v", err)
		}

		return id, nil
	})
	if err != nil {
		return "", fmt.Errorf("collecting new profile id: %v", err)
	}

	return fmt.Sprintf("%d", id), nil
}

func (p ProfilesProvider) getProfileInfo(ctx context.Context, profileID int64) (*profile, error) {
	var profiles []profile

	query := `
		select
			ps.id,
			name,
			surname,
			birthdate,
			sex,
			address,
			hobbies
		from scl.profiles as ps
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

func (p ProfilesProvider) getAllProfileInfo(ctx context.Context) ([]profile, error) {
	var profiles []profile

	query := `
		select
			ps.id,
			name,
			surname,
			birthdate,
			sex,
			address,
			hobbies
		from scl.profiles as ps
	`

	err := pgxscan.Select(ctx, p.querier, &profiles, query)
	if err != nil {
		return nil, fmt.Errorf("executing query `%s`: %v", query, err)
	}
	return profiles, nil
}

func (p ProfilesProvider) getProfilesInfoByParams(ctx context.Context, params *models.SearchParams) ([]profile, error) {
	var profiles []profile

	args := pgx.NamedArgs{
		"name":    fmt.Sprintf("%s%%", params.NamePrefix),
		"surname": fmt.Sprintf("%s%%", params.SurnamePrefix),
	}

	query := `
		select
			ps.id,
			name,
			surname,
			birthdate,
			sex,
			address,
			hobbies
		from scl.profiles as ps
		where
			name LIKE @name
			and
			surname LIKE @surname
		order by
			ps.id
	`

	err := pgxscan.Select(ctx, p.querier, &profiles, query, args)
	if err != nil {
		return nil, fmt.Errorf("executing query `%s`: %v", query, err)
	}

	if len(profiles) == 0 {
		return nil, fmt.Errorf("querying db: %v", models.ProfileNotFound)
	}

	return profiles, nil
}
