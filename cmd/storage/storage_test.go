package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddGauge(t *testing.T) {
    tests := []struct{
        name string
        val float64
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
    
    stor := NewEmptyStorage()

    for _, v := range tests {
        t.Run(v.name,func(t *testing.T){
        stor.AddGauge(v.name, v.val)
        assert.Equal(t, v.val, stor.gauge[v.name])
        })
    }
}

func TestAddCounter(t *testing.T) {
    tests := []struct{
        name string
        val int64
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
    
    stor := NewEmptyStorage()

    for _, v := range tests {
        t.Run(v.name,func(t *testing.T){
        stor.AddCounter(v.name, v.val)
        assert.Equal(t, v.val, stor.counter[v.name])
        })
    }
}
func TestGetGeoge(t *testing.T) {
    tests := []struct{
        name string
        val float64
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

    stor := NewEmptyStorage()

    for _, test := range tests {
        t.Run(test.name, func(t*testing.T){
        stor.gauge[test.name] = test.val
        v, err := stor.GetGauge(test.name)
        assert.Equal(t, test.val, v)
        assert.Equal(t, nil, err)
        })
    }
}
func TestGetCounter(t *testing.T) {
    tests := []struct{
        name string
        val int64
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

    stor := NewEmptyStorage()

    for _, test := range tests {
        t.Run(test.name, func(t*testing.T){
        stor.counter[test.name] = test.val
        v, err := stor.GetCounter(test.name)
        assert.Equal(t, test.val, v)
        assert.Equal(t, nil, err)
        })
    }
}
