package app

import (
	"context"

	domain "usersvc/internal/domain/user"
)

// Driven adapters (outbound) ports:
type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
}

type UserPublisher interface {
	PublishUserCreated(ctx context.Context, id domain.ID) error
}
