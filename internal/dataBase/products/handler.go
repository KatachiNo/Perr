package products

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KatachiNo/Perr/internal/handlers"
	"github.com/KatachiNo/Perr/pkg/logg"
	"github.com/gorilla/mux"
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
	//productsPriceStory        = "/products/PriceStory"
	productFindOne    = "/products/FindOne"
	productGetPicture = "/products/GetPicture"

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
	router.HandleFunc(productsAdd, h.productsAdd).Methods("POST")

	router.HandleFunc(productsChangeProductItem, h.productsChangeProductItem).Methods("PATCH")
	router.HandleFunc(productsDeleteItem, h.productsDeleteItem).Methods("DELETE")
	//router.HandleFunc(productsPriceStory, h.getProductsPriceStory).Methods("GET")

	router.HandleFunc(productGetPicture, h.productGetPicture).Methods("GET")
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

func (h *handler) productsAdd(w http.ResponseWriter, r *http.Request) {
	var arrPr []Products
	json.NewDecoder(r.Body).Decode(&arrPr)

	err := h.storage.ProductsAddItems(context.TODO(), arrPr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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
	}
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
		w.WriteHeader(http.StatusBadRequest)
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
		fmt.Print(errConv)
	}

	pr, errFind := h.storage.ProductFindOne(context.TODO(), intId)
	if errFind != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Print(errFind)
	}

	// картинка конец
	file, err := os.Open(pr.PictureAddress)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
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
func hey(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "добрый вечер")
}
