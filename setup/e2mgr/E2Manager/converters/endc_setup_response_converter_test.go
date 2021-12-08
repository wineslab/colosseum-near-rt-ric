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

func TestUnpackEndcX2SetupResponseAndExtract(t *testing.T) {
	logger, _ := logger.InitLogger(logger.InfoLevel)

	var testCases = []struct {
		key       *entities.GlobalNbId
		gnb       string
		packedPdu string
		failure   error
	}{
		{
			key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "4a952a0a"},
			gnb: "CONNECTED [served_nr_cell_information:<nr_pci:5 cell_id:\"1e3f27:1f2e3d4ff0\" stac5g:\"3d44d3\" configured_stac:\"4e4f\" served_plmns:\"3e4e5e\" nr_mode:TDD choice_nr_mode:<tdd:<nr_freq_info:<nr_ar_fcn:1 sulInformation:<sul_ar_fcn:2 sul_transmission_bandwidth:<nrscs:SCS60 ncnrb:NRB107 > > frequency_bands:<nr_frequency_band:22 supported_sul_bands:11 > > transmission_bandwidth:<nrscs:SCS30 ncnrb:NRB133 > > > >  served_nr_cell_information:<nr_pci:5 cell_id:\"1e3f27:1f2e3d4ff0\" stac5g:\"3d44d3\" configured_stac:\"4e4f\" served_plmns:\"3e4e5e\" nr_mode:TDD choice_nr_mode:<tdd:<nr_freq_info:<nr_ar_fcn:1 sulInformation:<sul_ar_fcn:2 sul_transmission_bandwidth:<nrscs:SCS120 ncnrb:NRB121 > > frequency_bands:<nr_frequency_band:22 supported_sul_bands:11 > > transmission_bandwidth:<nrscs:SCS15 ncnrb:NRB132 > > > > nr_neighbour_infos:<nr_pci:44 nr_cgi:\"1e3f27:1f2e3d4ff0\" nr_mode:TDD choice_nr_mode:<tdd:<ar_fcn_nr_freq_info:<nr_ar_fcn:1 sulInformation:<sul_ar_fcn:2 sul_transmission_bandwidth:<nrscs:SCS15 ncnrb:NRB11 > > frequency_bands:<nr_frequency_band:22 supported_sul_bands:11 > > > > > ]",
			/*
		E2AP-PDU:
			 successfulOutcome_t
			  procedureCode_t = 0x24
			  criticality_t = 0
			  ENDCX2SetupResponse
			   protocolIEs_t:
			    ProtocolIE_Container_elm
			     id_t = 0xf6
			     criticality_t = 0
			     RespondingNodeType-EndcX2Setup:
			      respond_en_gNB_t:
			       ProtocolIE_Container_elm
			        id_t = 0xfc
			        criticality_t = 0
			        GlobalGNB-ID
			         pLMN_Identity_t = 02 f8 29
			         gNB_ID_t:
			          gNB_ID_t = 4a 95 2a 0a (32 bits)
			       ProtocolIE_Container_elm
			        id_t = 0xfd
			        criticality_t = 0
			        ServedNRcellsENDCX2ManagementList:
			         ServedNRcellsENDCX2ManagementList_elm
			          servedNRCellInfo_t
			           nrpCI_t = 0x5
			           nrCellID_t
			            pLMN_Identity_t = 1e 3f 27
			            nRcellIdentifier_t = 1f 2e 3d 4f f0 (36 bits)
			           fiveGS_TAC_t = 3d 44 d3
			           configured_TAC_t = 4e 4f
			           broadcastPLMNs_t:
			            BroadcastPLMNs_Item_elm = 3e 4e 5e
			           nrModeInfo_t:
			            tdd_t
			             nRFreqInfo_t
			              nRARFCN_t = 0x1
			              freqBandListNr_t:
			               freqBandListNr_t_elm
			                freqBandIndicatorNr_t = 0x16
			                supportedSULBandList_t:
			                 supportedSULBandList_t_elm
			                  freqBandIndicatorNr_t = 0xb
			              sULInformation_t
			               sUL_ARFCN_t = 0x2
			               sUL_TxBW_t
			                nRSCS_t = 0x2
			                nRNRB_t = 0xf
			             nR_TxBW_t
			              nRSCS_t = 0x1
			              nRNRB_t = 0x12
			           measurementTimingConfiguration_t = 3e 4e 5e
			         ServedNRcellsENDCX2ManagementList_elm
			          servedNRCellInfo_t
			           nrpCI_t = 0x5
			           nrCellID_t
			            pLMN_Identity_t = 1e 3f 27
			            nRcellIdentifier_t = 1f 2e 3d 4f f0 (36 bits)
			           fiveGS_TAC_t = 3d 44 d3
			           configured_TAC_t = 4e 4f
			           broadcastPLMNs_t:
			            BroadcastPLMNs_Item_elm = 3e 4e 5e
			           nrModeInfo_t:
			            tdd_t
			             nRFreqInfo_t
			              nRARFCN_t = 0x1
			              freqBandListNr_t:
			               freqBandListNr_t_elm
			                freqBandIndicatorNr_t = 0x16
			                supportedSULBandList_t:
			                 supportedSULBandList_t_elm
			                  freqBandIndicatorNr_t = 0xb
			              sULInformation_t
			               sUL_ARFCN_t = 0x2
			               sUL_TxBW_t
			                nRSCS_t = 0x3
			                nRNRB_t = 0x10
			             nR_TxBW_t
			              nRSCS_t = 0
			              nRNRB_t = 0x11
			           measurementTimingConfiguration_t = 3e 4e 5e
			          nRNeighbourInfo_t:
			           NRNeighbour_Information_elm
			            nrpCI_t = 0x2c
			            nrCellID_t
			             pLMN_Identity_t = 1e 3f 27
			             nRcellIdentifier_t = 1f 2e 3d 4f f0 (36 bits)
			            measurementTimingConfiguration_t = 1e 3f 27
			            nRNeighbourModeInfo_t:
			             tdd_t
			              nRFreqInfo_t
			               nRARFCN_t = 0x1
			               freqBandListNr_t:
			                freqBandListNr_t_elm
			                 freqBandIndicatorNr_t = 0x16
			                 supportedSULBandList_t:
			                  supportedSULBandList_t_elm
			                   freqBandIndicatorNr_t = 0xb
			               sULInformation_t
			                sUL_ARFCN_t = 0x2
			                sUL_TxBW_t
			                 nRSCS_t = 0
			                 nRNRB_t = 0

			*/
			packedPdu: "202400808e00000100f600808640000200fc00090002f829504a952a0a00fd007200010c0005001e3f271f2e3d4ff03d44d34e4f003e4e5e4400010000150400000a000211e148033e4e5e4c0005001e3f271f2e3d4ff03d44d34e4f003e4e5e4400010000150400000a00021a0044033e4e5e000000002c001e3f271f2e3d4ff0031e3f274400010000150400000a00020000"},
		{
			key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "4a952a0a"},
			gnb: "CONNECTED [served_nr_cell_information:<nr_pci:5 cell_id:\"1e3f27:1f2e3d4ff0\" stac5g:\"3d44d3\" configured_stac:\"4e4f\" served_plmns:\"3e4e5e\" nr_mode:TDD choice_nr_mode:<tdd:<nr_freq_info:<nr_ar_fcn:1 sulInformation:<sul_ar_fcn:2 sul_transmission_bandwidth:<nrscs:SCS30 ncnrb:NRB107 > > frequency_bands:<nr_frequency_band:22 supported_sul_bands:11 > > transmission_bandwidth:<nrscs:SCS15 ncnrb:NRB121 > > > > nr_neighbour_infos:<nr_pci:44 nr_cgi:\"1e3f27:1f2e3d4ff0\" nr_mode:TDD choice_nr_mode:<tdd:<ar_fcn_nr_freq_info:<nr_ar_fcn:5 sulInformation:<sul_ar_fcn:6 sul_transmission_bandwidth:<nrscs:SCS120 ncnrb:NRB18 > > frequency_bands:<nr_frequency_band:22 supported_sul_bands:11 > > > > > ]",
			/*
			E2AP-PDU:
			 successfulOutcome_t
			  procedureCode_t = 0x24
			  criticality_t = 0
			  ENDCX2SetupResponse
			   protocolIEs_t:
			    ProtocolIE_Container_elm
			     id_t = 0xf6
			     criticality_t = 0
			     RespondingNodeType-EndcX2Setup:
			      respond_en_gNB_t:
			       ProtocolIE_Container_elm
			        id_t = 0xfc
			        criticality_t = 0
			        GlobalGNB-ID
			         pLMN_Identity_t = 02 f8 29
			         gNB_ID_t:
			          gNB_ID_t = 4a 95 2a 0a (32 bits)
			       ProtocolIE_Container_elm
			        id_t = 0xfd
			        criticality_t = 0
			        ServedNRcellsENDCX2ManagementList:
			         ServedNRcellsENDCX2ManagementList_elm
			          servedNRCellInfo_t
			           nrpCI_t = 0x5
			           nrCellID_t
			            pLMN_Identity_t = 1e 3f 27
			            nRcellIdentifier_t = 1f 2e 3d 4f f0 (36 bits)
			           fiveGS_TAC_t = 3d 44 d3
			           configured_TAC_t = 4e 4f
			           broadcastPLMNs_t:
			            BroadcastPLMNs_Item_elm = 3e 4e 5e
			           nrModeInfo_t:
			            tdd_t
			             nRFreqInfo_t
			              nRARFCN_t = 0x1
			              freqBandListNr_t:
			               freqBandListNr_t_elm
			                freqBandIndicatorNr_t = 0x16
			                supportedSULBandList_t:
			                 supportedSULBandList_t_elm
			                  freqBandIndicatorNr_t = 0xb
			              sULInformation_t
			               sUL_ARFCN_t = 0x2
			               sUL_TxBW_t
			                nRSCS_t = 0x1
			                nRNRB_t = 0xf
			             nR_TxBW_t
			              nRSCS_t = 0
			              nRNRB_t = 0x10
			           measurementTimingConfiguration_t = 3e 4e 5e
			          nRNeighbourInfo_t:
			           NRNeighbour_Information_elm
			            nrpCI_t = 0x2c
			            nrCellID_t
			             pLMN_Identity_t = 1e 3f 27
			             nRcellIdentifier_t = 1f 2e 3d 4f f0 (36 bits)
			            measurementTimingConfiguration_t = 1e 3f 27
			            nRNeighbourModeInfo_t:
			             tdd_t
			              nRFreqInfo_t
			               nRARFCN_t = 0x5
			               freqBandListNr_t:
			                freqBandListNr_t_elm
			                 freqBandIndicatorNr_t = 0x16
			                 supportedSULBandList_t:
			                  supportedSULBandList_t_elm
			                   freqBandIndicatorNr_t = 0xb
			               sULInformation_t
			                sUL_ARFCN_t = 0x6
			                sUL_TxBW_t
			                 nRSCS_t = 0x3
			                 nRNRB_t = 0x1

			*/
			packedPdu: "2024006500000100f6005e40000200fc00090002f829504a952a0a00fd004a00004c0005001e3f271f2e3d4ff03d44d34e4f003e4e5e4400010000150400000a000209e040033e4e5e000000002c001e3f271f2e3d4ff0031e3f274400050000150400000a00061820"},

		{
			key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "4a952a0a"},
			gnb: "CONNECTED [served_nr_cell_information:<nr_pci:5 cell_id:\"1e3f27:1f2e3d4ff0\" stac5g:\"3d44d3\" configured_stac:\"4e4f\" served_plmns:\"3e4e5e\" nr_mode:TDD choice_nr_mode:<tdd:<nr_freq_info:<nr_ar_fcn:1 sulInformation:<sul_ar_fcn:2 sul_transmission_bandwidth:<nrscs:SCS60 ncnrb:NRB107 > > frequency_bands:<nr_frequency_band:22 supported_sul_bands:11 > > transmission_bandwidth:<nrscs:SCS30 ncnrb:NRB133 > > > >  served_nr_cell_information:<nr_pci:8 cell_id:\"2e3f45:1f2e3d4ff0\" stac5g:\"4faa3c\" configured_stac:\"1a2f\" served_plmns:\"50321e\" nr_mode:TDD choice_nr_mode:<tdd:<nr_freq_info:<nr_ar_fcn:4 sulInformation:<sul_ar_fcn:8 sul_transmission_bandwidth:<nrscs:SCS120 ncnrb:NRB121 > > frequency_bands:<nr_frequency_band:7 supported_sul_bands:3 > > transmission_bandwidth:<nrscs:SCS15 ncnrb:NRB132 > > > > nr_neighbour_infos:<nr_pci:44 nr_cgi:\"1e3f27:1f2e3d4ff0\" nr_mode:TDD choice_nr_mode:<tdd:<ar_fcn_nr_freq_info:<nr_ar_fcn:1 sulInformation:<sul_ar_fcn:2 sul_transmission_bandwidth:<nrscs:SCS15 ncnrb:NRB11 > > frequency_bands:<nr_frequency_band:22 supported_sul_bands:11 > > > > > ]",
			/*
			E2AP-PDU:
			 successfulOutcome_t
			  procedureCode_t = 0x24
			  criticality_t = 0
			  ENDCX2SetupResponse
			   protocolIEs_t:
			    ProtocolIE_Container_elm
			     id_t = 0xf6
			     criticality_t = 0
			     RespondingNodeType-EndcX2Setup:
			      respond_en_gNB_t:
			       ProtocolIE_Container_elm
			        id_t = 0xfc
			        criticality_t = 0
			        GlobalGNB-ID
			         pLMN_Identity_t = 02 f8 29
			         gNB_ID_t:
			          gNB_ID_t = 4a 95 2a 0a (32 bits)
			       ProtocolIE_Container_elm
			        id_t = 0xfd
			        criticality_t = 0
			        ServedNRcellsENDCX2ManagementList:
			         ServedNRcellsENDCX2ManagementList_elm
			          servedNRCellInfo_t
			           nrpCI_t = 0x5
			           nrCellID_t
			            pLMN_Identity_t = 1e 3f 27
			            nRcellIdentifier_t = 1f 2e 3d 4f f0 (36 bits)
			           fiveGS_TAC_t = 3d 44 d3
			           configured_TAC_t = 4e 4f
			           broadcastPLMNs_t:
			            BroadcastPLMNs_Item_elm = 3e 4e 5e
			           nrModeInfo_t:
			            tdd_t
			             nRFreqInfo_t
			              nRARFCN_t = 0x1
			              freqBandListNr_t:
			               freqBandListNr_t_elm
			                freqBandIndicatorNr_t = 0x16
			                supportedSULBandList_t:
			                 supportedSULBandList_t_elm
			                  freqBandIndicatorNr_t = 0xb
			              sULInformation_t
			               sUL_ARFCN_t = 0x2
			               sUL_TxBW_t
			                nRSCS_t = 0x2
			                nRNRB_t = 0xf
			             nR_TxBW_t
			              nRSCS_t = 0x1
			              nRNRB_t = 0x12
			           measurementTimingConfiguration_t = 3e 4e 5e
			         ServedNRcellsENDCX2ManagementList_elm
			          servedNRCellInfo_t
			           nrpCI_t = 0x8
			           nrCellID_t
			            pLMN_Identity_t = 2e 3f 45
			            nRcellIdentifier_t = 1f 2e 3d 4f f0 (36 bits)
			           fiveGS_TAC_t = 4f aa 3c
			           configured_TAC_t = 1a 2f
			           broadcastPLMNs_t:
			            BroadcastPLMNs_Item_elm = 50 32 1e
			           nrModeInfo_t:
			            tdd_t
			             nRFreqInfo_t
			              nRARFCN_t = 0x4
			              freqBandListNr_t:
			               freqBandListNr_t_elm
			                freqBandIndicatorNr_t = 0x7
			                supportedSULBandList_t:
			                 supportedSULBandList_t_elm
			                  freqBandIndicatorNr_t = 0x3
			              sULInformation_t
			               sUL_ARFCN_t = 0x8
			               sUL_TxBW_t
			                nRSCS_t = 0x3
			                nRNRB_t = 0x10
			             nR_TxBW_t
			              nRSCS_t = 0
			              nRNRB_t = 0x11
			           measurementTimingConfiguration_t = 50 32 1e
			          nRNeighbourInfo_t:
			           NRNeighbour_Information_elm
			            nrpCI_t = 0x2c
			            nrCellID_t
			             pLMN_Identity_t = 1e 3f 27
			             nRcellIdentifier_t = 1f 2e 3d 4f f0 (36 bits)
			            measurementTimingConfiguration_t = 1e 3f 27
			            nRNeighbourModeInfo_t:
			             tdd_t
			              nRFreqInfo_t
			               nRARFCN_t = 0x1
			               freqBandListNr_t:
			                freqBandListNr_t_elm
			                 freqBandIndicatorNr_t = 0x16
			                 supportedSULBandList_t:
			                  supportedSULBandList_t_elm
			                   freqBandIndicatorNr_t = 0xb
			               sULInformation_t
			                sUL_ARFCN_t = 0x2
			                sUL_TxBW_t
			                 nRSCS_t = 0
			                 nRNRB_t = 0
			*/
			packedPdu: "202400808e00000100f600808640000200fc00090002f829504a952a0a00fd007200010c0005001e3f271f2e3d4ff03d44d34e4f003e4e5e4400010000150400000a000211e148033e4e5e4c0008002e3f451f2e3d4ff04faa3c1a2f0050321e4400040000060400000200081a00440350321e000000002c001e3f271f2e3d4ff0031e3f274400010000150400000a00020000"},

		{
			key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "4a952a0a"},
			gnb: "CONNECTED [served_nr_cell_information:<nr_pci:5 cell_id:\"1e3f27:1f2e3d4ff0\" served_plmns:\"3e4e5e\" nr_mode:FDD choice_nr_mode:<fdd:<ul_freq_info:<nr_ar_fcn:5 frequency_bands:<nr_frequency_band:44 supported_sul_bands:33 > > dl_freq_info:<nr_ar_fcn:1 frequency_bands:<nr_frequency_band:22 supported_sul_bands:11 > > ul_transmission_bandwidth:<nrscs:SCS120 ncnrb:NRB11 > dl_transmission_bandwidth:<nrscs:SCS15 ncnrb:NRB135 > > > > nr_neighbour_infos:<nr_pci:44 nr_cgi:\"1e3f27:1f2e3d4ff0\" nr_mode:FDD choice_nr_mode:<fdd:<ular_fcn_freq_info:<nr_ar_fcn:5 frequency_bands:<nr_frequency_band:22 supported_sul_bands:11 > > dlar_fcn_freq_info:<nr_ar_fcn:1 frequency_bands:<nr_frequency_band:22 supported_sul_bands:11 > > > > > ]",
			/*
			E2AP-PDU:
			 successfulOutcome_t
			  procedureCode_t = 0x24
			  criticality_t = 0
			  ENDCX2SetupResponse
			   protocolIEs_t:
			    ProtocolIE_Container_elm
			     id_t = 0xf6
			     criticality_t = 0
			     RespondingNodeType-EndcX2Setup:
			      respond_en_gNB_t:
			       ProtocolIE_Container_elm
			        id_t = 0xfc
			        criticality_t = 0
			        GlobalGNB-ID
			         pLMN_Identity_t = 02 f8 29
			         gNB_ID_t:
			          gNB_ID_t = 4a 95 2a 0a (32 bits)
			       ProtocolIE_Container_elm
			        id_t = 0xfd
			        criticality_t = 0
			        ServedNRcellsENDCX2ManagementList:
			         ServedNRcellsENDCX2ManagementList_elm
			          servedNRCellInfo_t
			           nrpCI_t = 0x5
			           nrCellID_t
			            pLMN_Identity_t = 1e 3f 27
			            nRcellIdentifier_t = 1f 2e 3d 4f f0 (36 bits)
			           broadcastPLMNs_t:
			            BroadcastPLMNs_Item_elm = 3e 4e 5e
			           nrModeInfo_t:
			            fdd_t
			             ul_NRFreqInfo_t
			              nRARFCN_t = 0x5
			              freqBandListNr_t:
			               freqBandListNr_t_elm
			                freqBandIndicatorNr_t = 0x2c
			                supportedSULBandList_t:
			                 supportedSULBandList_t_elm
			                  freqBandIndicatorNr_t = 0x21
			             dl_NRFreqInfo_t
			              nRARFCN_t = 0x1
			              freqBandListNr_t:
			               freqBandListNr_t_elm
			                freqBandIndicatorNr_t = 0x16
			                supportedSULBandList_t:
			                 supportedSULBandList_t_elm
			                  freqBandIndicatorNr_t = 0xb
			             ul_NR_TxBW_t
			              nRSCS_t = 0x3
			              nRNRB_t = 0
			             dl_NR_TxBW_t
			              nRSCS_t = 0
			              nRNRB_t = 0x13
			           measurementTimingConfiguration_t = 01 02 03
			          nRNeighbourInfo_t:
			           NRNeighbour_Information_elm
			            nrpCI_t = 0x2c
			            nrCellID_t
			             pLMN_Identity_t = 1e 3f 27
			             nRcellIdentifier_t = 1f 2e 3d 4f f0 (36 bits)
			            measurementTimingConfiguration_t = 01 02 03
			            nRNeighbourModeInfo_t:
			             fdd_t
			              ul_NRFreqInfo_t
			               nRARFCN_t = 0x5
			               freqBandListNr_t:
			                freqBandListNr_t_elm
			                 freqBandIndicatorNr_t = 0x16
			                 supportedSULBandList_t:
			                  supportedSULBandList_t_elm
			                   freqBandIndicatorNr_t = 0xb
			              dl_NRFreqInfo_t
			               nRARFCN_t = 0x1
			               freqBandListNr_t:
			                freqBandListNr_t_elm
			                 freqBandIndicatorNr_t = 0x16
			                 supportedSULBandList_t:
			                  supportedSULBandList_t_elm
			                   freqBandIndicatorNr_t = 0xb


			*/
			packedPdu: "2024006b00000100f6006440000200fc00090002f829504a952a0a00fd00500000400005001e3f271f2e3d4ff03e4e5e00000500002b0400002000010000150400000a18004c03010203000000002c001e3f271f2e3d4ff0030102030000050000150400000a00010000150400000a"},


		{
			key: &entities.GlobalNbId{PlmnId: "04a5c1", NbId: "4fc52bff"},
			gnb: "CONNECTED [served_nr_cell_information:<nr_pci:9 cell_id:\"aeafa7:2a3e3b4cd0\" stac5g:\"7d4773\" configured_stac:\"477f\" served_plmns:\"7e7e7e\" nr_mode:TDD choice_nr_mode:<tdd:<nr_freq_info:<nr_ar_fcn:8 sulInformation:<sul_ar_fcn:9 sul_transmission_bandwidth:<nrscs:SCS15 ncnrb:NRB121 > > frequency_bands:<nr_frequency_band:22 supported_sul_bands:11 > > transmission_bandwidth:<nrscs:SCS60 ncnrb:NRB18 > > > > nr_neighbour_infos:<nr_pci:44 nr_cgi:\"5a5ff1:2a3e3b4cd0\" nr_mode:TDD choice_nr_mode:<tdd:<ar_fcn_nr_freq_info:<nr_ar_fcn:5 sulInformation:<sul_ar_fcn:6 sul_transmission_bandwidth:<nrscs:SCS30 ncnrb:NRB18 > > frequency_bands:<nr_frequency_band:4 supported_sul_bands:3 > > > > > nr_neighbour_infos:<nr_pci:9 nr_cgi:\"5d5caa:af3e354ac0\" nr_mode:TDD choice_nr_mode:<tdd:<ar_fcn_nr_freq_info:<nr_ar_fcn:7 sulInformation:<sul_ar_fcn:8 sul_transmission_bandwidth:<nrscs:SCS120 ncnrb:NRB25 > > frequency_bands:<nr_frequency_band:3 supported_sul_bands:1 > > > > > ]",
			/*
			E2AP-PDU:
			 successfulOutcome_t
			  procedureCode_t = 0x24
			  criticality_t = 0
			  ENDCX2SetupResponse
			   protocolIEs_t:
			    ProtocolIE_Container_elm
			     id_t = 0xf6
			     criticality_t = 0
			     RespondingNodeType-EndcX2Setup:
			      respond_en_gNB_t:
			       ProtocolIE_Container_elm
			        id_t = 0xfc
			        criticality_t = 0
			        GlobalGNB-ID
			         pLMN_Identity_t = 04 a5 c1
			         gNB_ID_t:
			          gNB_ID_t = 4f c5 2b ff (32 bits)
			       ProtocolIE_Container_elm
			        id_t = 0xfd
			        criticality_t = 0
			        ServedNRcellsENDCX2ManagementList:
			         ServedNRcellsENDCX2ManagementList_elm
			          servedNRCellInfo_t
			           nrpCI_t = 0x9
			           nrCellID_t
			            pLMN_Identity_t = ae af a7
			            nRcellIdentifier_t = 2a 3e 3b 4c d0 (36 bits)
			           fiveGS_TAC_t = 7d 47 73
			           configured_TAC_t = 47 7f
			           broadcastPLMNs_t:
			            BroadcastPLMNs_Item_elm = 7e 7e 7e
			           nrModeInfo_t:
			            tdd_t
			             nRFreqInfo_t
			              nRARFCN_t = 0x8
			              freqBandListNr_t:
			               freqBandListNr_t_elm
			                freqBandIndicatorNr_t = 0x16
			                supportedSULBandList_t:
			                 supportedSULBandList_t_elm
			                  freqBandIndicatorNr_t = 0xb
			              sULInformation_t
			               sUL_ARFCN_t = 0x9
			               sUL_TxBW_t
			                nRSCS_t = 0
			                nRNRB_t = 0x10
			             nR_TxBW_t
			              nRSCS_t = 0x2
			              nRNRB_t = 0x1
			           measurementTimingConfiguration_t = 7e 7e 7e
			          nRNeighbourInfo_t:
			           NRNeighbour_Information_elm
			            nrpCI_t = 0x2c
			            nrCellID_t
			             pLMN_Identity_t = 5a 5f f1
			             nRcellIdentifier_t = 2a 3e 3b 4c d0 (36 bits)
			            measurementTimingConfiguration_t = 5a 5f f1
			            nRNeighbourModeInfo_t:
			             tdd_t
			              nRFreqInfo_t
			               nRARFCN_t = 0x5
			               freqBandListNr_t:
			                freqBandListNr_t_elm
			                 freqBandIndicatorNr_t = 0x4
			                 supportedSULBandList_t:
			                  supportedSULBandList_t_elm
			                   freqBandIndicatorNr_t = 0x3
			               sULInformation_t
			                sUL_ARFCN_t = 0x6
			                sUL_TxBW_t
			                 nRSCS_t = 0x1
			                 nRNRB_t = 0x1
			           NRNeighbour_Information_elm
			            nrpCI_t = 0x9
			            nrCellID_t
			             pLMN_Identity_t = 5d 5c aa
			             nRcellIdentifier_t = af 3e 35 4a c0 (36 bits)
			            measurementTimingConfiguration_t = 5d 5c aa
			            nRNeighbourModeInfo_t:
			             tdd_t
			              nRFreqInfo_t
			               nRARFCN_t = 0x7
			               freqBandListNr_t:
			                freqBandListNr_t_elm
			                 freqBandIndicatorNr_t = 0x3
			                 supportedSULBandList_t:
			                  supportedSULBandList_t_elm
			                   freqBandIndicatorNr_t = 0x1
			               sULInformation_t
			                sUL_ARFCN_t = 0x8
			                sUL_TxBW_t
			                 nRSCS_t = 0x3
			                 nRNRB_t = 0x3
			*/
			packedPdu: "202400808200000100f6007b40000200fc00090004a5c1504fc52bff00fd006700004c000900aeafa72a3e3b4cd07d4773477f007e7e7e4400080000150400000a0009020204037e7e7e000100002c005a5ff12a3e3b4cd0035a5ff144000500000304000002000608200009005d5caaaf3e354ac0035d5caa4400070000020400000000081860"},

		{key: &entities.GlobalNbId{PlmnId: "02f829", NbId: "4a952aaa"},
			/*
			E2AP-PDU:
			 successfulOutcome_t
			  procedureCode_t = 0x24
			  criticality_t = 0
			  ENDCX2SetupResponse
			   protocolIEs_t:
			    ProtocolIE_Container_elm
			     id_t = 0xf6
			     criticality_t = 0
			     RespondingNodeType-EndcX2Setup:
			      respond_en_gNB_t:
			       ProtocolIE_Container_elm
			        id_t = 0xfc
			        criticality_t = 0
			        GlobalGNB-ID
			         pLMN_Identity_t = 02 f8 29
			         gNB_ID_t:
			          gNB_ID_t = 4a 95 2a aa (32 bits)
			*/
			packedPdu: "2024001700000100f6001040000100fc00090002f829504a952aaa",

			failure: fmt.Errorf("getList for path [successfulOutcome_t ENDCX2SetupResponse protocolIEs_t ProtocolIE_Container_elm RespondingNodeType-EndcX2Setup respond_en_gNB_t ProtocolIE_Container_elm ServedNRcellsENDCX2ManagementList ServedNRcellsENDCX2ManagementList_elm servedNRCellInfo_t nrpCI_t] failed, rc = 1" /*NO_ITEMS*/),},
	}

	converter := NewEndcSetupResponseConverter(logger)

	for _, tc := range testCases {
		t.Run(tc.packedPdu, func(t *testing.T) {

			var payload []byte

			_, err := fmt.Sscanf(tc.packedPdu, "%x", &payload)

			if err != nil {
				t.Errorf("convert inputPayloadAsStr to payloadAsByte. Error: %v\n", err)
			}

			key, gnb, err := converter.UnpackEndcSetupResponseAndExtract(payload)

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
					t.Errorf("want: key=%v, got: %v", tc.key, key)
				}
			}

			if gnb == nil {
				if tc.failure == nil {
					t.Errorf("want: enb=%s, got: empty enb", tc.gnb)
				}
			} else {
				nb := &entities.NodebInfo{}
				nb.ConnectionStatus = entities.ConnectionStatus_CONNECTED
				nb.Configuration = &entities.NodebInfo_Gnb{Gnb: gnb}
				gnbStr := fmt.Sprintf("%s %s", nb.ConnectionStatus, gnb.ServedNrCells)
				if !strings.EqualFold(gnbStr, tc.gnb) {
					t.Errorf("want: enb=%s, got: %s", tc.gnb, gnbStr)
				}

			}
		})
	}
}
