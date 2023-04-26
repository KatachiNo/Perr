package userData

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KatachiNo/Perr/internal/Tokens"
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
	router.Handle(userDAll, Tokens.CheckAuthorizedUser(h.userDAll)).Methods("GET")
	router.Handle(userDAdd, Tokens.CheckAuthorizedUser(h.userDAdd)).Methods("POST")

	router.Handle(userDDelete, Tokens.CheckAuthorizedUser(h.userDDelete)).Methods("DELETE")

	router.Handle(userDFindOne, Tokens.CheckAuthorizedUser(h.userDFindOne)).Methods("GET")
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
		fmt.Print(errDel)
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
