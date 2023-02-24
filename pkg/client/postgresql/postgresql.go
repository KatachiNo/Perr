package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/KatachiNo/Perr/internal/config"
	_ "github.com/lib/pq"
	"strconv"
	"time"
)

type Client interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Begin() (*sql.Tx, error)
}

func NewClient(ctx context.Context, conf config.PostgresDb) (db *sql.DB, err error) {

	maxAttempts, err := strconv.Atoi(conf.MaxAttemptsForConnection)
	if err != nil {
		panic(err)
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.Username, conf.Password, conf.Dbname)

	if conf.Username == "" || conf.Password == "" {

	}
	//изменить
	for maxAttempts > 0 {
		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		db, err = sql.Open("postgres", connStr)

		if err != nil {
			return nil, err
		}

		err = db.Ping()
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
