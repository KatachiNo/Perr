package productsDb

import (
	"context"
	"github.com/KatachiNo/Perr/internal/dataBase/products"
	"github.com/KatachiNo/Perr/pkg/client/postgresql"
	"github.com/KatachiNo/Perr/pkg/logg"
)

type db struct {
	client postgresql.Client
	l      *logg.Logger
}

func (r db) ProductAddItem(ctx context.Context, p products.Products) error {
	q := `INSERT INTO "Products" 
    (product_name, category, quantity_of_goods, last_price, available_status, picture_address)
		VALUES ($1,$2,$3,$4,$5,$6)`

	_, err := r.client.ExecContext(ctx, q, p.ProductName, p.Category, p.QuantityOfGoods, p.LastPrice, p.AvailableStatus, p.PictureAddress)

	return err
}

func (r db) ProductsUpdateItem(ctx context.Context, p products.Products) error {
	//TODO implement me
	panic("implement me")
}

func (r db) ProductDeleteItem(ctx context.Context, p products.Products) error {
	//TODO implement me
	panic("implement me")
}

func (r db) ProductFind(ctx context.Context, id string) ([]products.Products, error) {
	//TODO implement me
	panic("implement me")
}

func (r db) ProductsGetAll(ctx context.Context) ([]products.Products, error) {
	q := `SELECT id,product_name,category,quantity_of_goods,last_price,available_status,picture_address FROM "Products"`
	rows, err := r.client.QueryContext(ctx, q)

	if err != nil {
		return nil, err
	}

	arrProd := make([]products.Products, 0)

	for rows.Next() {
		pr := products.Products{}

		err = rows.Scan(&pr.Id, &pr.ProductName, &pr.Category,
			&pr.QuantityOfGoods, &pr.LastPrice, &pr.AvailableStatus,
			&pr.PictureAddress)

		if err != nil {
			return nil, err
		}

		arrProd = append(arrProd, pr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return arrProd, nil
}

func NewStorage(client postgresql.Client, l *logg.Logger) products.Storage {
	return &db{
		client: client,
		l:      l,
	}
}
