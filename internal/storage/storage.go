package storage

import (
    "fmt"
    // "sort"
)

type Gauge float64
type Counter int64

type Getter interface{
    GetGauge(name string) (Gauge, error)
    GetCounter(name string) (Counter, error)
}

type Adder interface{
    AddGauge(name string, val Gauge)
    AddCounter(name string, val Counter)
}

type GetAdder interface{
    Adder
    Getter
}

type MemStorage struct{
    Gauge map[string]Gauge
    Counter map[string]Counter
}

func NewEmptyStorage() *MemStorage{
    return &MemStorage{ 
        Gauge: make(map[string]Gauge),
        Counter: make(map[string]Counter),
    }
}

func (m *MemStorage) GetGauge(name string) (Gauge, error) {
    val, ok := m.Gauge[name]
    if !ok {
        return 0, fmt.Errorf("no %s in Gauge", name)
    }
    return val, nil
}

func (m *MemStorage) GetCounter(name string) (Counter, error) {
    val, ok := m.Counter[name]
    if !ok {
        return 0, fmt.Errorf("no %s in Counter", name)
    }
    return val, nil
}

func (m *MemStorage) AddGauge(name string, val Gauge) {
    m.Gauge[name] = val
}

func (m *MemStorage) AddCounter(name string, val Counter) {
    m.Counter[name] += val
}

func (m *MemStorage) AddGaugeMap(mG map[string]Gauge) {
    m.Gauge = mG
}

func (m *MemStorage) AddCounterMap(mC map[string]Counter) {
    m.Counter = mC
}

func (m *MemStorage) GetGaugeNames() []string{
    arr := make([]string, 0, len(m.Gauge))
    for i := range m.Gauge {
        arr = append(arr, i)
    }
    return arr
}

func (m *MemStorage) GetCounterNames() []string{
    arr := make([]string, 0, len(m.Counter))
    for i := range m.Counter {
        arr = append(arr, i)
    }

    return arr
}
