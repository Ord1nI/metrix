package repo

import (
    "github.com/Ord1nI/metrix/internal/repo/metrics"
    "encoding/json"
)

type Repo interface {
    Adder
    Getter
    MetricGetAdder
    Closer
    json.Marshaler
}

type Adder interface {
    Add(name string, val interface{}) (error)
}

type Getter  interface {
    Get(name string, val interface{}) (error)
}

type MetricGetAdder interface {
    AddMetric(metrics.Metric) error
    GetMetric(string, string) (*metrics.Metric, bool)
}

type Closer interface {
    Close() error
}

