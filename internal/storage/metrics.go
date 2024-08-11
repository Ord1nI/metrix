package storage

import(
    "encoding/json"
)

type Gauge float64
type Counter int64

type mGauge map[string]Gauge
type mCounter map[string]Counter

type jMetric struct {
    ID string `json:"id"`
    Mtype string  `json:"type"`
    Delta *int64 `json:"delta,omitempty"`
    Value *float64 `json:"value,omitempty"`
}

func NewGaugeM() *mGauge{
    m := make(mGauge)
    return &m
}

func NewCounterM() *mCounter{
    m := make(mCounter)
    return &m
}

func (mG *mGauge) Get(name string) (Gauge, bool) {
    v, ok := (*mG)[name]
    return v, ok
}

func (mC *mCounter) Get(name string) (Counter, bool) {
    v, ok := (*mC)[name]
    return v, ok
}

func (mG *mGauge) Add(name string, val Gauge) {
    (*mG)[name] = val
}

func (mC *mCounter) Add(name string, val Counter) {
    (*mC)[name]+=val
}

func(mC *mCounter) tojMetrics() ([]jMetric){
    jm := make([]jMetric,0,len(*mC))

    for i, v := range (*mC) {
        fV := int64(v)
        jm = append(jm, jMetric{ID:i,Mtype:"counter",Delta:&fV})
    }
    return jm
}

func(mG *mGauge) tojMetrics() ([]jMetric) {
    jm := make([]jMetric,0,len(*mG))

    for i, v := range (*mG) {
        fV := float64(v)
        jm = append(jm, jMetric{ID:i,Mtype:"gauge",Value:&fV})
    }
    return jm
}

func (mG *mGauge) MarshalJSON() ([]byte, error){
    jm := mG.tojMetrics()

    r, err := json.Marshal(jm)
    return r, err
}

func (mC *mCounter) MarshalJSON() ([]byte, error){
    jm := mC.tojMetrics()

    r, err := json.Marshal(jm)
    return r, err
}
