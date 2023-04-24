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

	_, err := r.client.ExecContext(ctx, q, p.ProductName, p.Category, p.QuantityOfGoods, p.LastPrice, p.AvailableStatus)

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
	//TODO implement me
	q := `SELECT (id,product_name,category,quantity_of_goods,available_status,picture_address) FROM "Products"`
	rows, err := r.client.QueryContext(ctx, q)

	var aP []products.Products

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := products.Products{}
		rows.Scan(&p.Id, &p.ProductName, &p.Category, &p.QuantityOfGoods,
			&p.LastPrice, &p.AvailableStatus, &p.PictureAddress)
		aP = append(aP, p)
	}

	return aP, nil
}

func NewStorage(client postgresql.Client, l *logg.Logger) products.Storage {
	return &db{
		client: client,
		l:      l,
	}
}
