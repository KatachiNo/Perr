package db

import (
	"context"
	"github.com/KatachiNo/Perr/internal/dataBase/user"
	"github.com/KatachiNo/Perr/pkg/logg"
)

type db struct {
	l *logg.Logger
}

func (d db) UserCreate(ctx context.Context, logins user.User) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (d db) UserFind(ctx context.Context, id string) (user.User, error) {
	//TODO implement me
	panic("implement me")
}

func (d db) UserUpdate(ctx context.Context, logins user.User) error {
	//TODO implement me
	panic("implement me")
}

func (d db) UserDelete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func NewStorage(l *logg.Logger) user.Storage {
	return &db{}
}
