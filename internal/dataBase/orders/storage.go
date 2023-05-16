package orders

import "context"

type Storage interface {
	CreateOrder(ctx context.Context, order Orders) error
	ChangeOrder(ctx context.Context, order Orders) error
	OrderDelete(ctx context.Context, orderId int) error

	OrderFindOne(ctx context.Context, id int) (Orders, error)
	OrdersGetAll(ctx context.Context) ([]Orders, error)
}
