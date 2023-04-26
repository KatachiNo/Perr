package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/KatachiNo/Perr/internal/dataBase/categoryTable"
	categoryTableDb "github.com/KatachiNo/Perr/internal/dataBase/categoryTable/db"
	"github.com/KatachiNo/Perr/internal/dataBase/orders"
	ordersDb "github.com/KatachiNo/Perr/internal/dataBase/orders/db"
	"github.com/KatachiNo/Perr/internal/dataBase/productPriceStory"
	productPriceStoryDb "github.com/KatachiNo/Perr/internal/dataBase/productPriceStory/db"
	"github.com/KatachiNo/Perr/internal/dataBase/products"
	productsDb "github.com/KatachiNo/Perr/internal/dataBase/products/db"
	"github.com/KatachiNo/Perr/internal/dataBase/user"
	userDb "github.com/KatachiNo/Perr/internal/dataBase/user/db"
	"github.com/KatachiNo/Perr/internal/dataBase/userData"
	UserDataDb "github.com/KatachiNo/Perr/internal/dataBase/userData/db"
	"github.com/KatachiNo/Perr/pkg/client/postgresql"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/KatachiNo/Perr/internal/config"

	"github.com/KatachiNo/Perr/pkg/logg"
	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

func main() {
	l := logg.GetLogger()

	l.Info("create router")
	router := mux.NewRouter()

	conf := config.GetConfig()

	cli, _ := postgresql.NewClient(context.TODO(), l, conf.PostgresDb)
	l.Info("register products handler")
	st := productsDb.NewStorage(cli, l)
	h := products.NewRegister(st, l)
	h.Register(router)

	l.Info("register categorytable handler")
	st1 := categoryTableDb.NewStorage(cli, l)
	hh := categoryTable.NewRegister(st1, l)
	hh.Register(router)

	l.Info("register ProductPriceStory handler")
	st11 := productPriceStoryDb.NewStorage(cli, l)
	hhh := productPriceStory.NewRegister(st11, l)
	hhh.Register(router)

	l.Info("register Users handler")
	st22 := userDb.NewStorage(cli, l)
	hhhh := user.NewRegister(st22, l)
	hhhh.Register(router)

	l.Info("register UsersData handler")
	st222 := UserDataDb.NewStorage(cli, l)
	h1 := userData.NewRegister(st222, l)
	h1.Register(router)

	l.Info("register orders handler")
	stt := ordersDb.NewStorage(cli, l)
	h3 := orders.NewRegister(stt, l)
	h3.Register(router)

	start(router, conf)
}

func start(router *mux.Router, conf *config.Config) {
	l := logg.GetLogger()
	l.Info("starting application . . .")

	listener, err := getListener(conf)
	if err != nil {
		l.Fatal(err)
	}

	serv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	err = serv.Serve(listener)
	if err != nil {
		l.Fatal(err)
	}

	l.Info("application is started")

}

func getListener(conf *config.Config) (net.Listener, error) {
	const (
		socket = "socket"
		port   = "port"
	)

	l := logg.GetLogger()

	var listener net.Listener
	var listenErr error

	switch conf.Server.Type {
	case port:
		l.Info("listen tcp")
		listener, listenErr = net.Listen(
			"tcp",
			fmt.Sprintf(":%s", conf.Server.Port),
		)
		l.Infof("configuration server type %s, by port %s", conf.Server.Type, conf.Server.Port)

	case socket:
		directoryOfApp, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			l.Fatal(err)
		}
		// make socket
		l.Info("create socket")
		socketPath := path.Join(directoryOfApp, "perr.sock")

		l.Info("create unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		l.Infof("configuration server type %s", conf.Server.Type)

	default:
		listenErr = errors.New("unacceptable type in server configuration")
	}

	return listener, listenErr
}
