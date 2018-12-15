package ack

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	uuid "github.com/satori/go.uuid"
)

const VERSION = 7

var (
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

	durationTiming = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "ack_duration",
			Help: "The time it took between creating the Ack and sending it by type.",
		},
		[]string{"payload_type"},
	)
)

func Gin(c *gin.Context) Ack {
	t := time.Now()
	u, _ := uuid.NewV4()

	// get uuid from header
	ru := c.Request.Header.Get("uuid")

	ack := Ack{
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
		ginContext:  c,
	}

	// timer ends of SetPayload
	ack.StartTimer()

	return ack
}

// Ack
type Ack struct {
	Version      int         `json:"ack_version"`
	Agent        string      `json:"agent"`
	SrvEnv       string      `json:"srv_env"`
	SrvNS        string      `json:"srv_ns"`
	Uuid         string      `json:"ack_uuid"`
	RequestUuid  string      `json:"req_uuid"`
	DateTime     string      `json:"date_time"`
	Success      bool        `json:"success"`
	ErrorCode    string      `json:"error_code"`
	ErrorMessage string      `json:"error_message"`
	ServerCode   int         `json:"server_code"`
	Location     string      `json:"location"`
	PayloadType  string      `json:"payload_type"`
	Payload      interface{} `json:"payload"`
	Duration     string      `json:"duration"`
	instTime     time.Time
	ginContext   *gin.Context
}

// GinError aborts gin context with JSON error
func (a *Ack) GinErrorAbort(ServerCode int, errorCode string, errorMessage string) {
	a.MakeError(ServerCode, errorCode, errorMessage)
	a.Duration = fmt.Sprintf("%s", time.Since(a.instTime))

	ackCounter.With(prometheus.Labels{"code": strconv.Itoa(a.ServerCode)}).Inc()

	a.ginContext.AbortWithStatusJSON(a.ServerCode, a)
}

// UnmarshalPostAbort unmarshals raw data posted through gin
// or aborts.
func (a *Ack) UnmarshalPostAbort(v interface{}) error {
	rs, err := a.ginContext.GetRawData()
	if err != nil {
		a.SetPayloadType("ErrorMessage")
		a.SetPayload("There was a problem with the posted data")
		a.GinErrorAbort(500, "PostDataError", err.Error())
		return err
	}

	return a.UnmarshalAbort(rs, v)
}

// UnmarshalAbort unmarshals data and aborts if it can not.
func (a *Ack) UnmarshalAbort(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		a.SetPayloadType("ErrorMessage")
		a.SetPayload("There was a problem unmarshaling data")
		a.GinErrorAbort(500, "UnmarshalError", err.Error())
	}

	return err
}

// MakeError
func (a *Ack) MakeError(ServerCode int, errorCode string, errorMessage string) {
	a.ServerCode = ServerCode
	a.Success = false
	a.PayloadType = "ErrorMessage"
	a.ErrorCode = errorCode
	a.ErrorMessage = errorMessage

	// increment a counter for this error type
	ackErrorCounter.With(prometheus.Labels{"error_type": errorMessage}).Inc()
}

// StartTimer
func (a *Ack) StartTimer() {
	a.instTime = time.Now()
}

// SetPayload
func (a *Ack) SetPayload(payload interface{}) {
	a.Duration = fmt.Sprintf("%s", time.Since(a.instTime))
	a.Payload = payload
}

// SetPayloadType
func (a *Ack) SetPayloadType(payloadType string) {
	a.PayloadType = payloadType
}

// GinSend responds with JSON on the gin context
func (a *Ack) GinSend(payload interface{}) {

	a.SetPayload(payload)

	ackCounter.With(prometheus.Labels{"code": strconv.Itoa(a.ServerCode)}).Inc()
	durationTiming.With(prometheus.Labels{"payload_type": a.PayloadType}).Observe(float64(time.Since(a.instTime).Seconds()))

	a.ginContext.JSON(a.ServerCode, a)
}
