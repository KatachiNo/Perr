package products

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KatachiNo/Perr/internal/handlers"
	"github.com/KatachiNo/Perr/pkg/logg"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const (
	productsAll               = "/products/all"
	productsAdd               = "/products/add"
	productsChangeProductItem = "/products/changeProductItem"
	productsDeleteItem        = "/products/delete"
	productsPriceStory        = "/products/PriceStory"
	productFindOne            = "/products/FindOne"

	addPicture = "/products/addPicture"
	testHey    = "/test"
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
	router.HandleFunc(testHey, hey).Methods("GET")

	router.HandleFunc(productsAll, h.getProductsAll).Methods("GET")
	router.HandleFunc(productsAdd, h.productAdd).Methods("POST")

	router.HandleFunc(productsChangeProductItem, hey).Methods("PATCH")
	router.HandleFunc(productsDeleteItem, h.productsDeleteItem).Methods("DELETE")
	//router.HandleFunc(productsPriceStory, h.getProductsPriceStory).Methods("GET")

	router.HandleFunc(addPicture, h.addPicture).Methods("POST")
	router.HandleFunc(productFindOne, h.productFindOne).Methods("GET")
}

func (h *handler) addPicture(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	fmt.Println(id)
	// Принимаем изображение
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Print(err)
		return
	}
	defer file.Close()

	// Сохраняем изображение на диске
	fileName := fmt.Sprintf("%s.jpg", id)
	fmt.Print(fileName)
	filePath := filepath.Join("./pictureFiles/", fileName)
	targetFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Print(err)
		return
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (h *handler) getProductsAll(w http.ResponseWriter, r *http.Request) {
	arrP, err := h.storage.ProductsGetAll(context.TODO())

	if err != nil {
		w.WriteHeader(400)
	}

	js, errJs := json.Marshal(arrP)
	if errJs != nil {
		w.WriteHeader(400)
	}
	if err == nil && errJs == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	}

}

func (h *handler) productAdd(w http.ResponseWriter, r *http.Request) {

	//json.NewDecoder(r.Body).Decode()

	var t = decimal.NewFromInt(3)
	pp := Products{
		Id:              11,
		ProductName:     "Слива",
		Category:        34,
		QuantityOfGoods: 33,
		LastPrice:       t,
		AvailableStatus: "good",
		PictureAddress:  "/123",
	}

	err := h.storage.ProductAddItem(context.TODO(), pp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	var text = "hey33" + err.Error()
	io.WriteString(w, text)

}

func (h *handler) productsDeleteItem(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	intId, errConv := strconv.Atoi(id)
	if errConv != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Print(errConv)
	}

	errDel := h.storage.ProductDeleteItem(context.TODO(), intId)
	if errDel != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Print(errDel)
	}
}

func (h *handler) productFindOne(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	intId, errConv := strconv.Atoi(id)
	if errConv != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Print(errConv)
	}

	pr, errFind := h.storage.ProductFindOne(context.TODO(), intId)
	if errFind != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Print(errFind)
	}

	js, errJs := json.Marshal(pr)
	if errJs != nil {
		w.WriteHeader(400)
	}
	if errFind == nil && errJs == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	}
}
func hey(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "добрый вечер")
}
