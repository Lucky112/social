package postgres

import (
	"fmt"

	"github.com/guregu/null/v5"

	"github.com/Lucky112/social/internal/models"
	"github.com/Lucky112/social/internal/models/sex"
)

type profile struct {
	Id      int64       `db:"id"`
	Name    null.String `db:"name"`
	Surname null.String `db:"surname"`
	Sex     null.String `db:"sex"`
	Age     null.Int16  `db:"age"`
	Address null.String `db:"address"`
	Hobbies null.String `db:"hobbies"`
}

func (p *profile) toModel() (*models.Profile, error) {
	sex, err := sex.FromString(p.Sex.String)
	if err != nil {
		return nil, fmt.Errorf("parsing sex: %v", err)
	}

	return &models.Profile{
		Name:    p.Name.String,
		Surname: p.Surname.String,
		Age:     uint8(p.Age.Int16),
		Sex:     sex,
		Address: p.Address.String,
		Hobbies: p.Hobbies.String,
	}, nil
}
