// Package ack
package ack

import (
	"fmt"
	"time"
)

const VERSION = 8

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

func (a *Ack) MakeError(ServerCode int, errorCode string, errorMessage string) {
	a.ServerCode = ServerCode
	a.Success = false
	a.PayloadType = "ErrorMessage"
	a.ErrorCode = errorCode
	a.ErrorMessage = errorMessage
}

func (a *Ack) StartTimer() {
	a.instTime = time.Now()
}

func (a *Ack) SetPayload(payload interface{}) {
	a.Duration = fmt.Sprintf("%s", time.Since(a.instTime))
	a.Payload = payload
}

func (a *Ack) SetPayloadType(payloadType string) {
	a.PayloadType = payloadType
}
