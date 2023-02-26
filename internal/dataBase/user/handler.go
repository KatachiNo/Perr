package user

import (
	"io"
	"net/http"

	"github.com/KatachiNo/Perr/internal/handlers"

	"github.com/gorilla/mux"
)

const (
	loginAuthorization  = "/login/authorization"
	loginRegistration   = "/login/registration"
	loginChangeData     = "/login/changeData"
	loginForgetPassword = "/login/forgetPassword"

	productsAll               = "/db/all"
	productsAdd               = "/db/add"
	productsChangeProductItem = "/db/changeProductItem"
	productsDeleteItem        = "/db/delete"
	productsPriceStory        = "/db/PriceStory"

	testHey = "/test"
)

type handler struct {
}

func New() handlers.Handler {
	return &handler{}
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc(testHey, hey).Methods("GET")

	router.HandleFunc(productsAll, getProductsAll).Methods("GET")
	router.HandleFunc(productsAdd, prodAdd).Methods("POST")
	router.HandleFunc(productsChangeProductItem, hey).Methods("PATCH")
	router.HandleFunc(productsDeleteItem, hey).Methods("DELETE")
	router.HandleFunc(productsPriceStory, hey).Methods("GET")

}

func hey(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hey")
}
func getProductsAll(writer http.ResponseWriter, request *http.Request) {

}

func prodAdd(w http.ResponseWriter, r *http.Request) {

}
