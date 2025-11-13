package stdout

import (
	"context"
	"fmt"

	domain "usersvc/internal/domain/user"
)

type Publisher struct{}

func New() *Publisher { return &Publisher{} }

func (p *Publisher) PublishUserCreated(ctx context.Context, id domain.ID) error {
	fmt.Printf("event=user_created user_id=%s\n", id)
	return nil
}
