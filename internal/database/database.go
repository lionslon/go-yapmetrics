package database

import (
	"database/sql"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
)

type DBConnection struct {
	db *sql.DB
}

func New(dsn string) *DBConnection {
	dbc := &DBConnection{}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		dbc.db = nil
		zap.S().Error(err)
	} else {
		dbc.db = db
	}

	return dbc
}

func CheckConnection(dbc *DBConnection) error {
	if dbc.db != nil {
		err := dbc.db.Ping()
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("Empty connection string")
}
