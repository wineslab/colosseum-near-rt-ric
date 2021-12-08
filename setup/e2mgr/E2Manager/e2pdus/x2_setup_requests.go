/*
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
 */

/*
 * This source code is part of the near-RT RIC (RAN Intelligent Controller)
 * platform project (RICP).
 */

package e2pdus

// #cgo CFLAGS: -I../3rdparty/asn1codec/inc/ -I../3rdparty/asn1codec/e2ap_engine/
// #cgo LDFLAGS: -L ../3rdparty/asn1codec/lib/ -L../3rdparty/asn1codec/e2ap_engine/ -le2ap_codec -lasncodec
// #include <asn1codec_utils.h>
// #include <x2setup_request_wrapper.h>
import "C"
import (
	"fmt"
	"github.com/pkg/errors"
	"unsafe"
)

const (
	EnvRicId          = "RIC_ID"
	ShortMacro_eNB_ID = 18
	Macro_eNB_ID      = 20
	LongMacro_eNB_ID  = 21
	Home_eNB_ID       = 28
)

var PackedEndcX2setupRequest []byte
var PackedX2setupRequest []byte
var PackedEndcX2setupRequestAsString string
var PackedX2setupRequestAsString string

/*The Ric Id is the combination of pLMNId and ENBId*/
var pLMNId []byte
var eNBId []byte
var eNBIdBitqty uint
var ricFlag = []byte{0xbb, 0xbc, 0xcc} /*pLMNId [3]bytes*/

func preparePackedEndcX2SetupRequest(maxAsn1PackedBufferSize int, maxAsn1CodecMessageBufferSize int, pLMNId []byte, eNB_Id []byte /*18, 20, 21, 28 bits length*/, bitqty uint, ricFlag []byte) ([]byte, string, error) {
	packedBuf := make([]byte, maxAsn1PackedBufferSize)
	errBuf := make([]C.char, maxAsn1CodecMessageBufferSize)
	packedBufSize := C.ulong(len(packedBuf))
	pduAsString := ""

	if !C.build_pack_endc_x2setup_request(
		(*C.uchar)(unsafe.Pointer(&pLMNId[0])) /*pLMN_Identity*/,
		(*C.uchar)(unsafe.Pointer(&eNB_Id[0])),
		C.uint(bitqty),
		(*C.uchar)(unsafe.Pointer(&ricFlag[0])) /*pLMN_Identity*/,
		&packedBufSize,
		(*C.uchar)(unsafe.Pointer(&packedBuf[0])),
		C.ulong(len(errBuf)),
		&errBuf[0]) {
		return nil, "", errors.New(fmt.Sprintf("packing error: %s", C.GoString(&errBuf[0])))
	}

	pdu := C.new_pdu(C.size_t(1)) //TODO: change signature
	defer C.delete_pdu(pdu)
	if C.per_unpack_pdu(pdu, packedBufSize, (*C.uchar)(unsafe.Pointer(&packedBuf[0])), C.size_t(len(errBuf)), &errBuf[0]) {
		C.asn1_pdu_printer(pdu, C.size_t(len(errBuf)), &errBuf[0])
		pduAsString = C.GoString(&errBuf[0])
	}

	return packedBuf[:packedBufSize], pduAsString, nil
}

func preparePackedX2SetupRequest(maxAsn1PackedBufferSize int, maxAsn1CodecMessageBufferSize int, pLMNId []byte, eNB_Id []byte /*18, 20, 21, 28 bits length*/, bitqty uint, ricFlag []byte) ([]byte, string, error) {
	packedBuf := make([]byte, maxAsn1PackedBufferSize)
	errBuf := make([]C.char, maxAsn1CodecMessageBufferSize)
	packedBufSize := C.ulong(len(packedBuf))
	pduAsString := ""

	if !C.build_pack_x2setup_request(
		(*C.uchar)(unsafe.Pointer(&pLMNId[0])) /*pLMN_Identity*/,
		(*C.uchar)(unsafe.Pointer(&eNB_Id[0])),
		C.uint(bitqty),
		(*C.uchar)(unsafe.Pointer(&ricFlag[0])) /*pLMN_Identity*/,
		&packedBufSize,
		(*C.uchar)(unsafe.Pointer(&packedBuf[0])),
		C.ulong(len(errBuf)),
		&errBuf[0]) {
		return nil, "", errors.New(fmt.Sprintf("packing error: %s", C.GoString(&errBuf[0])))
	}

	pdu := C.new_pdu(C.size_t(1)) //TODO: change signature
	defer C.delete_pdu(pdu)
	if C.per_unpack_pdu(pdu, packedBufSize, (*C.uchar)(unsafe.Pointer(&packedBuf[0])), C.size_t(len(errBuf)), &errBuf[0]) {
		C.asn1_pdu_printer(pdu, C.size_t(len(errBuf)), &errBuf[0])
		pduAsString = C.GoString(&errBuf[0])
	}
	return packedBuf[:packedBufSize], pduAsString, nil
}

//Expected value in RIC_ID = pLMN_Identity-eNB_ID/<eNB_ID size in bits>
//<6 hex digits>-<6 or 8 hex digits>/<18|20|21|28>
//Each byte is represented by two hex digits, the value in the lowest byte of the eNB_ID must be assigned to the lowest bits
//For example, to get the value of ffffeab/28  the last byte must be 0x0b, not 0xb0 (-ffffea0b/28).
func parseRicID(ricId string) error {
	if _, err := fmt.Sscanf(ricId, "%6x-%8x/%2d", &pLMNId, &eNBId, &eNBIdBitqty); err != nil {
		return fmt.Errorf("unable to extract the value of %s: %s", EnvRicId, err)
	}

	if len(pLMNId) < 3 {
		return fmt.Errorf("invalid value for %s, len(pLMNId:%v) != 3", EnvRicId, pLMNId)
	}

	if len(eNBId) < 3 {
		return fmt.Errorf("invalid value for %s, len(eNBId:%v) != 3 or 4", EnvRicId, eNBId)
	}

	if eNBIdBitqty != ShortMacro_eNB_ID && eNBIdBitqty != Macro_eNB_ID && eNBIdBitqty != LongMacro_eNB_ID && eNBIdBitqty != Home_eNB_ID {
		return fmt.Errorf("invalid value for %s, eNBIdBitqty: %d", EnvRicId, eNBIdBitqty)
	}

	return nil
}

func init() {
	var err error
	ricId := "bbbccc-abcd0e/20"
	if err = parseRicID(ricId); err != nil {
		panic(err)
	}

	PackedEndcX2setupRequest, PackedEndcX2setupRequestAsString, err = preparePackedEndcX2SetupRequest(MaxAsn1PackedBufferSize, MaxAsn1CodecMessageBufferSize, pLMNId, eNBId, eNBIdBitqty, ricFlag)
	if err != nil {
		panic(err)
	}
	PackedX2setupRequest, PackedX2setupRequestAsString, err = preparePackedX2SetupRequest(MaxAsn1PackedBufferSize, MaxAsn1CodecMessageBufferSize, pLMNId, eNBId, eNBIdBitqty, ricFlag)
	if err != nil {
		panic(err)
	}
}
