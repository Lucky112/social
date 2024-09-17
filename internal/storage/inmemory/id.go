package inmemory

import "github.com/google/uuid"

func generateId() string {
	return uuid.NewString()
}
