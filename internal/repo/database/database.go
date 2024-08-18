package database 

import (
    "database/sql"
    _ "github.com/jackc/pgx/v5/stdlib"

    "encoding/json"
    "context"
    "reflect"
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
    Db *sql.DB
    ctx context.Context
}
func NewDB(ctx context.Context, dsn string,) (*Database, error){
    db, err := sql.Open("pgx", dsn)
    return &Database{db,ctx}, err
}

func(db *Database) CreateTable() error{
    _, err := db.Db.ExecContext(db.ctx,`CREATE TABLE if not EXISTS metrix(
        "name" TEXT PRIMARY KEY,
        "type" TEXT NOT NULL,
        "counter" INTEGER DEFAULT NULL,
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
        _, err = db.Db.ExecContext(db.ctx, insertMetric, name, "gauge", nil, val)
        return err
    case metrics.Counter:
        _, err = db.Db.ExecContext(db.ctx, insertMetric, name, "counter", val, nil)
        return err
    }
    return errors.New("incorect metric type")
}

func (db *Database) Get(name string, val interface{}) (error) {
    var err error

    v := reflect.ValueOf(val)
    if v.Kind() == reflect.Pointer {
        v = v.Elem()
        switch v.Type().Name(){
            case "Gauge":
                row := db.Db.QueryRowContext(db.ctx, selectMetric, name, "gauge")
                var gauge metrics.Gauge
                err = row.Scan(&gauge)
                return err
            case "Counter":
                row := db.Db.QueryRowContext(db.ctx, selectMetric, name, "counter")
                var counter metrics.Counter
                err = row.Scan(&counter)
                return err
            default:
                return errors.New("incorect val type")
        }
    }
    return errors.New("incorect val")
}
func (db *Database) AddMetric(m metrics.Metric) error {
    _, err := db.Db.ExecContext(db.ctx, insertMetric, m.ID, m.MType, m.Delta, m.Value)
    return err
}

func (db *Database) GetMetric(name string, t string) (*metrics.Metric, bool) {
    row := db.Db.QueryRowContext(db.ctx, selectMetric, name, t)

    var (
        counter sql.NullInt64
        gauge sql.NullFloat64
    )

    err := row.Scan(&counter, &gauge)

    if err != nil {
        panic(err)  //FIX as soon as possible
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
    rows, err := db.Db.QueryContext(db.ctx,"SELECT * FROM metrix")

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
    return db.Db.Close()
}
