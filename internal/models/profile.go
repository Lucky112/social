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
	Address Address
	Hobbies []Hobby
}

type Address struct {
	Country string
	City    string
}

type Hobby struct {
	Title string
}

var ProfileNotFound = errors.New("profile not found")
