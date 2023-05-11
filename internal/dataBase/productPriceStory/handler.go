package productPriceStory

import (
	"context"
	"encoding/json"
	"github.com/KatachiNo/Perr/internal/authCheck"
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
	go router.Handle(ppsAll, authCheck.UserAndAdmin(h.ppsAll)).Methods("GET")
	router.Handle(ppsAdd, authCheck.Admin(h.ppsAdd)).Methods("POST")

	router.Handle(ppsDelete, authCheck.Admin(h.ppsDelete)).Methods("DELETE")

	go router.Handle(ppsFindOne, authCheck.UserAndAdmin(h.ppsFindOne)).Methods("GET")
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
		h.l.Error(err)
	}
}

func (h *handler) ppsDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	intId, errConv := strconv.Atoi(id)
	if errConv != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.l.Error(errConv)
	}

	errDel := h.storage.ProductPriceTableDeleteItem(context.TODO(), intId)
	if errDel != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.l.Error(errDel)
	}
}

func (h *handler) ppsFindOne(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	intId, errConv := strconv.Atoi(id)
	if errConv != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.l.Error(errConv)
	}

	cT, errFind := h.storage.ProductPriceTableFindOne(context.TODO(), intId)
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
