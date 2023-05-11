package orders

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KatachiNo/Perr/internal/handlers"
	"github.com/KatachiNo/Perr/internal/tokens"
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
	go router.Handle(orderTAll, tokens.CheckAuthorizedAdmin(h.orderTAll)).Methods("GET")
	router.Handle(orderTCreate, tokens.CheckAuthorizedUser(h.orderTCreate)).Methods("POST")

	//router.Handle(catTChange, tokens.CheckAuthorizedAdmin(h.catTChange)).Methods("PATCH")
	router.Handle(orderTDelete, tokens.CheckAuthorizedAdmin(h.orderTDelete)).Methods("DELETE")

	go router.Handle(orderTFindOne, tokens.CheckAuthorizedAdmin(h.orderTFindOne)).Methods("GET")
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
		h.l.Error(errConv)
	}

	errDel := h.storage.OrderDelete(context.TODO(), intId)
	if errDel != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.l.Error(errDel)
	}
}

func (h *handler) orderTFindOne(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	intId, errConv := strconv.Atoi(id)
	if errConv != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.l.Error(errConv)
	}

	cT, errFind := h.storage.OrderFindOne(context.TODO(), intId)
	if errFind != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.l.Error(errFind)
	}

	js, errJs := json.Marshal(cT)
	if errJs != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.l.Error(errJs)
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
		h.l.Error(err)
	}
}
