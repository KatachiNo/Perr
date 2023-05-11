package userData

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KatachiNo/Perr/internal/authCheck"
	"github.com/KatachiNo/Perr/internal/handlers"
	"github.com/KatachiNo/Perr/pkg/logg"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	userDAll    = "/userD/all"
	userDAdd    = "/userD/add"
	userDDelete = "/userD/delete"

	userDFindOne = "/userD/FindOne"
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
	go router.Handle(userDAll, authCheck.UserAndAdmin(h.userDAll)).Methods("GET")
	router.Handle(userDAdd, authCheck.Admin(h.userDAdd)).Methods("POST")
	router.Handle(userDDelete, authCheck.Admin(h.userDDelete)).Methods("DELETE")
	go router.Handle(userDFindOne, authCheck.UserAndAdmin(h.userDFindOne)).Methods("GET")
}

func (h *handler) userDAll(w http.ResponseWriter, r *http.Request) {
	arrUD, err := h.storage.UserDataGetAll(context.TODO())

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	js, errJs := json.Marshal(arrUD)
	if errJs != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	if err == nil && errJs == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	}
}

func (h *handler) userDAdd(w http.ResponseWriter, r *http.Request) {
	var arrPr []UserData
	json.NewDecoder(r.Body).Decode(&arrPr)

	err := h.storage.UserDataAdd(context.TODO(), arrPr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *handler) userDDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	errDel := h.storage.UserDataDelete(context.TODO(), id)
	if errDel != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.l.Error(errDel)
		h.l.Error(id)
	}
}

func (h *handler) userDFindOne(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	cT, errFind := h.storage.UserDataFindOne(context.TODO(), id)
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
