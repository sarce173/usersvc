package user

import (
	"context"
	"errors"

	app "usersvc/internal/app"
	domain "usersvc/internal/domain/user"
)

var ErrPersist = errors.New("could not persist user")

type CreateUserCommand struct {
	Name  string
	Email string
}

type CreateUserResult struct {
	ID string
}

type CreateUserUseCase struct {
	Repo app.UserRepository
	Pub  app.UserPublisher
}

func NewCreateUserUseCase(repo app.UserRepository, pub app.UserPublisher) *CreateUserUseCase {
	return &CreateUserUseCase{Repo: repo, Pub: pub}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, cmd CreateUserCommand) (CreateUserResult, error) {
	u, err := domain.New(cmd.Name, cmd.Email)
	if err != nil {
		return CreateUserResult{}, err
	}
	if err := uc.Repo.Create(ctx, u); err != nil {
		return CreateUserResult{}, ErrPersist
	}
	// Best-effort publish; don't fail the command if publish errors.
	_ = uc.Pub.PublishUserCreated(ctx, u.ID)
	return CreateUserResult{ID: string(u.ID)}, nil
}
