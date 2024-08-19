package storage

import (
	"encoding/json"
    "github.com/Ord1nI/metrix/internal/repo/metrics"
)


type MGauge map[string]metrics.Gauge
type MCounter map[string]metrics.Counter

func NewGaugeM() *MGauge{
    m := make(MGauge)
    return &m
}

func NewCounterM() *MCounter{
    m := make(MCounter)
    return &m
}

func (mG *MGauge) Get(name string) (metrics.Gauge, bool) {
    v, ok := (*mG)[name]
    return v, ok
}

func (mC *MCounter) Get(name string) (metrics.Counter, bool) {
    v, ok := (*mC)[name]
    return v, ok
}

func (mG *MGauge) Add(name string, val metrics.Gauge) {
    (*mG)[name] = val
}

func (mC *MCounter) Add(name string, val metrics.Counter) {
    (*mC)[name] += val
}

func (mG *MGauge) Set(name string, val metrics.Gauge) {
    mG.Add(name,val)
}

func (mC *MCounter) Set(name string, val metrics.Counter) {
    (*mC)[name] = val
}

func(mC *MCounter) ToMetrics() ([]metrics.Metric){
    jm := make([]metrics.Metric,0,len(*mC))

    for i, v := range (*mC) {
        fV := int64(v)
        jm = append(jm, metrics.Metric{ID:i,MType:"counter",Delta:&fV})
    }
    return jm
}

func(mG *MGauge) ToMetrics() ([]metrics.Metric) {
    jm := make([]metrics.Metric,0,len(*mG))

    for i, v := range (*mG) {
        fV := float64(v)
        jm = append(jm, metrics.Metric{ID:i,MType:"gauge",Value:&fV})
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
