package user

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/KatachiNo/Perr/internal/authCheck"
	"github.com/KatachiNo/Perr/internal/handlers"
	"github.com/KatachiNo/Perr/pkg/logg"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"time"

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
	router.Handle(userChange, authCheck.Admin(h.userChange)).Methods("GET")
	router.Handle(userDelete, authCheck.Admin(h.userDelete)).Methods("DELETE")

}

func (h *handler) userReg(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	err := h.storage.UserCreate(context.TODO(), user)
	if err != nil {
		h.l.Error(err)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *handler) userAuth(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	fmt.Println(user)
	u, errFind := h.storage.UserFind(context.TODO(), user.Login)
	if errFind != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.l.Error(errFind)
	}

	saltTable := u.Salt
	hashTable := u.PasswordHash

	hashOp := sha512.New()
	hashOp.Write([]byte(user.Password))

	byteArraySalt, err := hex.DecodeString(saltTable)
	if err != nil {
		panic(err)
	}

	hashOp.Write(byteArraySalt)
	ha := fmt.Sprintf("%x", hashOp.Sum(nil))
	if user.Login != u.Login || hashTable != ha {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := genJWT(u.CategoryOfUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return

	}
	tkn := Token{
		TokenString: token,
	}

	jsTkn, errJs := json.Marshal(tkn)
	if errJs != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsTkn)
}

func genJWT(categoryOfUser string) (string, error) {
	l := logg.GetLogger()
	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Hour * 2160).Unix()

	var tokenString string
	var err error
	if categoryOfUser == "0" {
		tokenString, err = token.SignedString(authCheck.MySigningKeyAdmin)

	} else {
		tokenString, err = token.SignedString(authCheck.MySigningKeyUser)

	}

	if err != nil {
		l.Error(err)
	}
	return tokenString, nil
}

func (h *handler) userChange(w http.ResponseWriter, r *http.Request) {

}

func (h *handler) userDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	intId, errConv := strconv.Atoi(id)
	if errConv != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.l.Error(errConv)
	}

	errDel := h.storage.UserDelete(context.TODO(), intId)
	if errDel != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.l.Error(errDel)
	}
}
