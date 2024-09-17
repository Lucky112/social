package models

import "errors"

// Структура данных с информацией о пользователе
type User struct {
	Id       string
	Email    string
	Login    string
	Password []byte
}

var UserNotFound = errors.New("user not found")
