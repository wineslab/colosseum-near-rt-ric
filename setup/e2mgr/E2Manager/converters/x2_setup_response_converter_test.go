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

package converters

import (
	"e2mgr/logger"
	"fmt"
	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
	"strings"
	"testing"
)

/*
Test permutations of x2 setup response to protobuf enb
*/

func TestUnpackX2SetupResponseAndExtract(t *testing.T) {
	logger, _ := logger.InitLogger(logger.InfoLevel)

	var testCases = []struct {
		key       *entities.GlobalNbId
		enb       string
		packedPdu string
		failure   error
	}{
		{
			key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "007a80"},
			enb: "CONNECTED MACRO_ENB [pci:99 cell_id:\"02f829:0007ab50\" tac:\"0102\" broadcast_plmns:\"02f829\" choice_eutra_mode:<fdd:<ulear_fcn:1 dlear_fcn:1 ul_transmission_bandwidth:BW50 dl_transmission_bandwidth:BW50 > > eutra_mode:FDD ] []",
			/*
				X2AP-PDU:
				 successfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupResponse
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x15
				     criticality_t = 0
				     GlobalENB-ID
				      pLMN_Identity_t = 02 f8 29
				      eNB_ID_t:
				       macro_eNB_ID_t = 00 7a 80 (20 bits)
				    ProtocolIE_Container_elm
				     id_t = 0x14
				     criticality_t = 0
				     ServedCells:
				      ServedCells_elm
				       servedCellInfo_t
				        pCI_t = 0x63
				        cellId_t
				         pLMN_Identity_t = 02 f8 29
				         eUTRANcellIdentifier_t = 00 07 ab 50 (28 bits)
				        tAC_t = 01 02
				        broadcastPLMNs_t:
				         BroadcastPLMNs_Item_elm = 02 f8 29
				        eUTRA_Mode_Info_t:
				         fDD_t
				          uL_EARFCN_t = 0x1
				          dL_EARFCN_t = 0x1
				          uL_Transmission_Bandwidth_t = 0x3
				          dL_Transmission_Bandwidth_t = 0x3`
			*/
			packedPdu: "2006002a000002001500080002f82900007a8000140017000000630002f8290007ab50102002f829000001000133"},
		{
			key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "007a80"},
			enb: "CONNECTED MACRO_ENB [pci:99 cell_id:\"02f829:0007ab50\" tac:\"0102\" broadcast_plmns:\"02f829\" choice_eutra_mode:<fdd:<ulear_fcn:1 dlear_fcn:1 ul_transmission_bandwidth:BW50 dl_transmission_bandwidth:BW50 > > eutra_mode:FDD  pci:100 cell_id:\"02f929:0007ac50\" tac:\"0203\" broadcast_plmns:\"02f829\" broadcast_plmns:\"02f929\" choice_eutra_mode:<fdd:<ulear_fcn:2 dlear_fcn:3 ul_transmission_bandwidth:BW75 dl_transmission_bandwidth:BW75 > > eutra_mode:FDD ] []",
			/*
				X2AP-PDU:
				 successfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupResponse
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x15
				     criticality_t = 0
				     GlobalENB-ID
				      pLMN_Identity_t = 02 f8 29
				      eNB_ID_t:
				       macro_eNB_ID_t = 00 7a 80 (20 bits)
				    ProtocolIE_Container_elm
				     id_t = 0x14
				     criticality_t = 0
				     ServedCells:
				      ServedCells_elm
				       servedCellInfo_t
				        pCI_t = 0x63
				        cellId_t
				         pLMN_Identity_t = 02 f8 29
				         eUTRANcellIdentifier_t = 00 07 ab 50 (28 bits)
				        tAC_t = 01 02
				        broadcastPLMNs_t:
				         BroadcastPLMNs_Item_elm = 02 f8 29
				        eUTRA_Mode_Info_t:
				         fDD_t
				          uL_EARFCN_t = 0x1
				          dL_EARFCN_t = 0x1
				          uL_Transmission_Bandwidth_t = 0x3
				          dL_Transmission_Bandwidth_t = 0x3
				      ServedCells_elm
				       servedCellInfo_t
				        pCI_t = 0x64
				        cellId_t
				         pLMN_Identity_t = 02 f9 29
				         eUTRANcellIdentifier_t = 00 07 ac 50 (28 bits)
				        tAC_t = 02 03
				        broadcastPLMNs_t:
				         BroadcastPLMNs_Item_elm = 02 f8 29
				         BroadcastPLMNs_Item_elm = 02 f9 29
				        eUTRA_Mode_Info_t:
				         fDD_t
				          uL_EARFCN_t = 0x2
				          dL_EARFCN_t = 0x3
				          uL_Transmission_Bandwidth_t = 0x4
				          dL_Transmission_Bandwidth_t = 0x4
			*/
			packedPdu: "20060043000002001500080002f82900007a8000140030010000630002f8290007ab50102002f8290000010001330000640002f9290007ac50203202f82902f929000002000344"},

		{key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "007a80"},
			enb: "CONNECTED MACRO_ENB [pci:99 cell_id:\"02f829:0007ab50\" tac:\"0102\" broadcast_plmns:\"02f829\" choice_eutra_mode:<fdd:<ulear_fcn:2 dlear_fcn:1 ul_transmission_bandwidth:BW50 dl_transmission_bandwidth:BW50 > > eutra_mode:FDD ] []",
			/*
				X2AP-PDU:
				 successfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupResponse
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x15
				     criticality_t = 0
				     GlobalENB-ID
				      pLMN_Identity_t = 02 f8 29
				      eNB_ID_t:
				       macro_eNB_ID_t = 00 7a 80 (20 bits)
				    ProtocolIE_Container_elm
				     id_t = 0x14
				     criticality_t = 0
				     ServedCells:
				      ServedCells_elm
				       servedCellInfo_t
				        pCI_t = 0x63
				        cellId_t
				         pLMN_Identity_t = 02 f8 29
				         eUTRANcellIdentifier_t = 00 07 ab 50 (28 bits)
				        tAC_t = 01 02
				        broadcastPLMNs_t:
				         BroadcastPLMNs_Item_elm = 02 f8 29
				        eUTRA_Mode_Info_t:
				         fDD_t
				          uL_EARFCN_t = 0x1
				          dL_EARFCN_t = 0x1
				          uL_Transmission_Bandwidth_t = 0x3
				          dL_Transmission_Bandwidth_t = 0x3
				          iE_Extensions_t:
				           ProtocolExtensionContainer_elm
				            id_t = 0x5f  //ul_EARFCN
				            criticality_t = 0
				            EARFCNExtension = 0x2
			*/
			packedPdu: "20060033000002001500080002f82900007a8000140020000000630002f8290007ab50102002f8291000010001330000005f0003800102"},
		{key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "007a80"},
			enb: "CONNECTED MACRO_ENB [pci:99 cell_id:\"02f829:0007ab50\" tac:\"0102\" broadcast_plmns:\"02f829\" choice_eutra_mode:<fdd:<ulear_fcn:1 dlear_fcn:1 ul_transmission_bandwidth:BW50 dl_transmission_bandwidth:BW50 > > eutra_mode:FDD ] [02f729:0203 02f929:0304]",
			/*
				X2AP-PDU:
				 successfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupResponse
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x15
				     criticality_t = 0
				     GlobalENB-ID
				      pLMN_Identity_t = 02 f8 29
				      eNB_ID_t:
				       macro_eNB_ID_t = 00 7a 80 (20 bits)
				    ProtocolIE_Container_elm
				     id_t = 0x14
				     criticality_t = 0
				     ServedCells:
				      ServedCells_elm
				       servedCellInfo_t
				        pCI_t = 0x63
				        cellId_t
				         pLMN_Identity_t = 02 f8 29
				         eUTRANcellIdentifier_t = 00 07 ab 50 (28 bits)
				        tAC_t = 01 02
				        broadcastPLMNs_t:
				         BroadcastPLMNs_Item_elm = 02 f8 29
				        eUTRA_Mode_Info_t:
				         fDD_t
				          uL_EARFCN_t = 0x1
				          dL_EARFCN_t = 0x1
				          uL_Transmission_Bandwidth_t = 0x3
				          dL_Transmission_Bandwidth_t = 0x3
				    ProtocolIE_Container_elm
				     id_t = 0x18
				     criticality_t = 0
				     GUGroupIDList:
				      GUGroupIDList_elm
				       pLMN_Identity_t = 02 f7 29
				       mME_Group_ID_t = 02 03
				      GUGroupIDList_elm
				       pLMN_Identity_t = 02 f9 29
				       mME_Group_ID_t = 03 04
			*/
			packedPdu: "2006003a000003001500080002f82900007a8000140017000000630002f8290007ab50102002f8290000010001330018000c1002f72902030002f9290304"},

		{key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "007a80"},
			enb: "CONNECTED MACRO_ENB [pci:99 cell_id:\"02f829:0007ab50\" tac:\"0102\" broadcast_plmns:\"02f829\" choice_eutra_mode:<tdd:<ear_fcn:1 transmission_bandwidth:BW50 subframe_assignment:SA2 special_subframe_info:<special_subframe_patterns:SSP4 cyclic_prefix_dl:NORMAL cyclic_prefix_ul:EXTENDED > > > eutra_mode:TDD ] []",
			/*
				X2AP-PDU:
				 successfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupResponse
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x15
				     criticality_t = 0
				     GlobalENB-ID
				      pLMN_Identity_t = 02 f8 29
				      eNB_ID_t:
				       macro_eNB_ID_t = 00 7a 80 (20 bits)
				    ProtocolIE_Container_elm
				     id_t = 0x14
				     criticality_t = 0
				     ServedCells:
				      ServedCells_elm
				       servedCellInfo_t
				        pCI_t = 0x63
				        cellId_t
				         pLMN_Identity_t = 02 f8 29
				         eUTRANcellIdentifier_t = 00 07 ab 50 (28 bits)
				        tAC_t = 01 02
				        broadcastPLMNs_t:
				         BroadcastPLMNs_Item_elm = 02 f8 29
				        eUTRA_Mode_Info_t:
				         tDD_t
				          eARFCN_t = 0x1
				          transmission_Bandwidth_t = 0x3
				          subframeAssignment_t = 0x2
				          specialSubframe_Info_t
				           specialSubframePatterns_t = 0x4
				           cyclicPrefixDL_t = 0
				           cyclicPrefixUL_t = 0x1
			*/
			packedPdu: "2006002a000002001500080002f82900007a8000140017000000630002f8290007ab50102002f829400001320820"},
		{key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "007a80"},
			enb: "CONNECTED MACRO_ENB [pci:99 cell_id:\"02f829:0007ab50\" tac:\"0102\" broadcast_plmns:\"02f829\" choice_eutra_mode:<tdd:<ear_fcn:1 transmission_bandwidth:BW50 subframe_assignment:SA2 special_subframe_info:<special_subframe_patterns:SSP4 cyclic_prefix_dl:EXTENDED cyclic_prefix_ul:NORMAL > additional_special_subframe_info:<additional_special_subframe_patterns:SSP9 cyclic_prefix_dl:NORMAL cyclic_prefix_ul:EXTENDED > > > eutra_mode:TDD ] []",
			/*
				X2AP-PDU:
				 successfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupResponse
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x15
				     criticality_t = 0
				     GlobalENB-ID
				      pLMN_Identity_t = 02 f8 29
				      eNB_ID_t:
				       macro_eNB_ID_t = 00 7a 80 (20 bits)
				    ProtocolIE_Container_elm
				     id_t = 0x14
				     criticality_t = 0
				     ServedCells:
				      ServedCells_elm
				       servedCellInfo_t
				        pCI_t = 0x63
				        cellId_t
				         pLMN_Identity_t = 02 f8 29
				         eUTRANcellIdentifier_t = 00 07 ab 50 (28 bits)
				        tAC_t = 01 02
				        broadcastPLMNs_t:
				         BroadcastPLMNs_Item_elm = 02 f8 29
				        eUTRA_Mode_Info_t:
				         tDD_t
				          eARFCN_t = 0x1
				          transmission_Bandwidth_t = 0x3
				          subframeAssignment_t = 0x2
				          specialSubframe_Info_t
				           specialSubframePatterns_t = 0x4
				           cyclicPrefixDL_t = 0x1
				           cyclicPrefixUL_t = 0
				          iE_Extensions_t:
				           ProtocolExtensionContainer_elm
				            id_t = 0x61
				            criticality_t = 0x1
				            AdditionalSpecialSubframe-Info
				             additionalspecialSubframePatterns_t = 0x9
				             cyclicPrefixDL_t = 0
				             cyclicPrefixUL_t = 0x1
			*/
			packedPdu: "20060032000002001500080002f82900007a800014001f000000630002f8290007ab50102002f8295000013208800000006140021220"},

		{key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "007a80"},
			enb: "CONNECTED MACRO_ENB [pci:99 cell_id:\"02f829:0007ab50\" tac:\"0102\" broadcast_plmns:\"02f829\" choice_eutra_mode:<tdd:<ear_fcn:2 transmission_bandwidth:BW50 subframe_assignment:SA2 special_subframe_info:<special_subframe_patterns:SSP4 cyclic_prefix_dl:EXTENDED cyclic_prefix_ul:NORMAL > > > eutra_mode:TDD ] []",
			/*
				X2AP-PDU:
				 successfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupResponse
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x15
				     criticality_t = 0
				     GlobalENB-ID
				      pLMN_Identity_t = 02 f8 29
				      eNB_ID_t:
				       macro_eNB_ID_t = 00 7a 80 (20 bits)
				    ProtocolIE_Container_elm
				     id_t = 0x14
				     criticality_t = 0
				     ServedCells:
				      ServedCells_elm
				       servedCellInfo_t
				        pCI_t = 0x63
				        cellId_t
				         pLMN_Identity_t = 02 f8 29
				         eUTRANcellIdentifier_t = 00 07 ab 50 (28 bits)
				        tAC_t = 01 02
				        broadcastPLMNs_t:
				         BroadcastPLMNs_Item_elm = 02 f8 29
				        eUTRA_Mode_Info_t:
				         tDD_t
				          eARFCN_t = 0x1
				          transmission_Bandwidth_t = 0x3
				          subframeAssignment_t = 0x2
				          specialSubframe_Info_t
				           specialSubframePatterns_t = 0x4
				           cyclicPrefixDL_t = 0x1
				           cyclicPrefixUL_t = 0
				          iE_Extensions_t:
				           ProtocolExtensionContainer_elm
				            id_t = 0x5e
				            criticality_t = 0
				            EARFCNExtension = 0x2
			*/
			packedPdu: "20060033000002001500080002f82900007a8000140020000000630002f8290007ab50102002f8295000013208800000005e0003800102"},

		{key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "007a80"},
			enb: "CONNECTED MACRO_ENB [pci:99 cell_id:\"02f829:0007ab50\" tac:\"0102\" broadcast_plmns:\"02f829\" choice_eutra_mode:<tdd:<ear_fcn:1 transmission_bandwidth:BW50 subframe_assignment:SA2 special_subframe_info:<special_subframe_patterns:SSP4 cyclic_prefix_dl:EXTENDED cyclic_prefix_ul:NORMAL > additional_special_subframe_extension_info:<additional_special_subframe_patterns_extension:SSP10 cyclic_prefix_dl:NORMAL cyclic_prefix_ul:NORMAL > > > eutra_mode:TDD ] []",
			/*
				X2AP-PDU:
				 successfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupResponse
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x15
				     criticality_t = 0
				     GlobalENB-ID
				      pLMN_Identity_t = 02 f8 29
				      eNB_ID_t:
				       macro_eNB_ID_t = 00 7a 80 (20 bits)
				    ProtocolIE_Container_elm
				     id_t = 0x14
				     criticality_t = 0
				     ServedCells:
				      ServedCells_elm
				       servedCellInfo_t
				        pCI_t = 0x63
				        cellId_t
				         pLMN_Identity_t = 02 f8 29
				         eUTRANcellIdentifier_t = 00 07 ab 50 (28 bits)
				        tAC_t = 01 02
				        broadcastPLMNs_t:
				         BroadcastPLMNs_Item_elm = 02 f8 29
				        eUTRA_Mode_Info_t:
				         tDD_t
				          eARFCN_t = 0x1
				          transmission_Bandwidth_t = 0x3
				          subframeAssignment_t = 0x2
				          specialSubframe_Info_t
				           specialSubframePatterns_t = 0x4
				           cyclicPrefixDL_t = 0x1
				           cyclicPrefixUL_t = 0
				          iE_Extensions_t:
				           ProtocolExtensionContainer_elm
				            id_t = 0xb3
				            criticality_t = 0x1
				            AdditionalSpecialSubframeExtension-Info
				             additionalspecialSubframePatternsExtension_t = 0
				             cyclicPrefixDL_t = 0
				             cyclicPrefixUL_t = 0
			*/
			packedPdu: "20060031000002001500080002f82900007a800014001e000000630002f8290007ab50102002f829500001320880000000b3400100"},

		{key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "007a80"},
			enb: "CONNECTED MACRO_ENB [pci:99 cell_id:\"02f829:0007ab50\" tac:\"0102\" broadcast_plmns:\"02f829\" choice_eutra_mode:<tdd:<ear_fcn:2 transmission_bandwidth:BW50 subframe_assignment:SA2 special_subframe_info:<special_subframe_patterns:SSP4 cyclic_prefix_dl:EXTENDED cyclic_prefix_ul:NORMAL > additional_special_subframe_info:<additional_special_subframe_patterns:SSP9 cyclic_prefix_dl:NORMAL cyclic_prefix_ul:EXTENDED > additional_special_subframe_extension_info:<additional_special_subframe_patterns_extension:SSP10 cyclic_prefix_dl:NORMAL cyclic_prefix_ul:NORMAL > > > eutra_mode:TDD ] []",
			/*
				X2AP-PDU:
				 successfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupResponse
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x15
				     criticality_t = 0
				     GlobalENB-ID
				      pLMN_Identity_t = 02 f8 29
				      eNB_ID_t:
				       macro_eNB_ID_t = 00 7a 80 (20 bits)
				    ProtocolIE_Container_elm
				     id_t = 0x14
				     criticality_t = 0
				     ServedCells:
				      ServedCells_elm
				       servedCellInfo_t
				        pCI_t = 0x63
				        cellId_t
				         pLMN_Identity_t = 02 f8 29
				         eUTRANcellIdentifier_t = 00 07 ab 50 (28 bits)
				        tAC_t = 01 02
				        broadcastPLMNs_t:
				         BroadcastPLMNs_Item_elm = 02 f8 29
				        eUTRA_Mode_Info_t:
				         tDD_t
				          eARFCN_t = 0x1
				          transmission_Bandwidth_t = 0x3
				          subframeAssignment_t = 0x2
				          specialSubframe_Info_t
				           specialSubframePatterns_t = 0x4
				           cyclicPrefixDL_t = 0x1
				           cyclicPrefixUL_t = 0
				          iE_Extensions_t:
				           ProtocolExtensionContainer_elm
				            id_t = 0xb3
				            criticality_t = 0x1
				            AdditionalSpecialSubframeExtension-Info
				             additionalspecialSubframePatternsExtension_t = 0
				             cyclicPrefixDL_t = 0
				             cyclicPrefixUL_t = 0
				           ProtocolExtensionContainer_elm
				            id_t = 0x61
				            criticality_t = 0x1
				            AdditionalSpecialSubframe-Info
				             additionalspecialSubframePatterns_t = 0x9
				             cyclicPrefixDL_t = 0
				             cyclicPrefixUL_t = 0x1
				           ProtocolExtensionContainer_elm
				            id_t = 0x5e
				            criticality_t = 0
				            EARFCNExtension = 0x2

			*/
			packedPdu: "2006003e000002001500080002f82900007a800014002b000000630002f8290007ab50102002f829500001320880000200b3400100006140021220005e0003800102"},

		{key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "007a80b0"},
			enb: "CONNECTED HOME_ENB [pci:99 cell_id:\"02f829:0007ab50\" tac:\"0102\" broadcast_plmns:\"02f829\" choice_eutra_mode:<fdd:<ulear_fcn:1 dlear_fcn:1 ul_transmission_bandwidth:BW50 dl_transmission_bandwidth:BW50 > > eutra_mode:FDD number_of_antenna_ports:AN1  pci:100 cell_id:\"02f929:0007ac50\" tac:\"0203\" broadcast_plmns:\"02f829\" broadcast_plmns:\"02f929\" choice_eutra_mode:<fdd:<ulear_fcn:2 dlear_fcn:3 ul_transmission_bandwidth:BW75 dl_transmission_bandwidth:BW75 > > eutra_mode:FDD ] []",
			/*
				X2AP-PDU:
				 successfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupResponse
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x15
				     criticality_t = 0
				     GlobalENB-ID
				      pLMN_Identity_t = 02 f8 29
				      eNB_ID_t:
				       home_eNB_ID_t = 00 7a 80 b0 (28 bits)
				    ProtocolIE_Container_elm
				     id_t = 0x14
				     criticality_t = 0
				     ServedCells:
				      ServedCells_elm
				       servedCellInfo_t
				        pCI_t = 0x63
				        cellId_t
				         pLMN_Identity_t = 02 f8 29
				         eUTRANcellIdentifier_t = 00 07 ab 50 (28 bits)
				        tAC_t = 01 02
				        broadcastPLMNs_t:
				         BroadcastPLMNs_Item_elm = 02 f8 29
				        eUTRA_Mode_Info_t:
				         fDD_t
				          uL_EARFCN_t = 0x1
				          dL_EARFCN_t = 0x1
				          uL_Transmission_Bandwidth_t = 0x3
				          dL_Transmission_Bandwidth_t = 0x3
				        iE_Extensions_t:
				         ProtocolExtensionContainer_elm
				          id_t = 0x29
				          criticality_t = 0x1
				          Number-of-Antennaports = 0
				      ServedCells_elm
				       servedCellInfo_t
				        pCI_t = 0x64
				        cellId_t
				         pLMN_Identity_t = 02 f9 29
				         eUTRANcellIdentifier_t = 00 07 ac 50 (28 bits)
				        tAC_t = 02 03
				        broadcastPLMNs_t:
				         BroadcastPLMNs_Item_elm = 02 f8 29
				         BroadcastPLMNs_Item_elm = 02 f9 29
				        eUTRA_Mode_Info_t:nb_id
				         fDD_t
				          uL_EARFCN_t = 0x2
				          dL_EARFCN_t = 0x3
				          uL_Transmission_Bandwidth_t = 0x4
				          dL_Transmission_Bandwidth_t = 0x4

			*/
			packedPdu: "2006004b000002001500090002f82940007a80b000140037010800630002f8290007ab50102002f829000001000133000000294001000000640002f9290007ac50203202f82902f929000002000344"},

		{key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "007a40"},
			enb: "CONNECTED SHORT_MACRO_ENB [pci:99 cell_id:\"02f829:0007ab50\" tac:\"0102\" broadcast_plmns:\"02f829\" choice_eutra_mode:<fdd:<ulear_fcn:1 dlear_fcn:1 ul_transmission_bandwidth:BW50 dl_transmission_bandwidth:BW50 > > eutra_mode:FDD number_of_antenna_ports:AN1  pci:100 cell_id:\"02f929:0007ac50\" tac:\"0203\" broadcast_plmns:\"02f829\" broadcast_plmns:\"02f929\" choice_eutra_mode:<fdd:<ulear_fcn:2 dlear_fcn:3 ul_transmission_bandwidth:BW75 dl_transmission_bandwidth:BW75 > > eutra_mode:FDD prach_configuration:<root_sequence_index:15 zero_correlation_zone_configuration:7 high_speed_flag:true prach_frequency_offset:30 > ] []",
			/*
				X2AP-PDU:
				 successfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupResponse
				   protocolIEs_t:
					ProtocolIE_Container_elm
					 id_t = 0x15
					 criticality_t = 0
					 GlobalENB-ID
					  pLMN_Identity_t = 02 f8 29
					  eNB_ID_t:
					   short_Macro_eNB_ID_t = 00 7a 40 (18 bits)
					ProtocolIE_Container_elm
					 id_t = 0x14
					 criticality_t = 0
					 ServedCells:
					  ServedCells_elm
					   servedCellInfo_t
						pCI_t = 0x63
						cellId_t
						 pLMN_Identity_t = 02 f8 29
						 eUTRANcellIdentifier_t = 00 07 ab 50 (28 bits)
						tAC_t = 01 02
						broadcastPLMNs_t:
						 BroadcastPLMNs_Item_elm = 02 f8 29
						eUTRA_Mode_Info_t:
						 fDD_t
						  uL_EARFCN_t = 0x1
						  dL_EARFCN_t = 0x1
						  uL_Transmission_Bandwidth_t = 0x3
						  dL_Transmission_Bandwidth_t = 0x3
						iE_Extensions_t:
						 ProtocolExtensionContainer_elm
						  id_t = 0x29
						  criticality_t = 0x1
						  Number-of-Antennaports = 0
					  ServedCells_elm
					   servedCellInfo_t
						pCI_t = 0x64
						cellId_t
						 pLMN_Identity_t = 02 f9 29
						 eUTRANcellIdentifier_t = 00 07 ac 50 (28 bits)
						tAC_t = 02 03
						broadcastPLMNs_t:
						 BroadcastPLMNs_Item_elm = 02 f8 29
						 BroadcastPLMNs_Item_elm = 02 f9 29
						eUTRA_Mode_Info_t:
						 fDD_t
						  uL_EARFCN_t = 0x2
						  dL_EARFCN_t = 0x3
						  uL_Transmission_Bandwidth_t = 0x4
						  dL_Transmission_Bandwidth_t = 0x4
						iE_Extensions_t:
						 ProtocolExtensionContainer_elm
						  id_t = 0x37
						  criticality_t = 0x1
						  PRACH-Configuration
						   rootSequenceIndex_t = 0xf
						   zeroCorrelationIndex_t = 0x7
						   highSpeedFlag_t = true
						   prach_FreqOffset_t = 0x1e

			*/
			packedPdu: "20060056000002001500090002f8298003007a4000140042010800630002f8290007ab50102002f829000001000133000000294001000800640002f9290007ac50203202f82902f92900000200034400000037400500000f79e0"},

		{key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "007a40"},
			enb: "CONNECTED SHORT_MACRO_ENB [pci:99 cell_id:\"02f829:0007ab50\" tac:\"0102\" broadcast_plmns:\"02f829\" choice_eutra_mode:<fdd:<ulear_fcn:1 dlear_fcn:1 ul_transmission_bandwidth:BW50 dl_transmission_bandwidth:BW50 > > eutra_mode:FDD number_of_antenna_ports:AN1  pci:100 cell_id:\"02f929:0007ac50\" tac:\"0203\" broadcast_plmns:\"02f829\" broadcast_plmns:\"02f929\" choice_eutra_mode:<fdd:<ulear_fcn:2 dlear_fcn:3 ul_transmission_bandwidth:BW75 dl_transmission_bandwidth:BW75 > > eutra_mode:FDD mbsfn_subframe_infos:<radioframe_allocation_period:N8 radioframe_allocation_offset:3 subframe_allocation:\"28\" subframe_allocation_type:ONE_FRAME > ] []",
			/*
				X2AP-PDU:
				 successfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupResponse
				   protocolIEs_t:
					ProtocolIE_Container_elm
					 id_t = 0x15
					 criticality_t = 0
					 GlobalENB-ID
					  pLMN_Identity_t = 02 f8 29
					  eNB_ID_t:
					   short_Macro_eNB_ID_t = 00 7a 40 (18 bits)
					ProtocolIE_Container_elm
					 id_t = 0x14
					 criticality_t = 0
					 ServedCells:
					  ServedCells_elm
					   servedCellInfo_t
						pCI_t = 0x63
						cellId_t
						 pLMN_Identity_t = 02 f8 29
						 eUTRANcellIdentifier_t = 00 07 ab 50 (28 bits)
						tAC_t = 01 02
						broadcastPLMNs_t:
						 BroadcastPLMNs_Item_elm = 02 f8 29
						eUTRA_Mode_Info_t:
						 fDD_t
						  uL_EARFCN_t = 0x1
						  dL_EARFCN_t = 0x1
						  uL_Transmission_Bandwidth_t = 0x3
						  dL_Transmission_Bandwidth_t = 0x3
						iE_Extensions_t:
						 ProtocolExtensionContainer_elm
						  id_t = 0x29
						  criticality_t = 0x1
						  Number-of-Antennaports = 0
					  ServedCells_elm
					   servedCellInfo_t
						pCI_t = 0x64
						cellId_t
						 pLMN_Identity_t = 02 f9 29
						 eUTRANcellIdentifier_t = 00 07 ac 50 (28 bits)
						tAC_t = 02 03
						broadcastPLMNs_t:
						 BroadcastPLMNs_Item_elm = 02 f8 29
						 BroadcastPLMNs_Item_elm = 02 f9 29
						eUTRA_Mode_Info_t:
						 fDD_t
						  uL_EARFCN_t = 0x2
						  dL_EARFCN_t = 0x3
						  uL_Transmission_Bandwidth_t = 0x4
						  dL_Transmission_Bandwidth_t = 0x4
						iE_Extensions_t:
						 ProtocolExtensionContainer_elm
						  id_t = 0x38
						  criticality_t = 0x1
						  MBSFN-Subframe-Infolist:
						   MBSFN_Subframe_Infolist_elm
							radioframeAllocationPeriod_t = 0x3
							radioframeAllocationOffset_t = 0x3
							subframeAllocation_t:
							 oneframe_t = 28 (6 bits)
			*/
			packedPdu: "20060054000002001500090002f8298003007a4000140040010800630002f8290007ab50102002f829000001000133000000294001000800640002f9290007ac50203202f82902f929000002000344000000384003019850"},
		{key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "007a40"},
			enb: "CONNECTED SHORT_MACRO_ENB [pci:99 cell_id:\"02f829:0007ab50\" tac:\"0102\" broadcast_plmns:\"02f829\" choice_eutra_mode:<fdd:<ulear_fcn:1 dlear_fcn:1 ul_transmission_bandwidth:BW50 dl_transmission_bandwidth:BW50 > > eutra_mode:FDD number_of_antenna_ports:AN1 mbsfn_subframe_infos:<radioframe_allocation_period:N8 radioframe_allocation_offset:3 subframe_allocation:\"28\" subframe_allocation_type:ONE_FRAME >  pci:100 cell_id:\"02f929:0007ac50\" tac:\"0203\" broadcast_plmns:\"02f829\" broadcast_plmns:\"02f929\" choice_eutra_mode:<fdd:<ulear_fcn:2 dlear_fcn:3 ul_transmission_bandwidth:BW75 dl_transmission_bandwidth:BW75 > > eutra_mode:FDD ] []",
			/*
				X2AP-PDU:
				 successfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupResponse
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x15
				     criticality_t = 0
				     GlobalENB-ID
				      pLMN_Identity_t = 02 f8 29
				      eNB_ID_t:
				       short_Macro_eNB_ID_t = 00 7a 40 (18 bits)
				    ProtocolIE_Container_elm
				     id_t = 0x14
				     criticality_t = 0
				     ServedCells:
				      ServedCells_elm
				       servedCellInfo_t
				        pCI_t = 0x63
				        cellId_t
				         pLMN_Identity_t = 02 f8 29
				         eUTRANcellIdentifier_t = 00 07 ab 50 (28 bits)
				        tAC_t = 01 02
				        broadcastPLMNs_t:
				         BroadcastPLMNs_Item_elm = 02 f8 29
				        eUTRA_Mode_Info_t:
				         fDD_t
				          uL_EARFCN_t = 0x1
				          dL_EARFCN_t = 0x1
				          uL_Transmission_Bandwidth_t = 0x3
				          dL_Transmission_Bandwidth_t = 0x3
				        iE_Extensions_t:
				         ProtocolExtensionContainer_elm
				          id_t = 0x29
				          criticality_t = 0x1
				          Number-of-Antennaports = 0
				         ProtocolExtensionContainer_elm
				          id_t = 0x38
				          criticality_t = 0x1
				          MBSFN-Subframe-Infolist:
				           MBSFN_Subframe_Infolist_elm
				            radioframeAllocationPeriod_t = 0x3
				            radioframeAllocationOffset_t = 0x3
				            subframeAllocation_t:
				             oneframe_t = 28 (6 bits)
				      ServedCells_elm
				       servedCellInfo_t
				        pCI_t = 0x64
				        cellId_t
				         pLMN_Identity_t = 02 f9 29
				         eUTRANcellIdentifier_t = 00 07 ac 50 (28 bits)
				        tAC_t = 02 03
				        broadcastPLMNs_t:
				         BroadcastPLMNs_Item_elm = 02 f8 29
				         BroadcastPLMNs_Item_elm = 02 f9 29
				        eUTRA_Mode_Info_t:
				         fDD_t
				          uL_EARFCN_t = 0x2
				          dL_EARFCN_t = 0x3
				          uL_Transmission_Bandwidth_t = 0x4
				          dL_Transmission_Bandwidth_t = 0x4
			*/
			packedPdu: "20060052000002001500090002f8298003007a400014003e010800630002f8290007ab50102002f82900000100013300010029400100003840030198500000640002f9290007ac50203202f82902f929000002000344"},
		{key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "007a40"},
			enb: "CONNECTED SHORT_MACRO_ENB [pci:99 cell_id:\"02f829:0007ab50\" tac:\"0102\" broadcast_plmns:\"02f829\" choice_eutra_mode:<fdd:<ulear_fcn:1 dlear_fcn:1 ul_transmission_bandwidth:BW50 dl_transmission_bandwidth:BW50 > > eutra_mode:FDD number_of_antenna_ports:AN1 prach_configuration:<root_sequence_index:15 zero_correlation_zone_configuration:7 high_speed_flag:true prach_frequency_offset:30 >  pci:100 cell_id:\"02f929:0007ac50\" tac:\"0203\" broadcast_plmns:\"02f829\" broadcast_plmns:\"02f929\" choice_eutra_mode:<fdd:<ulear_fcn:2 dlear_fcn:3 ul_transmission_bandwidth:BW75 dl_transmission_bandwidth:BW75 > > eutra_mode:FDD ] []",
			/*
				X2AP-PDU:
				 successfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupResponse
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x15
				     criticality_t = 0
				     GlobalENB-ID
				      pLMN_Identity_t = 02 f8 29
				      eNB_ID_t:
				       short_Macro_eNB_ID_t = 00 7a 40 (18 bits)
				    ProtocolIE_Container_elm
				     id_t = 0x14
				     criticality_t = 0
				     ServedCells:
				      ServedCells_elm
				       servedCellInfo_t
				        pCI_t = 0x63
				        cellId_t
				         pLMN_Identity_t = 02 f8 29
				         eUTRANcellIdentifier_t = 00 07 ab 50 (28 bits)
				        tAC_t = 01 02
				        broadcastPLMNs_t:
				         BroadcastPLMNs_Item_elm = 02 f8 29
				        eUTRA_Mode_Info_t:
				         fDD_t
				          uL_EARFCN_t = 0x1
				          dL_EARFCN_t = 0x1
				          uL_Transmission_Bandwidth_t = 0x3
				          dL_Transmission_Bandwidth_t = 0x3
				        iE_Extensions_t:
				         ProtocolExtensionContainer_elm
				          id_t = 0x29
				          criticality_t = 0x1
				          Number-of-Antennaports = 0
				         ProtocolExtensionContainer_elm
				          id_t = 0x37
				          criticality_t = 0x1
				          PRACH-Configuration
				           rootSequenceIndex_t = 0xf
				           zeroCorrelationIndex_t = 0x7
				           highSpeedFlag_t = true
				           prach_FreqOffset_t = 0x1e
				      ServedCells_elm
				       servedCellInfo_t
				        pCI_t = 0x64
				        cellId_t
				         pLMN_Identity_t = 02 f9 29
				         eUTRANcellIdentifier_t = 00 07 ac 50 (28 bits)
				        tAC_t = 02 03
				        broadcastPLMNs_t:
				         BroadcastPLMNs_Item_elm = 02 f8 29
				         BroadcastPLMNs_Item_elm = 02 f9 29
				        eUTRA_Mode_Info_t:
				         fDD_t
				          uL_EARFCN_t = 0x2
				          dL_EARFCN_t = 0x3
				          uL_Transmission_Bandwidth_t = 0x4
				          dL_Transmission_Bandwidth_t = 0x4
			*/
			packedPdu: "20060054000002001500090002f8298003007a4000140040010800630002f8290007ab50102002f829000001000133000100294001000037400500000f79e00000640002f9290007ac50203202f82902f929000002000344"},
		{key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "007a08"},
			enb: "CONNECTED LONG_MACRO_ENB [pci:99 cell_id:\"02f829:0007ab50\" tac:\"0102\" broadcast_plmns:\"02f829\" choice_eutra_mode:<fdd:<ulear_fcn:1 dlear_fcn:1 ul_transmission_bandwidth:BW50 dl_transmission_bandwidth:BW50 > > eutra_mode:FDD prach_configuration:<root_sequence_index:15 zero_correlation_zone_configuration:7 high_speed_flag:true prach_frequency_offset:30 prach_configuration_index:60 >  pci:100 cell_id:\"02f929:0007ac50\" tac:\"0203\" broadcast_plmns:\"02f829\" broadcast_plmns:\"02f929\" choice_eutra_mode:<fdd:<ulear_fcn:2 dlear_fcn:3 ul_transmission_bandwidth:BW75 dl_transmission_bandwidth:BW75 > > eutra_mode:FDD ] []",
			/*
				X2AP-PDU:
				 successfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupResponse
				   protocolIEs_t:
					ProtocolIE_Container_elm
					 id_t = 0x15
					 criticality_t = 0
					 GlobalENB-ID
					  pLMN_Identity_t = 02 f8 29
					  eNB_ID_t:
					   long_Macro_eNB_ID_t = 00 7a 08 (21 bits)
					ProtocolIE_Container_elm
					 id_t = 0x14
					 criticality_t = 0
					 ServedCells:
					  ServedCells_elm
					   servedCellInfo_t
						pCI_t = 0x63
						cellId_t
						 pLMN_Identity_t = 02 f8 29
						 eUTRANcellIdentifier_t = 00 07 ab 50 (28 bits)
						tAC_t = 01 02
						broadcastPLMNs_t:
						 BroadcastPLMNs_Item_elm = 02 f8 29
						eUTRA_Mode_Info_t:
						 fDD_t
						  uL_EARFCN_t = 0x1
						  dL_EARFCN_t = 0x1
						  uL_Transmission_Bandwidth_t = 0x3
						  dL_Transmission_Bandwidth_t = 0x3
						iE_Extensions_t:
						 ProtocolExtensionContainer_elm
						  id_t = 0x37
						  criticality_t = 0x1
						  PRACH-Configuration
						   rootSequenceIndex_t = 0xf
						   zeroCorrelationIndex_t = 0x7
						   highSpeedFlag_t = true
						   prach_FreqOffset_t = 0x1e
						   prach_ConfigIndex_t = 0x3c
					  ServedCells_elm
					   servedCellInfo_t
						pCI_t = 0x64
						cellId_t
						 pLMN_Identity_t = 02 f9 29
						 eUTRANcellIdentifier_t = 00 07 ac 50 (28 bits)
						tAC_t = 02 03
						broadcastPLMNs_t:
						 BroadcastPLMNs_Item_elm = 02 f8 29
						 BroadcastPLMNs_Item_elm = 02 f9 29
						eUTRA_Mode_Info_t:
						 fDD_t
						  uL_EARFCN_t = 0x2
						  dL_EARFCN_t = 0x3
						  uL_Transmission_Bandwidth_t = 0x4
						  dL_Transmission_Bandwidth_t = 0x4
			*/
			//packedPdu: "20060050000002001500090002f8298103007a080014003c010800630002f8290007ab50102002f82900000100013300000037400640000f79ef000000640002f9290007ac50203202f82902f929000002000344"},
			packedPdu: "20060050000002001500090002f829c003007a080014003c010800630002f8290007ab50102002f82900000100013300000037400640000f79ef000000640002f9290007ac50203202f82902f929000002000344"},

		{key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "007a80"},
			enb: "CONNECTED MACRO_ENB [pci:99 cell_id:\"02f829:0007ab50\" tac:\"0102\" broadcast_plmns:\"02f829\" choice_eutra_mode:<fdd:<ulear_fcn:1 dlear_fcn:1 ul_transmission_bandwidth:BW50 dl_transmission_bandwidth:BW50 > > eutra_mode:FDD csg_id:\"0007aba0\" freq_band_indicator_priority:BROADCASTED bandwidth_reduced_si:SCHEDULED ] []",
			/*
				X2AP-PDU:
				 successfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupResponse
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x15
				     criticality_t = 0
				     GlobalENB-ID
				      pLMN_Identity_t = 02 f8 29
				      eNB_ID_t:
				       macro_eNB_ID_t = 00 7a 80 (20 bits)
				    ProtocolIE_Container_elm
				     id_t = 0x14
				     criticality_t = 0
				     ServedCells:
				      ServedCells_elm
				       servedCellInfo_t
				        pCI_t = 0x63
				        cellId_t
				         pLMN_Identity_t = 02 f8 29
				         eUTRANcellIdentifier_t = 00 07 ab 50 (28 bits)
				        tAC_t = 01 02
				        broadcastPLMNs_t:
				         BroadcastPLMNs_Item_elm = 02 f8 29
				        eUTRA_Mode_Info_t:
				         fDD_t
				          uL_EARFCN_t = 0x1
				          dL_EARFCN_t = 0x1
				          uL_Transmission_Bandwidth_t = 0x3
				          dL_Transmission_Bandwidth_t = 0x3
				        iE_Extensions_t:
				         ProtocolExtensionContainer_elm
				          id_t = 0x46
				          criticality_t = 0x1
				          CSG-Id = 00 07 ab a0 (27 bits)
				         ProtocolExtensionContainer_elm
				          id_t = 0xa0
				          criticality_t = 0x1
				          FreqBandIndicatorPriority = 0x1
				         ProtocolExtensionContainer_elm
				          id_t = 0xb4
				          criticality_t = 0x1
				          BandwidthReducedSI = 0
			*/
			packedPdu: "2006003e000002001500080002f82900007a800014002b000800630002f8290007ab50102002f8290000010001330002004640040007aba000a040014000b4400100"},

		{key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "007a80"},
			enb: "CONNECTED MACRO_ENB [pci:99 cell_id:\"02f829:0007ab50\" tac:\"0102\" broadcast_plmns:\"02f829\" choice_eutra_mode:<fdd:<ulear_fcn:1 dlear_fcn:1 ul_transmission_bandwidth:BW50 dl_transmission_bandwidth:BW50 > > eutra_mode:FDD mbms_service_area_identities:\"02f8\" mbms_service_area_identities:\"03f9\" multiband_infos:1 multiband_infos:2 multiband_infos:3 freq_band_indicator_priority:NOT_BROADCASTED ] []",
			/*
				X2AP-PDU:
				 successfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupResponse
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x15
				     criticality_t = 0
				     GlobalENB-ID
				      pLMN_Identity_t = 02 f8 29
				      eNB_ID_t:
				       macro_eNB_ID_t = 00 7a 80 (20 bits)
				    ProtocolIE_Container_elm
				     id_t = 0x14
				     criticality_t = 0
				     ServedCells:
				      ServedCells_elm
				       servedCellInfo_t
				        pCI_t = 0x63
				        cellId_t
				         pLMN_Identity_t = 02 f8 29
				         eUTRANcellIdentifier_t = 00 07 ab 50 (28 bits)
				        tAC_t = 01 02
				        broadcastPLMNs_t:
				         BroadcastPLMNs_Item_elm = 02 f8 29
				        eUTRA_Mode_Info_t:
				         fDD_t
				          uL_EARFCN_t = 0x1
				          dL_EARFCN_t = 0x1
				          uL_Transmission_Bandwidth_t = 0x3
				          dL_Transmission_Bandwidth_t = 0x3
				        iE_Extensions_t:
				         ProtocolExtensionContainer_elm
				          id_t = 0x4f
				          criticality_t = 0x1
				          MBMS-Service-Area-Identity-List:
				           MBMS_Service_Area_Identity_List_elm = 02 f8
				           MBMS_Service_Area_Identity_List_elm = 03 f9
				         ProtocolExtensionContainer_elm
				          id_t = 0xa0
				          criticality_t = 0x1
				          FreqBandIndicatorPriority = 0
				         ProtocolExtensionContainer_elm
				          id_t = 0x54
				          criticality_t = 0x1
				          MultibandInfoList:
				           MultibandInfoList_elm
				            freqBandIndicator_t = 0x1
				           MultibandInfoList_elm
				            freqBandIndicator_t = 0x2
				           MultibandInfoList_elm
				            freqBandIndicator_t = 0x3
			*/
			packedPdu: "20060044000002001500080002f82900007a8000140031000800630002f8290007ab50102002f8290000010001330002004f40050102f803f900a040010000544006200000010002"},
		{key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "007a80"},
			enb: "CONNECTED MACRO_ENB [pci:99 cell_id:\"02f829:0007ab50\" tac:\"0102\" broadcast_plmns:\"02f829\" choice_eutra_mode:<fdd:<ulear_fcn:1 dlear_fcn:1 ul_transmission_bandwidth:BW50 dl_transmission_bandwidth:BW50 > > eutra_mode:FDD neighbour_infos:<ecgi:\"02f829:0007ab50\" pci:99 ear_fcn:1 > neighbour_infos:<ecgi:\"03f930:0008bc50\" pci:100 ear_fcn:2 > ] []",
			/*
				X2AP-PDU:
				 successfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupResponse
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x15
				     criticality_t = 0
				     GlobalENB-ID
				      pLMN_Identity_t = 02 f8 29
				      eNB_ID_t:
				       macro_eNB_ID_t = 00 7a 80 (20 bits)
				    ProtocolIE_Container_elm
				     id_t = 0x14
				     criticality_t = 0
				     ServedCells:
				      ServedCells_elm
				       servedCellInfo_t
				        pCI_t = 0x63
				        cellId_t
				         pLMN_Identity_t = 02 f8 29
				         eUTRANcellIdentifier_t = 00 07 ab 50 (28 bits)
				        tAC_t = 01 02
				        broadcastPLMNs_t:
				         BroadcastPLMNs_Item_elm = 02 f8 29
				        eUTRA_Mode_Info_t:
				         fDD_t
				          uL_EARFCN_t = 0x1
				          dL_EARFCN_t = 0x1
				          uL_Transmission_Bandwidth_t = 0x3
				          dL_Transmission_Bandwidth_t = 0x3
				       neighbour_Info_t:
				        Neighbour_Information_elm
				         eCGI_t
				          pLMN_Identity_t = 02 f8 29
				          eUTRANcellIdentifier_t = 00 07 ab 50 (28 bits)
				         pCI_t = 0x63
				         eARFCN_t = 0x1
				        Neighbour_Information_elm
				         eCGI_t
				          pLMN_Identity_t = 03 f9 30
				          eUTRANcellIdentifier_t = 00 08 bc 50 (28 bits)
				         pCI_t = 0x64
				         eARFCN_t = 0x2
			*/
			packedPdu: "20060044000002001500080002f82900007a8000140031004000630002f8290007ab50102002f82900000100013300020002f8290007ab50006300010003f9300008bc5000640002"},

		{key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "007a80"},
			enb: "CONNECTED MACRO_ENB [pci:99 cell_id:\"02f829:0007ab50\" tac:\"0102\" broadcast_plmns:\"02f829\" choice_eutra_mode:<fdd:<ulear_fcn:1 dlear_fcn:1 ul_transmission_bandwidth:BW50 dl_transmission_bandwidth:BW50 > > eutra_mode:FDD neighbour_infos:<ecgi:\"02f829:0007ab50\" pci:99 ear_fcn:1 tac:\"0102\" > neighbour_infos:<ecgi:\"03f930:0008bc50\" pci:100 ear_fcn:3 > ] []",
			/*
				X2AP-PDU:
				 successfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupResponse
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x15
				     criticality_t = 0
				     GlobalENB-ID
				      pLMN_Identity_t = 02 f8 29
				      eNB_ID_t:
				       macro_eNB_ID_t = 00 7a 80 (20 bits)
				    ProtocolIE_Container_elm
				     id_t = 0x14
				     criticality_t = 0
				     ServedCells:
				      ServedCells_elm
				       servedCellInfo_t
				        pCI_t = 0x63
				        cellId_t
				         pLMN_Identity_t = 02 f8 29
				         eUTRANcellIdentifier_t = 00 07 ab 50 (28 bits)
				        tAC_t = 01 02
				        broadcastPLMNs_t:
				         BroadcastPLMNs_Item_elm = 02 f8 29
				        eUTRA_Mode_Info_t:
				         fDD_t
				          uL_EARFCN_t = 0x1
				          dL_EARFCN_t = 0x1
				          uL_Transmission_Bandwidth_t = 0x3
				          dL_Transmission_Bandwidth_t = 0x3
				       neighbour_Info_t:
				        Neighbour_Information_elm
				         eCGI_t
				          pLMN_Identity_t = 02 f8 29
				          eUTRANcellIdentifier_t = 00 07 ab 50 (28 bits)
				         pCI_t = 0x63
				         eARFCN_t = 0x1
				         iE_Extensions_t:
				          ProtocolExtensionContainer_elm
				           id_t = 0x4c
				           criticality_t = 0x1
				           TAC = 01 02
				        Neighbour_Information_elm
				         eCGI_t
				          pLMN_Identity_t = 03 f9 30
				          eUTRANcellIdentifier_t = 00 08 bc 50 (28 bits)
				         pCI_t = 0x64
				         eARFCN_t = 0x2
				         iE_Extensions_t:
				          ProtocolExtensionContainer_elm
				           id_t = 0x5e
				           criticality_t = 0
				           EARFCNExtension = 0x3
			*/
			packedPdu: "20060055000002001500080002f82900007a8000140042004000630002f8290007ab50102002f82900000100013300024002f8290007ab50006300010000004c400201024003f9300008bc50006400020000005e0003800103"},

		{key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "007a80"},
			enb: "CONNECTED MACRO_ENB [pci:99 cell_id:\"02f829:0007ab50\" tac:\"0102\" broadcast_plmns:\"02f829\" choice_eutra_mode:<fdd:<ulear_fcn:1 dlear_fcn:1 ul_transmission_bandwidth:BW50 dl_transmission_bandwidth:BW50 > > eutra_mode:FDD neighbour_infos:<ecgi:\"02f829:0007ab50\" pci:99 ear_fcn:1 tac:\"0102\" > neighbour_infos:<ecgi:\"03f930:0008bc50\" pci:100 ear_fcn:3 > ] []",
			/*
				X2AP-PDU:
				 successfulOutcome_t
				  procedureCode_t = 0x6
				  criticality_t = 0
				  X2SetupResponse
				   protocolIEs_t:
				    ProtocolIE_Container_elm
				     id_t = 0x15
				     criticality_t = 0
				     GlobalENB-ID
				      pLMN_Identity_t = 02 f8 29
				      eNB_ID_t:
				       macro_eNB_ID_t = 00 7a 80 (20 bits)
				    ProtocolIE_Container_elm
				     id_t = 0x14
				     criticality_t = 0
				     ServedCells:
				      ServedCells_elm
				       servedCellInfo_t
				        pCI_t = 0x63
				        cellId_t
				         pLMN_Identity_t = 02 f8 29
				         eUTRANcellIdentifier_t = 00 07 ab 50 (28 bits)
				        tAC_t = 01 02
				        broadcastPLMNs_t:
				         BroadcastPLMNs_Item_elm = 02 f8 29
				        eUTRA_Mode_Info_t:
				         fDD_t
				          uL_EARFCN_t = 0x1
				          dL_EARFCN_t = 0x1
				          uL_Transmission_Bandwidth_t = 0x3
				          dL_Transmission_Bandwidth_t = 0x3
				       neighbour_Info_t:
				        Neighbour_Information_elm
				         eCGI_t
				          pLMN_Identity_t = 02 f8 29
				          eUTRANcellIdentifier_t = 00 07 ab 50 (28 bits)
				         pCI_t = 0x63
				         eARFCN_t = 0x1
				         iE_Extensions_t:
				          ProtocolExtensionContainer_elm
				           id_t = 0x4c
				           criticality_t = 0x1
				           TAC = 01 02
				        Neighbour_Information_elm
				         eCGI_t
				          pLMN_Identity_t = 03 f9 30
				          eUTRANcellIdentifier_t = 00 08 bc 50 (28 bits)
				         pCI_t = 0x64
				         eARFCN_t = 0x2
				         iE_Extensions_t:
				          ProtocolExtensionContainer_elm
				           id_t = 0x5e
				           criticality_t = 0
				           EARFCNExtension = 0x3
			*/
			packedPdu: "20060055000002001500080002f82900007a8000140042004000630002f8290007ab50102002f82900000100013300024002f8290007ab50006300010000004c400201024003f9300008bc50006400020000005e0003800103",
			/*failure: fmt.Errorf("getAtom for path [successfulOutcome_t X2SetupResponse protocolIEs_t ProtocolIE_Container_elm GlobalENB-ID pLMN_Identity_t] failed, rc = 2" /NO_SPACE_LEFT),*/ },
	}

	converter := NewX2SetupResponseConverter(logger)

	for _, tc := range testCases {
		t.Run(tc.packedPdu, func(t *testing.T) {

			var payload []byte
			_, err := fmt.Sscanf(tc.packedPdu, "%x", &payload)
			if err != nil {
				t.Errorf("convert inputPayloadAsStr to payloadAsByte. Error: %v\n", err)
			}

			key, enb, err := converter.UnpackX2SetupResponseAndExtract(payload)

			if err != nil {
				if tc.failure == nil {
					t.Errorf("want: success, got: error: %v\n", err)
				} else {
					if strings.Compare(err.Error(), tc.failure.Error()) != 0 {
						t.Errorf("want: %s, got: %s", tc.failure, err)
					}
				}
			}

			if key == nil {
				if tc.failure == nil {
					t.Errorf("want: key=%v, got: empty key", tc.key)
				}
			} else {
				if strings.Compare(key.PlmnId, tc.key.PlmnId) != 0 || strings.Compare(key.NbId, tc.key.NbId) != 0 {
					t.Errorf("want: key=%s, got: %s", tc.key, key)
				}
			}

			if enb == nil {
				if tc.failure == nil {
					t.Errorf("want: enb=%s, got: empty enb", tc.enb)
				}
			} else {
				nb := &entities.NodebInfo{}
				nb.ConnectionStatus = entities.ConnectionStatus_CONNECTED
				nb.Configuration = &entities.NodebInfo_Enb{Enb: enb}
				embStr := fmt.Sprintf("%s %s %s %s", nb.ConnectionStatus, enb.EnbType, enb.ServedCells, enb.GuGroupIds)
				if !strings.EqualFold(embStr, tc.enb) {
					t.Errorf("want: enb=%s, got: %s", tc.enb, embStr)
				}
			}
		})
	}
}
