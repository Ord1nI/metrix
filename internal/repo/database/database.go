package database 

import (
    "database/sql"
    _ "github.com/jackc/pgx/v5/stdlib"

    "encoding/json"
    "context"
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
    DB *sql.DB
    ctx context.Context
}
func NewDB(ctx context.Context, dsn string,) (*Database, error){
    db, err := sql.Open("pgx", dsn)
    return &Database{db,ctx}, err
}

func(db *Database) CreateTable() error{
    _, err := db.DB.ExecContext(db.ctx,`CREATE TABLE if not EXISTS metrix(
        "name" TEXT PRIMARY KEY,
        "type" TEXT NOT NULL,
        "counter" BIGINT DEFAULT NULL,
        "gauge" DOUBLE PRECISION DEFAULT NULL);`)
    if err != nil {
        return err
    }
    return nil
}

func (db *Database) Add(name string, val interface{}) (error) {
    var err error

    switch val := val.(type) {
    case metrics.Gauge:
        _, err = db.DB.ExecContext(db.ctx, insertMetric, name, "gauge", nil, val)
        return err
    case metrics.Counter:
        _, err = db.DB.ExecContext(db.ctx, insertMetric, name, "counter", val, nil)
        return err
    case metrics.Metric:
        return db.AddMetric(val)
    case []metrics.Metric:
        return db.AddMetrics(val)
    }
    return errors.New("incorect metric type")
}

func (db *Database) Get(name string, val interface{}) (error) {
    switch value := val.(type){
        case *metrics.Gauge:
            row := db.DB.QueryRowContext(db.ctx, selectMetric, name, "gauge")
            return row.Scan(value)
        case *metrics.Counter:
            row := db.DB.QueryRowContext(db.ctx, selectMetric, name, "counter")
            return row.Scan(value)
        case *metrics.Metric:
            m, ok := db.GetMetric(name,value.MType)
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
    _, err := db.DB.ExecContext(db.ctx, insertMetric, m.ID, m.MType, m.Delta, m.Value)
    return err
}
func (db *Database) AddMetrics(m []metrics.Metric) error {
    tx, err := db.DB.BeginTx(db.ctx, nil)

    if err != nil {
        return err
    }
    
    for _, v := range m {
        _, err := tx.ExecContext(db.ctx, insertMetric, v.ID, v.MType, v.Delta, v.Value)
        if err != nil {
            tx.Rollback()
            return err
        }
    }
    
    return tx.Commit()
}

func (db *Database) GetMetric(name string, t string) (*metrics.Metric, bool) {
    row := db.DB.QueryRowContext(db.ctx, selectMetric, name, t)

    var (
        counter sql.NullInt64
        gauge sql.NullFloat64
    )

    err := row.Scan(&counter, &gauge)

    if err != nil {
        return nil, false 
    }

    switch true {
    case counter.Valid:
        return &metrics.Metric{
            ID: name,
            MType: t,
            Delta: &counter.Int64,
            Value: nil,
        },true
    case gauge.Valid:
        return &metrics.Metric{
            ID: name,
            MType: t,
            Delta: nil,
            Value: &gauge.Float64,
        },true
    }
    return nil, false
}

func (db *Database) toMetrics() ([]metrics.Metric, error){
    rows, err:= db.DB.QueryContext(db.ctx,"SELECT * FROM metrix")

    if err != nil {
        return nil, err 
    }
    defer rows.Close()

    metricArr := make([]metrics.Metric,0) //add limit or exact count in future

    for rows.Next() {
        var metric metrics.Metric

        err = rows.Scan(&metric.ID,&metric.MType,&metric.Delta,&metric.Value)
        if err != nil {
            return nil, err
        }
        metricArr = append(metricArr, metric)
    }

    if err = rows.Err(); err != nil {
        return nil , err
    }
    
    return metricArr, nil
}

func(db *Database) MarshalJSON() ([]byte, error){
    metricArr,err := db.toMetrics()

    if err != nil {
        return nil, err
    }

    b, err := json.Marshal(&metricArr)

    if err != nil {
        return nil, err
    }
    return b,nil
}

func (db *Database) Close() error {
    return db.DB.Close()
}
