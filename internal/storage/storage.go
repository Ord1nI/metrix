package storage

import (
    "errors"
    "reflect"
    "encoding/json"
)

type Adder interface {
    Add(name string, val interface{}) (error)
}

type Getter  interface {
    Get(name string, val interface{}) (error)
}

type GetAdder interface {
    Adder
    Getter
}

type MemStorage struct{
    Gauge *mGauge
    Counter *mCounter
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
    jm := m.Gauge.tojMetrics()
    jm = append(jm, m.Counter.tojMetrics()...)

    r, err := json.Marshal(jm)
    if err != nil {
        return nil,err
    }

    return r, nil
}

func (m *MemStorage) AddGauge(g mGauge) {
    m.Gauge = &g
}

func (m *MemStorage) AddCounter(c mCounter) {
    m.Counter = &c
}
