package httpagent

import (
	"crypto/hmac"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/Ord1nI/metrix/internal/compressor"

	"github.com/Ord1nI/metrix/internal/repo/metrics"
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

func (h *HttpAgent) SendGaugeMetrics() error {
	for i, v := range *h.Repo.Gauge {
		var builder strings.Builder
		builder.WriteString("/update/gauge/")
		builder.WriteString(i)
		builder.WriteRune('/')
		builder.WriteString(strconv.FormatFloat(float64(v), 'f', -1, 64))

		res, err := h.Client.R().
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

func (h *HttpAgent) SendMetricJSON(data metrics.Metric) error {
	Mdata, err := json.Marshal(data)
	if err != nil {
		return err
	}
	Mdata, err = compressor.ToGzip(Mdata)

	if err != nil {
		return err
	}

	req := h.Client.R().SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("X-Real-IP", h.Config.IP).
		SetBody(Mdata)

	res, err := backOff(req, "/update/", h.Config.BackoffSchedule)

	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusOK {
		return errors.New("doesnt sent")
	}

	return nil
}
func (h *HttpAgent) SendMetricJSONwithEncryption(keyPath string) func(data metrics.Metric) error {
	key, err := utils.ReadPublicPEM(keyPath)
	if err != nil {
		h.Logger.Fatal("Error reading public key",err)
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
			h.Logger.Error("Error while encrypting")
		}

		req := h.Client.R().SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetHeader("Accept-Encoding", "gzip").
			SetHeader("X-Real-IP", h.Config.IP).
			SetBody(MdataEncrypted)

		res, err := backOff(req, "/update/", h.Config.BackoffSchedule)

		if err != nil {
			return err
		}

		if res.StatusCode() != http.StatusOK {
			return errors.New("doesnt sent")
		}

		return nil
	}
}

func (h *HttpAgent) SendMetricsJSON() error {
	var metricArr []metrics.Metric
	h.Repo.Get("", &metricArr)
	h.Logger.Infoln(metricArr)

	for _, m := range metricArr {
		err := h.SendMetricJSON(m)
		if err != nil {
			h.Logger.Errorln(err)
		}
	}

	return nil
}

func (h *HttpAgent) SendMetricsArrJSON() error {
	metricsJSON, err := h.Repo.MarshalJSON()
	if err != nil {
		h.Logger.Error("Error while marshaling")
		return err
	}

	metricsJSON, err = compressor.ToGzip(metricsJSON)
	if err != nil {
		h.Logger.Error("Error while compressing")
		return err
	}

	req := h.Client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(metricsJSON)

	res, err := backOff(req, "/updates/", h.Config.BackoffSchedule)
	if err != nil {
		h.Logger.Error("Error while sending request")
		return err
	}
	if res.StatusCode() != http.StatusOK {
		h.Logger.Infoln("get status code", res.StatusCode())
		return errors.New("StatusCode != OK")
	}
	return nil
}

func (h *HttpAgent) SendMetricsArrJSONwithSign() error {
	metricsJSON, err := h.Repo.MarshalJSON()
	if err != nil {
		h.Logger.Error("Error while marshaling")
		return err
	}

	metricsJSON, err = compressor.ToGzip(metricsJSON)
	if err != nil {
		h.Logger.Error("Error while compressing")
		return err
	}

	req := h.Client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(metricsJSON)

	signer := hmac.New(sha256.New, []byte(h.Config.Key))

	_, err = signer.Write(metricsJSON)

	if err != nil {
		h.Logger.Error("Error while signing")
		return err
	}
	req.SetHeader("HashSHA256", hex.EncodeToString(signer.Sum(nil)))

	res, err := backOff(req, "/updates/", h.Config.BackoffSchedule)
	if err != nil {
		h.Logger.Error("Error while sending request")
		return err
	}
	if res.StatusCode() != http.StatusOK {
		h.Logger.Infoln("get status code", res.StatusCode())
		return errors.New("StatusCode != OK")
	}
	return nil

}

func (h *HttpAgent)SendMetricsArrJSONwithEncryption(key *rsa.PublicKey) error{
	metricsJSON, err := h.Repo.MarshalJSON()
	if err != nil {
		h.Logger.Error("Error while marshaling")
		return err
	}

	metricsJSON, err = compressor.ToGzip(metricsJSON)
	if err != nil {
		h.Logger.Error("Error while compressing")
		return err
	}

	encryptedMetricsJSON, err := rsa.EncryptPKCS1v15(crand.Reader,key,metricsJSON)
	if err != nil {
		h.Logger.Error("Error while encrypting")
		return err
	}

	req := h.Client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(encryptedMetricsJSON)

	res, err := backOff(req, "/updates/", h.Config.BackoffSchedule)
	if err != nil {
		h.Logger.Error("Error while sending request")
		return err
	}
	if res.StatusCode() != http.StatusOK {
		h.Logger.Infoln("get status code", res.StatusCode())
		return errors.New("StatusCode != OK")
	}
	return nil
}
