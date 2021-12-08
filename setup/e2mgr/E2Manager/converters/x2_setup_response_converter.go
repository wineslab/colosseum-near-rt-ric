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

// #cgo CFLAGS: -I../3rdparty/asn1codec/inc/ -I../3rdparty/asn1codec/e2ap_engine/
// #cgo LDFLAGS: -L ../3rdparty/asn1codec/lib/ -L../3rdparty/asn1codec/e2ap_engine/ -le2ap_codec -lasncodec
// #include <asn1codec_utils.h>
// #include <x2setup_response_wrapper.h>
import "C"
import (
	"e2mgr/e2pdus"
	"e2mgr/logger"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"github.com/pkg/errors"
	"unsafe"
)

const (
	maxnoofMBSFN                     = 8
	maxnoofBPLMNs                    = 6
	maxCellineNB                     = 256
	maxnoofMBMSServiceAreaIdentities = 256
	maxnoofBands                     = 16
	maxPools                         = 16
	maxnoofNeighbours                = 512
)

type X2SetupResponseConverter struct {
	logger *logger.Logger
}

type IX2SetupResponseConverter interface {
	UnpackX2SetupResponseAndExtract(packedBuf []byte) (*entities.GlobalNbId, *entities.Enb, error)
}

func NewX2SetupResponseConverter(logger *logger.Logger) *X2SetupResponseConverter {
	return &X2SetupResponseConverter{
		logger: logger,
	}
}

// The following are possible values of a choice field, find which the pdu contains.
func getENB_ID_choice(eNB_ID C.ENB_ID_t) (entities.EnbType, []byte) {

	enbIdAsBitString := (*C.BIT_STRING_t)(unsafe.Pointer(&eNB_ID.choice[0]))
	switch eNB_ID.present {
	case C.ENB_ID_PR_macro_eNB_ID:
		return entities.EnbType_MACRO_ENB, C.GoBytes(unsafe.Pointer(enbIdAsBitString.buf), C.int(enbIdAsBitString.size))
	case C.ENB_ID_PR_home_eNB_ID:
		return entities.EnbType_HOME_ENB, C.GoBytes(unsafe.Pointer(enbIdAsBitString.buf), C.int(enbIdAsBitString.size))
	case C.ENB_ID_PR_short_Macro_eNB_ID:
		return entities.EnbType_SHORT_MACRO_ENB, C.GoBytes(unsafe.Pointer(enbIdAsBitString.buf), C.int(enbIdAsBitString.size))
	case C.ENB_ID_PR_long_Macro_eNB_ID:
		return entities.EnbType_LONG_MACRO_ENB, C.GoBytes(unsafe.Pointer(enbIdAsBitString.buf), C.int(enbIdAsBitString.size))
	}

	return entities.EnbType_UNKNOWN_ENB_TYPE, nil
}

func getFDDInfo(fdd *C.FDD_Info_t) (*entities.FddInfo, error) {
	var fddInfo *entities.FddInfo

	if fdd != nil {
		fddInfo = &entities.FddInfo{UlearFcn: uint32(fdd.uL_EARFCN)}
		fddInfo.DlearFcn = uint32(fdd.dL_EARFCN)
		fddInfo.UlTransmissionBandwidth = entities.TransmissionBandwidth(1 + fdd.uL_Transmission_Bandwidth)
		fddInfo.DlTransmissionBandwidth = entities.TransmissionBandwidth(1 + fdd.dL_Transmission_Bandwidth)

		extIEs := (*C.ProtocolExtensionContainer_170P145_t)(unsafe.Pointer(fdd.iE_Extensions))
		if extIEs != nil && extIEs.list.count > 0 {
			count := int(extIEs.list.count)
			extIEs_slice := (*[1 << 30]*C.FDD_Info_ExtIEs_t)(unsafe.Pointer(extIEs.list.array))[:count:count]
			for _, member := range extIEs_slice {
				switch member.extensionValue.present {
				case C.FDD_Info_ExtIEs__extensionValue_PR_EARFCNExtension:
					eARFCNExtension := (*C.EARFCNExtension_t)(unsafe.Pointer(&member.extensionValue.choice[0]))
					if member.id == C.ProtocolIE_ID_id_UL_EARFCNExtension {
						fddInfo.UlearFcn = uint32(*eARFCNExtension)
					}
					if member.id == C.ProtocolIE_ID_id_DL_EARFCNExtension {
						fddInfo.DlearFcn = uint32(*eARFCNExtension)
					}
				case C.FDD_Info_ExtIEs__extensionValue_PR_OffsetOfNbiotChannelNumberToEARFCN:
					/*ignored*/
				case C.FDD_Info_ExtIEs__extensionValue_PR_NRS_NSSS_PowerOffset:
					/*ignored*/
				case C.FDD_Info_ExtIEs__extensionValue_PR_NSSS_NumOccasionDifferentPrecoder:
					/*ignored*/
				}
			}
		}
	}

	return fddInfo, nil
}

func getSpecialSubframeInfo(info C.SpecialSubframe_Info_t) *entities.SpecialSubframeInfo {
	specialSubframeInfo := entities.SpecialSubframeInfo{}

	specialSubframeInfo.SpecialSubframePatterns = entities.SpecialSubframe_Patterns(1 + info.specialSubframePatterns)
	specialSubframeInfo.CyclicPrefixDl = entities.CyclicPrefix(1 + info.cyclicPrefixDL)
	specialSubframeInfo.CyclicPrefixUl = entities.CyclicPrefix(1 + info.cyclicPrefixUL)

	return &specialSubframeInfo
}

func getAdditionalSpecialSubframeInfo(info *C.AdditionalSpecialSubframe_Info_t) *entities.AdditionalSpecialSubframeInfo {
	additionalSpecialSubframeInfo := &entities.AdditionalSpecialSubframeInfo{AdditionalSpecialSubframePatterns: entities.AdditionalSpecialSubframe_Patterns(1 + info.additionalspecialSubframePatterns)}

	additionalSpecialSubframeInfo.CyclicPrefixDl = entities.CyclicPrefix(1 + info.cyclicPrefixDL)
	additionalSpecialSubframeInfo.CyclicPrefixUl = entities.CyclicPrefix(1 + info.cyclicPrefixUL)

	return additionalSpecialSubframeInfo
}

func getAdditionalSpecialSubframeExtensionInfo(info *C.AdditionalSpecialSubframeExtension_Info_t) *entities.AdditionalSpecialSubframeExtensionInfo {
	additionalSpecialSubframeExtensionInfo := &entities.AdditionalSpecialSubframeExtensionInfo{AdditionalSpecialSubframePatternsExtension: entities.AdditionalSpecialSubframePatterns_Extension(1 + info.additionalspecialSubframePatternsExtension)}

	additionalSpecialSubframeExtensionInfo.CyclicPrefixDl = entities.CyclicPrefix(1 + info.cyclicPrefixDL)
	additionalSpecialSubframeExtensionInfo.CyclicPrefixUl = entities.CyclicPrefix(1 + info.cyclicPrefixUL)

	return additionalSpecialSubframeExtensionInfo
}

func getTDDInfo(tdd *C.TDD_Info_t) (*entities.TddInfo, error) {
	var tddInfo *entities.TddInfo

	if tdd != nil {
		tddInfo = &entities.TddInfo{EarFcn: uint32(tdd.eARFCN)}
		tddInfo.TransmissionBandwidth = entities.TransmissionBandwidth(1 + tdd.transmission_Bandwidth)
		tddInfo.SubframeAssignment = entities.SubframeAssignment(1 + tdd.subframeAssignment)

		tddInfo.SpecialSubframeInfo = getSpecialSubframeInfo(tdd.specialSubframe_Info)

		extIEs := (*C.ProtocolExtensionContainer_170P206_t)(unsafe.Pointer(tdd.iE_Extensions))
		if extIEs != nil && extIEs.list.count > 0 {
			count := int(extIEs.list.count)
			extIEs_slice := (*[1 << 30]*C.TDD_Info_ExtIEs_t)(unsafe.Pointer(extIEs.list.array))[:count:count]
			for _, member := range extIEs_slice {
				switch member.extensionValue.present {
				case C.TDD_Info_ExtIEs__extensionValue_PR_AdditionalSpecialSubframe_Info:
					tddInfo.AdditionalSpecialSubframeInfo = getAdditionalSpecialSubframeInfo((*C.AdditionalSpecialSubframe_Info_t)(unsafe.Pointer(&member.extensionValue.choice[0])))
				case C.TDD_Info_ExtIEs__extensionValue_PR_EARFCNExtension:
					tddInfo.EarFcn = uint32(*(*C.EARFCNExtension_t)(unsafe.Pointer(&member.extensionValue.choice[0])))
				case C.TDD_Info_ExtIEs__extensionValue_PR_AdditionalSpecialSubframeExtension_Info:
					tddInfo.AdditionalSpecialSubframeExtensionInfo = getAdditionalSpecialSubframeExtensionInfo((*C.AdditionalSpecialSubframeExtension_Info_t)(unsafe.Pointer(&member.extensionValue.choice[0])))
				}
			}
		}
	}

	return tddInfo, nil
}

// The following are possible values of a choice field, find which the pdu contains.
func getSubframeAllocation_choice(subframeAllocation C.SubframeAllocation_t) (entities.SubframeAllocationType, string, error) {

	switch subframeAllocation.present {
	case C.SubframeAllocation_PR_oneframe:
		frameAllocation := (*C.Oneframe_t)(unsafe.Pointer(&subframeAllocation.choice[0]))
		return entities.SubframeAllocationType_ONE_FRAME, fmt.Sprintf("%02x", C.GoBytes(unsafe.Pointer(frameAllocation.buf), C.int(frameAllocation.size))), nil
	case C.SubframeAllocation_PR_fourframes:
		frameAllocation := (*C.Fourframes_t)(unsafe.Pointer(&subframeAllocation.choice[0]))
		return entities.SubframeAllocationType_FOUR_FRAME, fmt.Sprintf("%02x", C.GoBytes(unsafe.Pointer(frameAllocation.buf), C.int(frameAllocation.size))), nil
	}
	return entities.SubframeAllocationType_UNKNOWN_SUBFRAME_ALLOCATION_TYPE, "", errors.Errorf("unexpected subframe allocation choice")
}

func getMBSFN_Subframe_Infolist(mBSFN_Subframe_Infolist *C.MBSFN_Subframe_Infolist_t) ([]*entities.MbsfnSubframe, error) {
	var mBSFNSubframes []*entities.MbsfnSubframe

	if mBSFN_Subframe_Infolist.list.count > 0 && mBSFN_Subframe_Infolist.list.count <= maxnoofMBSFN {
		count := int(mBSFN_Subframe_Infolist.list.count)
		BSFN_Subframe_Infolist_slice := (*[1 << 30]*C.MBSFN_Subframe_Info_t)(unsafe.Pointer(mBSFN_Subframe_Infolist.list.array))[:count:count]
		for _, member := range BSFN_Subframe_Infolist_slice {
			mBSFNSubframe := &entities.MbsfnSubframe{RadioframeAllocationPeriod: entities.RadioframeAllocationPeriod(1 + member.radioframeAllocationPeriod)}

			mBSFNSubframe.RadioframeAllocationOffset = uint32(member.radioframeAllocationOffset)

			allocType, subframeAllocation, err := getSubframeAllocation_choice(member.subframeAllocation)
			if err != nil {
				return nil, err
			}
			mBSFNSubframe.SubframeAllocation = subframeAllocation
			mBSFNSubframe.SubframeAllocationType = allocType

			mBSFNSubframes = append(mBSFNSubframes, mBSFNSubframe)
		}
	}

	return mBSFNSubframes, nil
}

func getPRACHConfiguration(prachConf *C.PRACH_Configuration_t) *entities.PrachConfiguration {

	var prachConfiguration *entities.PrachConfiguration

	prachConfiguration = &entities.PrachConfiguration{RootSequenceIndex: uint32(prachConf.rootSequenceIndex)}
	prachConfiguration.ZeroCorrelationZoneConfiguration = uint32(prachConf.zeroCorrelationIndex)
	prachConfiguration.HighSpeedFlag = prachConf.highSpeedFlag != 0
	prachConfiguration.PrachFrequencyOffset = uint32(prachConf.prach_FreqOffset)
	if prachConf.prach_ConfigIndex != nil {
		prachConfiguration.PrachConfigurationIndex = uint32(*prachConf.prach_ConfigIndex)
	}

	return prachConfiguration
}
func getServedCellsInfoExt(extIEs *C.ProtocolExtensionContainer_170P192_t, servedCellInfo *entities.ServedCellInfo) error {

	if extIEs != nil && extIEs.list.count > 0 {
		count := int(extIEs.list.count)
		extIEs_slice := (*[1 << 30]*C.ServedCell_Information_ExtIEs_t)(unsafe.Pointer(extIEs.list.array))[:count:count]
		for _, member := range extIEs_slice {
			switch member.extensionValue.present {
			case C.ServedCell_Information_ExtIEs__extensionValue_PR_Number_of_Antennaports:
				ports := (*C.Number_of_Antennaports_t)(unsafe.Pointer(&member.extensionValue.choice[0]))
				servedCellInfo.NumberOfAntennaPorts = entities.NumberOfAntennaPorts(1 + *ports)
			case C.ServedCell_Information_ExtIEs__extensionValue_PR_PRACH_Configuration:
				prachConfiguration := getPRACHConfiguration((*C.PRACH_Configuration_t)(unsafe.Pointer(&member.extensionValue.choice[0])))
				servedCellInfo.PrachConfiguration = prachConfiguration
			case C.ServedCell_Information_ExtIEs__extensionValue_PR_MBSFN_Subframe_Infolist:
				mBSFN_Subframe_Infolist, err := getMBSFN_Subframe_Infolist((*C.MBSFN_Subframe_Infolist_t)(unsafe.Pointer(&member.extensionValue.choice[0])))
				if err != nil {
					return err
				}
				servedCellInfo.MbsfnSubframeInfos = mBSFN_Subframe_Infolist
			case C.ServedCell_Information_ExtIEs__extensionValue_PR_CSG_Id:
				csgId := (*C.CSG_Id_t)(unsafe.Pointer(&member.extensionValue.choice[0]))
				servedCellInfo.CsgId = fmt.Sprintf("%02x", C.GoBytes(unsafe.Pointer(csgId.buf), C.int(csgId.size)))
			case C.ServedCell_Information_ExtIEs__extensionValue_PR_MBMS_Service_Area_Identity_List:
				mBMS_Service_Area_Identity_List := (*C.MBMS_Service_Area_Identity_List_t)(unsafe.Pointer(&member.extensionValue.choice[0]))
				if mBMS_Service_Area_Identity_List.list.count > 0 && mBMS_Service_Area_Identity_List.list.count < maxnoofMBMSServiceAreaIdentities {
					count := int(mBMS_Service_Area_Identity_List.list.count)
					mBMS_Service_Area_Identity_List_slice := (*[1 << 30]*C.MBMS_Service_Area_Identity_t)(unsafe.Pointer(mBMS_Service_Area_Identity_List.list.array))[:count:count]
					for _, identity := range mBMS_Service_Area_Identity_List_slice {
						servedCellInfo.MbmsServiceAreaIdentities = append(servedCellInfo.MbmsServiceAreaIdentities, fmt.Sprintf("%02x", C.GoBytes(unsafe.Pointer(identity.buf), C.int(identity.size))))
					}
				}
			case C.ServedCell_Information_ExtIEs__extensionValue_PR_MultibandInfoList:
				multibandInfoList := (*C.MultibandInfoList_t)(unsafe.Pointer(&member.extensionValue.choice[0]))
				if multibandInfoList.list.count > 0 && multibandInfoList.list.count < maxnoofBands {
					count := int(multibandInfoList.list.count)
					multibandInfoList_slice := (*[1 << 30]*C.BandInfo_t)(unsafe.Pointer(multibandInfoList.list.array))[:count:count]
					for _, bandInfo := range multibandInfoList_slice {
						servedCellInfo.MultibandInfos = append(servedCellInfo.MultibandInfos, uint32(bandInfo.freqBandIndicator))
					}
				}
			case C.ServedCell_Information_ExtIEs__extensionValue_PR_FreqBandIndicatorPriority:
				priority := (*C.FreqBandIndicatorPriority_t)(unsafe.Pointer(&member.extensionValue.choice[0]))
				servedCellInfo.FreqBandIndicatorPriority = entities.FreqBandIndicatorPriority(1 + *priority)
			case C.ServedCell_Information_ExtIEs__extensionValue_PR_BandwidthReducedSI:
				si := (*C.BandwidthReducedSI_t)(unsafe.Pointer(&member.extensionValue.choice[0]))
				servedCellInfo.BandwidthReducedSi = entities.BandwidthReducedSI(1 + *si)
			case C.ServedCell_Information_ExtIEs__extensionValue_PR_ProtectedEUTRAResourceIndication:
				/*ignored*/

			}

		}

	}

	return nil
}

func getServedCellsNeighbour_Info(neighbour_Information *C.Neighbour_Information_t) ([]*entities.NeighbourInformation, error) {
	var neighbours []*entities.NeighbourInformation

	if neighbour_Information != nil && neighbour_Information.list.count > 0 && neighbour_Information.list.count <= maxnoofNeighbours {
		count := int(neighbour_Information.list.count)
		neighbour_Information_slice := (*[1 << 30]*C.Neighbour_Information__Member)(unsafe.Pointer(neighbour_Information.list.array))[:count:count]
		for _, member := range neighbour_Information_slice {

			//pLMN_Identity:eUTRANcellIdentifier
			plmnId := C.GoBytes(unsafe.Pointer(member.eCGI.pLMN_Identity.buf), C.int(member.eCGI.pLMN_Identity.size))
			eUTRANcellIdentifier := C.GoBytes(unsafe.Pointer(member.eCGI.eUTRANcellIdentifier.buf), C.int(member.eCGI.eUTRANcellIdentifier.size))
			neighbourInfo := &entities.NeighbourInformation{Ecgi: fmt.Sprintf("%02x:%02x", plmnId, eUTRANcellIdentifier)}

			neighbourInfo.Pci = uint32(member.pCI)

			neighbourInfo.EarFcn = uint32(member.eARFCN)

			extIEs := (*C.ProtocolExtensionContainer_170P172_t)(unsafe.Pointer(member.iE_Extensions))
			if extIEs != nil && extIEs.list.count > 0 {
				count := int(extIEs.list.count)
				neighbour_Information_ExtIEs_slice := (*[1 << 30]*C.Neighbour_Information_ExtIEs_t)(unsafe.Pointer(extIEs.list.array))[:count:count]
				for _, neighbour_Information_ExtIE := range neighbour_Information_ExtIEs_slice {
					switch neighbour_Information_ExtIE.extensionValue.present {
					case C.Neighbour_Information_ExtIEs__extensionValue_PR_TAC:
						tac := (*C.TAC_t)(unsafe.Pointer(&neighbour_Information_ExtIE.extensionValue.choice[0]))
						neighbourInfo.Tac = fmt.Sprintf("%02x", C.GoBytes(unsafe.Pointer(tac.buf), C.int(tac.size)))
					case C.Neighbour_Information_ExtIEs__extensionValue_PR_EARFCNExtension:
						earFcn := (*C.EARFCNExtension_t)(unsafe.Pointer(&neighbour_Information_ExtIE.extensionValue.choice[0]))
						neighbourInfo.EarFcn = uint32(*earFcn)
					}
				}
			}

			neighbours = append(neighbours, neighbourInfo)
		}
	}

	return neighbours, nil
}

func getServedCells(servedCellsIE *C.ServedCells_t) ([]*entities.ServedCellInfo, error) {
	var servedCells []*entities.ServedCellInfo

	if servedCellsIE != nil && servedCellsIE.list.count > 0 && servedCellsIE.list.count < maxCellineNB {
		count := int(servedCellsIE.list.count)
		servedCells__Member_slice := (*[1 << 30]*C.ServedCells__Member)(unsafe.Pointer(servedCellsIE.list.array))[:count:count]
		for _, member := range servedCells__Member_slice {
			servedCellInfo := &entities.ServedCellInfo{Pci: uint32(member.servedCellInfo.pCI)}

			//pLMN_Identity:eUTRANcellIdentifier
			plmnId := C.GoBytes(unsafe.Pointer(member.servedCellInfo.cellId.pLMN_Identity.buf), C.int(member.servedCellInfo.cellId.pLMN_Identity.size))
			eUTRANcellIdentifier := C.GoBytes(unsafe.Pointer(member.servedCellInfo.cellId.eUTRANcellIdentifier.buf), C.int(member.servedCellInfo.cellId.eUTRANcellIdentifier.size))
			servedCellInfo.CellId = fmt.Sprintf("%02x:%02x", plmnId, eUTRANcellIdentifier)

			servedCellInfo.Tac = fmt.Sprintf("%02x", C.GoBytes(unsafe.Pointer(member.servedCellInfo.tAC.buf), C.int(member.servedCellInfo.tAC.size)))

			if member.servedCellInfo.broadcastPLMNs.list.count > 0 && member.servedCellInfo.broadcastPLMNs.list.count <= maxnoofBPLMNs {
				count := int(member.servedCellInfo.broadcastPLMNs.list.count)
				pLMN_Identity_slice := (*[1 << 30]*C.PLMN_Identity_t)(unsafe.Pointer(member.servedCellInfo.broadcastPLMNs.list.array))[:count:count]
				for _, pLMN_Identity := range pLMN_Identity_slice {
					servedCellInfo.BroadcastPlmns = append(servedCellInfo.BroadcastPlmns, fmt.Sprintf("%02x", C.GoBytes(unsafe.Pointer(pLMN_Identity.buf), C.int(pLMN_Identity.size))))
				}
			}

			switch member.servedCellInfo.eUTRA_Mode_Info.present {
			case C.EUTRA_Mode_Info_PR_fDD:
				if fdd, err := getFDDInfo(*(**C.FDD_Info_t)(unsafe.Pointer(&member.servedCellInfo.eUTRA_Mode_Info.choice[0]))); fdd != nil && err == nil {
					servedCellInfo.ChoiceEutraMode, servedCellInfo.EutraMode = &entities.ChoiceEUTRAMode{Fdd: fdd}, entities.Eutra_FDD
				} else {
					return nil, err
				}
			case C.EUTRA_Mode_Info_PR_tDD:
				if tdd, err := getTDDInfo(*(**C.TDD_Info_t)(unsafe.Pointer(&member.servedCellInfo.eUTRA_Mode_Info.choice[0]))); tdd != nil && err == nil {
					servedCellInfo.ChoiceEutraMode, servedCellInfo.EutraMode = &entities.ChoiceEUTRAMode{Tdd: tdd}, entities.Eutra_TDD
				} else {
					return nil, err
				}
			}

			neighbours, err := getServedCellsNeighbour_Info(member.neighbour_Info)
			if err != nil {
				return nil, err
			}
			servedCellInfo.NeighbourInfos = neighbours

			if err := getServedCellsInfoExt((*C.ProtocolExtensionContainer_170P192_t)(unsafe.Pointer(member.servedCellInfo.iE_Extensions)), servedCellInfo); err != nil {
				return nil, err
			}

			servedCells = append(servedCells, servedCellInfo)

		}
	}

	return servedCells, nil
}

func getGUGroupIDList(guGroupIDList *C.GUGroupIDList_t) []string {
	var ids []string

	if guGroupIDList != nil && guGroupIDList.list.count > 0 && guGroupIDList.list.count <= maxPools {
		count := int(guGroupIDList.list.count)
		guGroupIDList_slice := (*[1 << 30]*C.GU_Group_ID_t)(unsafe.Pointer(guGroupIDList.list.array))[:count:count]
		for _, guGroupID := range guGroupIDList_slice {
			plmnId := C.GoBytes(unsafe.Pointer(guGroupID.pLMN_Identity.buf), C.int(guGroupID.pLMN_Identity.size))
			mME_Group_ID := C.GoBytes(unsafe.Pointer(guGroupID.mME_Group_ID.buf), C.int(guGroupID.mME_Group_ID.size))
			ids = append(ids, fmt.Sprintf("%02x:%02x", plmnId, mME_Group_ID))
		}
	}

	return ids
}

// Populate  the ENB structure with data from the pdu
// Return the ENB and the associated key which can later be used to retrieve the ENB from the database.
func x2SetupResponseToProtobuf(pdu *C.E2AP_PDU_t) (*entities.GlobalNbId, *entities.Enb, error) {
	var globalNbId *entities.GlobalNbId

	enb := entities.Enb{}

	if pdu.present == C.E2AP_PDU_PR_successfulOutcome {
		//dereference a union of pointers (C union is represented as a byte array with the size of the largest member)
		successfulOutcome := *(**C.SuccessfulOutcome_t)(unsafe.Pointer(&pdu.choice[0]))
		if successfulOutcome != nil && successfulOutcome.value.present == C.SuccessfulOutcome__value_PR_X2SetupResponse {
			x2SetupResponse := (*C.X2SetupResponse_t)(unsafe.Pointer(&successfulOutcome.value.choice[0]))
			if x2SetupResponse != nil && x2SetupResponse.protocolIEs.list.count > 0 {
				count := int(x2SetupResponse.protocolIEs.list.count)
				x2SetupResponse_IEs_slice := (*[1 << 30]*C.X2SetupResponse_IEs_t)(unsafe.Pointer(x2SetupResponse.protocolIEs.list.array))[:count:count]
				for _, x2SetupResponse_IE := range x2SetupResponse_IEs_slice {
					switch x2SetupResponse_IE.value.present {
					case C.X2SetupResponse_IEs__value_PR_GlobalENB_ID:
						globalENB_ID := (*C.GlobalENB_ID_t)(unsafe.Pointer(&x2SetupResponse_IE.value.choice[0]))
						plmnId := C.GoBytes(unsafe.Pointer(globalENB_ID.pLMN_Identity.buf), C.int(globalENB_ID.pLMN_Identity.size))
						enbType, enbVal := getENB_ID_choice(globalENB_ID.eNB_ID)

						globalNbId = &entities.GlobalNbId{}
						globalNbId.NbId = fmt.Sprintf("%02x", enbVal)
						globalNbId.PlmnId = fmt.Sprintf("%02x", plmnId)
						enb.EnbType = enbType

					case C.X2SetupResponse_IEs__value_PR_ServedCells:
						ServedCells, err := getServedCells((*C.ServedCells_t)(unsafe.Pointer(&x2SetupResponse_IE.value.choice[0])))
						if err != nil {
							return globalNbId, nil, err
						}
						enb.ServedCells = ServedCells
					case C.X2SetupResponse_IEs__value_PR_GUGroupIDList:
						enb.GuGroupIds = getGUGroupIDList((*C.GUGroupIDList_t)(unsafe.Pointer(&x2SetupResponse_IE.value.choice[0])))
					case C.X2SetupResponse_IEs__value_PR_CriticalityDiagnostics:
						/*ignored*/
					case C.X2SetupResponse_IEs__value_PR_LHN_ID:
						/*ignored*/
					}
				}
			}
		}
	}

	return globalNbId, &enb, nil
}

func (c *X2SetupResponseConverter) UnpackX2SetupResponseAndExtract(packedBuf []byte) (*entities.GlobalNbId, *entities.Enb, error) {
	pdu, err := UnpackX2apPdu(c.logger, e2pdus.MaxAsn1CodecAllocationBufferSize, len(packedBuf), packedBuf, e2pdus.MaxAsn1CodecMessageBufferSize)
	if err != nil {
		return nil, nil, err
	}

	defer C.delete_pdu(pdu)
	return x2SetupResponseToProtobuf(pdu)
}
