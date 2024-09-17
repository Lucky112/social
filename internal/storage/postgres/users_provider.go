package postgres

import (
	"context"
	"fmt"

	"github.com/Lucky112/social/internal/models"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type UsersProvider struct {
	querier pgxscan.Querier
}

func (p UsersProvider) Exists(ctx context.Context, user *models.User) (bool, error) {
	emailExists, err := p.checkEmailExists(ctx, user.Email)
	if err != nil {
		return false, fmt.Errorf("checking email exists: %v", err)
	}
	if emailExists {
		return true, nil
	}

	loginExists, err := p.checkLoginExists(ctx, user.Login)
	if err != nil {
		return false, fmt.Errorf("checking login exists: %v", err)
	}
	if loginExists {
		return true, nil
	}

	return false, nil
}

func (p UsersProvider) Get(ctx context.Context, login string) (*models.User, error) {
	user, err := p.getUserInfo(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("getting user info of '%s': %v", login, err)
	}

	return &models.User{
		Email:    user.Email.String,
		Login:    user.Login.String,
		Password: user.Password.String,
	}, nil
}

func (p UsersProvider) Add(ctx context.Context, user *models.User) (string, error) {
	query := `
		insert into scl.users(email, login, password)
		values (@email, @login, @password)
		returning id
	`

	args := pgx.NamedArgs{
		"email":    user.Email,
		"login":    user.Login,
		"password": user.Password,
	}

	rows, err := p.querier.Query(ctx, query, args)
	if err != nil {
		return "", fmt.Errorf("inserting into db: %v", err)
	}

	id, err := pgx.CollectExactlyOneRow(rows, func(row pgx.CollectableRow) (int64, error) {
		var id int64
		err := row.Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("scanning user id: %v", err)
		}

		return id, nil
	})
	if err != nil {
		return "", fmt.Errorf("collecting new user id: %v", err)
	}

	return fmt.Sprintf("%d", id), nil
}

func (p UsersProvider) getUserInfo(ctx context.Context, login string) (*user, error) {
	var users []user

	query := `
		select
			id,
			login,
			password,
			email
		from scl.users
		where login = $1
	`

	err := pgxscan.Select(ctx, p.querier, &users, query, login)
	if err != nil {
		return nil, fmt.Errorf("executing query `%s`: %v", query, err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("querying db: %v", models.UserNotFound)
	}

	return &users[0], nil
}

func (p UsersProvider) checkEmailExists(ctx context.Context, email string) (bool, error) {
	var res []bool

	query := `
		select
			exists(
				select
					1
				from scl.users
				where email = $1
				limit 1
			) as exists
	`

	err := pgxscan.Select(ctx, p.querier, &res, query, email)
	if err != nil {
		return false, fmt.Errorf("executing query `%s`: %v", query, err)
	}

	if len(res) == 0 {
		return false, nil
	}

	return res[0], nil
}

func (p UsersProvider) checkLoginExists(ctx context.Context, login string) (bool, error) {
	var res []bool

	query := `
		select
			exists(
				select
					1
				from scl.users
				where login = $1
				limit 1
			) as exists
	`

	err := pgxscan.Select(ctx, p.querier, &res, query, login)
	if err != nil {
		return false, fmt.Errorf("executing query `%s`: %v", query, err)
	}

	if len(res) == 0 {
		return false, nil
	}

	return res[0], nil
}
