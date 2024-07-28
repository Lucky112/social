package postgres

import "github.com/guregu/null/v5"

type Profile struct {
	Id      int64       `db:"id"`
	Name    null.String `db:"name"`
	Surname null.String `db:"surname"`
	Sex     null.String `db:"sex"`
	Age     null.Int16  `db:"age"`
}
