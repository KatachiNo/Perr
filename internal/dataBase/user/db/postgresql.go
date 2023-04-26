package userDb

import (
	"context"
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"github.com/KatachiNo/Perr/internal/dataBase/user"
	"github.com/KatachiNo/Perr/pkg/client/postgresql"
	"github.com/KatachiNo/Perr/pkg/logg"
	"time"
)

type db struct {
	client postgresql.Client
	l      *logg.Logger
}

func (d db) UserCreate(ctx context.Context, data user.User) error {
	//TODO implement me
	hash := sha512.New()
	hash.Write([]byte(data.Password))

	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}

	h := fmt.Sprintf("%x", hash.Sum(salt))
	s := fmt.Sprintf("%x", salt)

	date := time.Now().Format("2006-01-02 15:04:05.000000")
	q := fmt.Sprintf(`INSERT INTO "Users" (login, "passwordHash", "categoryOfUser", "dateOfRegistration", salt, algorithm) VALUES ('%s','%s','%s','%s','%s','%s')`,
		data.Login, h, "0", date, s, "sha512")

	fmt.Println(q)
	_, err = d.client.ExecContext(ctx, q)

	if err != nil {
		fmt.Print("ошибка")
		fmt.Print(err)
		return err
	}

	return nil
}

func (d db) UserFind(ctx context.Context, id string) (user.User, error) {

	q := fmt.Sprintf(`Select id, login, "passwordHash", "categoryOfUser", "dateOfRegistration", salt, algorithm from "Users" WHERE id=%d`, id)

	row := d.client.QueryRowContext(ctx, q)
	u := user.User{}
	err := row.Scan(&u.Id, &u.Login, &u.PasswordHash, &u.CategoryOfUser, &u.DateOfRegistration, &u.Salt, &u.Algorithm)

	if err != nil {
		return u, err
	}

	return u, nil

}

func (d db) UserUpdate(ctx context.Context, data user.User) error {
	q := fmt.Sprintf(`Update "Users" SET "salt"='%s',"algorithm"='%s',"passwordHash"='%s',"categoryOfUser"='%s'
               where id=1`, data.Salt, data.Algorithm, data.PasswordHash, data.CategoryOfUser)

	fmt.Println(q)
	_, err := d.client.ExecContext(ctx, q)

	if err != nil {
		fmt.Print("ошибка")
		fmt.Print(err)
		return err
	}

	return nil
}

func (d db) UserDelete(ctx context.Context, id int) error {
	q := fmt.Sprintf(`DELETE FROM "Users" WHERE id=%d`, id)
	_, err := d.client.ExecContext(ctx, q)

	return err
}

func NewStorage(client postgresql.Client, l *logg.Logger) user.Storage {
	return &db{
		client: client,
		l:      l,
	}
}
