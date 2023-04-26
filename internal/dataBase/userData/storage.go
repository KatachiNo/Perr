package userData

import "context"

type Storage interface {
	UserDataAdd(ctx context.Context, arrUserData []UserData) error
	UserDataGetAll(ctx context.Context) ([]UserData, error)

	UserDataFindOne(ctx context.Context, id string) (UserData, error)
	UserDataDelete(ctx context.Context, id string) error
}
