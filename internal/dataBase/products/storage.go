package products

import "context"

type Storage interface {
	ProductAddItem(ctx context.Context, p Products) error
	ProductsUpdateItem(ctx context.Context, p Products) error
	ProductDeleteItem(ctx context.Context, p Products) error

	ProductFind(ctx context.Context, id string) (Products, error)
	ProductsGetAll(ctx context.Context) (Products, error)
}
