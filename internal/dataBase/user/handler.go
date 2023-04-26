package user

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KatachiNo/Perr/internal/handlers"
	"github.com/KatachiNo/Perr/pkg/logg"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	userAuth   = "/users/authorization"
	userReg    = "/users/registration"
	userChange = "/users/changeData"
	userDelete = "/users/delete"
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
	router.HandleFunc(userReg, h.userReg).Methods("POST")
	router.HandleFunc(userAuth, h.userAuth).Methods("POST")
	router.HandleFunc(userChange, h.userChange).Methods("GET")
	router.HandleFunc(userDelete, h.userDelete).Methods("GET")

}

func (h *handler) userReg(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	fmt.Println("reg")
	fmt.Println(user)
	err := h.storage.UserCreate(context.TODO(), user)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *handler) userAuth(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) userChange(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) userDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	intId, errConv := strconv.Atoi(id)
	if errConv != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Print(errConv)
	}

	errDel := h.storage.UserDelete(context.TODO(), intId)
	if errDel != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Print(errDel)
	}
}
