//
// Copyright 2019 AT&T Intellectual Property
// Copyright 2019 Nokia
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//  This source code is part of the near-RT RIC (RAN Intelligent Controller)
//  platform project (RICP).


package rmr

// #cgo LDFLAGS: -L/usr/local/lib -lrmr_nng -lnng
// #include <rmr/rmr.h>
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"github.com/pkg/errors"
	"time"
	"unsafe"
)

func (*Context) Init(port string, maxMsgSize int, maxRetries int, flags int) *Messenger {
	pp := C.CString(port)
	defer C.free(unsafe.Pointer(pp))
	ctx := NewContext(maxMsgSize, maxRetries, flags, C.rmr_init(pp, C.int(maxMsgSize), C.int(flags)))
	start := time.Now()
	for !ctx.IsReady() {
		time.Sleep(time.Second)
		if time.Since(start) >= time.Minute {
			start = time.Now()
		}
	}
	// Configure the rmr to make rounds of attempts to send a message before notifying the application that it should retry.
	// Each round is about 1000 attempts with a short sleep between each round.
	C.rmr_set_stimeout(ctx.RmrCtx, C.int(0))
	r := Messenger(ctx)
	return &r
}

func (ctx *Context) SendMsg(msg *MBuf) (*MBuf, error) {

	allocatedCMBuf, err := ctx.getAllocatedCRmrMBuf(msg, ctx.MaxMsgSize)
	if err != nil {
		return nil, err
	}
	if state := allocatedCMBuf.state; state != RMR_OK {
		errorMessage := fmt.Sprintf("#rmrCgoApi.SendMsg - Failed to get allocated message. state: %v - %s", state, states[int(state)])
		return nil, errors.New(errorMessage)
	}
	defer C.rmr_free_msg(allocatedCMBuf)

	for i := 0; i < ctx.MaxRetries; i++ {
		currCMBuf := C.rmr_send_msg(ctx.RmrCtx, allocatedCMBuf)
		if state := currCMBuf.state; state != RMR_OK {
			if state != RMR_ERR_RETRY {
				errorMessage := fmt.Sprintf("#rmrCgoApi.SendMsg - Failed to send message. state: %v - %s", state, states[int(state)])
				return nil, errors.New(errorMessage)
			}
			time.Sleep(100 * time.Millisecond)
			continue
		}
		return convertToMBuf(currCMBuf)
	}

	return nil, errors.New(fmt.Sprintf("#rmrCgoApi.SendMsg - Too many retries"))
}

func (ctx *Context) RecvMsg() (*MBuf, error) {
	allocatedCMBuf, err := C.rmr_alloc_msg(ctx.RmrCtx, C.int(ctx.MaxMsgSize))
	if err != nil {
		return nil, err
	}
	if state := allocatedCMBuf.state; state != RMR_OK {
		errorMessage := fmt.Sprintf("#rmrCgoApi.SendMsg - Failed to get allocated message. state: %v - %s", state, states[int(state)])
		return nil, errors.New(errorMessage)
	}
	defer C.rmr_free_msg(allocatedCMBuf)

	currCMBuf := C.rmr_rcv_msg(ctx.RmrCtx, allocatedCMBuf)
	if state := currCMBuf.state; state != RMR_OK {
		errorMessage := fmt.Sprintf("#rmrCgoApi.RecvMsg - Failed to receive message. state: %v - %s", state, states[int(state)])
		return nil, errors.New(errorMessage)
	}

	return convertToMBuf(currCMBuf)
}

func (ctx *Context) IsReady() bool {
	return int(C.rmr_ready(ctx.RmrCtx)) != 0
}

func (ctx *Context) Close() {
	C.rmr_close(ctx.RmrCtx)
}
