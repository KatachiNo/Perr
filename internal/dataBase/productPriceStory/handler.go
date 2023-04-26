package productPriceStory

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
	ppsAll    = "/pps/all"
	ppsAdd    = "/pps/add"
	ppsDelete = "/pps/delete"

	ppsFindOne = "/pps/FindOne"
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
	router.Handle(ppsAll, Tokens.CheckAuthorizedUser(h.ppsAll)).Methods("GET")
	router.Handle(ppsAdd, Tokens.CheckAuthorizedAdmin(h.ppsAdd)).Methods("POST")

	router.Handle(ppsDelete, Tokens.CheckAuthorizedAdmin(h.ppsDelete)).Methods("DELETE")

	router.Handle(ppsFindOne, Tokens.CheckAuthorizedUser(h.ppsFindOne)).Methods("GET")
}

func (h *handler) ppsAll(w http.ResponseWriter, r *http.Request) {
	arrPPS, err := h.storage.ProductPriceTableGetAll(context.TODO())

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	js, errJs := json.Marshal(arrPPS)
	if errJs != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	if err == nil && errJs == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	}
}

func (h *handler) ppsAdd(w http.ResponseWriter, r *http.Request) {
	var arrPr []ProductPriceStoryTable
	json.NewDecoder(r.Body).Decode(&arrPr)

	err := h.storage.ProductPriceAddItems(context.TODO(), arrPr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *handler) ppsDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	intId, errConv := strconv.Atoi(id)
	if errConv != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Print(errConv)
	}

	errDel := h.storage.ProductPriceTableDeleteItem(context.TODO(), intId)
	if errDel != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Print(errDel)
	}
}

func (h *handler) ppsFindOne(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	intId, errConv := strconv.Atoi(id)
	if errConv != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Print(errConv)
	}

	cT, errFind := h.storage.ProductPriceTableFindOne(context.TODO(), intId)
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