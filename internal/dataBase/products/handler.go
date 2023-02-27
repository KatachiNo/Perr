package products

import (
	"context"
	"github.com/KatachiNo/Perr/internal/handlers"
	"github.com/KatachiNo/Perr/pkg/logg"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"io"
	"net/http"
)

const (
	productsAll               = "/db/all"
	productsAdd               = "/db/add"
	productsChangeProductItem = "/db/changeProductItem"
	productsDeleteItem        = "/db/delete"
	productsPriceStory        = "/db/PriceStory"

	testHey = "/test"
)

type handler struct {
	storage Storage
	l       *logg.Logger
}

func NewRegister(storage Storage, l *logg.Logger) handlers.Handler {
	return &handler{
		storage: storage,
		l:       l,
	}
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc(testHey, hey).Methods("GET")

	//router.HandleFunc(productsAll, getProductsAll).Methods("GET")
	router.HandleFunc(productsAdd, h.productAdd).Methods("GET")
	//router.HandleFunc(productsChangeProductItem, hey).Methods("PATCH")
	//router.HandleFunc(productsDeleteItem, hey).Methods("DELETE")
	//router.HandleFunc(productsPriceStory, hey).Methods("GET")

}

func hey(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hey")
}

func getProductsAll(writer http.ResponseWriter, request *http.Request) {

}

func (h *handler) productAdd(w http.ResponseWriter, r *http.Request) {

	pp := Products{
		Id:              33,
		ProductName:     "слива",
		Category:        1,
		QuantityOfGoods: 3,
		LastPrice:       decimal.Decimal{},
		AvailableStatus: "good",
		PictureAddress:  "/123",
	}
	err := h.storage.ProductAddItem(context.TODO(), pp)
	if err != nil {
		w.WriteHeader(400)
	}
	io.WriteString(w, "hey???")

}
