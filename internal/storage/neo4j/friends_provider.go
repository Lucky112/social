package neo4j

import (
	"context"
	"fmt"

	"github.com/Lucky112/social/internal/models"
	"github.com/Lucky112/social/pkg/neo4j"
)

type FriendsProvider struct {
	driver *neo4j.Driver
}

func NewFriendsProvider(driver *neo4j.Driver) *FriendsProvider {
	return &FriendsProvider{driver: driver}
}

func (s *FriendsProvider) Get(ctx context.Context, id string) (*models.Friends, error) {
	session := s.driver.NewReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.Run(ctx, `
		MATCH (u:User {id: $id})-[:FRIEND]->(f:User)
		RETURN f.id`, map[string]interface{}{"id": id})
	if err != nil {
		return nil, fmt.Errorf("failed to get friends: %w", err)
	}

	friends := []string{}
	for result.Next(ctx) {
		if id, ok := result.Record().Get("f.id"); ok {
			friends = append(friends, id.(string))
		}
	}

	return &models.Friends{
		Id:        id,
		FriendIds: friends,
	}, nil
}

func (s *FriendsProvider) Add(ctx context.Context, id, friendId string) error {
	session := s.driver.NewWriteSession(ctx)
	defer session.Close(ctx)

	_, err := session.Run(ctx, `
		MERGE (u:User {id: $id})
		MERGE (f:User {id: $friendId})
		MERGE (u)-[:FRIEND]->(f)
		MERGE (f)-[:FRIEND]->(u)`,
		map[string]interface{}{"id": id, "friendId": friendId})
	if err != nil {
		return fmt.Errorf("failed to add friend: %w", err)
	}
	return nil
}

func (s *FriendsProvider) Delete(ctx context.Context, id, friendId string) error {
	session := s.driver.NewWriteSession(ctx)
	defer session.Close(ctx)

	_, err := session.Run(ctx, `
		MATCH (u:User {id: $id})-[r:FRIEND]-(f:User {id: $friendId})
		DELETE r`, map[string]interface{}{"id": id, "friendId": friendId})
	if err != nil {
		return fmt.Errorf("failed to delete friend: %w", err)
	}
	return nil
}
