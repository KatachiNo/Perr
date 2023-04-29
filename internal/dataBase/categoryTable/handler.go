package categoryTable

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
	catTAll    = "/catT/all"
	catTAdd    = "/catT/add"
	catTChange = "/catT/changeProductItem"
	catTDelete = "/catT/delete"

	catTFindOne = "/catT/FindOne"
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
	router.Handle(catTAll, Tokens.CheckAuthorizedUser(h.catTAll)).Methods("GET")
	router.Handle(catTAdd, Tokens.CheckAuthorizedAdmin(h.catTAdd)).Methods("POST")

	router.Handle(catTChange, Tokens.CheckAuthorizedAdmin(h.catTChange)).Methods("PATCH")
	router.Handle(catTDelete, Tokens.CheckAuthorizedAdmin(h.catTDelete)).Methods("DELETE")

	router.Handle(catTFindOne, Tokens.CheckAuthorizedUser(h.catTFindOne)).Methods("GET")
}

func (h *handler) catTAll(w http.ResponseWriter, r *http.Request) {
	arrCT, err := h.storage.CategoryTableGetAll(context.TODO())

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	js, errJs := json.Marshal(arrCT)
	if errJs != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	if err == nil && errJs == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	}
}

func (h *handler) catTFindOne(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	intId, errConv := strconv.Atoi(id)
	if errConv != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.l.Error(errConv)
	}

	cT, errFind := h.storage.CategoryTableFindOne(context.TODO(), intId)
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

func (h *handler) catTDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	intId, errConv := strconv.Atoi(id)
	if errConv != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.l.Error(errConv)
	}

	errDel := h.storage.CategoryTableDeleteItem(context.TODO(), intId)
	if errDel != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.l.Error(errDel)
	}
}

func (h *handler) catTChange(w http.ResponseWriter, r *http.Request) {
	var arrCT []CategoryTable
	json.NewDecoder(r.Body).Decode(&arrCT)

	fmt.Println("json")
	fmt.Println(arrCT[0])
	err := h.storage.CategoryTableUpdateItem(context.TODO(), arrCT[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *handler) catTAdd(w http.ResponseWriter, r *http.Request) {
	var arrPr []CategoryTable
	json.NewDecoder(r.Body).Decode(&arrPr)

	err := h.storage.CategoryTableAddItems(context.TODO(), arrPr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}
