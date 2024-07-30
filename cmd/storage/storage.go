package storage

import (
    "fmt"
    "sort"
)

type RepositoriesGetter interface{
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

type MemStorage struct{
    gauge map[string]float64
    counter map[string]int64
}

type metric struct {
    Name string
    Val float64
}

func NewEmptyStorage() *MemStorage{
    return &MemStorage{ 
        gauge: make(map[string]float64),
        counter: make(map[string]int64),
    }
}

func (m *MemStorage) GetGauge(name string) (float64, error) {
    val, ok := m.gauge[name]
    if !ok {
        return 0, fmt.Errorf("no %s in Gauge", name)
    }
    return val, nil
}

func (m *MemStorage) GetCounter(name string) (int64, error) {
    val, ok := m.counter[name]
    if !ok {
        return 0, fmt.Errorf("no %s in Counter", name)
    }
    return val, nil
}

func (m *MemStorage) AddGauge(name string, val float64){
    m.gauge[name] = val
}

func (m *MemStorage) AddCounter(name string, val int64){
    m.counter[name] += val
}

func (m *MemStorage) GetAll() []metric{
    r := make([]metric,0,len(m.gauge) +len(m.counter))

    for i, v := range m.gauge {
        r = append(r, metric{i,v})
    }

    for i, v := range m.counter {
        r = append(r, metric{i,float64(v)})
    }

    sort.SliceStable(r,func(i, j int) bool { return r[i].Name < r[j].Name})

    return r
}
