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


package converters

// #cgo CFLAGS: -I../3rdparty/asn1codec/inc/  -I../3rdparty/asn1codec/e2ap_engine/
// #cgo LDFLAGS: -L ../3rdparty/asn1codec/lib/ -L../3rdparty/asn1codec/e2ap_engine/ -le2ap_codec -lasncodec
// #include <asn1codec_utils.h>
// #include <x2setup_response_wrapper.h>
import "C"
import (
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"unsafe"
)

const (
	maxCellinengNB     = 16384
	maxofNRNeighbours  = 1024
	maxnoofNrCellBands = 32
)

type EndcSetupResponseConverter struct {
	logger *logger.Logger
}

type IEndcSetupResponseConverter interface {
	UnpackEndcSetupResponseAndExtract(packedBuf []byte) (*entities.GlobalNbId, *entities.Gnb, error)
}

func NewEndcSetupResponseConverter(logger *logger.Logger) *EndcSetupResponseConverter {
	return &EndcSetupResponseConverter{
		logger: logger,
	}
}

func getNRFreqInfo(freqInfo C.NRFreqInfo_t) (*entities.NrFrequencyInfo, error) {
	var info *entities.NrFrequencyInfo
	info = &entities.NrFrequencyInfo{NrArFcn: uint64(freqInfo.nRARFCN)}

	if freqInfo.sULInformation != nil {
		info.SulInformation = &entities.NrFrequencyInfo_SulInformation{SulArFcn: uint64((*C.SULInformation_t)(freqInfo.sULInformation).sUL_ARFCN)}

		if value, err := getNR_TxBW((*C.SULInformation_t)(freqInfo.sULInformation).sUL_TxBW); err == nil {
			info.SulInformation.SulTransmissionBandwidth = value
		} else {
			return nil, err
		}
	}

	if freqInfo.freqBandListNr.list.count > 0 && freqInfo.freqBandListNr.list.count <= maxnoofNrCellBands {
		count := int(freqInfo.freqBandListNr.list.count)
		freqBandListNr_slice := (*[1 << 30]*C.FreqBandNrItem_t)(unsafe.Pointer(freqInfo.freqBandListNr.list.array))[:count:count]
		for _, freqBandNrItem := range freqBandListNr_slice {
			frequencyBand := &entities.FrequencyBandItem{NrFrequencyBand: uint32(freqBandNrItem.freqBandIndicatorNr)}

			if freqBandNrItem.supportedSULBandList.list.count > 0 && freqBandNrItem.supportedSULBandList.list.count <= maxnoofNrCellBands {
				count := int(freqBandNrItem.supportedSULBandList.list.count)
				supportedSULBandList_slice := (*[1 << 30]*C.SupportedSULFreqBandItem_t)(unsafe.Pointer(freqBandNrItem.supportedSULBandList.list.array))[:count:count]
				for _, supportedSULFreqBandItem := range supportedSULBandList_slice {
					frequencyBand.SupportedSulBands = append(frequencyBand.SupportedSulBands, uint32(supportedSULFreqBandItem.freqBandIndicatorNr))
				}
			}

			info.FrequencyBands = append(info.FrequencyBands, frequencyBand)
		}
	}

	return info, nil
}

func getNR_TxBW(txBW C.NR_TxBW_t) (*entities.NrTransmissionBandwidth, error) {
	var bw *entities.NrTransmissionBandwidth

	bw = &entities.NrTransmissionBandwidth{Nrscs: entities.Nrscs(1 + int64(txBW.nRSCS))}
	bw.Ncnrb = entities.Ncnrb(1 + int64(txBW.nRNRB))

	return bw, nil
}

func getnrModeInfoFDDInfo(fdd *C.FDD_InfoServedNRCell_Information_t) (*entities.ServedNRCellInformation_ChoiceNRMode_FddInfo, error) {
	var fddInfo *entities.ServedNRCellInformation_ChoiceNRMode_FddInfo

	if info, err := getNRFreqInfo(fdd.ul_NRFreqInfo); err == nil {
		fddInfo = &entities.ServedNRCellInformation_ChoiceNRMode_FddInfo{UlFreqInfo: info}
	} else {
		return nil, err
	}

	if info, err := getNRFreqInfo(fdd.dl_NRFreqInfo); err == nil {
		fddInfo.DlFreqInfo = info
	} else {
		return nil, err
	}

	if bw, err := getNR_TxBW(fdd.ul_NR_TxBW); err == nil {
		fddInfo.UlTransmissionBandwidth = bw
	} else {
		return nil, err
	}

	if bw, err := getNR_TxBW(fdd.dl_NR_TxBW); err == nil {
		fddInfo.DlTransmissionBandwidth = bw
	} else {
		return nil, err
	}

	return fddInfo, nil
}

func getnrModeInfoTDDInfo(tdd *C.TDD_InfoServedNRCell_Information_t) (*entities.ServedNRCellInformation_ChoiceNRMode_TddInfo, error) {
	var tddInfo *entities.ServedNRCellInformation_ChoiceNRMode_TddInfo

	if info, err := getNRFreqInfo(tdd.nRFreqInfo); err == nil {
		tddInfo = &entities.ServedNRCellInformation_ChoiceNRMode_TddInfo{NrFreqInfo: info}
	} else {
		return nil, err

	}

	if bw, err := getNR_TxBW(tdd.nR_TxBW); err == nil {
		tddInfo.TransmissionBandwidth = bw
	} else {
		return nil, err
	}

	return tddInfo, nil
}

func getNRNeighbourInformation_ChoiceNRMode_FDDInfo(fdd *C.FDD_InfoNeighbourServedNRCell_Information_t) (*entities.NrNeighbourInformation_ChoiceNRMode_FddInfo, error) {
	var fddInfo *entities.NrNeighbourInformation_ChoiceNRMode_FddInfo

	if info, err := getNRFreqInfo(fdd.ul_NRFreqInfo); err == nil {
		fddInfo = &entities.NrNeighbourInformation_ChoiceNRMode_FddInfo{UlarFcnFreqInfo: info}
	} else {
		return nil, err
	}

	if info, err := getNRFreqInfo(fdd.dl_NRFreqInfo); err == nil {
		fddInfo.DlarFcnFreqInfo = info
	} else {
		return nil, err
	}

	return fddInfo, nil
}
func getNRNeighbourInformation_ChoiceNRMode_TDDInfo(tdd *C.TDD_InfoNeighbourServedNRCell_Information_t) (*entities.NrNeighbourInformation_ChoiceNRMode_TddInfo, error) {
	var tddInfo *entities.NrNeighbourInformation_ChoiceNRMode_TddInfo

	if info, err := getNRFreqInfo(tdd.nRFreqInfo); err == nil {
		tddInfo = &entities.NrNeighbourInformation_ChoiceNRMode_TddInfo{ArFcnNrFreqInfo: info}
	} else {
		return nil, err
	}

	return tddInfo, nil
}

func getnRNeighbourInfo(neighbour_Information *C.NRNeighbour_Information_t) ([]*entities.NrNeighbourInformation, error) {
	var neighbours []*entities.NrNeighbourInformation

	if neighbour_Information != nil && neighbour_Information.list.count > 0 && neighbour_Information.list.count <= maxofNRNeighbours {
		count := int(neighbour_Information.list.count)
		neighbour_Information_slice := (*[1 << 30]*C.NRNeighbour_Information__Member)(unsafe.Pointer(neighbour_Information.list.array))[:count:count]
		for _, member := range neighbour_Information_slice {
			info := &entities.NrNeighbourInformation{NrPci: uint32(member.nrpCI)}

			//pLMN_Identity:nRcellIdentifier
			plmnId := C.GoBytes(unsafe.Pointer(member.nrCellID.pLMN_Identity.buf), C.int(member.nrCellID.pLMN_Identity.size))
			nRcellIdentifier := C.GoBytes(unsafe.Pointer(member.nrCellID.nRcellIdentifier.buf), C.int(member.nrCellID.nRcellIdentifier.size))
			info.NrCgi = fmt.Sprintf("%02x:%02x", plmnId, nRcellIdentifier)

			if member.fiveGS_TAC != nil {
				info.Stac5G = fmt.Sprintf("%02x", C.GoBytes(unsafe.Pointer(member.fiveGS_TAC.buf), C.int(member.fiveGS_TAC.size)))

			}

			if member.configured_TAC != nil {
				info.ConfiguredStac = fmt.Sprintf("%02x", C.GoBytes(unsafe.Pointer(member.configured_TAC.buf), C.int(member.configured_TAC.size)))
			}
			switch member.nRNeighbourModeInfo.present {
			case C.NRNeighbour_Information__Member__nRNeighbourModeInfo_PR_fdd:
				if fdd, err := getNRNeighbourInformation_ChoiceNRMode_FDDInfo(*(**C.FDD_InfoNeighbourServedNRCell_Information_t)(unsafe.Pointer(&member.nRNeighbourModeInfo.choice[0]))); fdd != nil && err == nil {
					info.ChoiceNrMode, info.NrMode = &entities.NrNeighbourInformation_ChoiceNRMode{Fdd: fdd}, entities.Nr_FDD
				}

			case C.NRNeighbour_Information__Member__nRNeighbourModeInfo_PR_tdd:
				if tdd, err := getNRNeighbourInformation_ChoiceNRMode_TDDInfo(*(**C.TDD_InfoNeighbourServedNRCell_Information_t)(unsafe.Pointer(&member.nRNeighbourModeInfo.choice[0]))); tdd != nil && err == nil {
					info.ChoiceNrMode, info.NrMode = &entities.NrNeighbourInformation_ChoiceNRMode{Tdd: tdd}, entities.Nr_TDD
				}
			}
			neighbours = append(neighbours, info)
		}

	}

	return neighbours, nil
}

func getServedNRCells(servedNRcellsManagementList *C.ServedNRcellsENDCX2ManagementList_t) ([]*entities.ServedNRCell, error) {
	var servedNRCells []*entities.ServedNRCell

	if servedNRcellsManagementList != nil && servedNRcellsManagementList.list.count > 0 && servedNRcellsManagementList.list.count <= maxCellinengNB {
		count := int(servedNRcellsManagementList.list.count)
		servedNRcellsENDCX2ManagementList__Member_slice := (*[1 << 30]*C.ServedNRcellsENDCX2ManagementList__Member)(unsafe.Pointer(servedNRcellsManagementList.list.array))[:count:count]
		for _, servedNRcellsENDCX2ManagementList__Member := range servedNRcellsENDCX2ManagementList__Member_slice {
			servedNRCellInfo := servedNRcellsENDCX2ManagementList__Member.servedNRCellInfo
			servedNRCell := &entities.ServedNRCell{ServedNrCellInformation: &entities.ServedNRCellInformation{NrPci: uint32(servedNRCellInfo.nrpCI)}}

			//pLMN_Identity:nRcellIdentifier
			plmnId := C.GoBytes(unsafe.Pointer(servedNRCellInfo.nrCellID.pLMN_Identity.buf), C.int(servedNRCellInfo.nrCellID.pLMN_Identity.size))
			nRcellIdentifier := C.GoBytes(unsafe.Pointer(servedNRCellInfo.nrCellID.nRcellIdentifier.buf), C.int(servedNRCellInfo.nrCellID.nRcellIdentifier.size))
			servedNRCell.ServedNrCellInformation.CellId = fmt.Sprintf("%02x:%02x", plmnId, nRcellIdentifier)

			if servedNRCellInfo.fiveGS_TAC != nil {
				servedNRCell.ServedNrCellInformation.Stac5G = fmt.Sprintf("%02x", C.GoBytes(unsafe.Pointer(servedNRCellInfo.fiveGS_TAC.buf), C.int(servedNRCellInfo.fiveGS_TAC.size)))
			}

			if servedNRCellInfo.configured_TAC != nil {
				servedNRCell.ServedNrCellInformation.ConfiguredStac = fmt.Sprintf("%02x", C.GoBytes(unsafe.Pointer(servedNRCellInfo.configured_TAC.buf), C.int(servedNRCellInfo.configured_TAC.size)))
			}

			if servedNRCellInfo.broadcastPLMNs.list.count > 0 && servedNRCellInfo.broadcastPLMNs.list.count <= maxnoofBPLMNs {
				count := int(servedNRCellInfo.broadcastPLMNs.list.count)
				pLMN_Identity_slice := (*[1 << 30]*C.PLMN_Identity_t)(unsafe.Pointer(servedNRCellInfo.broadcastPLMNs.list.array))[:count:count]
				for _, pLMN_Identity := range pLMN_Identity_slice {
					servedNRCell.ServedNrCellInformation.ServedPlmns = append(servedNRCell.ServedNrCellInformation.ServedPlmns, fmt.Sprintf("%02x", C.GoBytes(unsafe.Pointer(pLMN_Identity.buf), C.int(pLMN_Identity.size))))
				}
			}
			switch servedNRCellInfo.nrModeInfo.present {
			case C.ServedNRCell_Information__nrModeInfo_PR_fdd:
				if fdd, err := getnrModeInfoFDDInfo(*(**C.FDD_InfoServedNRCell_Information_t)(unsafe.Pointer(&servedNRCellInfo.nrModeInfo.choice[0]))); fdd != nil && err == nil {
					servedNRCell.ServedNrCellInformation.ChoiceNrMode, servedNRCell.ServedNrCellInformation.NrMode = &entities.ServedNRCellInformation_ChoiceNRMode{Fdd: fdd}, entities.Nr_FDD
				} else {
					return nil, err
				}
			case C.ServedNRCell_Information__nrModeInfo_PR_tdd:
				if tdd, err := getnrModeInfoTDDInfo(*(**C.TDD_InfoServedNRCell_Information_t)(unsafe.Pointer(&servedNRCellInfo.nrModeInfo.choice[0]))); tdd != nil && err == nil {
					servedNRCell.ServedNrCellInformation.ChoiceNrMode, servedNRCell.ServedNrCellInformation.NrMode = &entities.ServedNRCellInformation_ChoiceNRMode{Tdd: tdd}, entities.Nr_TDD
				} else {
					return nil, err
				}
			}

			neighbours, err := getnRNeighbourInfo(servedNRcellsENDCX2ManagementList__Member.nRNeighbourInfo)
			if err != nil {
				return nil, err
			}
			servedNRCell.NrNeighbourInfos = neighbours

			servedNRCells = append(servedNRCells, servedNRCell)
		}
	}

	return servedNRCells, nil
}

// Populate  the GNB structure with data from the pdu
// Return the GNB and the associated key which can later be used to retrieve the GNB from the database.

func endcX2SetupResponseToProtobuf(pdu *C.E2AP_PDU_t) (*entities.GlobalNbId, *entities.Gnb, error) {

	var gnb *entities.Gnb
	var globalNbId *entities.GlobalNbId

	if pdu.present == C.E2AP_PDU_PR_successfulOutcome {
		//dereference a union of pointers (C union is represented as a byte array with the size of the largest member)
		successfulOutcome := *(**C.SuccessfulOutcome_t)(unsafe.Pointer(&pdu.choice[0]))
		if successfulOutcome != nil && successfulOutcome.value.present == C.SuccessfulOutcome__value_PR_ENDCX2SetupResponse {
			endcX2SetupResponse := (*C.ENDCX2SetupResponse_t)(unsafe.Pointer(&successfulOutcome.value.choice[0]))
			if endcX2SetupResponse != nil && endcX2SetupResponse.protocolIEs.list.count > 0 {
				count := int(endcX2SetupResponse.protocolIEs.list.count)
				endcX2SetupResponse_IEs_slice := (*[1 << 30]*C.ENDCX2SetupResponse_IEs_t)(unsafe.Pointer(endcX2SetupResponse.protocolIEs.list.array))[:count:count]
				for _, endcX2SetupResponse_IE := range endcX2SetupResponse_IEs_slice {
					if endcX2SetupResponse_IE.value.present == C.ENDCX2SetupResponse_IEs__value_PR_RespondingNodeType_EndcX2Setup {
						respondingNodeType := (*C.RespondingNodeType_EndcX2Setup_t)(unsafe.Pointer(&endcX2SetupResponse_IE.value.choice[0]))
						switch respondingNodeType.present {
						case C.RespondingNodeType_EndcX2Setup_PR_respond_en_gNB:
							en_gNB_ENDCX2SetupReqAckIEs_Container := *(**C.ProtocolIE_Container_119P89_t)(unsafe.Pointer(&respondingNodeType.choice[0]))
							if en_gNB_ENDCX2SetupReqAckIEs_Container != nil && en_gNB_ENDCX2SetupReqAckIEs_Container.list.count > 0 {
								count := int(en_gNB_ENDCX2SetupReqAckIEs_Container.list.count)
								en_gNB_ENDCX2SetupReqAckIEs_slice := (*[1 << 30]*C.En_gNB_ENDCX2SetupReqAckIEs_t)(unsafe.Pointer(en_gNB_ENDCX2SetupReqAckIEs_Container.list.array))[:count:count]
								for _, en_gNB_ENDCX2SetupReqAckIE := range en_gNB_ENDCX2SetupReqAckIEs_slice {
									switch en_gNB_ENDCX2SetupReqAckIE.value.present {
									case C.En_gNB_ENDCX2SetupReqAckIEs__value_PR_GlobalGNB_ID:
										globalGNB_ID := (*C.GlobalGNB_ID_t)(unsafe.Pointer(&en_gNB_ENDCX2SetupReqAckIE.value.choice[0]))
										plmnId := C.GoBytes(unsafe.Pointer(globalGNB_ID.pLMN_Identity.buf), C.int(globalGNB_ID.pLMN_Identity.size))
										if globalGNB_ID.gNB_ID.present == C.GNB_ID_PR_gNB_ID {
											gnbIdAsBitString := (*C.BIT_STRING_t)(unsafe.Pointer(&globalGNB_ID.gNB_ID.choice[0]))
											globalNbId = &entities.GlobalNbId{}
											globalNbId.NbId = fmt.Sprintf("%02x", C.GoBytes(unsafe.Pointer(gnbIdAsBitString.buf), C.int(gnbIdAsBitString.size)))
											globalNbId.PlmnId = fmt.Sprintf("%02x", plmnId)
										}
									case C.En_gNB_ENDCX2SetupReqAckIEs__value_PR_ServedNRcellsENDCX2ManagementList:
										servedCells, err := getServedNRCells((*C.ServedNRcellsENDCX2ManagementList_t)(unsafe.Pointer(&en_gNB_ENDCX2SetupReqAckIE.value.choice[0])))
										if err != nil {
											return globalNbId, nil, err
										}
										gnb = &entities.Gnb{}
										gnb.ServedNrCells = servedCells
									}
								}
							}
						case C.RespondingNodeType_EndcX2Setup_PR_respond_eNB:
							/*ignored*/
						}
					}
				}
			}
		}
	}

	return globalNbId, gnb, nil
}

func (c *EndcSetupResponseConverter) UnpackEndcSetupResponseAndExtract(packedBuf []byte) (*entities.GlobalNbId, *entities.Gnb, error) {
	pdu, err := UnpackX2apPdu(c.logger, e2pdus.MaxAsn1CodecAllocationBufferSize, len(packedBuf), packedBuf, e2pdus.MaxAsn1CodecMessageBufferSize)

	if err != nil {
		return nil, nil, err
	}

	defer C.delete_pdu(pdu)
	return endcX2SetupResponseToProtobuf(pdu)
}
