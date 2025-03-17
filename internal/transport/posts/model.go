package posts

import (
	"github.com/Lucky112/social/internal/models"
)

const birthdateFormat = "2006-01-02"

type post struct {
	userId  string
	Content string `json:"content"      validate:"required"`
}

type postResponse struct {
	Id string `json:"id"`
}
type postError struct {
	Message string `json:"msg"`
}

func (p *post) toModel() (*models.Post, error) {
	return &models.Post{
		UserId:  p.userId,
		Content: p.Content,
	}, nil
}

func fromModel(mp *models.Post) *post {
	return &post{
		mp.UserId,
		mp.Content,
	}
}
