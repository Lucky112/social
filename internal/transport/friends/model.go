package friends

import "github.com/Lucky112/social/internal/models"

type friendsResponse struct {
	ID      string   `json:"id"`
	Friends []string `json:"friend_ids"`
}

type friendError struct {
	Message string `json:"msg"`
}

func fromModel(mf *models.Friends) *friendsResponse {
	return &friendsResponse{
		ID:      mf.Id,
		Friends: mf.FriendIds,
	}
}
