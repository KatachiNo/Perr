package orders

type Orders struct {
	OrderId            int    `json:"order-id"`
	UserId             int    `json:"user-id"`
	DataOfOrder        string `json:"data-of-order"`
	OrderedProductsIds string `json:"ordered-products-ids"`
	FinalPrice         string `json:"final-price"`
	DeliveredStatus    string `json:"delivered-status"`
}
