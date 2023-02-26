package products

import (
	"context"
	"github.com/KatachiNo/Perr/internal/dataBase/products"
	"github.com/KatachiNo/Perr/pkg/client/postgresql"
	"github.com/KatachiNo/Perr/pkg/logg"
)

type repo struct {
	client postgresql.Client
	l      *logg.Logger
}

func (r repo) ProductAddItem(ctx context.Context, p products.Products) error {
	q := `INSERT INTO "Products" 
    (product_name, category, quantity_of_goods, last_price, available_status, picture_address)
		VALUES ($1,$2,$3,$4,$5,$6)`

	_, error := r.client.ExecContext(ctx, q, p.ProductName, p.Category, p.QuantityOfGoods, p.LastPrice, p.AvailableStatus)

	return error
}

func (r repo) ProductsUpdateItem(ctx context.Context, p products.Products) error {
	//TODO implement me
	panic("implement me")
}

func (r repo) ProductDeleteItem(ctx context.Context, p products.Products) error {
	//TODO implement me
	panic("implement me")
}

func (r repo) ProductFind(ctx context.Context, id string) (products.Products, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) ProductsGetAll(ctx context.Context) (products.Products, error) {
	// q := `SELECT * FROM "Products"`

	//rows, err := r.client.QueryContext(ctx, q)

	/// scan + array Products
	//TODO implement me
	panic("implement me")
}

func NewRepo(client postgresql.Client, l *logg.Logger) products.Storage {
	return &repo{
		client: client,
		l:      l,
	}
}
