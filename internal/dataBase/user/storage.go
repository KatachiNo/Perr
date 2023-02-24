package user

import "context"

type Storage interface {
	UserCreate(ctx context.Context, logins User) (string, error)
	UserFind(ctx context.Context, id string) (User, error)
	UserUpdate(ctx context.Context, logins User) error
	UserDelete(ctx context.Context, id string) error
}
