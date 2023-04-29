package main

import (
	"context"
	"crypto/sha512"
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
	conf := config.GetConfig()

	beforeStart(router, conf)
	start(router, conf)
}

func beforeStart(router *mux.Router, conf *config.Config) {
	l := logg.GetLogger()
	cli, _ := postgresql.NewClient(context.TODO(), l, conf.PostgresDb)

	l.Info("MakeTables(if they don't exist)")
	makeTables(cli, conf)
	makeAdmins(cli, conf)
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

// изменено
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

// изменено
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

	if conf.MakeStartTables == "true" {
		_, err := client.ExecContext(context.TODO(), SqlReq)
		if err != nil {
			//l.Error(resp)
			l.Error(err)
		}
		l.Info("Tables have made successful")
	}

}

const SqlReq = `
create table "Products"
(
    category          integer not null,
    picture_address   varchar(100),
    id                serial
        primary key
        unique,
    quantity_of_goods integer not null,
    last_price        numeric,
    product_name      varchar(100),
    available_status  varchar(50)
);

comment on column "Products".category is 'number of category ';

comment on column "Products".picture_address is 'Address of picture in file system';

comment on column "Products".product_name is 'Name of product';

comment on column "Products".available_status is 'available/not available/in stock - vars';

alter table "Products"
    owner to postgres;

create table "CategoryTable"
(
    categoryid   serial
        unique,
    categoryname varchar(100) not null,
    id           serial
        constraint "CategoryTable_pk"
            primary key
);

alter table "CategoryTable"
    owner to postgres;

create table "ProductPriceStory"
(
    id      serial
        constraint "ProductPrice_pkey"
            primary key
        unique,
    "Price" numeric   not null,
    "Date"  timestamp not null
);

alter table "ProductPriceStory"
    owner to postgres;

create table "Users"
(
    id                   serial
        constraint "Logins_pk"
            primary key
        unique,
    login                varchar(100)  not null
        unique,
    "passwordHash"       varchar(1000) not null,
    "categoryOfUser"     varchar(100)  not null,
    "dateOfRegistration" timestamp     not null,
    salt                 varchar(1000) not null,
    algorithm            varchar(200)  not null
);

alter table "Users"
    owner to postgres;

create table "UserData"
(
    id           serial
        constraint id
            primary key
        unique,
    email        varchar(50),
    phone_number varchar(30),
    country      varchar(30),
    city         varchar(30),
    index        varchar(30),
    street       varchar(30),
    number_house varchar(10),
    note         varchar(100),
    first_name   varchar(50),
    middle_name  varchar(50),
    last_name    varchar(50)
);

alter table "UserData"
    owner to postgres;

create table "Orders"
(
    "orderId"            serial
        constraint "Orders_pk"
            primary key
        unique,
    user_id              integer   not null,
    data_of_order        timestamp not null,
    ordered_products_ids integer[] not null,
    final_price          numeric   not null,
    delivery_status      varchar(50)
);

comment on column "Orders"."orderId" is 'order id';

comment on column "Orders".user_id is 'who ordered';

comment on column "Orders".ordered_products_ids is 'ids of products which was ordered';

alter table "Orders"
    owner to postgres;

`
