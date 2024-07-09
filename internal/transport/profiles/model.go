package profiles

import (
	"fmt"
	"strings"

	"github.com/Lucky112/social/internal/models"
)

type profile struct {
	Name    string   `json:"name"`
	Surname string   `json:"surname"`
	Sex     string   `json:"sex"`
	Age     uint8    `json:"age"`
	City    string   `json:"city"`
	Hobbies []string `json:"hobbies"`
}

type profileResponse struct {
	Id string `json:"id"`
}

func (p *profile) toModel() (*models.Profile, error) {
	hobbies := make([]models.Hobby, len(p.Hobbies))
	for i, h := range p.Hobbies {
		hobbies[i] = models.Hobby{
			Title: h,
		}
	}

	var sex models.Sex
	switch strings.ToLower(p.Sex) {
	case "":
		sex = models.Unknown
	case "male":
		sex = models.Male
	case "female":
		sex = models.Female
	default:
		return nil, fmt.Errorf("unknown sex '%s': only %v are available", p.Sex, []string{"Male", "Female"})
	}

	return &models.Profile{
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

	var sex string
	switch mp.Sex {
	case models.Male:
		sex = "Male"
	case models.Female:
		sex = "Female"
	}

	return &profile{
		Name:    mp.Name,
		Surname: mp.Surname,
		Sex:     sex,
		Age:     mp.Age,
		City:    mp.Address.City,
		Hobbies: hobbies,
	}
}
