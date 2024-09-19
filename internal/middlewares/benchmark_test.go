package middlewares

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

const html = `<html>
<body>
<p>Alloc = 3277184</p>
<p>BuckHashSys = 7321</p>
<p>CPUutilization1 = 7.591093117408124</p>
<p>CPUutilization10 = 5.04540867810759</p>
<p>CPUutilization11 = 1.5060240963570322</p>
<p>CPUutilization12 = 1.6080402009874488</p>
<p>CPUutilization13 = 1.612903225825306</p>
<p>CPUutilization14 = 1.3052208835444783</p>
<p>CPUutilization15 = 1.0020040080736698</p>
<p>CPUutilization16 = 0.6006006006136142</p>
<p>CPUutilization2 = 3.121852970773174</p>
<p>CPUutilization3 = 1.713709677438098</p>
<p>CPUutilization4 = 1.306532663326725</p>
<p>CPUutilization5 = 2.8311425682034113</p>
<p>CPUutilization6 = 1.604814443328343</p>
<p>CPUutilization7 = 0.9018036071798025</p>
<p>CPUutilization8 = 0.401606425711466</p>
<p>CPUutilization9 = 1.7085427135707307</p>
<p>FreeMemory = 23678066688</p>
<p>Frees = 1417574</p>
<p>GCCPUFraction = 0.0002109945283273709</p>
<p>GCSys = 2669728</p>
<p>HeapAlloc = 3277184</p>
<p>HeapIdle = 6832128</p>
<p>HeapInuse = 4521984</p>
<p>HeapObjects = 2319</p>
<p>HeapReleased = 4841472</p>
<p>HeapSys = 11354112</p>
<p>LastGC = 1726766280589277400</p>
<p>Lookups = 0</p>
<p>MCacheInuse = 19200</p>
<p>MCacheSys = 31200</p>
<p>MSpanInuse = 92640</p>
<p>MSpanInuse = 248800</p>
<p>MSpanSys = 277440</p>
<p>MSpanSys = 114240</p>
<p>Mallocs = 1419893</p>
<p>Mallocs = 6959</p>
<p>NextGC = 4194304</p>
<p>NumForcedGC = 0</p>
<p>NumGC = 4747</p>
<p>OtherSys = 2100535</p>
<p>PauseTotalNs = 375776865</p>
<p>PollCount = 1775</p>
<p>RandomValue = 0.9853556300161874</p>
<p>StackInuse = 1081344</p>
<p>StackSys = 1081344</p>
<p>Sys = 17521680</p>
<p>TotalAlloc = 9726861992</p>
<p>TotalMemory = 32554754048</p>
</body>
</html>`

const json = `{  
        "employee": {  
            "name":       "sonoo",   
            "salary":      56000,   
            "married":    true  
        }  
    }`

func HandlerMock(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func BenchmarkCompressor(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	gz := gzip.NewWriter(buf)
	gz.Write([]byte(json))
	defer gz.Close()

	req := httptest.NewRequest(http.MethodGet, "/ping", buf)
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Accept", "html")

	testF := CompressorMW(zap.NewNop().Sugar())(http.HandlerFunc(HandlerMock))

	recorder := httptest.NewRecorder()

	b.Run("compressing test", func(b *testing.B) {
		testF.ServeHTTP(recorder, req)
	})
}
