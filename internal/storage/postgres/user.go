package postgres

import "github.com/guregu/null/v5"

type user struct {
	Id       int64       `db:"id"`
	Login    null.String `db:"login"`
	Password null.String `db:"password"`
	Email    null.String `db:"email"`
}
