package database

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"context"
	"encoding/json"
	"errors"

	"github.com/Ord1nI/metrix/internal/repo/metrics"
)

const (
	insertMetric = `INSERT INTO metrix AS a(name, type, counter, gauge)
                    VALUES($1,$2,$3,$4)
                    ON CONFLICT (name) 
                    DO 
                        UPDATE 
                            SET type=EXCLUDED.type, 
                                counter= a.counter+$3,
                                gauge=EXCLUDED.gauge;`
	selectMetric = `SELECT counter, gauge FROM metrix WHERE name = $1 AND type = $2;`
)

type Database struct {
	DB       *sql.DB
	WaitTime time.Duration
}

func (db *Database) Ping() error {
	childCtx, cancel := context.WithTimeout(context.Background(), db.WaitTime)
	defer cancel()
	return db.DB.PingContext(childCtx)
}

func NewDB(dsn string, waitTime time.Duration) (*Database, error) {
	db, err := sql.Open("pgx", dsn)
	return &Database{db, waitTime}, err
}

func (db *Database) CreateTable() error {
	childCtx, cancel := context.WithTimeout(context.Background(), db.WaitTime)
	defer cancel()
	_, err := db.DB.ExecContext(childCtx, `CREATE TABLE if not EXISTS metrix(
        "name" TEXT PRIMARY KEY,
        "type" TEXT NOT NULL,
        "counter" BIGINT DEFAULT NULL,
        "gauge" DOUBLE PRECISION DEFAULT NULL);`)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) Add(name string, val interface{}) error {
	var err error

	switch val := val.(type) {
	case metrics.Gauge:
		childCtx, cancel := context.WithTimeout(context.Background(), db.WaitTime)
		defer cancel()
		_, err = db.DB.ExecContext(childCtx, insertMetric, name, "gauge", nil, val)
		return err
	case metrics.Counter:
		childCtx, cancel := context.WithTimeout(context.Background(), db.WaitTime)
		defer cancel()
		_, err = db.DB.ExecContext(childCtx, insertMetric, name, "counter", val, nil)
		return err
	case metrics.Metric:
		return db.AddMetric(val)
	case []metrics.Metric:
		return db.AddMetrics(val)
	}
	return errors.New("incorect metric type")
}

func (db *Database) Get(name string, val interface{}) error {
	switch value := val.(type) {
	case *metrics.Gauge:
		childCtx, cancel := context.WithTimeout(context.Background(), db.WaitTime)
		defer cancel()
		row := db.DB.QueryRowContext(childCtx, selectMetric, name, "gauge")
		return row.Scan(value)
	case *metrics.Counter:
		childCtx, cancel := context.WithTimeout(context.Background(), db.WaitTime)
		defer cancel()
		row := db.DB.QueryRowContext(childCtx, selectMetric, name, "counter")
		return row.Scan(value)
	case *metrics.Metric:
		m, ok := db.GetMetric(name, value.MType)
		if !ok {
			return errors.New("metric not found")
		}
		*value = *m
		return nil
	case *[]metrics.Metric:
		v, err := db.toMetrics()
		*value = v
		return err

	}
	return errors.New("incorect val")
}
func (db *Database) AddMetric(m metrics.Metric) error {
	childCtx, cancel := context.WithTimeout(context.Background(), db.WaitTime)
	defer cancel()
	_, err := db.DB.ExecContext(childCtx, insertMetric, m.ID, m.MType, m.Delta, m.Value)
	return err
}
func (db *Database) AddMetrics(m []metrics.Metric) error {
	childCtx, cancel := context.WithTimeout(context.Background(), db.WaitTime)
	defer cancel()
	tx, err := db.DB.BeginTx(childCtx, nil)

	if err != nil {
		return err
	}

	for _, v := range m {
		childCtx, cancel := context.WithTimeout(context.Background(), db.WaitTime)
		defer cancel()
		_, err := tx.ExecContext(childCtx, insertMetric, v.ID, v.MType, v.Delta, v.Value)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (db *Database) GetMetric(name string, t string) (*metrics.Metric, bool) {
	childCtx, cancel := context.WithTimeout(context.Background(), db.WaitTime)
	defer cancel()
	row := db.DB.QueryRowContext(childCtx, selectMetric, name, t)

	var (
		counter sql.NullInt64
		gauge   sql.NullFloat64
	)

	err := row.Scan(&counter, &gauge)

	if err != nil {
		return nil, false
	}

	switch true {
	case counter.Valid:
		return &metrics.Metric{
			ID:    name,
			MType: t,
			Delta: &counter.Int64,
			Value: nil,
		}, true
	case gauge.Valid:
		return &metrics.Metric{
			ID:    name,
			MType: t,
			Delta: nil,
			Value: &gauge.Float64,
		}, true
	}
	return nil, false
}

func (db *Database) toMetrics() ([]metrics.Metric, error) {
	childCtx, cancel := context.WithTimeout(context.Background(), db.WaitTime)
	defer cancel()
	rows, err := db.DB.QueryContext(childCtx, "SELECT * FROM metrix")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metricArr := make([]metrics.Metric, 0) //add limit or exact count in future

	for rows.Next() {
		var metric metrics.Metric

		err = rows.Scan(&metric.ID, &metric.MType, &metric.Delta, &metric.Value)
		if err != nil {
			return nil, err
		}
		metricArr = append(metricArr, metric)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return metricArr, nil
}

func (db *Database) MarshalJSON() ([]byte, error) {
	metricArr, err := db.toMetrics()

	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(&metricArr)

	if err != nil {
		return nil, err
	}
	return b, nil
}

func (db *Database) Close() error {
	return db.DB.Close()
}
