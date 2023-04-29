package categoryTableDb

import (
	"context"
	"fmt"
	"github.com/KatachiNo/Perr/internal/dataBase/categoryTable"
	"github.com/KatachiNo/Perr/pkg/client/postgresql"
	"github.com/KatachiNo/Perr/pkg/logg"
)

type db struct {
	client postgresql.Client
	l      *logg.Logger
}

func (d db) CategoryTableAddItems(ctx context.Context, arrCT []categoryTable.CategoryTable) error {
	for _, c := range arrCT {
		q := fmt.Sprintf(`INSERT INTO "CategoryTable" (categoryid, categoryname) VALUES ('%d','%s')`, c.CategoryId, c.CategoryName)
		fmt.Println(q)
		_, err := d.client.ExecContext(ctx, q)

		if err != nil {
			d.l.Error(err)
			d.l.Error(q)
			return err
		}
	}
	return nil
}

func (d db) CategoryTableUpdateItem(ctx context.Context, cT categoryTable.CategoryTable) error {
	q := fmt.Sprintf(`UPDATE "CategoryTable" SET categoryname='%s' WHERE id=%d`, cT.CategoryName, cT.Id)

	_, err := d.client.ExecContext(ctx, q)

	if err != nil {
		d.l.Error(err)
		d.l.Error(q)
		return err
	}

	return nil
}

func (d db) CategoryTableDeleteItem(ctx context.Context, id int) error {
	q := fmt.Sprintf(`DELETE FROM "CategoryTable" WHERE id=%d`, id)
	_, err := d.client.ExecContext(ctx, q)

	return err
}

func (d db) CategoryTableFindOne(ctx context.Context, id int) (categoryTable.CategoryTable, error) {
	q := fmt.Sprintf(`SELECT id,categoryid,categoryname FROM "CategoryTable" WHERE id=%d`, id)

	row := d.client.QueryRowContext(ctx, q)
	cT := categoryTable.CategoryTable{}
	err := row.Scan(&cT.Id, &cT.CategoryId, &cT.CategoryName)

	if err != nil {
		return cT, err
	}

	return cT, nil
}

func (d db) CategoryTableGetAll(ctx context.Context) ([]categoryTable.CategoryTable, error) {
	q := `SELECT id,categoryid,categoryname FROM "CategoryTable"`
	rows, err := d.client.QueryContext(ctx, q)

	if err != nil {
		d.l.Error(err)
		d.l.Error(q)
		return nil, err
	}

	arrCT := make([]categoryTable.CategoryTable, 0)

	for rows.Next() {
		cT := categoryTable.CategoryTable{}

		err = rows.Scan(&cT.Id, &cT.CategoryId, &cT.CategoryName)

		if err != nil {
			return nil, err
		}

		arrCT = append(arrCT, cT)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return arrCT, nil
}

func NewStorage(client postgresql.Client, l *logg.Logger) categoryTable.Storage {
	return &db{
		client: client,
		l:      l,
	}
}
