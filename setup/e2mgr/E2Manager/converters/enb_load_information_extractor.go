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
/*
// #cgo CFLAGS: -I../3rdparty/asn1codec/inc/  -I../3rdparty/asn1codec/e2ap_engine/
// #cgo LDFLAGS: -L ../3rdparty/asn1codec/lib/ -L../3rdparty/asn1codec/e2ap_engine/ -le2ap_codec -lasncodec
// #include <asn1codec_utils.h>
// #include <load_information_wrapper.h>
import "C"
import (
	"e2mgr/logger"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"unsafe"
)

const (
	MaxCellsInEnb            = 256
	MaxNoOfPrbs              = 110
	NaxNoOfCompHypothesisSet = 256
	NaxNoOfCompCells         = 32
	MaxNoOfPa                = 3
)

type EnbLoadInformationExtractor struct {
	logger *logger.Logger
}

type IEnbLoadInformationExtractor interface {
	ExtractAndBuildRanLoadInformation(pdu *C.E2AP_PDU_t, ranLoadInformation *entities.RanLoadInformation) error
}

func NewEnbLoadInformationExtractor(logger *logger.Logger) *EnbLoadInformationExtractor {
	return &EnbLoadInformationExtractor{
		logger: logger,
	}
}

var populators = map[C.CellInformation_Item_ExtIEs__extensionValue_PR]func(string, *entities.CellLoadInformation, *C.CellInformation_Item_ExtIEs_t) error{
	C.CellInformation_Item_ExtIEs__extensionValue_PR_ABSInformation:                     populateAbsInformation,
	C.CellInformation_Item_ExtIEs__extensionValue_PR_InvokeIndication:                   populateInvokeIndication,
	C.CellInformation_Item_ExtIEs__extensionValue_PR_SubframeAssignment:                 populateIntendedUlDlConfiguration,
	C.CellInformation_Item_ExtIEs__extensionValue_PR_ExtendedULInterferenceOverloadInfo: populateExtendedUlInterferenceOverloadInfo,
	C.CellInformation_Item_ExtIEs__extensionValue_PR_CoMPInformation:                    populateCompInformation,
	C.CellInformation_Item_ExtIEs__extensionValue_PR_DynamicDLTransmissionInformation:   populateDynamicDLTransmissionInformation,
}

func extractPduCellInformationItemIEs(pdu *C.E2AP_PDU_t) ([]*C.CellInformation_ItemIEs_t, error) {

	if pdu.present != C.E2AP_PDU_PR_initiatingMessage {
		return nil, fmt.Errorf("#extractPduCellInformationItemIEs - Invalid E2AP_PDU value")
	}

	initiatingMessage := *(**C.InitiatingMessage_t)(unsafe.Pointer(&pdu.choice[0]))

	if initiatingMessage == nil || initiatingMessage.value.present != C.InitiatingMessage__value_PR_LoadInformation {
		return nil, fmt.Errorf("#extractPduCellInformationItemIEs - Invalid InitiatingMessage value")
	}

	loadInformationMessage := (*C.LoadInformation_t)(unsafe.Pointer(&initiatingMessage.value.choice[0]))

	if loadInformationMessage == nil {
		return nil, fmt.Errorf("#extractPduCellInformationItemIEs - Invalid LoadInformation container")
	}

	protocolIEsListCount := loadInformationMessage.protocolIEs.list.count

	if protocolIEsListCount != 1 {
		return nil, fmt.Errorf("#extractPduCellInformationItemIEs - Invalid protocolIEs list count")
	}

	loadInformationIEs := (*[1 << 30]*C.LoadInformation_IEs_t)(unsafe.Pointer(loadInformationMessage.protocolIEs.list.array))[:int(protocolIEsListCount):int(protocolIEsListCount)]

	loadInformationIE := loadInformationIEs[0]

	if loadInformationIE.value.present != C.LoadInformation_IEs__value_PR_CellInformation_List {
		return nil, fmt.Errorf("#extractPduCellInformationItemIEs - Invalid protocolIEs value")
	}

	loadInformationCellList := (*C.CellInformation_List_t)(unsafe.Pointer(&loadInformationIE.value.choice[0]))

	loadInformationCellListCount := loadInformationCellList.list.count

	if loadInformationCellListCount == 0 || loadInformationCellListCount > MaxCellsInEnb {
		return nil, fmt.Errorf("#extractPduCellInformationItemIEs - Invalid CellInformation list count")
	}

	return (*[1 << 30]*C.CellInformation_ItemIEs_t)(unsafe.Pointer(loadInformationCellList.list.array))[:int(loadInformationCellListCount):int(loadInformationCellListCount)], nil
}

func populateCellLoadInformation(pduCellInformationItemIE *C.CellInformation_ItemIEs_t, ranLoadInformation *entities.RanLoadInformation) error {

	pduCellInformationItem, err := extractPduCellInformationItem(pduCellInformationItemIE)

	if err != nil {
		return err
	}

	cellLoadInformation := entities.CellLoadInformation{}
	err = buildCellLoadInformation(&cellLoadInformation, pduCellInformationItem)

	if err != nil {
		return err
	}

	ranLoadInformation.CellLoadInfos = append(ranLoadInformation.CellLoadInfos, &cellLoadInformation)

	return nil
}

func extractPduCellInformationItem(cell *C.CellInformation_ItemIEs_t) (*C.CellInformation_Item_t, error) {

	if cell.value.present != C.CellInformation_ItemIEs__value_PR_CellInformation_Item {
		return nil, fmt.Errorf("#extractPduCellInformationItem - Failed extracting pdu cell information item")
	}

	return (*C.CellInformation_Item_t)(unsafe.Pointer(&cell.value.choice[0])), nil
}

func buildCellLoadInformation(cellLoadInformation *entities.CellLoadInformation, cellInformationItem *C.CellInformation_Item_t) error {

	cellId := buildCellId(cellInformationItem.cell_ID)
	cellLoadInformation.CellId = cellId

	pduUlInterferenceOverloadIndicationItems, err := extractPduUlInterferenceOverloadIndicationItems(cellId, cellInformationItem)

	if err != nil {
		return err
	}

	if (pduUlInterferenceOverloadIndicationItems != nil) {
		cellLoadInformation.UlInterferenceOverloadIndications = buildUlInterferenceOverloadIndicationList(pduUlInterferenceOverloadIndicationItems)
	}

	pduUlHighInterferenceIndicationInfoItems, err := extractPduUlHighInterferenceIndicationInfoItems(cellId, cellInformationItem)

	if err != nil {
		return err
	}

	if pduUlHighInterferenceIndicationInfoItems != nil {
		cellLoadInformation.UlHighInterferenceInfos = buildUlHighInterferenceIndicationInfoList(pduUlHighInterferenceIndicationInfoItems)
	}

	pduRelativeNarrowbandTxPower := extractPduRelativeNarrowbandTxPower(cellInformationItem)

	if pduRelativeNarrowbandTxPower != nil {
		pduEnhancedRntp, err := extractPduEnhancedRntp(cellId, pduRelativeNarrowbandTxPower)

		if err != nil {
			return err
		}

		cellLoadInformation.RelativeNarrowbandTxPower = buildRelativeNarrowbandTxPower(pduRelativeNarrowbandTxPower, pduEnhancedRntp)
	}

	pduCellInformationItemExtensionIEs := extractPduCellInformationItemExtensionIEs(cellInformationItem)

	if (pduCellInformationItemExtensionIEs == nil) {
		return nil
	}

	err = populateCellLoadInformationExtensionIEs(cellId, cellLoadInformation, pduCellInformationItemExtensionIEs)

	if err != nil {
		return err
	}

	return nil

}

func buildCellId(cellId C.ECGI_t) string {
	plmnId := cUcharArrayToGoByteSlice(cellId.pLMN_Identity.buf, cellId.pLMN_Identity.size)
	eutranCellIdentifier := cUcharArrayToGoByteSlice(cellId.eUTRANcellIdentifier.buf, cellId.eUTRANcellIdentifier.size)
	return fmt.Sprintf("%x:%x", plmnId, eutranCellIdentifier)
}

func buildUlInterferenceOverloadIndicationList(pduUlInterferenceOverloadIndicationItems []*C.UL_InterferenceOverloadIndication_Item_t) []entities.UlInterferenceOverloadIndication {
	indications := make([]entities.UlInterferenceOverloadIndication, len(pduUlInterferenceOverloadIndicationItems))
	for i, ci := range pduUlInterferenceOverloadIndicationItems {
		indications[i] = entities.UlInterferenceOverloadIndication(*ci + 1)
	}

	return indications
}

func extractPduUlInterferenceOverloadIndicationItems(cellId string, cellInformationItem *C.CellInformation_Item_t) ([]*C.UL_InterferenceOverloadIndication_Item_t, error) {

	if cellInformationItem.ul_InterferenceOverloadIndication == nil {
		return nil, nil
	}

	ulInterferenceOverLoadIndicationCount := cellInformationItem.ul_InterferenceOverloadIndication.list.count

	if ulInterferenceOverLoadIndicationCount == 0 || ulInterferenceOverLoadIndicationCount > MaxNoOfPrbs {
		return nil, fmt.Errorf("#extractPduUlInterferenceOverloadIndicationItems - cellId: %s - Invalid UL Interference OverLoad Indication list count", cellId)
	}

	pduUlInterferenceOverloadIndicationItems := (*[1 << 30]*C.UL_InterferenceOverloadIndication_Item_t)(unsafe.Pointer(cellInformationItem.ul_InterferenceOverloadIndication.list.array))[:int(ulInterferenceOverLoadIndicationCount):int(ulInterferenceOverLoadIndicationCount)]

	return pduUlInterferenceOverloadIndicationItems, nil
}

func NewStartTime(startSfn C.long, startSubframeNumber C.long) *entities.StartTime {
	return &entities.StartTime{
		StartSfn:            int32(startSfn),
		StartSubframeNumber: int32(startSubframeNumber),
	}
}

func buildEnhancedRntp(pduEnhancedRntp *C.EnhancedRNTP_t) *entities.EnhancedRntp {

	enhancedRntp := entities.EnhancedRntp{
		EnhancedRntpBitmap:     NewHexString(pduEnhancedRntp.enhancedRNTPBitmap.buf, pduEnhancedRntp.enhancedRNTPBitmap.size),
		RntpHighPowerThreshold: entities.RntpThreshold(pduEnhancedRntp.rNTP_High_Power_Threshold + 1),
	}

	pduEnhancedRntpStartTime := (*C.EnhancedRNTPStartTime_t)(unsafe.Pointer(pduEnhancedRntp.enhancedRNTPStartTime))

	if pduEnhancedRntpStartTime != nil {
		enhancedRntp.EnhancedRntpStartTime = NewStartTime(pduEnhancedRntpStartTime.startSFN, pduEnhancedRntpStartTime.startSubframeNumber)
	}

	return &enhancedRntp
}

func cUcharArrayToGoByteSlice(buf *C.uchar, size C.ulong) []byte {
	return C.GoBytes(unsafe.Pointer(buf), C.int(size))
}

func NewHexString(buf *C.uchar, size C.ulong) string {
	return fmt.Sprintf("%x", cUcharArrayToGoByteSlice(buf, size))
}

func buildRelativeNarrowbandTxPower(pduRelativeNarrowbandTxPower *C.RelativeNarrowbandTxPower_t, pduEnhancedRntp *C.EnhancedRNTP_t) *entities.RelativeNarrowbandTxPower {

	relativeNarrowbandTxPower := entities.RelativeNarrowbandTxPower{
		RntpThreshold:                    entities.RntpThreshold(pduRelativeNarrowbandTxPower.rNTP_Threshold + 1),
		NumberOfCellSpecificAntennaPorts: entities.NumberOfCellSpecificAntennaPorts(pduRelativeNarrowbandTxPower.numberOfCellSpecificAntennaPorts + 1),
		PB:                               uint32(pduRelativeNarrowbandTxPower.p_B),
		PdcchInterferenceImpact:          uint32(pduRelativeNarrowbandTxPower.pDCCH_InterferenceImpact),
		RntpPerPrb:                       NewHexString(pduRelativeNarrowbandTxPower.rNTP_PerPRB.buf, pduRelativeNarrowbandTxPower.rNTP_PerPRB.size),
	}

	if pduEnhancedRntp != nil {
		relativeNarrowbandTxPower.EnhancedRntp = buildEnhancedRntp(pduEnhancedRntp)
	}

	return &relativeNarrowbandTxPower
}

func extractPduRelativeNarrowbandTxPower(cellInformationItem *C.CellInformation_Item_t) *C.RelativeNarrowbandTxPower_t {

	if cellInformationItem.relativeNarrowbandTxPower == nil {
		return nil
	}

	return (*C.RelativeNarrowbandTxPower_t)(unsafe.Pointer(cellInformationItem.relativeNarrowbandTxPower))
}

func buildUlHighInterferenceIndicationInfoList(ulHighInterferenceIndicationInfoList []*C.UL_HighInterferenceIndicationInfo_Item_t) []*entities.UlHighInterferenceInformation {
	infos := make([]*entities.UlHighInterferenceInformation, len(ulHighInterferenceIndicationInfoList))

	for i, v := range ulHighInterferenceIndicationInfoList {

		infos[i] = &entities.UlHighInterferenceInformation{
			TargetCellId:                 buildCellId(v.target_Cell_ID),
			UlHighInterferenceIndication: NewHexString(v.ul_interferenceindication.buf, v.ul_interferenceindication.size),
		}
	}

	return infos
}

func extractPduUlHighInterferenceIndicationInfoItems(cellId string, cellInformationItem *C.CellInformation_Item_t) ([]*C.UL_HighInterferenceIndicationInfo_Item_t, error) {
	pduUlHighInterferenceIndicationInfo := (*C.UL_HighInterferenceIndicationInfo_t)(unsafe.Pointer(cellInformationItem.ul_HighInterferenceIndicationInfo))

	if (pduUlHighInterferenceIndicationInfo == nil) {
		return nil, nil
	}

	pduUlHighInterferenceIndicationInfoListCount := pduUlHighInterferenceIndicationInfo.list.count

	if pduUlHighInterferenceIndicationInfoListCount == 0 || pduUlHighInterferenceIndicationInfoListCount > MaxCellsInEnb {
		return nil, fmt.Errorf("#extractPduUlHighInterferenceIndicationInfoItems - cellId: %s - Invalid UL High Interference Indication info list count", cellId)
	}

	pduUlHighInterferenceIndicationInfoItems := (*[1 << 30]*C.UL_HighInterferenceIndicationInfo_Item_t)(unsafe.Pointer(cellInformationItem.ul_HighInterferenceIndicationInfo.list.array))[:int(pduUlHighInterferenceIndicationInfoListCount):int(pduUlHighInterferenceIndicationInfoListCount)]

	return pduUlHighInterferenceIndicationInfoItems, nil
}

func extractPduCellInformationItemExtensionIEs(cellInformationItem *C.CellInformation_Item_t) []*C.CellInformation_Item_ExtIEs_t {
	extIEs := (*C.ProtocolExtensionContainer_170P7_t)(unsafe.Pointer(cellInformationItem.iE_Extensions))

	if extIEs == nil {
		return nil
	}

	extIEsCount := int(extIEs.list.count)

	if extIEsCount == 0 {
		return nil
	}

	return (*[1 << 30]*C.CellInformation_Item_ExtIEs_t)(unsafe.Pointer(extIEs.list.array))[:extIEsCount:extIEsCount]
}

func extractPduEnhancedRntp(cellId string, pduRelativeNarrowbandTxPower *C.RelativeNarrowbandTxPower_t) (*C.EnhancedRNTP_t, error) {

	extIEs := (*C.ProtocolExtensionContainer_170P184_t)(unsafe.Pointer(pduRelativeNarrowbandTxPower.iE_Extensions))

	if extIEs == nil {
		return nil, nil
	}

	extIEsCount := int(extIEs.list.count)

	if extIEsCount != 1 {
		return nil, fmt.Errorf("#extractPduEnhancedRntp - cellId: %s - Invalid Enhanced RNTP container", cellId)
	}

	enhancedRntpExtIEs := (*[1 << 30]*C.RelativeNarrowbandTxPower_ExtIEs_t)(unsafe.Pointer(extIEs.list.array))[:extIEsCount:extIEsCount]

	enhancedRntpExtIE := enhancedRntpExtIEs[0]

	if enhancedRntpExtIE.extensionValue.present != C.RelativeNarrowbandTxPower_ExtIEs__extensionValue_PR_EnhancedRNTP {
		return nil, fmt.Errorf("#extractPduEnhancedRntp - cellId: %s - Invalid Enhanced RNTP container", cellId)
	}

	return (*C.EnhancedRNTP_t)(unsafe.Pointer(&enhancedRntpExtIE.extensionValue.choice[0])), nil
}

func buildAbsInformationFdd(cellId string, pduAbsInformation *C.ABSInformation_t) (*entities.AbsInformation, error) {
	pduAbsInformationFdd := *(**C.ABSInformationFDD_t)(unsafe.Pointer(&pduAbsInformation.choice[0]))

	if pduAbsInformationFdd == nil {
		return nil, fmt.Errorf("#buildAbsInformationFdd - cellId: %s - Invalid FDD Abs Information", cellId)
	}

	absInformation := entities.AbsInformation{
		Mode:                             entities.AbsInformationMode_ABS_INFO_FDD,
		AbsPatternInfo:                   NewHexString(pduAbsInformationFdd.abs_pattern_info.buf, pduAbsInformationFdd.abs_pattern_info.size),
		NumberOfCellSpecificAntennaPorts: entities.NumberOfCellSpecificAntennaPorts(pduAbsInformationFdd.numberOfCellSpecificAntennaPorts + 1),
		MeasurementSubset:                NewHexString(pduAbsInformationFdd.measurement_subset.buf, pduAbsInformationFdd.measurement_subset.size),
	}

	return &absInformation, nil
}

func buildAbsInformationTdd(cellId string, pduAbsInformation *C.ABSInformation_t) (*entities.AbsInformation, error) {
	pduAbsInformationTdd := *(**C.ABSInformationTDD_t)(unsafe.Pointer(&pduAbsInformation.choice[0]))

	if pduAbsInformationTdd == nil {
		return nil, fmt.Errorf("#buildAbsInformationTdd - cellId: %s - Invalid TDD Abs Information", cellId)
	}

	absInformation := entities.AbsInformation{
		Mode:                             entities.AbsInformationMode_ABS_INFO_TDD,
		AbsPatternInfo:                   NewHexString(pduAbsInformationTdd.abs_pattern_info.buf, pduAbsInformationTdd.abs_pattern_info.size),
		NumberOfCellSpecificAntennaPorts: entities.NumberOfCellSpecificAntennaPorts(pduAbsInformationTdd.numberOfCellSpecificAntennaPorts + 1),
		MeasurementSubset:                NewHexString(pduAbsInformationTdd.measurement_subset.buf, pduAbsInformationTdd.measurement_subset.size),
	}

	return &absInformation, nil
}

func extractAndBuildCellLoadInformationAbsInformation(cellId string, extIE *C.CellInformation_Item_ExtIEs_t) (*entities.AbsInformation, error) {
	pduAbsInformation := (*C.ABSInformation_t)(unsafe.Pointer(&extIE.extensionValue.choice[0]))

	switch pduAbsInformation.present {
	case C.ABSInformation_PR_fdd:
		return buildAbsInformationFdd(cellId, pduAbsInformation)
	case C.ABSInformation_PR_tdd:
		return buildAbsInformationTdd(cellId, pduAbsInformation)
	case C.ABSInformation_PR_abs_inactive:
		return &entities.AbsInformation{Mode: entities.AbsInformationMode_ABS_INACTIVE}, nil
	}

	return nil, fmt.Errorf("#extractAndBuildCellLoadInformationAbsInformation - cellId: %s - Failed extracting AbsInformation", cellId)
}

func extractExtendedUlInterferenceOverloadIndicationList(extendedULInterferenceOverloadInfo *C.ExtendedULInterferenceOverloadInfo_t) []*C.UL_InterferenceOverloadIndication_Item_t {

	extendedUlInterferenceOverLoadIndicationCount := extendedULInterferenceOverloadInfo.extended_ul_InterferenceOverloadIndication.list.count

	if extendedUlInterferenceOverLoadIndicationCount == 0 {
		return nil
	}

	extendedUlInterferenceOverLoadIndicationList := (*[1 << 30]*C.UL_InterferenceOverloadIndication_Item_t)(unsafe.Pointer(extendedULInterferenceOverloadInfo.extended_ul_InterferenceOverloadIndication.list.array))[:int(extendedUlInterferenceOverLoadIndicationCount):int(extendedUlInterferenceOverLoadIndicationCount)]

	if (extendedUlInterferenceOverLoadIndicationList == nil) {
		return nil
	}

	return extendedUlInterferenceOverLoadIndicationList
}

func buildExtendedULInterferenceOverloadInfo(pduExtendedULInterferenceOverloadInfo *C.ExtendedULInterferenceOverloadInfo_t) *entities.ExtendedUlInterferenceOverloadInfo {
	associatedSubframes := NewHexString(pduExtendedULInterferenceOverloadInfo.associatedSubframes.buf, pduExtendedULInterferenceOverloadInfo.associatedSubframes.size)
	indications := extractExtendedUlInterferenceOverloadIndicationList(pduExtendedULInterferenceOverloadInfo)

	return &entities.ExtendedUlInterferenceOverloadInfo{
		AssociatedSubframes:                       associatedSubframes,
		ExtendedUlInterferenceOverloadIndications: buildUlInterferenceOverloadIndicationList(indications),
	}
}

func extractPaListFromDynamicNaicsInformation(cellId string, dynamicNAICSInformation *C.DynamicNAICSInformation_t) ([]*C.long, error) {

	paListCount := dynamicNAICSInformation.pA_list.list.count

	if paListCount == 0 {
		return nil, nil
	}

	if paListCount > MaxNoOfPa {
		return nil, fmt.Errorf("#extractPaListFromDynamicNaicsInformation - cellId: %s - Invalid PA list count", cellId)
	}

	extendedUlInterferenceOverLoadIndicationList := (*[1 << 30]*C.long)(unsafe.Pointer(dynamicNAICSInformation.pA_list.list.array))[:int(paListCount):int(paListCount)]

	if (extendedUlInterferenceOverLoadIndicationList == nil) {
		return nil, fmt.Errorf("#extractPaListFromDynamicNaicsInformation - cellId: %s - Extended Ul Interference OverLoad Indication List is nil", cellId)
	}

	return extendedUlInterferenceOverLoadIndicationList, nil
}

func buildPaList(paList []*C.long) []entities.PA {
	pas := make([]entities.PA, len(paList))
	for i, pi := range paList {
		pas[i] = entities.PA(*pi + 1)
	}

	return pas
}

func extractAndBuildActiveDynamicDlTransmissionInformation(cellId string, pduDynamicDlTransmissionInformation *C.DynamicDLTransmissionInformation_t) (*entities.DynamicDlTransmissionInformation, error) {
	dynamicNaicsInformation := *(**C.DynamicNAICSInformation_t)(unsafe.Pointer(&pduDynamicDlTransmissionInformation.choice[0]))

	if dynamicNaicsInformation == nil {
		return nil, fmt.Errorf("#extractAndBuildActiveDynamicDlTransmissionInformation - cellId: %s - Invalid NAICS Information value", cellId)
	}

	dynamicDlTransmissionInformation := entities.DynamicDlTransmissionInformation{State: entities.NaicsState_NAICS_ACTIVE}

	if dynamicNaicsInformation.transmissionModes != nil {
		transmissionModes := NewHexString(dynamicNaicsInformation.transmissionModes.buf, dynamicNaicsInformation.transmissionModes.size)
		dynamicDlTransmissionInformation.TransmissionModes = transmissionModes
	}

	if dynamicNaicsInformation.pB_information != nil {
		dynamicDlTransmissionInformation.PB = uint32(*dynamicNaicsInformation.pB_information)
	}

	paList, err := extractPaListFromDynamicNaicsInformation(cellId, dynamicNaicsInformation)

	if err != nil {
		return nil, err
	}

	if (paList != nil) {
		dynamicDlTransmissionInformation.PAList = buildPaList(paList)
	}

	return &dynamicDlTransmissionInformation, nil
}

func extractAndBuildDynamicDlTransmissionInformation(cellId string, extIE *C.CellInformation_Item_ExtIEs_t) (*entities.DynamicDlTransmissionInformation, error) {
	pduDynamicDlTransmissionInformation := (*C.DynamicDLTransmissionInformation_t)(unsafe.Pointer(&extIE.extensionValue.choice[0]))

	if pduDynamicDlTransmissionInformation.present == C.DynamicDLTransmissionInformation_PR_naics_inactive {
		return &entities.DynamicDlTransmissionInformation{State: entities.NaicsState_NAICS_INACTIVE}, nil
	}

	if pduDynamicDlTransmissionInformation.present != C.DynamicDLTransmissionInformation_PR_naics_active {
		return nil, fmt.Errorf("#extractAndBuildDynamicDlTransmissionInformation - cellId: %s - Invalid Dynamic Dl Transmission Information value", cellId)
	}

	return extractAndBuildActiveDynamicDlTransmissionInformation(cellId, pduDynamicDlTransmissionInformation)
}

func extractCompInformationStartTime(cellId string, compInformation *C.CoMPInformation_t) (*C.CoMPInformationStartTime__Member, error) {
	compInformationStartTimeListCount := compInformation.coMPInformationStartTime.list.count

	if compInformationStartTimeListCount == 0 {
		return nil, nil
	}

	compInformationStartTimeList := (*[1 << 30]*C.CoMPInformationStartTime__Member)(unsafe.Pointer(compInformation.coMPInformationStartTime.list.array))[:int(compInformationStartTimeListCount):int(compInformationStartTimeListCount)]

	if len(compInformationStartTimeList) != 1 {
		return nil, fmt.Errorf("#extractCompInformationStartTime - cellId: %s - Invalid Comp Information StartTime list count", cellId)
	}

	return compInformationStartTimeList[0], nil
}

func buildCompHypothesisSet(pduCompHypothesisSetItem *C.CoMPHypothesisSetItem_t) *entities.CompHypothesisSet {
	return &entities.CompHypothesisSet{
		CellId:         buildCellId(pduCompHypothesisSetItem.coMPCellID),
		CompHypothesis: NewHexString(pduCompHypothesisSetItem.coMPHypothesis.buf, pduCompHypothesisSetItem.coMPHypothesis.size),
	}
}

func buildCompHypothesisSets(cellId string, pduCompInformationItemMember *C.CoMPInformationItem__Member) ([]*entities.CompHypothesisSet, error) {

	compHypothesisSets := []*entities.CompHypothesisSet{}

	pduCompHypothesisSetItemsListCount := pduCompInformationItemMember.coMPHypothesisSet.list.count

	if pduCompHypothesisSetItemsListCount == 0 || pduCompHypothesisSetItemsListCount > NaxNoOfCompCells {
		return nil, fmt.Errorf("#buildCompHypothesisSets - cellId: %s - Invalid Comp Hypothesis Set Items list count", cellId)
	}

	pduCompHypothesisSetItems := (*[1 << 30]*C.CoMPHypothesisSetItem_t)(unsafe.Pointer(pduCompInformationItemMember.coMPHypothesisSet.list.array))[:int(pduCompHypothesisSetItemsListCount):int(pduCompHypothesisSetItemsListCount)]

	for _, pduCompHypothesisSetItem := range pduCompHypothesisSetItems {
		compHypothesisSet := buildCompHypothesisSet(pduCompHypothesisSetItem)
		compHypothesisSets = append(compHypothesisSets, compHypothesisSet)
	}

	return compHypothesisSets, nil
}

func buildCompInformationItem(cellId string, pduCompInformationItemMember *C.CoMPInformationItem__Member) (*entities.CompInformationItem, error) {

	compHypothesisSets, err := buildCompHypothesisSets(cellId, pduCompInformationItemMember)

	if err != nil {
		return nil, err
	}

	compInformation := entities.CompInformationItem{
		CompHypothesisSets: compHypothesisSets,
		BenefitMetric:      int32(pduCompInformationItemMember.benefitMetric),
	}

	return &compInformation, nil
}

func extractAndBuildCompInformation(cellId string, extIE *C.CellInformation_Item_ExtIEs_t) (*entities.CompInformation, error) {

	pduCompInformation := (*C.CoMPInformation_t)(unsafe.Pointer(&extIE.extensionValue.choice[0]))

	compInformation := entities.CompInformation{}
	pduCompInformationStartTime, err := extractCompInformationStartTime(cellId, pduCompInformation)

	if err != nil {
		return nil, err
	}

	if pduCompInformationStartTime != nil {
		compInformation.CompInformationStartTime = NewStartTime(pduCompInformationStartTime.startSFN, pduCompInformationStartTime.startSubframeNumber)
	}

	pduCompInformationItemsListCount := pduCompInformation.coMPInformationItem.list.count

	if pduCompInformationItemsListCount == 0 || pduCompInformationItemsListCount > NaxNoOfCompHypothesisSet {
		return nil, fmt.Errorf("#extractAndBuildCompInformation - cellId: %s - Invalid Comp Information Items list count", cellId)
	}

	pduCompInformationItems := (*[1 << 30]*C.CoMPInformationItem__Member)(unsafe.Pointer(pduCompInformation.coMPInformationItem.list.array))[:int(pduCompInformationItemsListCount):int(pduCompInformationItemsListCount)]

	for _, pduCompInformationItem := range pduCompInformationItems {
		compInformationItem, err := buildCompInformationItem(cellId, pduCompInformationItem)

		if err != nil {
			return nil, err
		}

		compInformation.CompInformationItems = append(compInformation.CompInformationItems, compInformationItem)
	}

	return &compInformation, nil
}

func populateAbsInformation(cellId string, cellLoadInformation *entities.CellLoadInformation, extIE *C.CellInformation_Item_ExtIEs_t) error {
	absInformation, err := extractAndBuildCellLoadInformationAbsInformation(cellId, extIE)

	if err != nil {
		return err
	}

	cellLoadInformation.AbsInformation = absInformation
	return nil
}

func populateInvokeIndication(cellId string, cellLoadInformation *entities.CellLoadInformation, extIE *C.CellInformation_Item_ExtIEs_t) error {
	pduInvokeIndication := (*C.InvokeIndication_t)(unsafe.Pointer(&extIE.extensionValue.choice[0]))
	cellLoadInformation.InvokeIndication = entities.InvokeIndication(*pduInvokeIndication + 1)
	return nil
}

func populateIntendedUlDlConfiguration(cellId string, cellLoadInformation *entities.CellLoadInformation, extIE *C.CellInformation_Item_ExtIEs_t) error {
	pduSubframeAssignment := (*C.SubframeAssignment_t)(unsafe.Pointer(&extIE.extensionValue.choice[0]))
	cellLoadInformation.IntendedUlDlConfiguration = entities.SubframeAssignment(*pduSubframeAssignment + 1)
	return nil
}

func populateExtendedUlInterferenceOverloadInfo(cellId string, cellLoadInformation *entities.CellLoadInformation, extIE *C.CellInformation_Item_ExtIEs_t) error {
	pduExtendedULInterferenceOverloadInfo := (*C.ExtendedULInterferenceOverloadInfo_t)(unsafe.Pointer(&extIE.extensionValue.choice[0]))
	cellLoadInformation.ExtendedUlInterferenceOverloadInfo = buildExtendedULInterferenceOverloadInfo(pduExtendedULInterferenceOverloadInfo)
	return nil
}

func populateCompInformation(cellId string, cellLoadInformation *entities.CellLoadInformation, extIE *C.CellInformation_Item_ExtIEs_t) error {
	pduCompInformation, err := extractAndBuildCompInformation(cellId, extIE)

	if err != nil {
		return err
	}

	cellLoadInformation.CompInformation = pduCompInformation
	return nil
}

func populateDynamicDLTransmissionInformation(cellId string, cellLoadInformation *entities.CellLoadInformation, extIE *C.CellInformation_Item_ExtIEs_t) error {
	dynamicDLTransmissionInformation, err := extractAndBuildDynamicDlTransmissionInformation(cellId, extIE)

	if err != nil {
		return err
	}

	cellLoadInformation.DynamicDlTransmissionInformation = dynamicDLTransmissionInformation
	return nil
}

func populateCellLoadInformationExtensionIEs(cellId string, cellLoadInformation *entities.CellLoadInformation, extIEs []*C.CellInformation_Item_ExtIEs_t) error {
	for _, extIE := range extIEs {

		populator, ok := populators[extIE.extensionValue.present]

		if (!ok) {
			continue
		}

		err := populator(cellId, cellLoadInformation, extIE)

		if err != nil {
			return err
		}

	}

	return nil
}

func (*EnbLoadInformationExtractor) ExtractAndBuildRanLoadInformation(pdu *C.E2AP_PDU_t, ranLoadInformation *entities.RanLoadInformation) error {

	defer C.delete_pdu(pdu)

	pduCellInformationItemIEs, err := extractPduCellInformationItemIEs(pdu)

	if err != nil {
		return err
	}

	for _, pduCellInformationItemIE := range pduCellInformationItemIEs {
		err = populateCellLoadInformation(pduCellInformationItemIE, ranLoadInformation)

		if err != nil {
			return err
		}
	}

	return nil
}
*/