package postgres

import (
	"context"
	"fmt"

	"github.com/Lucky112/social/internal/models"
	"github.com/Lucky112/social/internal/models/sex"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type ProfilesProvider struct {
	querier pgxscan.Querier
}

// TODO : add pagination
func (p ProfilesProvider) GetAll(ctx context.Context) ([]models.Profile, error) {
	var res []models.Profile

	profilesInfo, err := p.getAllProfileInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting all profiles info: %v", err)
	}

	hobbies, err := p.getAllHobbies(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting all hobbies: %v", err)
	}

	for _, profileInfo := range profilesInfo {
		sex, err := sex.FromString(profileInfo.Sex.String)
		if err != nil {
			return nil, fmt.Errorf("parsing sex: %v", err)
		}

		profile := models.Profile{
			Name:    profileInfo.Name.String,
			Surname: profileInfo.Surname.String,
			Age:     uint8(profileInfo.Age.Int16),
			Sex:     sex,
			Address: models.Address{
				Country: profileInfo.Country.String,
				City:    profileInfo.City.String,
			},
		}

		for _, h := range hobbies[profileInfo.Id] {
			profile.Hobbies = append(profile.Hobbies, models.Hobby{
				Title: h.Title.String,
			})
		}

		res = append(res, profile)
	}

	return res, nil
}

func (p ProfilesProvider) Get(ctx context.Context, profileID int64) (*models.Profile, error) {
	profileInfo, err := p.getProfileInfo(ctx, profileID)
	if err != nil {
		return nil, fmt.Errorf("getting profile info of '%d': %v", profileID, err)
	}

	hobbies, err := p.getHobbies(ctx, profileID)
	if err != nil {
		return nil, fmt.Errorf("getting hobbies of '%d': %v", profileID, err)
	}

	sex, err := sex.FromString(profileInfo.Sex.String)
	if err != nil {
		return nil, fmt.Errorf("parsing sex: %v", err)
	}

	profile := &models.Profile{
		Name:    profileInfo.Name.String,
		Surname: profileInfo.Surname.String,
		Age:     uint8(profileInfo.Age.Int16),
		Sex:     sex,
		Address: models.Address{
			Country: profileInfo.Country.String,
			City:    profileInfo.City.String,
		},
	}

	for _, h := range hobbies {
		profile.Hobbies = append(profile.Hobbies, models.Hobby{
			Title: h.Title.String,
		})
	}

	return profile, nil
}

func (p ProfilesProvider) Add(ctx context.Context, profile *models.Profile) (string, error) {
	query := `
		insert into scl.profiles(user_id, name, surname, age, sex)
		values (@user, @name, @surname, @age, @sex)
		returning id
	`

	args := pgx.NamedArgs{
		"user":    profile.UserId,
		"name":    profile.Name,
		"surname": profile.Surname,
		"age":     profile.Age,
		"sex":     profile.Sex.String(),
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
			age,
			sex,
			city,
			country
		from scl.profiles as ps
		left join scl.addresses as ads
			on (ps.address_id = ads.id)
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

func (p ProfilesProvider) getHobbies(ctx context.Context, profileID int64) ([]hobby, error) {
	var hobbies []hobby

	query := `
		select
			h.id,
			title
		from scl.profiles as ps
		join scl.profile_hobbies as phs
			on (ps.id = phs.profile_id)
		join scl.hobbies as hs
			on (hs.id = phs.hobby_id)
		where ps.id = $1
	`

	err := pgxscan.Select(ctx, p.querier, &hobbies, query, profileID)
	if err != nil {
		return nil, fmt.Errorf("executing query `%s`: %v", query, err)
	}

	return hobbies, nil
}

func (p ProfilesProvider) getAllProfileInfo(ctx context.Context) ([]profile, error) {
	var profiles []profile

	query := `
		select
			ps.id,
			name,
			surname,
			age,
			sex,
			city,
			country
		from scl.profiles as ps
		left join scl.addresses as ads
			on (ps.address_id = ads.id)
	`

	err := pgxscan.Select(ctx, p.querier, &profiles, query)
	if err != nil {
		return nil, fmt.Errorf("executing query `%s`: %v", query, err)
	}
	return profiles, nil
}

func (p ProfilesProvider) getAllHobbies(ctx context.Context) (map[int64][]hobby, error) {
	var hobbies []hobby

	query := `
	select
			h.id,
			phs.profile_id,
			title
		from scl.profile_hobbies as phs
		join scl.hobbies as hs
			on (hs.id = phs.hobby_id)
	`

	err := pgxscan.Select(ctx, p.querier, &hobbies, query)
	if err != nil {
		return nil, fmt.Errorf("executing query `%s`: %v", query, err)
	}

	res := make(map[int64][]hobby)
	for _, h := range hobbies {
		res[h.ProfileID] = append(res[h.ProfileID], h)
	}

	return res, nil
}
