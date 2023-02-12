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
	loginforgetPassword = "/login/forgetPassword"

	productsAll               = "/products/all"
	productsAdd               = "/products/add"
	productsChangeProductItem = "/products/changeProductItem"
	productsDeleteItem        = "/products/delete"
	productsPriceStory        = "/products/PriceStory"

	testHey = "/test"
)

type handler struct {
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc(testHey, hey).Methods("GET")
}

func NewHandler() handlers.Handler {
	return &handler{}
}

func hey(w http.ResponseWriter, r *http.Request) {

	io.WriteString(w, "hey")

}
