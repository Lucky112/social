package postgres

import "github.com/guregu/null/v5"

type profile struct {
	Id      int64       `db:"id"`
	Name    null.String `db:"name"`
	Surname null.String `db:"surname"`
	Sex     null.String `db:"sex"`
	Age     null.Int16  `db:"age"`
	City    null.String `db:"city"`
	Country null.String `db:"country"`
}

type hobby struct {
	Id        int64       `db:"id"`
	ProfileID int64       `db:"profile_id"`
	Title     null.String `db:"title"`
}
