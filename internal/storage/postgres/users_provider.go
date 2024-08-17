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

// TODO : add pagination
func (p UsersProvider) GetAll(ctx context.Context) ([]models.User, error) {
	var res []models.User

	users, err := p.getAllUserInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting all users info: %v", err)
	}

	for _, user := range users {
		res = append(res, models.User{
			Email:    user.Email.String,
			Login:    user.Login.String,
			Password: user.Password.String,
		})
	}

	return res, nil

}

func (p UsersProvider) Get(ctx context.Context, userId int64) (*models.User, error) {
	user, err := p.getUserInfo(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("getting user info of '%d': %v", userId, err)
	}

	return &models.User{
		Email:    user.Email.String,
		Login:    user.Login.String,
		Password: user.Password.String,
	}, nil
}

func (p UsersProvider) Add(ctx context.Context, user *models.User) error {
	query := `
		insert into scl.users(email, login, password)
		values (@email, @login, @password)
	`

	args := pgx.NamedArgs{
		"email":    user.Email,
		"login":    user.Login,
		"password": user.Password,
	}

	_, err := p.querier.Query(ctx, query, args)
	if err != nil {
		return fmt.Errorf("inserting into db: %v", err)
	}

	return nil
}

func (p UsersProvider) getUserInfo(ctx context.Context, userID int64) (*user, error) {
	var users []user

	query := `
		select
			id,
			login,
			password,
			email
		from scl.users
		where id = $1
	`

	err := pgxscan.Select(ctx, p.querier, &users, query, userID)
	if err != nil {
		return nil, fmt.Errorf("executing query `%s`: %v", query, err)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("querying db: %v", models.UserNotFound)
	}

	return &users[0], nil
}

func (p UsersProvider) getAllUserInfo(ctx context.Context) ([]user, error) {
	var users []user

	query := `
		select
			id,
			login,
			password,
			email
		from scl.users
	`

	err := pgxscan.Select(ctx, p.querier, &users, query)
	if err != nil {
		return nil, fmt.Errorf("executing query `%s`: %v", query, err)
	}

	return users, nil
}
