package inmemory

import (
	"context"
	"sync"

	"github.com/Lucky112/social/internal/models"
)

type FriendsStorage struct {
	mu      *sync.RWMutex
	friends map[string][]string
}

func NewFriendsStorage() FriendsStorage {
	return FriendsStorage{
		mu:      &sync.RWMutex{},
		friends: make(map[string][]string),
	}
}

func (s FriendsStorage) Get(ctx context.Context, id string) (*models.Friends, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	friendsList, exists := s.friends[id]
	if !exists {
		return &models.Friends{
			Id: id,
		}, nil
	}

	return &models.Friends{
		Id:        id,
		FriendIds: friendsList,
	}, nil
}

func (s FriendsStorage) Add(ctx context.Context, id, friendId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.friends[id] = append(s.friends[id], friendId)
	s.friends[friendId] = append(s.friends[friendId], id)
	return nil
}

func (s FriendsStorage) Delete(ctx context.Context, id, friendId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	removeFriend := func(friends []string, target string) []string {
		result := []string{}
		for _, f := range friends {
			if f != target {
				result = append(result, f)
			}
		}
		return result
	}

	s.friends[id] = removeFriend(s.friends[id], friendId)
	s.friends[friendId] = removeFriend(s.friends[friendId], id)

	return nil
}
