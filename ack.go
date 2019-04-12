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
	"fmt"
	"time"
)

const VERSION = 8

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

// MakeError
func (a *Ack) MakeError(ServerCode int, errorCode string, errorMessage string) {
	a.ServerCode = ServerCode
	a.Success = false
	a.PayloadType = "ErrorMessage"
	a.ErrorCode = errorCode
	a.ErrorMessage = errorMessage
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
