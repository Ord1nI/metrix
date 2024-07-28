package storage

import (
	"testing"
    "errors"

	"github.com/stretchr/testify/assert"
)

func TestGetGaugeE(t *testing.T) {
    tests := []struct{
        name string
        err error
    }{
        {
        name:"test",
        err: errors.New("no test in Gauge"),
        },
        {
        name:"test1",
        err: errors.New("no test1 in Gauge"),
        },
    }

    stor := NewEmptyStorage()

    for _, test := range tests {
        t.Run(test.name,func(t *testing.T){
            v, err := stor.GetGauge(test.name)
            assert.Equal(t, test.err,err)
            assert.Equal(t, float64(0), v)
        })
    }
}
func TestGetCounterE(t *testing.T) {
    tests := []struct{
        name string
        err error
    }{
        {
        name: "test",
        err: errors.New("no test in Counter"),
        },
        {
        name:"test1",
        err: errors.New("no test1 in Counter"),
        },
    }
    stor := NewEmptyStorage()

    for _, test := range tests {
        t.Run(test.name,func(t *testing.T){
            v, err := stor.GetCounter(test.name)
            assert.Equal(t, test.err,err)
            assert.Equal(t, int64(0), v)
        })
    }
}
