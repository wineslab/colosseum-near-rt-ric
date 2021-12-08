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
// #include <stdlib.h>
import "C"
import (
	"e2mgr/logger"
	"bytes"
	"encoding/binary"
	"strings"
	"unsafe"
)

func convertToMBuf(logger *logger.Logger, m *C.rmr_mbuf_t) *MBuf {
	payloadArr := C.GoBytes(unsafe.Pointer(m.payload),C.int(m.len))
	xActionArr := C.GoBytes(unsafe.Pointer(m.xaction),C.int(RMR_MAX_XACTION_LEN))

	// Trim padding (space and 0)
	xActionStr :=  strings.TrimRight(string(xActionArr),"\040\000")
	xActionArr = []byte(xActionStr)

	mbuf := &MBuf{
		MType: int(m.mtype),
		Len:   int(m.len),
		Payload: &payloadArr,
		XAction: &xActionArr,
		msgSrc: C.CBytes(make([]byte, RMR_MAX_SRC_LEN)),
	}

	C.rmr_get_src(m, (*C.uchar)(mbuf.msgSrc)) // Capture message source

	meidBuf := make([]byte, RMR_MAX_MEID_LEN)
	if meidCstr := C.rmr_get_meid(m, (*C.uchar)(unsafe.Pointer(&meidBuf[0]))); meidCstr != nil {
		mbuf.Meid =	strings.TrimRight(string(meidBuf), "\000")
	}

	return mbuf
}

func (ctx *Context) getAllocatedCRmrMBuf(logger *logger.Logger, mBuf *MBuf, maxMsgSize int) (cMBuf *C.rmr_mbuf_t) {
	var xActionBuf [RMR_MAX_XACTION_LEN]byte
	var meidBuf[RMR_MAX_MEID_LEN]byte

	cMBuf = C.rmr_alloc_msg(ctx.RmrCtx, C.int(maxMsgSize))
	cMBuf.mtype = C.int(mBuf.MType)
	cMBuf.len = C.int(mBuf.Len)

	payloadLen := len(*mBuf.Payload)
	xActionLen := len(*mBuf.XAction)

	//Add padding
	copy(xActionBuf[:], *mBuf.XAction)
	for i:= xActionLen; i < RMR_MAX_XACTION_LEN; i++{
		xActionBuf[i] = '\040' //space
	}

	// Add padding
	copy(meidBuf[:], mBuf.Meid)
	for i:= len(mBuf.Meid); i < RMR_MAX_MEID_LEN; i++{
		meidBuf[i] = 0
	}

	payloadArr := (*[1 << 30]byte)(unsafe.Pointer(cMBuf.payload))[:payloadLen:payloadLen]
	xActionArr := (*[1 << 30]byte)(unsafe.Pointer(cMBuf.xaction))[:RMR_MAX_XACTION_LEN:RMR_MAX_XACTION_LEN]

	err := binary.Read(bytes.NewReader(*mBuf.Payload), binary.LittleEndian, payloadArr)
	if err != nil {
		ctx.Logger.Errorf(
			"#rmrCgoUtils.getAllocatedCRmrMBuf - Failed to read payload to allocated RMR message buffer")
	}
	err = binary.Read(bytes.NewReader(xActionBuf[:]), binary.LittleEndian, xActionArr)
	if err != nil {
		ctx.Logger.Errorf(
			"#rmrCgoUtils.getAllocatedCRmrMBuf - Failed to read xAction data to allocated RMR message buffer")
	}
	len := C.rmr_bytes2meid(cMBuf,  (*C.uchar)(unsafe.Pointer(&meidBuf[0])), C.int(RMR_MAX_MEID_LEN))
	if int(len) != RMR_MAX_MEID_LEN {
		ctx.Logger.Errorf(
			"#rmrCgoUtils.getAllocatedCRmrMBuf - Failed to copy meid data to allocated RMR message buffer")
	}
	return cMBuf
}

//TODO: change to assert or return error
func (ctx *Context) checkContextInitialized() {
	if ctx.RmrCtx == nil {
		if ctx.Logger != nil {
			ctx.Logger.DPanicf("#rmrCgoUtils.checkContextInitialized - The RMR router has not been initialized")
		}
		panic("#rmrCgoUtils.checkContextInitialized - The RMR router has not been initialized. To initialize router please call Init() method")
	}
}
