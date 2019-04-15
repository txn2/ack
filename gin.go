/*
   Copyright 2019 txn2

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package ack

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	io_prometheus_client "github.com/prometheus/client_model/go"
)

var (
	ackVersion = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "ack_version",
			Help: "Ack Version",
		},
	)

	ackCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ack_code_count",
			Help: "Total Acks for a code.",
		},
		[]string{"code"},
	)

	ackErrorCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ack_error_count",
			Help: "Total Error Acks for an error type.",
		},
		[]string{"error_type"},
	)

	ackUnmarshalErrorCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "ack_unmarshal_error_count",
			Help: "Total Unmarshal Error Acks.",
		},
	)

	ackPostBodyErrorCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "ack_postbody_error_count",
			Help: "The number of error from attempting to retrieve a POST body.",
		},
	)

	payloadTypeCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ack_payload_type_count",
			Help: "Total Acks for a payload type.",
		},
		[]string{"payload_type"},
	)

	durationTiming = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "ack_duration",
			Help: "The time it took between creating the Ack and sending it by type.",
		},
		[]string{"payload_type"},
	)
)

// GinAck
type GinAck struct {
	Ack         Ack
	ctx         *gin.Context
	sendHeaders bool
}

// Gin Ack
func Gin(c *gin.Context) GinAck {
	t := time.Now()
	u, _ := uuid.NewV4()

	// get uuid from header
	ru := c.Request.Header.Get("uuid")

	ak := GinAck{
		Ack: Ack{
			Agent:       os.Getenv("AGENT"),
			SrvEnv:      os.Getenv("SERVICE_ENV"),
			SrvNS:       os.Getenv("SERVICE_NS"),
			Uuid:        u.String(),
			RequestUuid: ru,
			ServerCode:  200,
			Success:     true,
			Version:     VERSION,
			DateTime:    t.Format(time.RFC3339),
			Location:    c.Request.URL.String(),
		},
		ctx: c,
	}

	ackVersion.Set(float64(ak.Ack.Version))

	// timer ends of SetPayload
	ak.Ack.StartTimer()

	return ak
}

// SetPayload
func (ga *GinAck) SetPayload(payload interface{}) {
	ga.Ack.SetPayload(payload)
}

// SetPayloadType
func (ga *GinAck) SetPayloadType(payloadType string) {
	ga.Ack.SetPayloadType(payloadType)
}

// MakeError
func (ga *GinAck) MakeError(ServerCode int, errorCode string, errorMessage string) {
	ga.Ack.MakeError(ServerCode, errorCode, errorMessage)

	// increment a counter for this error type
	ackErrorCounter.With(prometheus.Labels{"error_type": errorCode}).Inc()
}

// UnmarshalPostAbort unmarshal raw data posted through gin
// or aborts.
func (ga *GinAck) UnmarshalPostAbort(v interface{}) error {
	rs, err := ga.ctx.GetRawData()
	if err != nil {
		ga.Ack.SetPayloadType("ErrorMessage")
		ga.Ack.SetPayload("There was a problem with the posted data")
		ga.GinErrorAbort(500, "PostDataError", err.Error())
		ackPostBodyErrorCounter.Inc()
		return err
	}

	return ga.UnmarshalAbort(rs, v)
}

// UnmarshalAbort unmarshals data and aborts if it can not.
func (ga *GinAck) UnmarshalAbort(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		ga.Ack.SetPayloadType("ErrorMessage")
		ga.Ack.SetPayload("There was a problem unmarshaling data")
		ga.GinErrorAbort(500, "UnmarshalError", err.Error())
		ackUnmarshalErrorCounter.Inc()
		payloadTypeCounter.With(prometheus.Labels{"payload_type": ga.Ack.PayloadType}).Inc()
	}

	return err
}

// GinError aborts gin context with JSON error
func (ga *GinAck) GinErrorAbort(ServerCode int, errorCode string, errorMessage string) {
	ga.Ack.MakeError(ServerCode, errorCode, errorMessage)
	ga.Ack.Duration = fmt.Sprintf("%s", time.Since(ga.Ack.instTime))

	ackCounter.With(prometheus.Labels{"code": strconv.Itoa(ga.Ack.ServerCode)}).Inc()
	payloadTypeCounter.With(prometheus.Labels{"payload_type": ga.Ack.PayloadType}).Inc()

	ga.setHeaders()
	ga.ctx.AbortWithStatusJSON(ga.Ack.ServerCode, ga.Ack)
}

// GinSend responds with JSON on the gin context
func (ga *GinAck) GinSend(payload interface{}) {

	ga.Ack.SetPayload(payload)

	ackCounter.With(prometheus.Labels{"code": strconv.Itoa(ga.Ack.ServerCode)}).Inc()
	durationTiming.With(prometheus.Labels{"payload_type": ga.Ack.PayloadType}).Observe(float64(time.Since(ga.Ack.instTime).Seconds()))
	payloadTypeCounter.With(prometheus.Labels{"payload_type": ga.Ack.PayloadType}).Inc()

	ga.setHeaders()
	ga.ctx.JSON(ga.Ack.ServerCode, ga.Ack)
}

// setHeaders
func (ga *GinAck) setHeaders() {
	ga.ctx.Header("X-Ack-Version", strconv.Itoa(ga.Ack.Version))
	ga.ctx.Header("X-Ack-Agent", ga.Ack.Agent)
	ga.ctx.Header("X-Ack-Srv-Env", ga.Ack.SrvEnv)
	ga.ctx.Header("X-Ack-Srv-NS", ga.Ack.SrvNS)
	ga.ctx.Header("X-Ack-Srv-Env", ga.Ack.SrvEnv)
	ga.ctx.Header("X-Ack-Uuid", ga.Ack.Uuid)
	ga.ctx.Header("X-Ack-Req-Uuid", ga.Ack.RequestUuid)
	ga.ctx.Header("X-Ack-Payload-Type", ga.Ack.PayloadType)
	ga.ctx.Header("X-Ack-Duration", ga.Ack.Duration)
}

// MappedMetricFamily
type MappedMetricFamily map[string]*io_prometheus_client.MetricFamily
