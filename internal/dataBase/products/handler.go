package products

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KatachiNo/Perr/internal/authCheck"
	"github.com/KatachiNo/Perr/internal/handlers"
	"github.com/KatachiNo/Perr/pkg/logg"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// изменено
const (
	productsAll               = "/products/all"
	productsAdd               = "/products/add"
	productsChangeProductItem = "/products/changeProductItem"
	productsDeleteItem        = "/products/delete"
	productFindOne            = "/products/FindOne"
	productGetPicture         = "/products/GetPicture"

	addPicture = "/products/addPicture"
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

	router.Handle(productsAll, authCheck.UserAndAdmin(h.getProductsAll)).Methods("GET")
	router.Handle(productsAdd, authCheck.Admin(h.productsAdd)).Methods("POST")

	router.Handle(productsChangeProductItem, authCheck.Admin(h.productsChangeProductItem)).Methods("PATCH")
	router.Handle(productsDeleteItem, authCheck.Admin(h.productsDeleteItem)).Methods("DELETE")

	router.Handle(productGetPicture, authCheck.UserAndAdmin(h.productGetPicture)).Methods("GET")
	router.Handle(addPicture, authCheck.Admin(h.addPicture)).Methods("POST")
	router.Handle(productFindOne, authCheck.UserAndAdmin(h.productFindOne)).Methods("GET")
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

	intId, errConv := strconv.Atoi(id)
	if errConv != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Print(errConv)
	}

	pr := Products{
		Id:              intId,
		ProductName:     "null",
		Category:        -1,
		QuantityOfGoods: -1,
		LastPrice:       "null",
		AvailableStatus: "null",
		PictureAddress:  filePath,
	}
	fmt.Println(pr)
	err = h.storage.ProductUpdateItem(context.TODO(), pr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func (h *handler) getProductsAll(w http.ResponseWriter, r *http.Request) {
	arrP, err := h.storage.ProductsGetAll(context.TODO())

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	js, errJs := json.Marshal(arrP)
	if errJs != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	if err == nil && errJs == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	}
}

func (h *handler) productsAdd(w http.ResponseWriter, r *http.Request) {
	var arrPr []Products
	json.NewDecoder(r.Body).Decode(&arrPr)

	err := h.storage.ProductsAddItems(context.TODO(), arrPr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.l.Error(err)
	}
}

func (h *handler) productsChangeProductItem(w http.ResponseWriter, r *http.Request) {
	var pr []Products
	json.NewDecoder(r.Body).Decode(&pr)

	fmt.Println("json")
	fmt.Println(pr[0])
	err := h.storage.ProductUpdateItem(context.TODO(), pr[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.l.Error(err)
	}
}

func (h *handler) productsDeleteItem(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	intId, errConv := strconv.Atoi(id)
	if errConv != nil {
		w.WriteHeader(http.StatusBadRequest)

		h.l.Error(errConv)
	}

	errDel := h.storage.ProductDeleteItem(context.TODO(), intId)
	if errDel != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.l.Error(errDel)
	}
}

func (h *handler) productFindOne(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	intId, errConv := strconv.Atoi(id)
	if errConv != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.l.Error(errConv)
	}

	pr, errFind := h.storage.ProductFindOne(context.TODO(), intId)
	if errFind != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.l.Error(errFind)
	}

	js, errJs := json.Marshal(pr)
	if errJs != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.l.Error(errJs)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(js)

}

func (h *handler) productGetPicture(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	intId, errConv := strconv.Atoi(id)
	if errConv != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.l.Error(errConv)
	}

	pr, errFind := h.storage.ProductFindOne(context.TODO(), intId)
	if errFind != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.l.Error(errFind)
	}

	// картинка конец
	file, err := os.Open("pictureFiles" + pr.PictureAddress)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.l.Error(err)
		return
	}
	defer file.Close()

	// Определение Content-Type и отправка картинки в ответ на запрос
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)
	file.Read(buffer)
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileSize))
	w.Write(buffer)
}
