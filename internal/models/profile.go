package models

import (
	"errors"
	"time"

	"github.com/Lucky112/social/internal/models/sex"
)

type Profile struct {
	id        string
	UserId    string
	Name      string
	Surname   string
	Sex       sex.Sex
	Birthdate time.Time
	Address   string
	Hobbies   string
}

type SearchParams struct {
	NamePrefix    string
	SurnamePrefix string
}

var ProfileNotFound = errors.New("profile not found")
