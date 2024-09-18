package profiles

import (
	"fmt"

	"github.com/Lucky112/social/internal/models"
	"github.com/Lucky112/social/internal/models/sex"
)

type profile struct {
	userId  string
	Name    string   `json:"name"    validate:"required"`
	Surname string   `json:"surname"`
	Sex     string   `json:"sex"`
	Age     uint8    `json:"age"`
	City    string   `json:"city"`
	Hobbies []string `json:"hobbies"`
}

type profileResponse struct {
	Id string `json:"id"`
}
type profileError struct {
	Message string `json:"msg"`
}

func (p *profile) toModel() (*models.Profile, error) {
	hobbies := make([]models.Hobby, len(p.Hobbies))
	for i, h := range p.Hobbies {
		hobbies[i] = models.Hobby{
			Title: h,
		}
	}

	sex, err := sex.FromString(p.Sex)
	if err != nil {
		return nil, fmt.Errorf("extracting sex: %v", err)
	}

	return &models.Profile{
		UserId:  p.userId,
		Name:    p.Name,
		Surname: p.Surname,
		Sex:     sex,
		Age:     p.Age,
		Address: models.Address{
			City: p.City,
		},
		Hobbies: hobbies,
	}, nil
}

func fromModel(mp *models.Profile) *profile {
	hobbies := make([]string, len(mp.Hobbies))
	for i, h := range mp.Hobbies {
		hobbies[i] = h.Title
	}

	return &profile{
		Name:    mp.Name,
		Surname: mp.Surname,
		Sex:     mp.Sex.String(),
		Age:     mp.Age,
		City:    mp.Address.City,
		Hobbies: hobbies,
	}
}
