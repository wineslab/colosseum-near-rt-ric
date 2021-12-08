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

package models

import (
	"encoding/xml"
)

type TimeToWait = int

var TimeToWaitEnum = struct {
	V60s TimeToWait
	V20s TimeToWait
	V10s TimeToWait
	V5s  TimeToWait
	V2s  TimeToWait
	V1s  TimeToWait
}{60, 20, 10, 5, 2, 1}

var timeToWaitMap = map[TimeToWait]interface{}{
	TimeToWaitEnum.V60s: struct {
		XMLName xml.Name `xml:"TimeToWait"`
		Text    string   `xml:",chardata"`
		V60s    string   `xml:"v60s"`
	}{},
	TimeToWaitEnum.V20s: struct {
		XMLName xml.Name `xml:"TimeToWait"`
		Text    string   `xml:",chardata"`
		V20s    string   `xml:"v20s"`
	}{},
	TimeToWaitEnum.V10s: struct {
		XMLName xml.Name `xml:"TimeToWait"`
		Text    string   `xml:",chardata"`
		V10s    string   `xml:"v10s"`
	}{},
	TimeToWaitEnum.V5s: struct {
		XMLName xml.Name `xml:"TimeToWait"`
		Text    string   `xml:",chardata"`
		V5s     string   `xml:"v5s"`
	}{},
	TimeToWaitEnum.V2s: struct {
		XMLName xml.Name `xml:"TimeToWait"`
		Text    string   `xml:",chardata"`
		V2s     string   `xml:"v2s"`
	}{},
	TimeToWaitEnum.V1s: struct {
		XMLName xml.Name `xml:"TimeToWait"`
		Text    string   `xml:",chardata"`
		V1s     string   `xml:"v1s"`
	}{},
}

func NewE2SetupSuccessResponseMessage(plmnId string, ricId string, request *E2SetupRequestMessage) E2SetupResponseMessage {
	outcome := SuccessfulOutcome{}
	outcome.ProcedureCode = "1"

	setupRequestIes := request.E2APPDU.InitiatingMessage.Value.E2setupRequest.ProtocolIEs.E2setupRequestIEs

	outcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs = make([]E2setupResponseIEs, len(setupRequestIes))
	outcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs[0].ID = "4"
	outcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs[0].Value = GlobalRICID{GlobalRICID: struct {
		Text         string `xml:",chardata"`
		PLMNIdentity string `xml:"pLMN-Identity"`
		RicID        string `xml:"ric-ID"`
	}{PLMNIdentity: plmnId, RicID: ricId}}

	if len(setupRequestIes) < 2 {
		return E2SetupResponseMessage{E2APPDU: E2APPDU{Outcome: outcome}}
	}

	functionsIdList := extractRanFunctionsIDList(request)

	outcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs[1].ID = "9"
	outcome.Value.E2setupResponse.ProtocolIEs.E2setupResponseIEs[1].Value = RANfunctionsIDList{RANfunctionsIDList: struct {
		Text                      string                      `xml:",chardata"`
		ProtocolIESingleContainer []ProtocolIESingleContainer `xml:"ProtocolIE-SingleContainer"`
	}{ProtocolIESingleContainer: functionsIdList}}

	return E2SetupResponseMessage{E2APPDU: E2APPDU{Outcome: outcome}}
}

func NewE2SetupFailureResponseMessage(timeToWait TimeToWait) E2SetupResponseMessage {
	outcome := UnsuccessfulOutcome{}
	outcome.Value.E2setupFailure.ProtocolIEs.E2setupFailureIEs = make([]E2setupFailureIEs, 2)
	outcome.ProcedureCode = "1"
	outcome.Value.E2setupFailure.ProtocolIEs.E2setupFailureIEs[0].ID = "1"
	outcome.Value.E2setupFailure.ProtocolIEs.E2setupFailureIEs[0].Value.Value = Cause{}
	outcome.Value.E2setupFailure.ProtocolIEs.E2setupFailureIEs[1].ID = "31"
	outcome.Value.E2setupFailure.ProtocolIEs.E2setupFailureIEs[1].Value.Value = timeToWaitMap[timeToWait]
	return E2SetupResponseMessage{E2APPDU: E2APPDU{Outcome: outcome}}
}

type E2SetupResponseMessage struct {
	XMLName xml.Name `xml:"E2SetupSuccessResponseMessage"`
	Text    string   `xml:",chardata"`
	E2APPDU E2APPDU
}

type E2APPDU struct {
	XMLName xml.Name `xml:"E2AP-PDU"`
	Text    string   `xml:",chardata"`
	Outcome interface{}
}

type SuccessfulOutcome struct {
	XMLName       xml.Name `xml:"successfulOutcome"`
	Text          string   `xml:",chardata"`
	ProcedureCode string   `xml:"procedureCode"`
	Criticality   struct {
		Text   string `xml:",chardata"`
		Reject string `xml:"reject"`
	} `xml:"criticality"`
	Value struct {
		Text            string `xml:",chardata"`
		E2setupResponse struct {
			Text        string `xml:",chardata"`
			ProtocolIEs struct {
				Text               string               `xml:",chardata"`
				E2setupResponseIEs []E2setupResponseIEs `xml:"E2setupResponseIEs"`
			} `xml:"protocolIEs"`
		} `xml:"E2setupResponse"`
	} `xml:"value"`
}

type E2setupResponseIEs struct {
	Text        string `xml:",chardata"`
	ID          string `xml:"id"`
	Criticality struct {
		Text   string `xml:",chardata"`
		Reject string `xml:"reject"`
	} `xml:"criticality"`
	Value interface{} `xml:"value"`
}

type GlobalRICID struct {
	Text        string `xml:",chardata"`
	GlobalRICID struct {
		Text         string `xml:",chardata"`
		PLMNIdentity string `xml:"pLMN-Identity"`
		RicID        string `xml:"ric-ID"`
	} `xml:"GlobalRIC-ID"`
}

type RANfunctionsIDList struct {
	Text               string `xml:",chardata"`
	RANfunctionsIDList struct {
		Text                      string                      `xml:",chardata"`
		ProtocolIESingleContainer []ProtocolIESingleContainer `xml:"ProtocolIE-SingleContainer"`
	} `xml:"RANfunctionsID-List"`
}

type ProtocolIESingleContainer struct {
	Text        string `xml:",chardata"`
	ID          string `xml:"id"`
	Criticality struct {
		Text   string `xml:",chardata"`
		Ignore string `xml:"ignore"`
	} `xml:"criticality"`
	Value struct {
		Text              string `xml:",chardata"`
		RANfunctionIDItem struct {
			Text                string `xml:",chardata"`
			RanFunctionID       string `xml:"ranFunctionID"`
			RanFunctionRevision string `xml:"ranFunctionRevision"`
		} `xml:"RANfunctionID-Item"`
	} `xml:"value"`
}

type UnsuccessfulOutcome struct {
	XMLName       xml.Name `xml:"unsuccessfulOutcome"`
	Text          string   `xml:",chardata"`
	ProcedureCode string   `xml:"procedureCode"`
	Criticality   struct {
		Text   string `xml:",chardata"`
		Reject string `xml:"reject"`
	} `xml:"criticality"`
	Value struct {
		Text           string `xml:",chardata"`
		E2setupFailure struct {
			Text        string `xml:",chardata"`
			ProtocolIEs struct {
				Text              string              `xml:",chardata"`
				E2setupFailureIEs []E2setupFailureIEs `xml:"E2setupFailureIEs"`
			} `xml:"protocolIEs"`
		} `xml:"E2setupFailure"`
	} `xml:"value"`
}

type E2setupFailureIEs struct {
	Text        string `xml:",chardata"`
	ID          string `xml:"id"`
	Criticality struct {
		Text   string `xml:",chardata"`
		Ignore string `xml:"ignore"`
	} `xml:"criticality"`
	Value struct {
		Text  string `xml:",chardata"`
		Value interface{}
	} `xml:"value"`
}

type Cause struct {
	XMLName   xml.Name `xml:"Cause"`
	Text      string   `xml:",chardata"`
	Transport struct {
		Text                         string `xml:",chardata"`
		TransportResourceUnavailable string `xml:"transport-resource-unavailable"`
	} `xml:"transport"`
}

func extractRanFunctionsIDList(request *E2SetupRequestMessage) []ProtocolIESingleContainer {

	list := &request.E2APPDU.InitiatingMessage.Value.E2setupRequest.ProtocolIEs.E2setupRequestIEs[1].Value.RANfunctionsList
	ids := make([]ProtocolIESingleContainer, len(list.ProtocolIESingleContainer))
	for i := 0; i < len(ids); i++ {
		ids[i] = convertToRANfunctionID(list, i)
	}
	return ids
}

func convertToRANfunctionID(list *RANfunctionsList, i int) ProtocolIESingleContainer {
	id := ProtocolIESingleContainer{}
	id.ID = "6"
	id.Value.RANfunctionIDItem.RanFunctionID = list.ProtocolIESingleContainer[i].Value.RANfunctionItem.RanFunctionID
	id.Value.RANfunctionIDItem.RanFunctionRevision = list.ProtocolIESingleContainer[i].Value.RANfunctionItem.RanFunctionRevision
	return id
}
