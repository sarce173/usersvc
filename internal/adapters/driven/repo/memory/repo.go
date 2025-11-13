package memory

import (
	"context"
	"sync"

	domain "usersvc/internal/domain/user"
)

type Repo struct {
	mu   sync.Mutex
	byID map[domain.ID]domain.User
	byEmail map[string]domain.ID
}

func New() *Repo {
	return &Repo{
		byID: make(map[domain.ID]domain.User),
		byEmail: make(map[string]domain.ID),
	}
}

func (r *Repo) Create(ctx context.Context, u domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.byEmail[u.Email]; exists {
		// For demo simplicity, treat duplicate as error.
		return ErrDuplicateEmail
	}
	r.byID[u.ID] = u
	r.byEmail[u.Email] = u.ID
	return nil
}

// Local error type (adapter-level); the use case collapses persistence errors.
type duplicateEmail string
func (e duplicateEmail) Error() string { return "duplicate email" }
var ErrDuplicateEmail = duplicateEmail("duplicate email")
