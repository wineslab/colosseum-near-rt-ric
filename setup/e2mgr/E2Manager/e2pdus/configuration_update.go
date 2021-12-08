/*******************************************************************************
 *
 *   Copyright (c) 2019 AT&T Intellectual Property.
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 *
 *******************************************************************************/

/*
* This source code is part of the near-RT RIC (RAN Intelligent Controller)
* platform project (RICP).
*/

package e2pdus

// #cgo CFLAGS: -I../asn1codec/inc/  -I../asn1codec/e2ap_engine/
// #cgo LDFLAGS: -L ../asn1codec/lib/ -L../asn1codec/e2ap_engine/ -le2ap_codec -lasncodec
// #include <configuration_update_wrapper.h>
import "C"
import (
	"fmt"
	"unsafe"
)

var PackedEndcConfigurationUpdateFailure []byte
var PackedEndcConfigurationUpdateAck []byte
var PackedX2EnbConfigurationUpdateFailure []byte
var PackedX2EnbConfigurationUpdateAck []byte

func prepareEndcConfigurationUpdateFailurePDU(maxAsn1PackedBufferSize int, maxAsn1CodecMessageBufferSize int) error {

	packedBuffer := make([]C.uchar, maxAsn1PackedBufferSize)
	errorBuffer := make([]C.char, maxAsn1CodecMessageBufferSize)
	var payloadSize = C.ulong(maxAsn1PackedBufferSize)

	if status := C.build_pack_endc_configuration_update_failure(&payloadSize, &packedBuffer[0], C.ulong(maxAsn1CodecMessageBufferSize), &errorBuffer[0]); !status {
		return fmt.Errorf("#configuration_update.prepareEndcConfigurationUpdateFailurePDU - failed to build and pack the endc configuration update failure message %s ", C.GoString(&errorBuffer[0]))

	}
	PackedEndcConfigurationUpdateFailure = C.GoBytes(unsafe.Pointer(&packedBuffer[0]), C.int(payloadSize))

	return nil
}

func prepareX2EnbConfigurationUpdateFailurePDU(maxAsn1PackedBufferSize int, maxAsn1CodecMessageBufferSize int) error {

	packedBuffer := make([]C.uchar, maxAsn1PackedBufferSize)
	errorBuffer := make([]C.char, maxAsn1CodecMessageBufferSize)
	var payloadSize = C.ulong(maxAsn1PackedBufferSize)

	if status := C.build_pack_x2enb_configuration_update_failure(&payloadSize, &packedBuffer[0], C.ulong(maxAsn1CodecMessageBufferSize), &errorBuffer[0]); !status {
		return fmt.Errorf("#configuration_update.prepareX2EnbConfigurationUpdateFailurePDU - failed to build and pack the x2 configuration update failure message %s ", C.GoString(&errorBuffer[0]))

	}
	PackedX2EnbConfigurationUpdateFailure = C.GoBytes(unsafe.Pointer(&packedBuffer[0]), C.int(payloadSize))

	return nil
}

func prepareEndcConfigurationUpdateAckPDU(maxAsn1PackedBufferSize int, maxAsn1CodecMessageBufferSize int) error {

	packedBuffer := make([]C.uchar, maxAsn1PackedBufferSize)
	errorBuffer := make([]C.char, maxAsn1CodecMessageBufferSize)
	var payloadSize = C.ulong(maxAsn1PackedBufferSize)

	if status := C.build_pack_endc_configuration_update_ack(&payloadSize, &packedBuffer[0], C.ulong(maxAsn1CodecMessageBufferSize), &errorBuffer[0]); !status {
		return fmt.Errorf("#configuration_update.prepareEndcConfigurationUpdateAckPDU - failed to build and pack the endc configuration update ack message %s ", C.GoString(&errorBuffer[0]))

	}
	PackedEndcConfigurationUpdateAck = C.GoBytes(unsafe.Pointer(&packedBuffer[0]), C.int(payloadSize))

	return nil
}

func prepareX2EnbConfigurationUpdateAckPDU(maxAsn1PackedBufferSize int, maxAsn1CodecMessageBufferSize int) error {

	packedBuffer := make([]C.uchar, maxAsn1PackedBufferSize)
	errorBuffer := make([]C.char, maxAsn1CodecMessageBufferSize)
	var payloadSize = C.ulong(maxAsn1PackedBufferSize)

	if status := C.build_pack_x2enb_configuration_update_ack(&payloadSize, &packedBuffer[0], C.ulong(maxAsn1CodecMessageBufferSize), &errorBuffer[0]); !status {
		return fmt.Errorf("#configuration_update.prepareX2EnbConfigurationUpdateAckPDU - failed to build and pack the x2 configuration update ack message %s ", C.GoString(&errorBuffer[0]))

	}
	PackedX2EnbConfigurationUpdateAck = C.GoBytes(unsafe.Pointer(&packedBuffer[0]), C.int(payloadSize))

	return nil
}

func init() {
	if err := prepareEndcConfigurationUpdateFailurePDU(MaxAsn1PackedBufferSize, MaxAsn1CodecMessageBufferSize); err != nil {
		panic(err)
	}
	if err := prepareEndcConfigurationUpdateAckPDU(MaxAsn1PackedBufferSize, MaxAsn1CodecMessageBufferSize); err != nil {
		panic(err)
	}
	if err := prepareX2EnbConfigurationUpdateFailurePDU(MaxAsn1PackedBufferSize, MaxAsn1CodecMessageBufferSize); err != nil {
		panic(err)
	}
	if err := prepareX2EnbConfigurationUpdateAckPDU(MaxAsn1PackedBufferSize, MaxAsn1CodecMessageBufferSize); err != nil {
		panic(err)
	}
}
