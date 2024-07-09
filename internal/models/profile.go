package models

import "errors"

type Sex uint8

const (
	Unknown Sex = iota
	Male
	Female
)

type Profile struct {
	id      string
	Name    string
	Surname string
	Sex     Sex
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
