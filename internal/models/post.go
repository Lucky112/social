package models

import (
	"errors"
	"time"
)

type Post struct {
	id        string
	UserId    string
	CreatedAt time.Time
	Content   string
}

var PostNotFound = errors.New("post not found")
