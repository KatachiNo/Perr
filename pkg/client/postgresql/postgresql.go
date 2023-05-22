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

func NewClient(ctx context.Context, l *logg.Logger, conf config.PostgresDb, envConf *config.EnvConf) (db *sql.DB, err error) {

	l.Info("Entrance to NewClient Postgresql")

	dbPort := conf.Port
	host := conf.Host
	if conf.Port == "docker" {
		dbPort = "5432"
		host = "postgres"
	}
	fmt.Print(dbPort)
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, dbPort, conf.Username, envConf.Password, conf.Dbname)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		l.Error(err)
		return nil, err
	}

	return db, nil
}
