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
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"unsafe"
)

/*
Allocates an mBuf and initialize it with the content of C.rmr_mbuf_t.
The xAction field is assigned a a value without trailing spaces.
*/
func convertToMBuf(m *C.rmr_mbuf_t) (*MBuf, error) {
	payloadArr := C.GoBytes(unsafe.Pointer(m.payload), C.int(m.len))
	xActionArr := C.GoBytes(unsafe.Pointer(m.xaction), C.int(RMR_MAX_XACTION_LEN))

	// Trim padding (space and 0)
	xActionStr := strings.TrimRight(string(xActionArr), "\040\000")
	xActionArr = []byte(xActionStr)

	mbuf := &MBuf{
		MType: int(m.mtype),
		Len:   int(m.len),
		//Payload: (*[]byte)(unsafe.Pointer(m.payload)),
		Payload: payloadArr,
		//XAction: (*[]byte)(unsafe.Pointer(m.xaction)),
		XAction: xActionArr,
	}

	meidBuf := make([]byte, RMR_MAX_MEID_LEN)
	if meidCstr := C.rmr_get_meid(m, (*C.uchar)(unsafe.Pointer(&meidBuf[0]))); meidCstr != nil {
		mbuf.Meid = strings.TrimRight(string(meidBuf), "\000")
	}

	return mbuf, nil
}

/*
Allocates an C.rmr_mbuf_t and initialize it with the content of mBuf.
The xAction field is padded with trailing spaces upto capacity
*/
func (ctx *Context) getAllocatedCRmrMBuf(mBuf *MBuf, maxMsgSize int) (cMBuf *C.rmr_mbuf_t, rc error) {
	var xActionBuf [RMR_MAX_XACTION_LEN]byte
	var meidBuf [RMR_MAX_MEID_LEN]byte

	cMBuf = C.rmr_alloc_msg(ctx.RmrCtx, C.int(maxMsgSize))
	cMBuf.mtype = C.int(mBuf.MType)
	cMBuf.len = C.int(mBuf.Len)

	payloadLen := len(mBuf.Payload)
	xActionLen := len(mBuf.XAction)

	copy(xActionBuf[:], mBuf.XAction)
	for i := xActionLen; i < RMR_MAX_XACTION_LEN; i++ {
		xActionBuf[i] = '\040' //space
	}

	// Add padding
	copy(meidBuf[:], mBuf.Meid)
	for i := len(mBuf.Meid); i < RMR_MAX_MEID_LEN; i++ {
		meidBuf[i] = 0
	}

	payloadArr := (*[1 << 30]byte)(unsafe.Pointer(cMBuf.payload))[:payloadLen:payloadLen]
	xActionArr := (*[1 << 30]byte)(unsafe.Pointer(cMBuf.xaction))[:RMR_MAX_XACTION_LEN:RMR_MAX_XACTION_LEN]

	err := binary.Read(bytes.NewReader(mBuf.Payload), binary.LittleEndian, payloadArr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("#rmrCgoUtils.getAllocatedCRmrMBuf - Failed to read payload to allocated RMR message buffer, %s", err))
	}
	err = binary.Read(bytes.NewReader(xActionBuf[:]), binary.LittleEndian, xActionArr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("#rmrCgoUtils.getAllocatedCRmrMBuf - Failed to read xAction data to allocated RMR message buffer, %s", err))
	}

	len := C.rmr_bytes2meid(cMBuf, (*C.uchar)(unsafe.Pointer(&meidBuf[0])), C.int(RMR_MAX_XACTION_LEN))
	if int(len) != RMR_MAX_MEID_LEN {
		return nil, errors.New(
			"#rmrCgoUtils.getAllocatedCRmrMBuf - Failed to copy meid data to allocated RMR message buffer")
	}
	return cMBuf, nil
}

func MessageIdToUint(id string) (msgId uint64, err error) {
	if len(id) == 0 {
		msgId, err = 0, nil
	} else {
		msgId, err = strconv.ParseUint(id, 10, 16)
	}
	return
}
