package products

import "github.com/shopspring/decimal"

type Products struct {
	Id              int             `json:"id"`
	ProductName     string          `json:"product-name"`
	Category        int             `json:"category"`
	QuantityOfGoods int             `json:"quantity-of-goods"`
	LastPrice       decimal.Decimal `json:"last-price"`
	AvailableStatus string          `json:"available-status"`
	PictureAddress  string          `json:"picture-address"`
}
