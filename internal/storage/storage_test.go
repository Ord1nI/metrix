package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddGauge(t *testing.T) {
    tests := []struct{
        name string
        val Gauge
    }{
        {
        name:"test",
        val: 23.43,
        },
        {
        name:"test1",
        val: 23,
        },
        {
        name:"test2",
        val: -23.32,
        },
        {
        name:"test3",
        val: 0,
        },
        {
        name:"test3",
        val: -0,
        },
    }
    
    stor := NewMemStorage()

    for _, v := range tests {
        t.Run(v.name,func(t *testing.T){
        stor.Add(v.name, v.val)
        val, ok := stor.Gauge.Get(v.name)
        assert.Equal(t, ok, true)
        assert.Equal(t, v.val, val)
        })
    }
}

func TestAddCounter(t *testing.T) {
    tests := []struct{
        name string
        val Counter
    }{
        {
        name:"test",
        val: 2343,
        },
        {
        name:"test1",
        val: 23,
        },
        {
        name:"test2",
        val: -2332,
        },
        {
        name:"test3",
        val: 0,
        },
    }
    
    stor := NewMemStorage()

    for _, v := range tests {
        t.Run(v.name,func(t *testing.T){
        stor.Add(v.name, v.val)
        
        val, ok := stor.Counter.Get(v.name)
        assert.Equal(t, ok, true)

        assert.Equal(t, v.val, val)
        })
    }
}
func TestGetGeoge(t *testing.T) {
    tests := []struct{
        name string
        val Gauge
    }{
        {
        name:"test",
        val: 23.43,
        },
        {
        name:"test1",
        val: 23,
        },
        {
        name:"test2",
        val: -23.32,
        },
        {
        name:"test3",
        val: 0,
        },
        {
        name:"test3",
        val: -0,
        },
    }

    stor := NewMemStorage()

    for _, test := range tests {
        t.Run(test.name, func(t*testing.T){
        stor.Gauge.Add(test.name, test.val)
        var v Gauge
        err := stor.Get(test.name, &v)
        assert.Equal(t, test.val, v)
        assert.Equal(t, nil, err)
        })
    }
}
func TestGetCounter(t *testing.T) {
    tests := []struct{
        name string
        val Counter
    }{
        {
        name:"test",
        val: 2343,
        },
        {
        name:"test1",
        val: 23,
        },
        {
        name:"test2",
        val: -2332,
        },
        {
        name:"test3",
        val: 0,
        },
        {
        name:"test3",
        val: -0,
        },
    }

    stor := NewMemStorage()

    for _, test := range tests {
        t.Run(test.name, func(t*testing.T){
        stor.Counter.Add(test.name, test.val)
        var v Counter
        err := stor.Get(test.name, &v)
        assert.Equal(t, test.val, v)
        assert.Equal(t, nil, err)
        })
    }
}
