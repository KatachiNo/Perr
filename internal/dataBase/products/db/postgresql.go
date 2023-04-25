package productsDb

import (
	"context"
	"fmt"
	"github.com/KatachiNo/Perr/internal/dataBase/products"
	"github.com/KatachiNo/Perr/pkg/client/postgresql"
	"github.com/KatachiNo/Perr/pkg/logg"
)

type db struct {
	client postgresql.Client
	l      *logg.Logger
}

func (r db) ProductsAddItems(ctx context.Context, arrP []products.Products) error {
	/*q := `INSERT INTO "Products"
	    (product_name, category, quantity_of_goods, last_price, available_status, picture_address)
			VALUES ($1,$2,$3,$4,$5,$6)`
	*/

	fmt.Println(arrP)
	for _, p := range arrP {
		fmt.Println("in for")
		fmt.Println(p)

		q := fmt.Sprintf(`INSERT INTO "Products" (product_name, category, quantity_of_goods, last_price, available_status, picture_address) VALUES ('%s','%d','%d','%s','%s','%s')`, p.ProductName, p.Category,
			p.QuantityOfGoods, p.LastPrice.String(), p.AvailableStatus, p.PictureAddress)

		fmt.Println(q)
		//_, err := r.client.ExecContext(ctx, q, p.ProductName, p.Category,
		//	p.QuantityOfGoods, p.LastPrice, p.AvailableStatus, p.PictureAddress)

		_, err := r.client.ExecContext(ctx, q)

		if err != nil {
			fmt.Print("ошибка")
			fmt.Print(err)
			return err
		}
	}
	return nil
	//return nil
}

func (r db) ProductUpdateItem(ctx context.Context, p products.Products) error {
	//TODO implement me
	panic("implement me")
}

func (r db) ProductFindOne(ctx context.Context, id int) (products.Products, error) {

	q := fmt.Sprintf(`SELECT * FROM "Products" WHERE id=%d`, id)

	row := r.client.QueryRowContext(ctx, q)
	pr := products.Products{}
	err := row.Scan(&pr.Id, &pr.ProductName, &pr.Category,
		&pr.QuantityOfGoods, &pr.LastPrice, &pr.AvailableStatus,
		&pr.PictureAddress)

	if err != nil {
		return pr, err
	}

	return pr, nil

}

func (r db) ProductDeleteItem(ctx context.Context, id int) error {
	//TODO implement me

	q := fmt.Sprintf(`DELETE FROM "Products" WHERE id=%d`, id)
	_, err := r.client.ExecContext(ctx, q)

	return err

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
