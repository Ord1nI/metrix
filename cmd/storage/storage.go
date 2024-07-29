package storage

import (
    "fmt"
)

type RepositoriesGeter interface{
    GetGauge(name string) (float64, error)
    GetCounter(name string) (int64, error)
}

type RepositoriesAdder interface{
    AddCounter(name string, val int64)
    AddGauge(name string, val float64)
}

type RepositoriesGetterAdder interface{
    AddCounter(name string, val int64)
    AddGauge(name string, val float64)
}

type memStorage struct{
    gauge map[string]float64
    counter map[string]int64
}

func NewEmptyStorage() *memStorage{
    return &memStorage{ 
        gauge: make(map[string]float64),
        counter: make(map[string]int64),
    }
}

func (m *memStorage) GetGauge(name string) (float64, error) {
    val, ok := m.gauge[name]
    if !ok {
        return 0, fmt.Errorf("no %s in Gauge", name)
    }
    return val, nil
}

func (m *memStorage) GetCounter(name string) (int64, error) {
    val, ok := m.counter[name]
    if !ok {
        return 0, fmt.Errorf("no %s in Counter", name)
    }
    return val, nil
}

func (m *memStorage) AddGauge(name string, val float64){
    m.gauge[name] = val
}

func (m *memStorage) AddCounter(name string, val int64){
    m.counter[name] += val
}

func (m *memStorage) PrintAll() {
    fmt.Println(m.counter)
    fmt.Println(m.gauge)
}
