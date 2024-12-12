package agent

import (
	"crypto/hmac"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	mrand "math/rand/v2"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"

	"github.com/go-resty/resty/v2"

	"github.com/Ord1nI/metrix/internal/compressor"

	"github.com/Ord1nI/metrix/internal/repo/metrics"
	"github.com/Ord1nI/metrix/internal/repo/storage"
	"github.com/Ord1nI/metrix/internal/utils"
)

func backOff(r *resty.Request, URI string, BackoffSchedule []time.Duration) (res *resty.Response, err error) {
	for _, backoff := range BackoffSchedule {
		res, err = r.Post(URI)

		if err == nil && res.StatusCode() == http.StatusOK {
			break
		}

		time.Sleep(backoff)
	}
	return res, err
}

func (a *Agent) getSubnetIP() {

}

func (a *Agent) CollectMetrics() {
	var mS runtime.MemStats
	runtime.ReadMemStats(&mS)
	memory, _ := mem.VirtualMemory()
	mGauge := storage.MGauge{
		"Alloc":         metrics.Gauge(mS.Alloc),
		"BuckHashSys":   metrics.Gauge(mS.BuckHashSys),
		"Frees":         metrics.Gauge(mS.Frees),
		"GCCPUFraction": metrics.Gauge(mS.GCCPUFraction),
		"GCSys":         metrics.Gauge(mS.GCSys),
		"HeapAlloc":     metrics.Gauge(mS.HeapAlloc),
		"HeapIdle":      metrics.Gauge(mS.HeapIdle),
		"HeapInuse":     metrics.Gauge(mS.HeapInuse),
		"HeapObjects":   metrics.Gauge(mS.HeapObjects),
		"HeapReleased":  metrics.Gauge(mS.HeapReleased),
		"HeapSys":       metrics.Gauge(mS.HeapSys),
		"LastGC":        metrics.Gauge(mS.LastGC),
		"Lookups":       metrics.Gauge(mS.Lookups),
		"MCacheInuse":   metrics.Gauge(mS.MCacheInuse),
		"MCacheSys":     metrics.Gauge(mS.MCacheSys),
		"MSpanInuse":    metrics.Gauge(mS.MSpanInuse),
		"MSpanSys":      metrics.Gauge(mS.MSpanSys),
		"Mallocs":       metrics.Gauge(mS.Mallocs),
		"NextGC":        metrics.Gauge(mS.NextGC),
		"NumForcedGC":   metrics.Gauge(mS.NumForcedGC),
		"NumGC":         metrics.Gauge(mS.NumGC),
		"OtherSys":      metrics.Gauge(mS.OtherSys),
		"PauseTotalNs":  metrics.Gauge(mS.PauseTotalNs),
		"StackInuse":    metrics.Gauge(mS.StackInuse),
		"StackSys":      metrics.Gauge(mS.StackSys),
		"Sys":           metrics.Gauge(mS.Sys),
		"TotalAlloc":    metrics.Gauge(mS.TotalAlloc),
		"RandomValue":   metrics.Gauge(mrand.Float64()),

		"TotalMemory": metrics.Gauge(memory.Total),
		"FreeMemory":  metrics.Gauge(memory.Free),
	}
	cpuUtil, _ := cpu.Percent(0, true)

	for i, v := range cpuUtil {
		mGauge.Add(fmt.Sprintf("CPUutilization%d", i+1), metrics.Gauge(v))
	}

	a.Repo.AddGauge(mGauge)

	a.Repo.Set("PollCount", metrics.Counter(1))
}

func (a *Agent) SendGaugeMetrics() error {
	for i, v := range *a.Repo.Gauge {
		var builder strings.Builder
		builder.WriteString("/update/gauge/")
		builder.WriteString(i)
		builder.WriteRune('/')
		builder.WriteString(strconv.FormatFloat(float64(v), 'f', -1, 64))

		res, err := a.Client.R().
			SetHeader("Content-Type", "text/plain").
			Post(builder.String())

		if err != nil {
			return err
		}

		if res.StatusCode() != http.StatusOK {
			return errors.New("doesnt sent")
		}
	}
	return nil
}

func (a *Agent) SendMetricJSON(data metrics.Metric) error {
	Mdata, err := json.Marshal(data)
	if err != nil {
		return err
	}
	Mdata, err = compressor.ToGzip(Mdata)

	if err != nil {
		return err
	}

	req := a.Client.R().SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("X-Real-IP", a.Config.IP).
		SetBody(Mdata)

	res, err := backOff(req, "/update/", a.Config.BackoffSchedule)

	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusOK {
		return errors.New("doesnt sent")
	}

	return nil
}
func (a *Agent) SendMetricJSONwithEncryption(keyPath string) func(data metrics.Metric) error {
	key, err := utils.ReadPublicPEM(keyPath)
	if err != nil {
		a.Logger.Fatal("Error reading public key",err)
	}
	return func(data metrics.Metric) error {
		Mdata, err := json.Marshal(data)
		if err != nil {
			return err
		}
		Mdata, err = compressor.ToGzip(Mdata)

		if err != nil {
			return err
		}

		MdataEncrypted, err := rsa.EncryptPKCS1v15(crand.Reader,key,Mdata)
		if err != nil {
			a.Logger.Error("Error while encrypting")
		}

		req := a.Client.R().SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetHeader("Accept-Encoding", "gzip").
			SetHeader("X-Real-IP", a.Config.IP).
			SetBody(MdataEncrypted)

		res, err := backOff(req, "/update/", a.Config.BackoffSchedule)

		if err != nil {
			return err
		}

		if res.StatusCode() != http.StatusOK {
			return errors.New("doesnt sent")
		}

		return nil
	}
}

func (a *Agent) SendMetricsJSON() error {
	var metricArr []metrics.Metric
	a.Repo.Get("", &metricArr)
	a.Logger.Infoln(metricArr)

	for _, m := range metricArr {
		err := a.SendMetricJSON(m)
		if err != nil {
			a.Logger.Errorln(err)
		}
	}

	return nil
}

func (a *Agent) SendMetricsArrJSON() error {
	metricsJSON, err := a.Repo.MarshalJSON()
	if err != nil {
		a.Logger.Error("Error while marshaling")
		return err
	}

	metricsJSON, err = compressor.ToGzip(metricsJSON)
	if err != nil {
		a.Logger.Error("Error while compressing")
		return err
	}

	req := a.Client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(metricsJSON)

	res, err := backOff(req, "/updates/", a.Config.BackoffSchedule)
	if err != nil {
		a.Logger.Error("Error while sending request")
		return err
	}
	if res.StatusCode() != http.StatusOK {
		a.Logger.Infoln("get status code", res.StatusCode())
		return errors.New("StatusCode != OK")
	}
	return nil
}

func (a *Agent) SendMetricsArrJSONwithSign() error {
	metricsJSON, err := a.Repo.MarshalJSON()
	if err != nil {
		a.Logger.Error("Error while marshaling")
		return err
	}

	metricsJSON, err = compressor.ToGzip(metricsJSON)
	if err != nil {
		a.Logger.Error("Error while compressing")
		return err
	}

	req := a.Client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(metricsJSON)

	signer := hmac.New(sha256.New, []byte(a.Config.Key))

	_, err = signer.Write(metricsJSON)

	if err != nil {
		a.Logger.Error("Error while signing")
		return err
	}
	req.SetHeader("HashSHA256", hex.EncodeToString(signer.Sum(nil)))

	res, err := backOff(req, "/updates/", a.Config.BackoffSchedule)
	if err != nil {
		a.Logger.Error("Error while sending request")
		return err
	}
	if res.StatusCode() != http.StatusOK {
		a.Logger.Infoln("get status code", res.StatusCode())
		return errors.New("StatusCode != OK")
	}
	return nil

}

func (a *Agent)SendMetricsArrJSONwithEncryption(key *rsa.PublicKey) error{
	metricsJSON, err := a.Repo.MarshalJSON()
	if err != nil {
		a.Logger.Error("Error while marshaling")
		return err
	}

	metricsJSON, err = compressor.ToGzip(metricsJSON)
	if err != nil {
		a.Logger.Error("Error while compressing")
		return err
	}

	encryptedMetricsJSON, err := rsa.EncryptPKCS1v15(crand.Reader,key,metricsJSON)
	if err != nil {
		a.Logger.Error("Error while encrypting")
		return err
	}

	req := a.Client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(encryptedMetricsJSON)

	res, err := backOff(req, "/updates/", a.Config.BackoffSchedule)
	if err != nil {
		a.Logger.Error("Error while sending request")
		return err
	}
	if res.StatusCode() != http.StatusOK {
		a.Logger.Infoln("get status code", res.StatusCode())
		return errors.New("StatusCode != OK")
	}
	return nil
}
