package storage

import (
    "errors"
    "reflect"
    "encoding/json"
    "github.com/Ord1nI/metrix/internal/myjson"
    "os"
)

type Adder interface {
    Add(name string, val interface{}) (error)
}

type MetricAdder interface {
    AddMetric(myjson.Metric)
}
type MetricGetter interface {
    GetMetric(string, string) (*myjson.Metric, bool)
}


type MetricGetAdder interface {
    MetricAdder
    MetricGetter
}

type Getter  interface {
    Get(name string, val interface{}) (error)
}

type GetAdder interface {
    Adder
    Getter
}

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
    case Gauge:
        m.Gauge.Add(name, val)
        return nil
    case Counter:
        m.Counter.Add(name, val)
        return nil
    }
    return errors.New("incorect metric type")
}

func (m *MemStorage)AddMetric(metric myjson.Metric) {
    switch metric.MType {
    case "gauge":
        m.Gauge.Add(metric.ID, Gauge(*metric.Value))
    case "counter":
        m.Counter.Add(metric.ID, Counter(*metric.Delta))
    }
}

func (m *MemStorage)GetMetric(name, mType string) (*myjson.Metric, bool){
    switch mType {
    case "gauge":
        val, ok := m.Gauge.Get(name)
        fval := float64(val)

        if !ok {
            return nil, false
        }

        mj := myjson.Metric{
            ID:name,
            MType:mType,
            Value:&fval,
        }
        return &mj, true

    case "counter":
        val, ok := m.Counter.Get(name)
        ival := int64(val)

        if !ok {
            return nil, false
        }

        mj := myjson.Metric{
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
    var metrics []myjson.Metric

    err := json.Unmarshal(d, &metrics)

    if err != nil {
        return err
    }

    for _, v := range metrics {
        m.AddMetric(v)
    }

    return nil
}

func (m *MemStorage) ToMetrics() ([]myjson.Metric){
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
