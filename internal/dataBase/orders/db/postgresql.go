package ordersDb

import (
	"context"
	"fmt"
	"github.com/KatachiNo/Perr/internal/dataBase/orders"
	"github.com/KatachiNo/Perr/pkg/client/postgresql"
	"github.com/KatachiNo/Perr/pkg/logg"
	"strings"
)

type db struct {
	client postgresql.Client
	l      *logg.Logger
}

func (d db) CreateOrder(ctx context.Context, order orders.Orders) error {

	orderedIds := fmt.Sprintf("{%s}", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(order.OrderedProductsIds)), ","), "[]"))

	q := fmt.Sprintf(`INSERT INTO "Orders" (user_id, data_of_order, ordered_products_ids, final_price, delivery_status)
VALUES ('%d','%s','%s','%s','%s')`, order.UserId, order.DataOfOrder, orderedIds,
		order.FinalPrice, "Заказ создан")
	fmt.Println(q)
	_, err := d.client.ExecContext(ctx, q)

	if err != nil {
		fmt.Print("ошибка")
		fmt.Print(err)
		return err
	}

	return nil
}

func (d db) ChangeOrder(ctx context.Context, order orders.Orders) error {
	//TODO implement me
	panic("implement me")
}

func (d db) OrderDelete(ctx context.Context, orderId int) error {
	q := fmt.Sprintf(`DELETE FROM "Orders" WHERE id=%d`, orderId)
	_, err := d.client.ExecContext(ctx, q)

	return err
}

func (d db) OrderFindOne(ctx context.Context, id int) (orders.Orders, error) {
	q := fmt.Sprintf(`SELECT "orderId", user_id, data_of_order, ordered_products_ids, final_price, delivery_status FROM "Orders" WHERE id=%d`, id)

	row := d.client.QueryRowContext(ctx, q)
	ord := orders.Orders{}
	err := row.Scan(&ord.OrderId, &ord.UserId, &ord.DataOfOrder, &ord.OrderedProductsIds, &ord.FinalPrice, &ord.DeliveredStatus)

	if err != nil {
		return ord, err
	}
	return ord, nil
}

func (d db) OrdersGetAll(ctx context.Context) ([]orders.Orders, error) {
	q := `SELECT "orderId", user_id, data_of_order, ordered_products_ids, final_price, delivery_status FROM "Orders"`
	rows, err := d.client.QueryContext(ctx, q)

	if err != nil {
		return nil, err
	}

	arrOrders := make([]orders.Orders, 0)

	for rows.Next() {
		ord := orders.Orders{}

		err = rows.Scan(&ord.OrderId, &ord.UserId, &ord.DataOfOrder, &ord.OrderedProductsIds, &ord.FinalPrice, &ord.DeliveredStatus)

		if err != nil {
			return nil, err
		}

		arrOrders = append(arrOrders, ord)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return arrOrders, nil
}

func NewStorage(client postgresql.Client, l *logg.Logger) orders.Storage {
	return &db{
		client: client,
		l:      l,
	}
}
