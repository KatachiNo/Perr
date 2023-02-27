package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/KatachiNo/Perr/internal/config"
	"github.com/KatachiNo/Perr/pkg/logg"
	_ "github.com/lib/pq"
)

type Client interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func NewClient(ctx context.Context, l *logg.Logger, conf config.PostgresDb) (db *sql.DB, err error) {

	l.Info("Entrance to NewClient Postgresql ")

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.Username, conf.Password, conf.Dbname)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}