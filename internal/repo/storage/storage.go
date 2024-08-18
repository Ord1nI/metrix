package storage

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
    "time"

    "github.com/Ord1nI/metrix/internal/repo/metrics"
)

type MemStorage struct{
    Gauge *MGauge
    Counter *MCounter
}

func NewMemStorage() *MemStorage{
    return &MemStorage{
        Gauge: NewGaugeM(),
        Counter: NewCounterM(),
    }
}

func (m *MemStorage) Add(name string, val interface{}) error {
    switch val := val.(type) {
    case metrics.Gauge:
        m.Gauge.Add(name, val)
        return nil
    case metrics.Counter:
        m.Counter.Add(name, val)
        return nil
    }
    return errors.New("incorect metric type")
}

func (m *MemStorage)AddMetric(metric metrics.Metric) error{
    if metric.ID == "" {
        return errors.New("metric must hame name")
    }
    switch metric.MType {
    case "gauge":
        m.Gauge.Add(metric.ID, metrics.Gauge(*metric.Value))
        return nil
    case "counter":
        m.Counter.Add(metric.ID, metrics.Counter(*metric.Delta))
        return nil
    }
    return errors.New("bad type")
}

func (m *MemStorage)GetMetric(name, mType string) (*metrics.Metric, bool){
    switch mType {
    case "gauge":
        val, ok := m.Gauge.Get(name)
        fval := float64(val)

        if !ok {
            return nil, false
        }

        mj := metrics.Metric{
            ID : name,
            MType : mType,
            Value : &fval,
        }
        return &mj, true

    case "counter":
        val, ok := m.Counter.Get(name)
        ival := int64(val)

        if !ok {
            return nil, false
        }

        mj := metrics.Metric{
            ID:name,
            MType:mType,
            Delta:&ival,
        }
        return &mj, true
    }
    return nil, false
}

func (m *MemStorage) Get(name string, val interface{}) error{
    v := reflect.ValueOf(val)
    if v.Kind() == reflect.Pointer {
        v = v.Elem()
        switch v.Type().Name(){
            case "Gauge":
                val, ok := m.Gauge.Get(name)
                if ok {
                    v.SetFloat(float64(val))
                    return nil
                } else {
                    return errors.New("no variable by this name")
                }
            case "Counter":
                val, ok := m.Counter.Get(name)
                if ok {
                    v.SetInt(int64(val))
                    return nil
                } else {
                    return errors.New("no variable by this name")
                }
            default:
                return errors.New("incorect val type")
        }
    }
    return errors.New("incorect val")
}

func (m *MemStorage) MarshalJSON() ([]byte, error){
    jm := m.Gauge.ToMetrics()
    jm = append(jm, m.Counter.ToMetrics()...)

    r, err := json.Marshal(jm)
    if err != nil {
        return nil,err
    }

    return r, nil
}

func (m *MemStorage) UnmarshalJSON(d []byte) error {
    var metrics []metrics.Metric

    err := json.Unmarshal(d, &metrics)

    if err != nil {
        return err
    }

    for _, v := range metrics {
        m.AddMetric(v)
    }

    return nil
}

func (m *MemStorage) ToMetrics() ([]metrics.Metric){
    jm := m.Gauge.ToMetrics()
    jm = append(jm, m.Counter.ToMetrics()...)
    return jm
}

func (m *MemStorage) AddGauge(g MGauge) {
    m.Gauge = &g
}

func (m *MemStorage) AddCounter(c MCounter) {
    m.Counter = &c
}

func (m *MemStorage) WriteToFile(f string) error {
    json, err := m.MarshalJSON()

    if err != nil {
        return err
     }

     err = os.WriteFile(f, json, 0666)
     return err
}

func (m *MemStorage) GetFromFile(f string) (error) {
    buf, err := os.ReadFile(f)

    if err != nil {
        return err
    }

    err = m.UnmarshalJSON(buf)

    if err != nil {
        return err
    }

    return nil
}

func (m *MemStorage) Close() error {
    return nil
}

func (m *MemStorage)StartDataSaver(StoreInterval int, FileStoragePath string) error{
    for {
        time.Sleep(time.Duration(StoreInterval) * time.Second)
        err := m.WriteToFile(FileStoragePath)
        if err != nil {
            return err
        }
        return nil
    }
}
