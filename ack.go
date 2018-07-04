package ack

import (
	"fmt"
	"time"

	"os"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

const VERSION = 6

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
