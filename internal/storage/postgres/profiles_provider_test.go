package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/guregu/null/v5"
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
		expected := []Profile{
			{
				Id:      1,
				Name:    null.StringFrom("user1"),
				Surname: null.StringFrom("surname1"),
				Sex:     null.StringFrom("male"),
				Age:     null.Int16From(18),
			},
			{
				Id:      2,
				Name:    null.StringFrom("user2"),
				Surname: null.StringFrom("surname2"),
				Sex:     null.StringFrom("female"),
				Age:     null.Int16From(21),
			},
		}

		rows := mock.NewRows([]string{"id", "name", "surname", "age", "sex"}).
			AddRow(int64(1), "user1", "surname1", 18, "male").
			AddRow(int64(2), "user2", "surname2", 21, "female")

		mock.ExpectQuery("select").WithArgs().WillReturnRows(rows)

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

func TestSingleProfiles(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	p := ProfilesProvider{mock}

	t.Run("Select successfully", func(t *testing.T) {
		expected := Profile{
			Id:      1,
			Name:    null.StringFrom("user1"),
			Surname: null.StringFrom("surname1"),
			Sex:     null.StringFrom("male"),
			Age:     null.Int16From(18),
		}

		rows := mock.NewRows([]string{"id", "name", "surname", "age", "sex"}).
			AddRow(
				expected.Id,
				expected.Name.String,
				expected.Surname.String,
				expected.Age.Int16,
				expected.Sex.String,
			)

		mock.ExpectQuery("select").WithArgs(expected.Id).WillReturnRows(rows)

		actual, err := p.Get(context.Background(), expected.Id)
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
