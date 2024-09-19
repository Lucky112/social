package models

import (
	"errors"

	"github.com/Lucky112/social/internal/models/sex"
)

type Profile struct {
	id      string
	UserId  string
	Name    string
	Surname string
	Sex     sex.Sex
	Age     uint8
	Address string
	Hobbies string
}

var ProfileNotFound = errors.New("profile not found")
