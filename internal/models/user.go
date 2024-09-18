package models

import "errors"

// Структура данных с информацией о пользователе
type User struct {
	Id             string
	Email          string
	Login          string
	Password       string
	HashedPassword []byte
}

var UserNotFound = errors.New("user not found")
var UserAlreadyExists = errors.New("user already exists")
var UserBadCredentials = errors.New("invalid credentials for user")
