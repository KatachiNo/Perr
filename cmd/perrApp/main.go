package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/KatachiNo/Perr/internal/dataBase/products"
	productsDb "github.com/KatachiNo/Perr/internal/dataBase/products/db"
	"github.com/KatachiNo/Perr/internal/dataBase/user"
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

	l.Info("register user handler")
	handler := user.NewRegister()
	handler.Register(router)

	cli, _ := postgresql.NewClient(context.TODO(), l, conf.PostgresDb)
	l.Info("register products handler")
	st := productsDb.NewStorage(cli, l)
	h := products.NewRegister(st, l)
	h.Register(router)

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
