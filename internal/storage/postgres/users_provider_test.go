package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/Lucky112/social/internal/models"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestAllUsers(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	p := UsersProvider{mock}

	t.Run("Select successfully", func(t *testing.T) {
		expected := []models.User{
			{
				Email:    "myemail@index.com",
				Login:    "mylogin",
				Password: "pwd",
			},
			{
				Email:    "newemail@index.com",
				Login:    "newlogin",
				Password: "pwd2",
			},
		}

		users := mock.NewRows([]string{"id", "email", "login", "password"}).
			AddRow(int64(1), "myemail@index.com", "mylogin", "pwd").
			AddRow(int64(2), "newemail@index.com", "newlogin", "pwd2")

		mock.ExpectQuery("select").WithArgs().WillReturnRows(users)

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

func TestSingleUser(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	p := UsersProvider{mock}

	t.Run("Select successfully", func(t *testing.T) {
		expected := models.User{
			Email:    "myemail@index.com",
			Login:    "mylogin",
			Password: "pwd",
		}

		users := mock.NewRows([]string{"id", "email", "login", "password"}).
			AddRow(int64(1), "myemail@index.com", "mylogin", "pwd")

		profileId := int64(0)
		mock.ExpectQuery("select").WithArgs(profileId).WillReturnRows(users)

		actual, err := p.Get(context.Background(), profileId)
		require.NoError(t, err)
		require.Equal(t, expected, *actual)
	})

	t.Run("select nothing found", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "email", "login", "password"})

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

func TestInsertUser(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	p := UsersProvider{mock}

	t.Run("Insert successfully", func(t *testing.T) {
		user := models.User{
			Email:    "myemail@index.com",
			Login:    "mylogin",
			Password: "pwd",
		}

		rows := mock.NewRows([]string{})

		mock.ExpectQuery("insert").WithArgs(user.Email, user.Login, user.Password).WillReturnRows(rows)

		err := p.Add(context.Background(), &user)
		require.NoError(t, err)
	})

	t.Run("insert with error", func(t *testing.T) {
		user := models.User{
			Email:    "myemail@index.com",
			Login:    "mylogin",
			Password: "pwd",
		}

		mock.ExpectQuery("insert").WithArgs(user.Email, user.Login, user.Password).WillReturnError(errors.New("db error"))

		err := p.Add(context.Background(), &user)
		require.Error(t, err)
	})

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
