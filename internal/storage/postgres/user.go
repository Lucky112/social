package postgres

type user struct {
	Id       int64  `db:"id"`
	Login    string `db:"login"`
	Password []byte `db:"password"`
	Email    string `db:"email"`
}
