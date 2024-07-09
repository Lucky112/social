package models

import "errors"

// Структура данных с информацией о пользователе
type User struct {
	Email    string
	Name     string
	Password string
}

var UserNotFound = errors.New("user not found")
