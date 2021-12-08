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

//import (
//	"gerrit.o-ran-sc.org/r/ric-plt/nodeb-rnib.git/entities"
//)
//
///*
//Test permutations of eNB Load Information to protobuf
//*/
//
//type LoadInformationTestCaseName string
//
//const LoadTimestamp = 1257894000000000000
//
//const (
//	SingleCellWithCellIdOnly     LoadInformationTestCaseName = "SINGLE CELL WITH CELL ID ONLY"   //base
//	SingleCellPartiallyPopulated LoadInformationTestCaseName = "SINGLE CELL PARTIALLY POPULATED" //8
//	TwoCellsFullInfo             LoadInformationTestCaseName = "TWO CELLS FULLY POPULATED"       //13
//)
//
//type LoadInformationTestCase struct {
//	loadInformationTestCaseName LoadInformationTestCaseName
//	packedUperPdu               string
//	packedAperPdu               string
//	expectedLoadInformation     *entities.RanLoadInformation
//}
//
//var testCases = []LoadInformationTestCase{
//	{
//		loadInformationTestCaseName: SingleCellWithCellIdOnly,
//		packedAperPdu:               "000240140000010006400d00000740080002f8290007ab50",
//		packedUperPdu:               "004898000400190d0000074200017c148003d5a80000",
//		expectedLoadInformation:     GenerateSingleCellWithCellIdOnlyRanLoadInformation(),
//	},
//	{
//		loadInformationTestCaseName: SingleCellPartiallyPopulated,
//		packedAperPdu:               "", //TODO: populate and USE
//		packedUperPdu:               "004b380004001961000007571e017c148003d5a8205000017c180003d5a875555003331420008007a85801f07c1f07c41f07c1e07801f2020000c680b0003220664102800d8908020000be0c4001ead4016e007ab50100002f8320067ab5005b8c1ead5070190c000000",
//		expectedLoadInformation:     GenerateSingleCellPartiallyPopulatedLoadInformation(),
//	},
//	{
//		loadInformationTestCaseName: TwoCellsFullInfo,
//		packedAperPdu:               "", //TODO: populate and USE
//		packedUperPdu:               "004c07080004001980da0100075bde017c148003d5a8205000017c180003d5a875555403331420000012883a0003547400cd20002801ea16007c1f07c1f107c1f0781e007c80800031a02c000c88199040a00352083669190000d8908020000be0c4001ead4016e007ab50100002f8320067ab5005b8c1ead5070190c00001d637805f220000f56a081400005f020000f56a1d555400ccc508002801ea16007c1f07c1f107c1f0781e007c80800031a02c000c88199040a00352083669190000d8908020000be044001ead4016e007ab50100002f8120067ab5005b8c1ead5070190c00000",
//		expectedLoadInformation:     GenerateTwoCellsFullyPopulatedRanLoadInformation(),
//	},
//}

//func TestExtractAndBuildRanLoadInformation(t *testing.T) {
//	logger, _ := logger.InitLogger(logger.InfoLevel)
//
//	for _, tc := range testCases {
//		t.Run(string(tc.loadInformationTestCaseName), func(t *testing.T) {
//
//			var payload []byte
//			_, err := fmt.Sscanf(tc.packedUperPdu, "%x", &payload)
//
//			if err != nil {
//				t.Errorf("convert inputPayloadAsStr to payloadAsByte. Error: %v\n", err)
//			}
//
//			pdu, err := unpackX2apPduUPer(logger, MaxAsn1CodecAllocationBufferSize, len(payload), payload, MaxAsn1CodecMessageBufferSize)
//
//			actualRanLoadInformation := &entities.RanLoadInformation{LoadTimestamp: LoadTimestamp}
//
//			err = ExtractAndBuildRanLoadInformation(pdu, actualRanLoadInformation)
//
//			if err != nil {
//				t.Errorf("want: success, got: error: %v\n", err)
//			}
//
//			if !assert.Equal(t, tc.expectedLoadInformation, actualRanLoadInformation) {
//				t.Errorf("want: %v, got: %v", tc.expectedLoadInformation, actualRanLoadInformation)
//			}
//		})
//	}
//}

/*func GenerateSingleCellWithCellIdOnlyRanLoadInformation() *entities.RanLoadInformation {
	return &entities.RanLoadInformation{
		LoadTimestamp: LoadTimestamp,
		CellLoadInfos: []*entities.CellLoadInformation{
			{CellId: "02f829:0007ab50"},
		},
	}
}

func GenerateSingleCellPartiallyPopulatedLoadInformation() *entities.RanLoadInformation {

	ulInterferenceOverloadIndications := []entities.UlInterferenceOverloadIndication{
		entities.UlInterferenceOverloadIndication_HIGH_INTERFERENCE,
		entities.UlInterferenceOverloadIndication_MEDIUM_INTERFERENCE,
		entities.UlInterferenceOverloadIndication_LOW_INTERFERENCE,
	}

	rntp := entities.RelativeNarrowbandTxPower{
		RntpPerPrb:                       "cc",
		RntpThreshold:                    entities.RntpThreshold_NEG_6,
		NumberOfCellSpecificAntennaPorts: entities.NumberOfCellSpecificAntennaPorts_V2_ANT_PRT,
		PB:                               2,
		PdcchInterferenceImpact:          1,
	}

	absInformation := entities.AbsInformation{
		Mode:                             entities.AbsInformationMode_ABS_INFO_FDD,
		AbsPatternInfo:                   "07c1f07c1f",
		NumberOfCellSpecificAntennaPorts: entities.NumberOfCellSpecificAntennaPorts_V1_ANT_PRT,
		MeasurementSubset:                "83e0f83c0f",
	}

	extendedUlInterferenceOverloadInfo := entities.ExtendedUlInterferenceOverloadInfo{
		AssociatedSubframes:                       "c8",
		ExtendedUlInterferenceOverloadIndications: ulInterferenceOverloadIndications,
	}

	compInformationStartTime := entities.StartTime{
		StartSfn:            50,
		StartSubframeNumber: 3,
	}

	return &entities.RanLoadInformation{
		LoadTimestamp: LoadTimestamp,
		CellLoadInfos: []*entities.CellLoadInformation{
			{
				CellId:                             "02f829:0007ab50",
				UlInterferenceOverloadIndications:  ulInterferenceOverloadIndications,
				UlHighInterferenceInfos:            []*entities.UlHighInterferenceInformation{{TargetCellId: "02f830:0007ab50", UlHighInterferenceIndication: "aaaa"}},
				RelativeNarrowbandTxPower:          &rntp,
				AbsInformation:                     &absInformation,
				InvokeIndication:                   entities.InvokeIndication_ABS_INFORMATION,
				IntendedUlDlConfiguration:          entities.SubframeAssignment_SA6,
				ExtendedUlInterferenceOverloadInfo: &extendedUlInterferenceOverloadInfo,
				CompInformation: &entities.CompInformation{
					CompInformationItems: []*entities.CompInformationItem{
						{
							CompHypothesisSets: []*entities.CompHypothesisSet{{CellId: "02f831:0007ab50", CompHypothesis: "e007ab50"}},
							BenefitMetric:      -99,
						},
						{
							CompHypothesisSets: []*entities.CompHypothesisSet{{CellId: "02f832:0067ab50", CompHypothesis: "e307ab50"}},
							BenefitMetric:      30,
						},
					},
					CompInformationStartTime: &compInformationStartTime,
				},
			},
		},
	}

}

func GenerateTwoCellsFullyPopulatedRanLoadInformation() *entities.RanLoadInformation {

	ulInterferenceOverloadIndications := []entities.UlInterferenceOverloadIndication{
		entities.UlInterferenceOverloadIndication_HIGH_INTERFERENCE,
		entities.UlInterferenceOverloadIndication_MEDIUM_INTERFERENCE,
		entities.UlInterferenceOverloadIndication_LOW_INTERFERENCE,
	}

	rntp := entities.RelativeNarrowbandTxPower{
		RntpPerPrb:                       "cc",
		RntpThreshold:                    entities.RntpThreshold_NEG_6,
		NumberOfCellSpecificAntennaPorts: entities.NumberOfCellSpecificAntennaPorts_V2_ANT_PRT,
		PB:                               2,
		PdcchInterferenceImpact:          1,
	}

	enhancedRntp := entities.EnhancedRntp{
		EnhancedRntpBitmap:     "aa38",
		RntpHighPowerThreshold: entities.RntpThreshold_NEG_4,
		EnhancedRntpStartTime:  &entities.StartTime{StartSfn: 51, StartSubframeNumber: 9},
	}

	rntpWithEnhanced := rntp
	rntpWithEnhanced.EnhancedRntp = &enhancedRntp

	absInformation := entities.AbsInformation{
		Mode:                             entities.AbsInformationMode_ABS_INFO_FDD,
		AbsPatternInfo:                   "07c1f07c1f",
		NumberOfCellSpecificAntennaPorts: entities.NumberOfCellSpecificAntennaPorts_V1_ANT_PRT,
		MeasurementSubset:                "83e0f83c0f",
	}

	extendedUlInterferenceOverloadInfo := entities.ExtendedUlInterferenceOverloadInfo{
		AssociatedSubframes:                       "c8",
		ExtendedUlInterferenceOverloadIndications: ulInterferenceOverloadIndications,
	}

	compInformationStartTime := entities.StartTime{
		StartSfn:            50,
		StartSubframeNumber: 3,
	}

	dynamicDlTransmissionInformation := entities.DynamicDlTransmissionInformation{
		State:             entities.NaicsState_NAICS_ACTIVE,
		TransmissionModes: "cd",
		PB:                0,
		PAList:            []entities.PA{entities.PA_DB_NEG_1_DOT_77, entities.PA_DB_NEG_3},
	}

	return &entities.RanLoadInformation{
		LoadTimestamp: LoadTimestamp,
		CellLoadInfos: []*entities.CellLoadInformation{
			{
				CellId:                             "02f829:0007ab50",
				UlInterferenceOverloadIndications:  ulInterferenceOverloadIndications,
				UlHighInterferenceInfos:            []*entities.UlHighInterferenceInformation{{TargetCellId: "02f830:0007ab50", UlHighInterferenceIndication: "aaaa"}},
				RelativeNarrowbandTxPower:          &rntpWithEnhanced,
				AbsInformation:                     &absInformation,
				InvokeIndication:                   entities.InvokeIndication_ABS_INFORMATION,
				IntendedUlDlConfiguration:          entities.SubframeAssignment_SA6,
				ExtendedUlInterferenceOverloadInfo: &extendedUlInterferenceOverloadInfo,
				CompInformation: &entities.CompInformation{
					CompInformationItems: []*entities.CompInformationItem{
						{
							CompHypothesisSets: []*entities.CompHypothesisSet{{CellId: "02f831:0007ab50", CompHypothesis: "e007ab50"}},
							BenefitMetric:      -99,
						},
						{
							CompHypothesisSets: []*entities.CompHypothesisSet{{CellId: "02f832:0067ab50", CompHypothesis: "e307ab50"}},
							BenefitMetric:      30,
						},
					},
					CompInformationStartTime: &compInformationStartTime,
				},
				DynamicDlTransmissionInformation: &dynamicDlTransmissionInformation,
			},
			{
				CellId:                             "02f910:0007ab50",
				UlInterferenceOverloadIndications:  ulInterferenceOverloadIndications,
				UlHighInterferenceInfos:            []*entities.UlHighInterferenceInformation{{TargetCellId: "02f810:0007ab50", UlHighInterferenceIndication: "aaaa"}},
				RelativeNarrowbandTxPower:          &rntp,
				AbsInformation:                     &absInformation,
				InvokeIndication:                   entities.InvokeIndication_ABS_INFORMATION,
				IntendedUlDlConfiguration:          entities.SubframeAssignment_SA6,
				ExtendedUlInterferenceOverloadInfo: &extendedUlInterferenceOverloadInfo,
				CompInformation: &entities.CompInformation{
					CompInformationItems: []*entities.CompInformationItem{
						{
							CompHypothesisSets: []*entities.CompHypothesisSet{{CellId: "02f811:0007ab50", CompHypothesis: "e007ab50"}},
							BenefitMetric:      -99,
						},
						{
							CompHypothesisSets: []*entities.CompHypothesisSet{{CellId: "02f812:0067ab50", CompHypothesis: "e307ab50"}},
							BenefitMetric:      30,
						},
					},
					CompInformationStartTime: &compInformationStartTime,
				},
				DynamicDlTransmissionInformation: &dynamicDlTransmissionInformation,
			},
		},
	}
}*/
