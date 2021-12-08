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


package rmrCgo

// #cgo LDFLAGS: -L/usr/local/lib -lrmr_si
// #include <rmr/rmr.h>
// #include <rmr/RIC_message_types.h>
// #include <stdlib.h>
import "C"
import (
	"e2mgr/logger"
	"fmt"
	"unsafe"
)

func NewMBuf(mType int, len int, meid string, payload *[]byte, xAction *[]byte, msgSrc unsafe.Pointer) *MBuf {
	return &MBuf{
		mType,
		len,
		meid,
		payload,
		xAction,
		msgSrc,
	}
}

func NewContext(maxMsgSize int, flags int, ctx unsafe.Pointer, logger *logger.Logger) *Context {
	return &Context{
		MaxMsgSize: maxMsgSize,
		Flags:      flags,
		RmrCtx:     ctx,
		Logger:     logger,
	}
}

//TODO: consider declaring using its own type
const (
	// messages
	RIC_X2_SETUP_REQ                     = C.RIC_X2_SETUP_REQ
	RIC_X2_SETUP_RESP                    = C.RIC_X2_SETUP_RESP
	RIC_X2_SETUP_FAILURE                 = C.RIC_X2_SETUP_FAILURE
	RIC_ENDC_X2_SETUP_REQ                = C.RIC_ENDC_X2_SETUP_REQ
	RIC_ENDC_X2_SETUP_RESP               = C.RIC_ENDC_X2_SETUP_RESP
	RIC_ENDC_X2_SETUP_FAILURE            = C.RIC_ENDC_X2_SETUP_FAILURE
	RIC_SCTP_CONNECTION_FAILURE          = C.RIC_SCTP_CONNECTION_FAILURE
	RIC_ENB_LOAD_INFORMATION             = C.RIC_ENB_LOAD_INFORMATION
	RIC_ENB_CONF_UPDATE                  = C.RIC_ENB_CONF_UPDATE
	RIC_ENB_CONFIGURATION_UPDATE_ACK     = C.RIC_ENB_CONF_UPDATE_ACK
	RIC_ENB_CONFIGURATION_UPDATE_FAILURE = C.RIC_ENB_CONF_UPDATE_FAILURE
	RIC_ENDC_CONF_UPDATE                 = C.RIC_ENDC_CONF_UPDATE
	RIC_ENDC_CONF_UPDATE_ACK             = C.RIC_ENDC_CONF_UPDATE_ACK
	RIC_ENDC_CONF_UPDATE_FAILURE         = C.RIC_ENDC_CONF_UPDATE_FAILURE
	RIC_SCTP_CLEAR_ALL                   = C.RIC_SCTP_CLEAR_ALL
	RIC_X2_RESET_RESP                    = C.RIC_X2_RESET_RESP
	RIC_X2_RESET                         = C.RIC_X2_RESET
	RIC_E2_TERM_INIT 					 = C.E2_TERM_INIT
	RAN_CONNECTED						 = C.RAN_CONNECTED
	RAN_RESTARTED						 = C.RAN_RESTARTED
	RAN_RECONFIGURED					 = C.RAN_RECONFIGURED
	E2_TERM_KEEP_ALIVE_REQ				 = C.E2_TERM_KEEP_ALIVE_REQ
	E2_TERM_KEEP_ALIVE_RESP				 = C.E2_TERM_KEEP_ALIVE_RESP
	RIC_E2_SETUP_REQ					 = C.RIC_E2_SETUP_REQ
	RIC_E2_SETUP_RESP                    = C.RIC_E2_SETUP_RESP
	RIC_E2_SETUP_FAILURE                 = C.RIC_E2_SETUP_FAILURE
)

const (
	RMR_MAX_XACTION_LEN = int(C.RMR_MAX_XID)
	RMR_MAX_MEID_LEN    = int(C.RMR_MAX_MEID)
	RMR_MAX_SRC_LEN			= int(C.RMR_MAX_SRC)

	//states
	RMR_OK             = C.RMR_OK
	RMR_ERR_BADARG     = C.RMR_ERR_BADARG
	RMR_ERR_NOENDPT    = C.RMR_ERR_NOENDPT
	RMR_ERR_EMPTY      = C.RMR_ERR_EMPTY
	RMR_ERR_NOHDR      = C.RMR_ERR_NOHDR
	RMR_ERR_SENDFAILED = C.RMR_ERR_SENDFAILED
	RMR_ERR_CALLFAILED = C.RMR_ERR_CALLFAILED
	RMR_ERR_NOWHOPEN   = C.RMR_ERR_NOWHOPEN
	RMR_ERR_WHID       = C.RMR_ERR_WHID
	RMR_ERR_OVERFLOW   = C.RMR_ERR_OVERFLOW
	RMR_ERR_RETRY      = C.RMR_ERR_RETRY
	RMR_ERR_RCVFAILED  = C.RMR_ERR_RCVFAILED
	RMR_ERR_TIMEOUT    = C.RMR_ERR_TIMEOUT
	RMR_ERR_UNSET      = C.RMR_ERR_UNSET
	RMR_ERR_TRUNC      = C.RMR_ERR_TRUNC
	RMR_ERR_INITFAILED = C.RMR_ERR_INITFAILED
)

var states = map[int]string{
	RMR_OK:             "state is good",
	RMR_ERR_BADARG:     "argument passd to function was unusable",
	RMR_ERR_NOENDPT:    "send/call could not find an endpoint based on msg type",
	RMR_ERR_EMPTY:      "msg received had no payload; attempt to send an empty message",
	RMR_ERR_NOHDR:      "message didn't contain a valid header",
	RMR_ERR_SENDFAILED: "send failed; errno has nano reason",
	RMR_ERR_CALLFAILED: "unable to send call() message",
	RMR_ERR_NOWHOPEN:   "no wormholes are open",
	RMR_ERR_WHID:       "wormhole id was invalid",
	RMR_ERR_OVERFLOW:   "operation would have busted through a buffer/field size",
	RMR_ERR_RETRY:      "request (send/call/rts) failed, but caller should retry (EAGAIN for wrappers)",
	RMR_ERR_RCVFAILED:  "receive failed (hard error)",
	RMR_ERR_TIMEOUT:    "message processing call timed out",
	RMR_ERR_UNSET:      "the message hasn't been populated with a transport buffer",
	RMR_ERR_TRUNC:      "received message likely truncated",
	RMR_ERR_INITFAILED: "initialisation of something (probably message) failed",
}

type MBuf struct {
	MType   int
	Len     int
	Meid    string //Managed entity id (RAN name)
	Payload *[]byte
	XAction *[]byte
	msgSrc  unsafe.Pointer
}

func (m MBuf) String() string {
	return fmt.Sprintf("{ MType: %d, Len: %d, Meid: %q, Xaction: %q, Payload: [%x] }", m.MType, m.Len, m.Meid, m.XAction, m.Payload)
}

func (m MBuf) GetMsgSrc() unsafe.Pointer {
	return m.msgSrc
}

type Context struct {
	MaxMsgSize int
	Flags      int
	RmrCtx     unsafe.Pointer
	Logger     *logger.Logger
}

type RmrMessenger interface {
	Init(port string, maxMsgSize int, flags int, logger *logger.Logger) RmrMessenger
	SendMsg(msg *MBuf, printLogs bool) (*MBuf, error)
	WhSendMsg(msg *MBuf, printLogs bool) (*MBuf, error)
	RecvMsg() (*MBuf, error)
	IsReady() bool
	Close()
}
