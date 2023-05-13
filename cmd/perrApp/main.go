package main

import (
	"bufio"
	"context"
	"crypto/sha512"
	"errors"
	"fmt"
	"github.com/KatachiNo/Perr/internal/authCheck"
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
	"math/rand"
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
	conf, envConf := config.GetConfig()

	beforeStart(router, conf, envConf)
	start(router, conf)
}

func beforeStart(router *mux.Router, conf *config.Config, envConf *config.EnvConf) {
	l := logg.GetLogger()
	cli, _ := postgresql.NewClient(context.TODO(), l, conf.PostgresDb, envConf)

	l.Info("Set data from .env")
	setKeysJWT(envConf)

	l.Info("MakeTables(if they don't exist)")
	go makeTables(cli, conf)
	go makeAdmins(cli, conf)
	l.Info("register products handler")
	st1 := productsDb.NewStorage(cli, l)
	h1 := products.NewRegister(st1, l)
	h1.Register(router)

	l.Info("register categorytable handler")
	st2 := categoryTableDb.NewStorage(cli, l)
	h2 := categoryTable.NewRegister(st2, l)
	h2.Register(router)

	l.Info("register ProductPriceStory handler")
	st3 := productPriceStoryDb.NewStorage(cli, l)
	h3 := productPriceStory.NewRegister(st3, l)
	h3.Register(router)

	l.Info("register Users handler")
	st4 := userDb.NewStorage(cli, l)
	h4 := user.NewRegister(st4, l)
	h4.Register(router)

	l.Info("register UsersData handler")
	st5 := UserDataDb.NewStorage(cli, l)
	h5 := userData.NewRegister(st5, l)
	h5.Register(router)

	l.Info("register orders handler")
	st6 := ordersDb.NewStorage(cli, l)
	h6 := orders.NewRegister(st6, l)
	h6.Register(router)
}

func start(router *mux.Router, conf *config.Config) {
	l := logg.GetLogger()
	l.Info("starting application . . .")

	listener, err := getListener(conf)
	if err != nil {
		l.Fatal(err)
	}

	serv := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		MaxHeaderBytes: 1 << 20,
		WriteTimeout:   10 * time.Second,
		ReadTimeout:    10 * time.Second,
	}

	err = serv.Serve(listener)
	if err != nil {
		l.Fatal(err)
	}

	l.Info("server part is started")

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

func setKeysJWT(envConf *config.EnvConf) {
	authCheck.MySigningKeyAdmin = []byte(envConf.SigningKeyAdmin)
	authCheck.MySigningKeyUser = []byte(envConf.SigningKeyUser)
}

func makeAdmins(client postgresql.Client, conf *config.Config) {
	if conf.MakeStartAdmin == "true" {
		l := logg.GetLogger()
		l.Info("Try to make admin")
		pswd := generatePassword(14)
		u := user.User{
			Login:          "adminStart",
			CategoryOfUser: "0",
			Password:       pswd,
		}

		hash := sha512.New()
		hash.Write([]byte(u.Password))

		salt := make([]byte, 128)
		_, err := rand.Read(salt)
		if err != nil {
			l.Fatal(err)
		}

		hash.Write(salt)
		h := fmt.Sprintf("%x", hash.Sum(nil))
		s := fmt.Sprintf("%x", salt)

		date := time.Now().Format("2006-01-02 15:04:05.000000")
		q := fmt.Sprintf(`INSERT INTO "Users" (login, "passwordHash", "categoryOfUser", "dateOfRegistration", salt, algorithm)
							 VALUES ('%s','%s','%s','%s','%s','%s')`,
			u.Login, h, u.CategoryOfUser, date, s, "sha512")

		_, err = client.ExecContext(context.TODO(), q)
		if err != nil {
			l.Error(err)
		} else {
			text := fmt.Sprintf("Ваш логин %s  // ваш пароль %s", u.Login, u.Password)
			l.Info(text)
		}

	}

}

func generatePassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	password := make([]byte, length)
	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}
	return string(password)
}
func makeTables(client postgresql.Client, conf *config.Config) {
	l := logg.GetLogger()
	l.Info("Make tables func")
	file, err := os.Open("tables.sql")
	if err != nil {
		l.Fatal(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var sqlReq string
	for scanner.Scan() {
		sqlReq += scanner.Text()
	}

	if conf.MakeStartTables == "true" {
		_, err = client.ExecContext(context.TODO(), sqlReq)
		if err != nil {
			l.Error(err)
			return
		}
		l.Info("Tables have made successful")
	}

}
