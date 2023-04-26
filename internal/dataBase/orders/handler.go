package orders

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KatachiNo/Perr/internal/Tokens"
	"github.com/KatachiNo/Perr/internal/handlers"
	"github.com/KatachiNo/Perr/pkg/logg"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const (
	orderTAll    = "/orderT/all"
	orderTCreate = "/orderT/create"
	//orderTChange = "/orderT/changeProductItem"
	orderTDelete = "/orderT/delete"

	orderTFindOne = "/orderT/FindOne"
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
	router.Handle(orderTAll, Tokens.CheckAuthorizedAdmin(h.orderTAll)).Methods("GET")
	router.Handle(orderTCreate, Tokens.CheckAuthorizedUser(h.orderTCreate)).Methods("POST")

	//router.Handle(catTChange, Tokens.CheckAuthorizedAdmin(h.catTChange)).Methods("PATCH")
	router.Handle(orderTDelete, Tokens.CheckAuthorizedAdmin(h.orderTDelete)).Methods("DELETE")

	router.Handle(orderTFindOne, Tokens.CheckAuthorizedAdmin(h.orderTFindOne)).Methods("GET")
}

func (h *handler) orderTAll(w http.ResponseWriter, r *http.Request) {
	arrO, err := h.storage.OrdersGetAll(context.TODO())

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	js, errJs := json.Marshal(arrO)
	if errJs != nil {
		fmt.Println(errJs)

		w.WriteHeader(http.StatusBadRequest)
	}
	if err == nil && errJs == nil {

		w.WriteHeader(http.StatusOK)
		w.Write(js)
	}
}

func (h *handler) orderTDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	intId, errConv := strconv.Atoi(id)
	if errConv != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Print(errConv)
	}

	errDel := h.storage.OrderDelete(context.TODO(), intId)
	if errDel != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Print(errDel)
	}
}

func (h *handler) orderTFindOne(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	intId, errConv := strconv.Atoi(id)
	if errConv != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Print(errConv)
	}

	cT, errFind := h.storage.OrderFindOne(context.TODO(), intId)
	if errFind != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Print(errFind)
	}

	js, errJs := json.Marshal(cT)
	if errJs != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func (h *handler) orderTCreate(w http.ResponseWriter, r *http.Request) {
	var order Orders
	json.NewDecoder(r.Body).Decode(&order)

	err := h.storage.CreateOrder(context.TODO(), order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}
