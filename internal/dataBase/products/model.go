package products

type Products struct {
	Id              int    `json:"id"`
	ProductName     string `json:"product-name"`
	Category        int    `json:"category"`
	QuantityOfGoods int    `json:"quantity-of-goods"`
	LastPrice       string `json:"last-price"`
	AvailableStatus string `json:"available-status"`
	PictureAddress  string `json:"picture-address"`
}
