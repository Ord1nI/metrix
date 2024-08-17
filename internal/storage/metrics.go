package storage

import (
	"encoding/json"
)
type Metric struct {
    ID string `json:"id"`
    MType string `json:"type"`
    Delta *int64 `json:"delta,omitempty"`
    Value *float64 `json:"value,omitempty"`
}

type Gauge float64
type Counter int64

type MGauge map[string]Gauge
type MCounter map[string]Counter

func NewGaugeM() *MGauge{
    m := make(MGauge)
    return &m
}

func NewCounterM() *MCounter{
    m := make(MCounter)
    return &m
}

func (mG *MGauge) Get(name string) (Gauge, bool) {
    v, ok := (*mG)[name]
    return v, ok
}

func (mC *MCounter) Get(name string) (Counter, bool) {
    v, ok := (*mC)[name]
    return v, ok
}

func (mG *MGauge) Add(name string, val Gauge) {
    (*mG)[name] = val
}

func (mC *MCounter) Add(name string, val Counter) {
    (*mC)[name] += val
}

func(mC *MCounter) ToMetrics() ([]Metric){
    jm := make([]Metric,0,len(*mC))

    for i, v := range (*mC) {
        fV := int64(v)
        jm = append(jm, Metric{ID:i,MType:"counter",Delta:&fV})
    }
    return jm
}

func(mG *MGauge) ToMetrics() ([]Metric) {
    jm := make([]Metric,0,len(*mG))

    for i, v := range (*mG) {
        fV := float64(v)
        jm = append(jm, Metric{ID:i,MType:"gauge",Value:&fV})
    }
    return jm
}

func (mG *MGauge) MarshalJSON() ([]byte, error){
    jm := mG.ToMetrics()

    r, err := json.Marshal(jm)
    return r, err
}

func (mC *MCounter) MarshalJSON() ([]byte, error){
    jm := mC.ToMetrics()

    r, err := json.Marshal(jm)
    return r, err
}
