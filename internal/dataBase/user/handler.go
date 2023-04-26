package user

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/KatachiNo/Perr/internal/Tokens"
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
	router.Handle(userChange, Tokens.CheckAuthorizedAdmin(h.userChange)).Methods("GET")
	router.Handle(userDelete, Tokens.CheckAuthorizedAdmin(h.userDelete)).Methods("GET")

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
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	fmt.Println(user)
	u, errFind := h.storage.UserFind(context.TODO(), user.Login)
	if errFind != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(errFind)
	}

	saltTable := u.Salt
	hashTable := u.PasswordHash

	hashOp := sha512.New()
	hashOp.Write([]byte(user.Password))

	byteArraySalt, err := hex.DecodeString(saltTable)
	if err != nil {
		panic(err)
	}
	fmt.Println("salt user auth saltTable", byteArraySalt)
	fmt.Println("salt user auth saltTable hex", saltTable)

	hashOp.Write(byteArraySalt)
	ha := fmt.Sprintf("%x", hashOp.Sum(nil))
	fmt.Println("ново сгенерированный хеш", ha)
	if user.Login != u.Login || hashTable != ha {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := genJWT(u.CategoryOfUser)
	if err != nil {
		fmt.Println("jsss", err)
	}
	fmt.Println("tokennnnn", token)
	tkn := Token{
		TokenString: token,
	}

	jsTkn, errJs := json.Marshal(tkn)
	if errJs != nil {
		fmt.Println("jsonTokenError", errJs)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsTkn)
}

func genJWT(categoryOfUser string) (string, error) {

	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)

	//claims["authorized"] = true
	//claims["user"] = "VASYA221"
	claims["exp"] = time.Now().Add(time.Hour * 2160).Unix()

	var tokenString string
	var err error
	if categoryOfUser == "0" {
		tokenString, err = token.SignedString(Tokens.MySigningKeyAdmin)

	} else {
		tokenString, err = token.SignedString(Tokens.MySigningKeyUser)

	}

	if err != nil {
		fmt.Errorf("Something went wrong: %s", err.Error())
	}
	fmt.Println("tokenString", tokenString)
	return tokenString, nil
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
