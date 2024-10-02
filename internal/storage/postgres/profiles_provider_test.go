package postgres

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

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
		birthdate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
		expected := []*models.Profile{
			{
				Name:      "user1",
				Surname:   "surname1",
				Sex:       sex.Male,
				Birthdate: birthdate,
				Address:   "Moscow",
				Hobbies:   "reading, dancing",
			},
			{
				Name:      "user2",
				Surname:   "surname2",
				Sex:       sex.Female,
				Birthdate: birthdate,
				Address:   "Los-Angeles",
				Hobbies:   "youtube",
			},
		}

		profiles := mock.NewRows([]string{"id", "name", "surname", "birthdate", "sex", "address", "hobbies"}).
			AddRow(int64(1), "user1", "surname1", birthdate, "male", "Moscow", "reading, dancing").
			AddRow(int64(2), "user2", "surname2", birthdate, "female", "Los-Angeles", "youtube")

		mock.ExpectQuery("select").WithArgs().WillReturnRows(profiles)

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
			Name:      "user1",
			Surname:   "surname1",
			Sex:       sex.Male,
			Birthdate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			Address:   "Moscow",
			Hobbies:   "reading, dancing",
		}

		profile := mock.NewRows([]string{"name", "surname", "birthdate", "sex", "address", "hobbies"}).
			AddRow(
				expected.Name,
				expected.Surname,
				expected.Birthdate,
				expected.Sex.String(),
				expected.Address,
				expected.Hobbies,
			)

		profileId := int64(0)
		mock.ExpectQuery("select").WithArgs(profileId).WillReturnRows(profile)

		actual, err := p.Get(context.Background(), fmt.Sprintf("%d", profileId))
		require.NoError(t, err)
		require.Equal(t, expected, *actual)
	})

	t.Run("select nothing found", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "name", "surname", "birthdate", "sex", "address", "hobbies"})

		id := int64(1)
		mock.ExpectQuery("select").WithArgs(id).WillReturnRows(rows)

		actual, err := p.Get(context.Background(), fmt.Sprintf("%d", id))
		require.Error(t, err)
		require.Nil(t, actual)
	})

	t.Run("select with error", func(t *testing.T) {
		id := int64(1)
		mock.ExpectQuery("select").WithArgs(id).WillReturnError(errors.New("db error"))

		actual, err := p.Get(context.Background(), fmt.Sprintf("%d", id))
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
			UserId:    "1",
			Name:      "user1",
			Surname:   "surname1",
			Sex:       sex.Male,
			Birthdate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			Address:   "Moscow",
			Hobbies:   "reading, dancing",
		}

		rows := mock.NewRows([]string{"id"}).AddRow(int64(1))

		mock.ExpectQuery("insert").WithArgs(
			prof.UserId,
			prof.Name,
			prof.Surname,
			prof.Birthdate,
			prof.Sex.String(),
			prof.Address,
			prof.Hobbies,
		).WillReturnRows(rows)

		id, err := p.Add(context.Background(), &prof)
		require.NoError(t, err)
		require.Equal(t, "1", id)
	})

	t.Run("insert with error", func(t *testing.T) {
		prof := models.Profile{
			UserId:    "1",
			Name:      "user1",
			Surname:   "surname1",
			Sex:       sex.Male,
			Birthdate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
			Address:   "Moscow",
			Hobbies:   "reading, dancing",
		}

		mock.ExpectQuery("insert").WithArgs(
			prof.UserId,
			prof.Name,
			prof.Surname,
			prof.Birthdate,
			prof.Sex.String(),
			prof.Address,
			prof.Hobbies,
		).WillReturnError(errors.New("db error"))

		id, err := p.Add(context.Background(), &prof)
		require.Error(t, err)
		require.Equal(t, "", id)
	})

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
