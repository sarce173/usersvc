package user_test

import (
	"context"
	"errors"
	"testing"

	app "usersvc/internal/app"
	appuser "usersvc/internal/app/user"
	domain "usersvc/internal/domain/user"
)

type fakeRepo struct{ fail bool }
func (f *fakeRepo) Create(ctx context.Context, u domain.User) error {
	if f.fail { return errors.New("persist fail") }
	return nil
}

type fakePub struct{ called bool }
func (f *fakePub) PublishUserCreated(ctx context.Context, id domain.ID) error { f.called = true; return nil }

var _ app.UserRepository = (*fakeRepo)(nil)
var _ app.UserPublisher = (*fakePub)(nil)

func TestCreateUser_OK(t *testing.T) {
	repo := &fakeRepo{}
	pub := &fakePub{}
	uc := appuser.NewCreateUserUseCase(repo, pub)

	res, err := uc.Execute(context.Background(), appuser.CreateUserCommand{Name: "Alice", Email: "alice@example.com"})
	if err != nil { t.Fatalf("unexpected err: %v", err) }
	if res.ID == "" { t.Fatalf("expected id") }
	if !pub.called { t.Fatalf("expected publish to be called") }
}

func TestCreateUser_InvalidEmail(t *testing.T) {
	repo := &fakeRepo{}
	pub := &fakePub{}
	uc := appuser.NewCreateUserUseCase(repo, pub)

	_, err := uc.Execute(context.Background(), appuser.CreateUserCommand{Name: "Alice", Email: "bad"})
	if err == nil { t.Fatalf("expected error") }
}
