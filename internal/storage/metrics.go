package storage

import (
	"encoding/json"

	"github.com/Ord1nI/metrix/internal/myjson"
)

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

func(mC *MCounter) ToMetrics() ([]myjson.Metric){
    jm := make([]myjson.Metric,0,len(*mC))

    for i, v := range (*mC) {
        fV := int64(v)
        jm = append(jm, myjson.Metric{ID:i,MType:"counter",Delta:&fV})
    }
    return jm
}

func(mG *MGauge) ToMetrics() ([]myjson.Metric) {
    jm := make([]myjson.Metric,0,len(*mG))

    for i, v := range (*mG) {
        fV := float64(v)
        jm = append(jm, myjson.Metric{ID:i,MType:"gauge",Value:&fV})
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
