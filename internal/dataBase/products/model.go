package products

import "github.com/shopspring/decimal"

type Products struct {
	Id              int             `json:"id"`
	Category        int             `json:"category"`
	QuantityOfGoods int             `json:"quantity-of-goods"`
	LastPrice       decimal.Decimal `json:"last-price"`
	PictureAddress  string          `json:"picture-address"`
}
