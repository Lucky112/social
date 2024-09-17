package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/Lucky112/social/internal/models"
	"github.com/Lucky112/social/internal/models/sex"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestAllProfiles(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	p := ProfilesProvider{mock}

	t.Run("Select successfully", func(t *testing.T) {
		expected := []models.Profile{
			{
				Name:    "user1",
				Surname: "surname1",
				Sex:     sex.Male,
				Age:     18,
				Address: models.Address{
					City:    "city1",
					Country: "country1",
				},
				Hobbies: []models.Hobby{
					{Title: "reading"},
					{Title: "dancing"},
				},
			},
			{
				Name:    "user2",
				Surname: "surname2",
				Sex:     sex.Female,
				Age:     21,
				Address: models.Address{
					City:    "city2",
					Country: "country2",
				},
				Hobbies: []models.Hobby{
					{Title: "youtube"},
				},
			},
		}

		profiles := mock.NewRows([]string{"id", "name", "surname", "age", "sex", "city", "country"}).
			AddRow(int64(1), "user1", "surname1", 18, "male", "city1", "country1").
			AddRow(int64(2), "user2", "surname2", 21, "female", "city2", "country2")

		hobbies := mock.NewRows([]string{"id", "profile_id", "title"}).
			AddRow(int64(100), int64(1), "reading").
			AddRow(int64(101), int64(1), "dancing").
			AddRow(int64(102), int64(2), "youtube")

		mock.ExpectQuery("select").WithArgs().WillReturnRows(profiles)
		mock.ExpectQuery("select").WithArgs().WillReturnRows(hobbies)

		actual, err := p.GetAll(context.Background())
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("select with error", func(t *testing.T) {
		mock.ExpectQuery("select").WillReturnError(errors.New("db error"))

		actual, err := p.GetAll(context.Background())
		require.Error(t, err)
		require.Nil(t, actual)
	})

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestSingleProfile(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	p := ProfilesProvider{mock}

	t.Run("Select successfully", func(t *testing.T) {
		expected := models.Profile{
			Name:    "user1",
			Surname: "surname1",
			Sex:     sex.Male,
			Age:     18,
			Address: models.Address{
				City:    "city1",
				Country: "country1",
			},
			Hobbies: []models.Hobby{
				{Title: "reading"},
				{Title: "dancing"},
			},
		}

		profile := mock.NewRows([]string{"name", "surname", "age", "sex", "country", "city"}).
			AddRow(
				expected.Name,
				expected.Surname,
				expected.Age,
				expected.Sex.String(),
				expected.Address.Country,
				expected.Address.City,
			)

		hobbies := mock.NewRows([]string{"id", "title"}).
			AddRow(int64(100), "reading").
			AddRow(int64(101), "dancing")

		profileId := int64(0)
		mock.ExpectQuery("select").WithArgs(profileId).WillReturnRows(profile)
		mock.ExpectQuery("select").WithArgs(profileId).WillReturnRows(hobbies)

		actual, err := p.Get(context.Background(), 0)
		require.NoError(t, err)
		require.Equal(t, expected, *actual)
	})

	t.Run("select nothing found", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "name", "surname", "age", "sex"})

		id := int64(1)
		mock.ExpectQuery("select").WithArgs(id).WillReturnRows(rows)

		actual, err := p.Get(context.Background(), id)
		require.Error(t, err)
		require.Nil(t, actual)
	})

	t.Run("select with error", func(t *testing.T) {
		id := int64(1)
		mock.ExpectQuery("select").WithArgs(id).WillReturnError(errors.New("db error"))

		actual, err := p.Get(context.Background(), id)
		require.Error(t, err)
		require.Nil(t, actual)
	})

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestInsertProfile(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	p := ProfilesProvider{mock}

	t.Run("Insert successfully", func(t *testing.T) {
		prof := models.Profile{
			Name:    "user1",
			Surname: "surname1",
			Sex:     sex.Male,
			Age:     18,
		}

		rows := mock.NewRows([]string{"id"}).AddRow(int64(1))

		mock.ExpectQuery("insert").WithArgs(prof.Name, prof.Surname, prof.Age, prof.Sex.String()).WillReturnRows(rows)

		id, err := p.Add(context.Background(), &prof)
		require.NoError(t, err)
		require.Equal(t, "1", id)
	})

	t.Run("insert with error", func(t *testing.T) {
		prof := models.Profile{
			Name:    "user1",
			Surname: "surname1",
			Sex:     sex.Male,
			Age:     18,
		}
		mock.ExpectQuery("insert").WithArgs(prof.Name, prof.Surname, prof.Age, prof.Sex.String()).WillReturnError(errors.New("db error"))

		id, err := p.Add(context.Background(), &prof)
		require.Error(t, err)
		require.Equal(t, "", id)
	})

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
