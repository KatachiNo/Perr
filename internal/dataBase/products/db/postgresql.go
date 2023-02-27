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

	_, error := r.client.ExecContext(ctx, q, p.ProductName, p.Category, p.QuantityOfGoods, p.LastPrice, p.AvailableStatus)

	return error

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
	panic("implement me")
}

func NewStorage(client postgresql.Client, l *logg.Logger) products.Storage {
	return &db{
		client: client,
		l:      l,
	}
}
