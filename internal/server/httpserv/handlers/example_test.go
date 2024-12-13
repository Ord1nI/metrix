package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/Ord1nI/metrix/internal/repo/metrics"
	"github.com/Ord1nI/metrix/internal/repo/storage"
)


func Example() {
    value := float64(1.312)

    //JSON tha we wil use in handler.
    metricJSON, _ := json.Marshal(metrics.Metric{
        ID: "name",
        MType:"gauge",
        Value: &value,
    })

    //Create storage.
    storage := storage.NewMemStorage()

    //Endpint
    mux := http.NewServeMux()
    mux.Handle("/update/", APIFunc(UpdateJSON(storage)))

    //Test server
    server := httptest.NewServer(mux)
    defer server.Close()

    //Send request
    res, _  := http.Post(server.URL+"/update/","application/json",bytes.NewBuffer(metricJSON))

    //Read response
    byteRes, _ := io.ReadAll(res.Body)
    defer res.Body.Close()


    //Result
    var resultMetric metrics.Metric
    json.Unmarshal(byteRes, &resultMetric)
    fmt.Println(resultMetric.ID)
    fmt.Println(resultMetric.MType)
    fmt.Println(*resultMetric.Value)

    //Output:
    //name
    //gauge
    //1.312
}
