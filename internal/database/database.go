package database

import (
	"context"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lionslon/go-yapmetrics/internal/storage"
	"github.com/pkg/errors"
)

type DBConnection struct {
	DB *sql.DB
}

type counterMetric struct {
	name  string
	value int64
}

type gaugeMetric struct {
	name  string
	value float64
}

func New(dsn string) *DBConnection {
	dbc := &DBConnection{}

	if dsn == "" {
		dbc.DB = nil
		return dbc
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		zap.S().Error(err)
		dbc.DB = nil
		return dbc
	} else {
		dbc.DB = db
	}

	if dbc.DB != nil {
		dbc.DB.Exec("CREATE TABLE IF NOT EXISTS counter_metrics (name char(30) UNIQUE, value integer);")
		dbc.DB.Exec("CREATE TABLE IF NOT EXISTS gauge_metrics (name char(30) UNIQUE, value double precision);")
	}
	return dbc
}

func CheckConnection(dbc *DBConnection) error {
	if dbc.DB != nil {
		err := dbc.DB.Ping()
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("Empty connection string")
}

func Restore(s *storage.MemStorage, dbc *DBConnection) {
	if dbc.DB == nil {
		return
	}

	ctx := context.Background()
	rowsCounter, err := dbc.DB.QueryContext(ctx, "SELECT name, value FROM counter_metrics;")
	if err != nil {
		zap.S().Error(err)
	}
	if err := rowsCounter.Err(); err != nil {
		zap.S().Error(err)
	}
	defer rowsCounter.Close()

	for rowsCounter.Next() {
		var cm counterMetric
		err = rowsCounter.Scan(&cm.name, &cm.value)
		if err != nil {
			zap.S().Error(err)
		}
		s.UpdateCounter(cm.name, cm.value)
	}

	rowsGauge, err := dbc.DB.QueryContext(ctx, "SELECT name, value FROM gauge_metrics;")
	if err != nil {
		zap.S().Error(err)
	}
	if err := rowsGauge.Err(); err != nil {
		zap.S().Error(err)
	}
	defer rowsGauge.Close()

	for rowsGauge.Next() {
		var gm gaugeMetric
		err = rowsGauge.Scan(&gm.name, &gm.value)
		if err != nil {
			zap.S().Error(err)
		}
		s.UpdateGauge(gm.name, gm.value)
	}
}

func Dump(s *storage.MemStorage, dbc *DBConnection, storeInterval int) {
	pollTicker := time.NewTicker(time.Duration(storeInterval) * time.Second)
	defer pollTicker.Stop()
	for range pollTicker.C {
		saveMetrics(s, dbc)
	}
}

func saveMetrics(s *storage.MemStorage, dbc *DBConnection) error {
	tx, err := dbc.DB.Begin()
	if err != nil {
		return err
	}
	var query string
	query = "TRUNCATE counter_metrics, gauge_metrics; "
	for k, v := range s.GetCounterData() {
		query += fmt.Sprintf("INSERT INTO counter_metrics (name, value) VALUES ('%s', %d); ", k, v)
	}

	for k, v := range s.GetGaugeData() {
		query += fmt.Sprintf("INSERT INTO gauge_metrics (name, value) VALUES ('%s', %f); ", k, v)
	}

	_, err = tx.Exec(query)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
