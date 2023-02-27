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

	testHey = "/test"
)

type handler struct {
}

func NewRegister() handlers.Handler {
	return &handler{}
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc(testHey, hey).Methods("GET")

}

func hey(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hey")
}
