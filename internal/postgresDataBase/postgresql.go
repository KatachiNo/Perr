package postgresDataBase

import (
	"database/sql"
	"fmt"
)

type ConfigDB struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func ConnectToDB(confData ConfigDB) (*sql.DB, error) {

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", confData.Username, confData.Password, confData.DBName, confData.SSLMode)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
