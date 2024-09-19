package profiles

import (
	"fmt"
	"time"

	"github.com/Lucky112/social/internal/models"
	"github.com/Lucky112/social/internal/models/sex"
)

const birthdateFormat = "2006-01-02"

type profile struct {
	userId    string
	Name      string `json:"name"      validate:"required"`
	Surname   string `json:"surname"`
	Sex       string `json:"sex"`
	Birthdate string `json:"birthdate" validate:"datetime=2006-01-02"`
	City      string `json:"city"`
	Hobbies   string `json:"hobbies"`
}

type profileResponse struct {
	Id string `json:"id"`
}
type profileError struct {
	Message string `json:"msg"`
}

func (p *profile) toModel() (*models.Profile, error) {
	sex, err := sex.FromString(p.Sex)
	if err != nil {
		return nil, fmt.Errorf("extracting sex: %v", err)
	}

	birthdate, err := time.Parse(birthdateFormat, p.Birthdate)
	if err != nil {
		return nil, fmt.Errorf("extracting birthdate: %v", err)
	}

	return &models.Profile{
		UserId:    p.userId,
		Name:      p.Name,
		Surname:   p.Surname,
		Sex:       sex,
		Birthdate: birthdate,
		Address:   p.City,
		Hobbies:   p.Hobbies,
	}, nil
}

func fromModel(mp *models.Profile) *profile {
	return &profile{
		Name:      mp.Name,
		Surname:   mp.Surname,
		Sex:       mp.Sex.String(),
		Birthdate: mp.Birthdate.Format(birthdateFormat),
		City:      mp.Address,
		Hobbies:   mp.Hobbies,
	}
}
