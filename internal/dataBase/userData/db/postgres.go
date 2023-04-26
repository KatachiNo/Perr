package UserDataDb

import (
	"context"
	"fmt"
	"github.com/KatachiNo/Perr/internal/dataBase/userData"
	"github.com/KatachiNo/Perr/pkg/client/postgresql"
	"github.com/KatachiNo/Perr/pkg/logg"
)

type db struct {
	client postgresql.Client
	l      *logg.Logger
}

func (d db) UserDataAdd(ctx context.Context, arrUserData []userData.UserData) error {
	for _, uD := range arrUserData {
		q := fmt.Sprintf(`INSERT INTO "UserData" (email,phone_number,country,city,index,street,number_house,note,first_name,middle_name,last_name)
VALUES ('%s','%s','%s','%s','%s','%s','%s','%s','%s','%s','%s',)`, uD.Email, uD.PhoneNumber, uD.Country, uD.City, uD.Index, uD.Street, uD.NumberHouse, uD.Note, uD.FirstName, uD.MiddleName, uD.LastName)
		fmt.Println(q)
		_, err := d.client.ExecContext(ctx, q)

		if err != nil {
			fmt.Print("ошибка")
			fmt.Print(err)
			return err
		}
	}
	return nil
}

func (d db) UserDataFindOne(ctx context.Context, id string) (userData.UserData, error) {
	q := fmt.Sprintf(`SELECT id,email,phone_number,country,city,index,street,number_house,note,
       first_name,middle_name,last_name FROM "UserData" WHERE id=%s`, id)

	row := d.client.QueryRowContext(ctx, q)
	uD := userData.UserData{}
	err := row.Scan(&uD.Id, &uD.Email, &uD.PhoneNumber, &uD.Country, &uD.City, &uD.Index, &uD.Street, &uD.NumberHouse,
		&uD.Note, &uD.FirstName, &uD.MiddleName, &uD.LastName)

	if err != nil {
		return uD, err
	}

	return uD, nil
}

func (d db) UserDataGetAll(ctx context.Context) ([]userData.UserData, error) {
	q := `SELECT id,email,phone_number,country,city,index,street,number_house,note,first_name,middle_name,last_name FROM "UserData"`
	rows, err := d.client.QueryContext(ctx, q)

	if err != nil {
		return nil, err
	}

	arrUD := make([]userData.UserData, 0)

	for rows.Next() {
		uD := userData.UserData{}

		err = rows.Scan(&uD.Id, &uD.Email, &uD.PhoneNumber, &uD.Country, &uD.City, &uD.Index, &uD.Street, &uD.NumberHouse, &uD.Note, &uD.FirstName, &uD.MiddleName, &uD.LastName)

		if err != nil {
			return nil, err
		}

		arrUD = append(arrUD, uD)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return arrUD, nil
}

func (d db) UserDataDelete(ctx context.Context, id string) error {
	q := fmt.Sprintf(`DELETE FROM "UserData" WHERE id=%d`, id)
	_, err := d.client.ExecContext(ctx, q)

	return err
}

func NewStorage(client postgresql.Client, l *logg.Logger) userData.Storage {
	return &db{
		client: client,
		l:      l,
	}
}
