package productPriceStoryDb

import (
	"context"
	"fmt"
	"github.com/KatachiNo/Perr/internal/dataBase/productPriceStory"
	"github.com/KatachiNo/Perr/pkg/client/postgresql"
	"github.com/KatachiNo/Perr/pkg/logg"
)

type db struct {
	client postgresql.Client
	l      *logg.Logger
}

func (d db) ProductPriceAddItems(ctx context.Context, arrPPS []productPriceStory.ProductPriceStoryTable) error {
	for _, ppS := range arrPPS {
		q := fmt.Sprintf(`INSERT INTO "ProductPriceStory" ("Price", "Date") VALUES ('%s','%s')`, ppS.Price, ppS.Date)
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

func (d db) ProductPriceTableDeleteItem(ctx context.Context, id int) error {
	q := fmt.Sprintf(`DELETE FROM "ProductPriceStory" WHERE id=%d`, id)
	_, err := d.client.ExecContext(ctx, q)

	return err
}

func (d db) ProductPriceTableFindOne(ctx context.Context, id int) (productPriceStory.ProductPriceStoryTable, error) {
	q := fmt.Sprintf(`SELECT id, "Price","Date" FROM "ProductPriceStory" WHERE id=%d`, id)

	row := d.client.QueryRowContext(ctx, q)
	ppS := productPriceStory.ProductPriceStoryTable{}
	err := row.Scan(&ppS.Id, &ppS.Price, &ppS.Date)

	if err != nil {
		return ppS, err
	}
	return ppS, nil
}

func (d db) ProductPriceTableGetAll(ctx context.Context) ([]productPriceStory.ProductPriceStoryTable, error) {
	q := `SELECT id,"Price","Date" FROM "ProductPriceStory"`
	rows, err := d.client.QueryContext(ctx, q)

	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	arrPPS := make([]productPriceStory.ProductPriceStoryTable, 0)

	for rows.Next() {
		cT := productPriceStory.ProductPriceStoryTable{}

		err = rows.Scan(&cT.Id, &cT.Price, &cT.Date)

		if err != nil {
			return nil, err
		}

		arrPPS = append(arrPPS, cT)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return arrPPS, nil
}

func NewStorage(client postgresql.Client, l *logg.Logger) productPriceStory.Storage {
	return &db{
		client: client,
		l:      l,
	}
}
