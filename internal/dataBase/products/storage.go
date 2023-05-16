package products

import "context"

type Storage interface {
	ProductsAddItems(ctx context.Context, p []Products) error
	ProductUpdateItem(ctx context.Context, p Products) error
	ProductDeleteItem(ctx context.Context, id int) error

	ProductFindOne(ctx context.Context, id int) (Products, error)
	ProductsGetAll(ctx context.Context) ([]Products, error)
}
