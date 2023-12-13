package storage

import (
	"context"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

type dbProvider struct {
	st            *MemStorage
	DB            *sqlx.DB
	storeInterval int
}

func NewDBProvider(dsn string, storeInterval int, m *MemStorage) (StorageWorker, error) {
	var err error
	dbc := &dbProvider{
		st:            m,
		storeInterval: storeInterval,
	}

	if dsn == "" {
		return dbc, errors.New("Empty dsn string")
	}
	dbc.DB, err = sqlx.Open("postgres", dsn)
	if err != nil {
		return dbc, err
	}

	if dbc.DB != nil {
		dbc.DB.MustExec("CREATE TABLE IF NOT EXISTS counter_metrics (name char(30) UNIQUE, value integer);")
		dbc.DB.MustExec("CREATE TABLE IF NOT EXISTS gauge_metrics (name char(30) UNIQUE, value double precision);")
	}
	return dbc, nil
}

func (d *dbProvider) Restore() error {
	ctx := context.Background()
	rowsCounter, err := d.DB.QueryContext(ctx, "SELECT name, value FROM counter_metrics;")
	if err != nil {
		return err
	}
	if err := rowsCounter.Err(); err != nil {
		return err
	}
	defer rowsCounter.Close()

	for rowsCounter.Next() {
		var cm counterMetric
		err = rowsCounter.Scan(&cm.name, &cm.value)
		if err != nil {
			return err
		}
		d.st.UpdateCounter(cm.name, cm.value)
	}

	rowsGauge, err := d.DB.QueryContext(ctx, "SELECT name, value FROM gauge_metrics;")
	if err != nil {
		return err
	}
	if err := rowsGauge.Err(); err != nil {
		return err
	}
	defer rowsGauge.Close()

	for rowsGauge.Next() {
		var gm gaugeMetric
		err = rowsGauge.Scan(&gm.name, &gm.value)
		if err != nil {
			return err
		}
		d.st.UpdateGauge(gm.name, gm.value)
	}
	return nil
}

func (d *dbProvider) IntervalDump() {
	pollTicker := time.NewTicker(time.Duration(d.storeInterval) * time.Second)
	defer pollTicker.Stop()
	for range pollTicker.C {
		err := d.Dump()
		if err != nil {
			zap.S().Error(err)
		}
	}
}

func (d *dbProvider) Check() error {
	err := d.DB.Ping()
	if err != nil {
		return errors.New("Empty connection string")
	}
	return nil
}

func (d *dbProvider) Dump() error {
	tx := d.DB.MustBegin()

	tx.MustExec("TRUNCATE counter_metrics, gauge_metrics; ")
	for k, v := range d.st.GetCounterData() {
		tx.MustExec("INSERT INTO counter_metrics (name, value) VALUES ($1, $2); ", k, v)
	}

	for k, v := range d.st.GetGaugeData() {
		tx.MustExec("INSERT INTO gauge_metrics (name, value) VALUES ($1, $2); ", k, v)
	}

	return tx.Commit()
}
