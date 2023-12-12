package database

import (
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/lionslon/go-yapmetrics/internal/storage"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

type DBConnection struct {
	DB *sqlx.DB
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

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		zap.S().Error(err)
		dbc.DB = nil
		return dbc
	} else {
		dbc.DB = db
	}

	if dbc.DB != nil {
		dbc.DB.MustExec("CREATE TABLE IF NOT EXISTS counter_metrics (name char(30) UNIQUE, value integer);")
		dbc.DB.MustExec("CREATE TABLE IF NOT EXISTS gauge_metrics (name char(30) UNIQUE, value double precision);")
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
	tx := dbc.DB.MustBegin()

	tx.MustExec("TRUNCATE counter_metrics, gauge_metrics; ")
	for k, v := range s.GetCounterData() {
		tx.MustExec("INSERT INTO counter_metrics (name, value) VALUES ($1, $2); ", k, v)
	}

	for k, v := range s.GetGaugeData() {
		tx.MustExec("INSERT INTO gauge_metrics (name, value) VALUES ($1, $2); ", k, v)
	}

	return tx.Commit()
}
