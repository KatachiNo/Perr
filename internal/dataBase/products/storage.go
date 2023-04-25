package products

import "context"

type Storage interface {
	ProductAddItem(ctx context.Context, p Products) error
	ProductsUpdateItem(ctx context.Context, p Products) error
	ProductDeleteItem(ctx context.Context, id int) error

	ProductFindOne(ctx context.Context, id int) (Products, error)
	ProductsGetAll(ctx context.Context) ([]Products, error)
}
