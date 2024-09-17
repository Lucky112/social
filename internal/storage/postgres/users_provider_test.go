package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/Lucky112/social/internal/models"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestUserExists(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	p := UsersProvider{mock}

	t.Run("Email exists", func(t *testing.T) {
		user := &models.User{
			Email:    "myemail@index.com",
			Login:    "mylogin",
			Password: []byte("pwd"),
		}

		exists := mock.NewRows([]string{"exists"}).
			AddRow(true)

		mock.ExpectQuery("select").WithArgs(user.Email).WillReturnRows(exists)

		actual, err := p.Exists(context.Background(), user)
		require.NoError(t, err)
		require.True(t, actual)
	})

	t.Run("Login exists", func(t *testing.T) {
		user := &models.User{
			Email:    "myemail@index.com",
			Login:    "mylogin",
			Password: []byte("pwd"),
		}

		notExists := mock.NewRows([]string{"exists"}).
			AddRow(false)
		exists := mock.NewRows([]string{"exists"}).
			AddRow(true)

		mock.ExpectQuery("select").WithArgs(user.Email).WillReturnRows(notExists)
		mock.ExpectQuery("select").WithArgs(user.Login).WillReturnRows(exists)

		actual, err := p.Exists(context.Background(), user)
		require.NoError(t, err)
		require.True(t, actual)
	})

	t.Run("Nothing exists", func(t *testing.T) {
		user := &models.User{
			Email:    "myemail@index.com",
			Login:    "mylogin",
			Password: []byte("pwd"),
		}

		notExists := mock.NewRows([]string{"exists"}).
			AddRow(false)

		mock.ExpectQuery("select").WithArgs(user.Email).WillReturnRows(notExists)
		mock.ExpectQuery("select").WithArgs(user.Login).WillReturnRows(notExists)

		actual, err := p.Exists(context.Background(), user)
		require.NoError(t, err)
		require.False(t, actual)
	})

	t.Run("exists with error", func(t *testing.T) {
		mock.ExpectQuery("select").WithArgs("").WillReturnError(errors.New("db error"))

		actual, err := p.Exists(context.Background(), &models.User{})
		require.Error(t, err)
		require.False(t, actual)
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
			Password: []byte("pwd"),
		}

		users := mock.NewRows([]string{"id", "email", "login", "password"}).
			AddRow(int64(1), "myemail@index.com", "mylogin", []byte("pwd"))

		mock.ExpectQuery("select").WithArgs("login").WillReturnRows(users)

		actual, err := p.Get(context.Background(), "login")
		require.NoError(t, err)
		require.Equal(t, expected, *actual)
	})

	t.Run("select nothing found", func(t *testing.T) {
		rows := mock.NewRows([]string{"id", "email", "login", "password"})

		mock.ExpectQuery("select").WithArgs("login").WillReturnRows(rows)

		actual, err := p.Get(context.Background(), "login")
		require.Error(t, err)
		require.Nil(t, actual)
	})

	t.Run("select with error", func(t *testing.T) {
		mock.ExpectQuery("select").WithArgs("login").WillReturnError(errors.New("db error"))

		actual, err := p.Get(context.Background(), "login")
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
			Password: []byte("pwd"),
		}

		rows := mock.NewRows([]string{"id"}).AddRow(int64(1))

		mock.ExpectQuery("insert").WithArgs(user.Email, user.Login, user.Password).WillReturnRows(rows)

		id, err := p.Add(context.Background(), &user)
		require.NoError(t, err)
		require.Equal(t, "1", id)
	})

	t.Run("insert with error", func(t *testing.T) {
		user := models.User{
			Email:    "myemail@index.com",
			Login:    "mylogin",
			Password: []byte("pwd"),
		}

		mock.ExpectQuery("insert").WithArgs(user.Email, user.Login, user.Password).WillReturnError(errors.New("db error"))

		id, err := p.Add(context.Background(), &user)
		require.Error(t, err)
		require.Equal(t, "", id)
	})

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
