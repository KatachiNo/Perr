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
	for _, p := range arrP {
		q := fmt.Sprintf(`INSERT INTO "Products" (product_name, category, quantity_of_goods, last_price, available_status, picture_address) VALUES ('%s','%d','%d','%s','%s','%s')`, p.ProductName, p.Category,
			p.QuantityOfGoods, p.LastPrice, p.AvailableStatus, p.PictureAddress)
		fmt.Println(q)
		_, err := r.client.ExecContext(ctx, q)

		if err != nil {
			fmt.Print("ошибка")
			fmt.Print(err)
			return err
		}
	}
	return nil
}

func (r db) ProductUpdateItem(ctx context.Context, p products.Products) error {

	fmt.Println(p)
	var pN, lP, aS, pA string
	c := -1
	qOG := -1
	var cc, qqOG string

	pN = `'` + p.ProductName + `'`
	c = p.Category
	qOG = p.QuantityOfGoods
	lP = `'` + p.LastPrice + `'`
	aS = `'` + p.AvailableStatus + `'`
	pA = `'` + p.PictureAddress + `'`

	if p.ProductName == "null" {
		pN = "product_name"
	}
	if p.Category == -1 {
		cc = "category"
	}
	if p.QuantityOfGoods == -1 {
		qqOG = "quantity_of_goods"
	}
	if p.LastPrice == "null" {
		lP = "last_price"
	}
	if p.AvailableStatus == "null" {
		aS = "available_status"
	}
	if p.PictureAddress == "null" {
		pA = "picture_address"
	}

	var q string
	if c != -1 && qOG != -1 {
		q = fmt.Sprintf(`UPDATE "Products" SET product_name=%s, category=%d, quantity_of_goods=%d, last_price=%s, available_status=%s, picture_address=%s WHERE id=%d`, pN, c, qOG, lP, aS, pA, p.Id)
	}
	if c == -1 && qOG != -1 {
		q = fmt.Sprintf(`UPDATE "Products" SET product_name=%s, category=%s, quantity_of_goods=%d, last_price=%s, available_status=%s, picture_address=%s WHERE id=%d`, pN, cc, qOG, lP, aS, pA, p.Id)

	}
	if c != -1 && qOG == -1 {
		q = fmt.Sprintf(`UPDATE "Products" SET product_name=%s, category=%d, quantity_of_goods=%s, last_price=%s, available_status=%s, picture_address=%s WHERE id=%d`, pN, c, qqOG, lP, aS, pA, p.Id)
	}
	if c == -1 && qOG == -1 {
		q = fmt.Sprintf(`UPDATE "Products" SET product_name=%s, category=%s, quantity_of_goods=%s, last_price=%s, available_status=%s, picture_address=%s WHERE id=%d`, pN, cc, qqOG, lP, aS, pA, p.Id)

	}

	fmt.Println(q)
	_, err := r.client.ExecContext(ctx, q)

	if err != nil {
		fmt.Print("ошибка")
		fmt.Print(err)
		return err
	}

	return nil
}

func (r db) ProductFindOne(ctx context.Context, id int) (products.Products, error) {
	q := fmt.Sprintf(`SELECT id,product_name,category,quantity_of_goods,last_price,available_status,picture_address FROM "Products" WHERE id=%d`, id)

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
