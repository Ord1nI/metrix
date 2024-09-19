package storage

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Ord1nI/metrix/internal/repo/metrics"
)

func ptrInt(val int64) *int64 {
	return &val
}

func ptrFloat(val float64) *float64 {
	return &val
}

func TestAddGauge(t *testing.T) {
	tests := []struct {
		name string
		val  metrics.Gauge
	}{
		{
			name: "test",
			val:  23.43,
		},
		{
			name: "test1",
			val:  23,
		},
		{
			name: "test2",
			val:  -23.32,
		},
		{
			name: "test3",
			val:  0,
		},
		{
			name: "test3",
			val:  -0,
		},
	}

	stor := NewMemStorage()

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			stor.Add(v.name, v.val)
			val, ok := stor.Gauge.Get(v.name)
			assert.Equal(t, ok, true)
			assert.Equal(t, v.val, val)
		})
	}
}

func TestAddCounter(t *testing.T) {
	tests := []struct {
		name string
		val  metrics.Counter
	}{
		{
			name: "test",
			val:  2343,
		},
		{
			name: "test1",
			val:  23,
		},
		{
			name: "test2",
			val:  -2332,
		},
		{
			name: "test3",
			val:  0,
		},
	}

	stor := NewMemStorage()

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			stor.Add(v.name, v.val)

			val, ok := stor.Counter.Get(v.name)
			assert.Equal(t, ok, true)

			assert.Equal(t, v.val, val)
		})
	}
}

func TestGetGeoge(t *testing.T) {
	tests := []struct {
		name string
		val  metrics.Gauge
	}{
		{
			name: "test",
			val:  23.43,
		},
		{
			name: "test1",
			val:  23,
		},
		{
			name: "test2",
			val:  -23.32,
		},
		{
			name: "test3",
			val:  0,
		},
		{
			name: "test3",
			val:  -0,
		},
	}

	stor := NewMemStorage()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stor.Gauge.Add(test.name, test.val)
			var v metrics.Gauge
			err := stor.Get(test.name, &v)
			assert.Equal(t, test.val, v)
			assert.Equal(t, nil, err)
		})
	}
}

func TestGetCounter(t *testing.T) {
	tests := []struct {
		name string
		val  metrics.Counter
	}{
		{
			name: "test",
			val:  2343,
		},
		{
			name: "test1",
			val:  23,
		},
		{
			name: "test2",
			val:  -2332,
		},
		{
			name: "test3",
			val:  0,
		},
		{
			name: "test3",
			val:  -0,
		},
	}

	stor := NewMemStorage()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stor.Counter.Add(test.name, test.val)
			var v metrics.Counter
			err := stor.Get(test.name, &v)
			assert.Equal(t, test.val, v)
			assert.Equal(t, nil, err)
		})
	}
}

func TestAddMetric(t *testing.T) {
	type want struct {
		metric metrics.Metric
	}
	tests := []struct {
		want   want
		metric metrics.Metric
	}{
		{
			want: want{
				metric: metrics.Metric{
					ID:    "name",
					MType: "gauge",
					Value: ptrFloat(1.5),
				},
			},
			metric: metrics.Metric{
				ID:    "name",
				MType: "gauge",
				Value: ptrFloat(1.5),
			},
		},
		{
			want: want{
				metric: metrics.Metric{
					ID:    "name2",
					MType: "gauge",
					Value: ptrFloat(1.6),
				},
			},
			metric: metrics.Metric{
				ID:    "name2",
				MType: "gauge",
				Value: ptrFloat(1.6),
			},
		},
		{
			want: want{
				metric: metrics.Metric{
					ID:    "Cname",
					MType: "counter",
					Delta: ptrInt(1),
				},
			},
			metric: metrics.Metric{
				ID:    "Cname",
				MType: "counter",
				Delta: ptrInt(1),
			},
		},
		{
			want: want{
				metric: metrics.Metric{
					ID:    "Cname",
					MType: "counter",
					Delta: ptrInt(2),
				},
			},
			metric: metrics.Metric{
				ID:    "Cname",
				MType: "counter",
				Delta: ptrInt(1),
			},
		},
	}

	stor := NewMemStorage()

	for v, test := range tests {
		t.Run(fmt.Sprintf("test %d", v), func(t *testing.T) {
			err := stor.Add("", test.metric)
			require.NoError(t, err)

			get := metrics.Metric{
				MType: test.metric.MType,
			}

			err = stor.Get(test.want.metric.ID, &get)

			require.NoError(t, err)
			assert.Equal(t, test.want.metric, get)
		})
	}
}

func TestGetMetrics(t *testing.T) {
	arr := []metrics.Metric{
		metrics.Metric{
			ID:    "cname",
			MType: "counter",
			Delta: ptrInt(1),
		},
		metrics.Metric{
			ID:    "name",
			MType: "gauge",
			Value: ptrFloat(1.5),
		},
		metrics.Metric{
			ID:    "name1",
			MType: "gauge",
			Value: ptrFloat(2.5),
		},
		metrics.Metric{
			ID:    "name2",
			MType: "gauge",
			Value: ptrFloat(3.5),
		},
		metrics.Metric{
			ID:    "name3",
			MType: "gauge",
			Value: ptrFloat(4.5),
		},
		metrics.Metric{
			ID:    "cname",
			MType: "counter",
			Delta: ptrInt(1),
		},
	}
	arrGet := []metrics.Metric{
		metrics.Metric{
			ID:    "cname",
			MType: "counter",
			Delta: ptrInt(2),
		},
		metrics.Metric{
			ID:    "name",
			MType: "gauge",
			Value: ptrFloat(1.5),
		},
		metrics.Metric{
			ID:    "name1",
			MType: "gauge",
			Value: ptrFloat(2.5),
		},
		metrics.Metric{
			ID:    "name2",
			MType: "gauge",
			Value: ptrFloat(3.5),
		},
		metrics.Metric{
			ID:    "name3",
			MType: "gauge",
			Value: ptrFloat(4.5),
		},
	}

	stor := NewMemStorage()

	stor.Add("", arr)

	var getArr []metrics.Metric

	stor.Get("", &getArr)

	sort.Slice(getArr, func(i, j int) bool { return getArr[i].ID < getArr[j].ID })

	assert.Equal(t, arrGet, getArr)
}
