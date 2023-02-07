/*
==================================================================================

        Copyright (c) 2019-2020 AT&T Intellectual Property.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
==================================================================================
*/
/*
 * msgs_proc.cc
 * Created on: 2019
 * Author: Ashwin Shridharan, Shraboni Jana, Leonardo Bonati, Andrea Lacava
 */


#include "msgs_proc.hpp"
#include <stdio.h>

extern "C" {
    #include "E2SM-KPM-IndicationMessage.h"
    #include "E2SM-KPM-IndicationMessage-Format1.h"
    #include "RAN-Container.h"
}

bool XappMsgHandler::encode_subscription_delete_request(unsigned char* buffer, size_t *buf_len){

	subscription_helper sub_helper;
	sub_helper.set_request(0); // requirement of subscription manager ... ?
	sub_helper.set_function_id(0);

	subscription_delete e2ap_sub_req_del;

	  // generate the delete request pdu

	  bool res = e2ap_sub_req_del.encode_e2ap_subscription(&buffer[0], buf_len, sub_helper);
	  if(!res){
	    mdclog_write(MDCLOG_ERR, "%s, %d: Error encoding subscription delete request pdu. Reason = %s", __FILE__, __LINE__, e2ap_sub_req_del.get_error().c_str());
	    return false;
	  }

	return true;

}

bool XappMsgHandler::decode_subscription_response(unsigned char* data_buf, size_t data_size){

	subscription_helper subhelper;
	subscription_response subresponse;
	bool res = true;
	E2AP_PDU_t *e2pdu = 0;

	asn_dec_rval_t rval;

	ASN_STRUCT_RESET(asn_DEF_E2AP_PDU, e2pdu);

	rval = asn_decode(0,ATS_ALIGNED_BASIC_PER, &asn_DEF_E2AP_PDU, (void**)&e2pdu, data_buf, data_size);
	switch(rval.code)
	{
		case RC_OK:
			   //Put in Subscription Response Object.
			   //asn_fprint(stdout, &asn_DEF_E2AP_PDU, e2pdu);
			   break;
		case RC_WMORE:
				mdclog_write(MDCLOG_ERR, "RC_WMORE");
				res = false;
				break;
		case RC_FAIL:
				mdclog_write(MDCLOG_ERR, "RC_FAIL");
				res = false;
				break;
		default:
				break;
	 }
	ASN_STRUCT_FREE(asn_DEF_E2AP_PDU, e2pdu);
	return res;

}

bool  XappMsgHandler::a1_policy_handler(char * message, int *message_len, a1_policy_helper &helper){

  rapidjson::Document doc;
  if (doc.Parse(message).HasParseError()){
    mdclog_write(MDCLOG_ERR, "Error: %s, %d :: Could not decode A1 JSON message %s\n", __FILE__, __LINE__, message);
    return false;
  }

  //Extract Operation
  rapidjson::Pointer temp1("/operation");
    rapidjson::Value * ref1 = temp1.Get(doc);
    if (ref1 == NULL){
      mdclog_write(MDCLOG_ERR, "Error : %s, %d:: Could not extract policy type id from %s\n", __FILE__, __LINE__, message);
      return false;
    }

   helper.operation = ref1->GetString();

  // Extract policy id type
  rapidjson::Pointer temp2("/policy_type_id");
  rapidjson::Value * ref2 = temp2.Get(doc);
  if (ref2 == NULL){
    mdclog_write(MDCLOG_ERR, "Error : %s, %d:: Could not extract policy type id from %s\n", __FILE__, __LINE__, message);
    return false;
  }
   //helper.policy_type_id = ref2->GetString();
    helper.policy_type_id = to_string(ref2->GetInt());

    // Extract policy instance id
    rapidjson::Pointer temp("/policy_instance_id");
    rapidjson::Value * ref = temp.Get(doc);
    if (ref == NULL){
      mdclog_write(MDCLOG_ERR, "Error : %s, %d:: Could not extract policy type id from %s\n", __FILE__, __LINE__, message);
      return false;
    }
    helper.policy_instance_id = ref->GetString();

    if (helper.policy_type_id == "1" && helper.operation == "CREATE"){
    	helper.status = "OK";
    	Document::AllocatorType& alloc = doc.GetAllocator();

    	Value handler_id;
    	handler_id.SetString(helper.handler_id.c_str(), helper.handler_id.length(), alloc);

    	Value status;
    	status.SetString(helper.status.c_str(), helper.status.length(), alloc);


    	doc.AddMember("handler_id", handler_id, alloc);
    	doc.AddMember("status",status, alloc);
    	doc.RemoveMember("operation");
    	StringBuffer buffer;
    	Writer<StringBuffer> writer(buffer);
    	doc.Accept(writer);
    	strncpy(message,buffer.GetString(), buffer.GetLength());
    	*message_len = buffer.GetLength();
    	return true;
    }
    return false;
 }


//For processing received messages.XappMsgHandler should mention if resend is required or not.
void XappMsgHandler::operator()(rmr_mbuf_t *message, bool *resend){

	if (message->len > MAX_RMR_RECV_SIZE){
		mdclog_write(MDCLOG_ERR, "Error : %s, %d, RMR message larger than %d. Ignoring ...", __FILE__, __LINE__, MAX_RMR_RECV_SIZE);
		return;
	}
	a1_policy_helper helper;
	bool res=false;
	switch(message->mtype){
		//need to fix the health check.
		case (RIC_HEALTH_CHECK_REQ):
				message->mtype = RIC_HEALTH_CHECK_RESP;        // if we're here we are running and all is ok
				message->sub_id = -1;
				strncpy( (char*)message->payload, "HELLOWORLD OK\n", rmr_payload_size( message) );
				*resend = true;
				break;

		case (RIC_INDICATION): {

			mdclog_write(MDCLOG_INFO, "Received RIC indication message of type = %d", message->mtype);

			unsigned char *me_id_null;
			unsigned char *me_id = rmr_get_meid(message, me_id_null);
			mdclog_write(MDCLOG_INFO,"RMR Received MEID: %s",me_id);

			process_ric_indication(message->mtype, me_id, message->payload, message->len);

			break;
		}

		case (RIC_SUB_RESP): {
        		mdclog_write(MDCLOG_INFO, "Received subscription message of type = %d", message->mtype);

				unsigned char *me_id_null;
				unsigned char *me_id = rmr_get_meid(message, me_id_null);
				mdclog_write(MDCLOG_INFO,"RMR Received MEID: %s",me_id);

				if(_ref_sub_handler !=NULL){
					_ref_sub_handler->manage_subscription_response(message->mtype, me_id, message->payload, message->len);
				} else {
					mdclog_write(MDCLOG_ERR, " Error :: %s, %d : Subscription handler not assigned in message processor !", __FILE__, __LINE__);
				}
				*resend = false;
				break;
		 }

	case A1_POLICY_REQ:

		    mdclog_write(MDCLOG_INFO, "In Message Handler: Received A1_POLICY_REQ.");
			helper.handler_id = xapp_id;

			res = a1_policy_handler((char*)message->payload, &message->len, helper);
			if(res){
				message->mtype = A1_POLICY_RESP;        // if we're here we are running and all is ok
				message->sub_id = -1;
				*resend = true;
			}
			break;

	default:
		{
			mdclog_write(MDCLOG_ERR, "Error :: Unknown message type %d received from RMR", message->mtype);
			*resend = false;
		}
	}

	return;

};

void process_ric_indication(int message_type, transaction_identifier id, const void *message_payload, size_t message_len) {

	std::cout << "In Process RIC indication" << std::endl;
	std::cout << "ID " << id << std::endl;

	// decode received message payload
  E2AP_PDU_t *pdu = nullptr;
  auto retval = asn_decode(nullptr, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2AP_PDU, (void **) &pdu, message_payload, message_len);

  // print decoded payload
  if (retval.code == RC_OK) {
    char *printBuffer;
    size_t size;
    FILE *stream = open_memstream(&printBuffer, &size);
    asn_fprint(stream, &asn_DEF_E2AP_PDU, pdu);
    mdclog_write(MDCLOG_DEBUG, "Decoded E2AP PDU: %s", printBuffer);

    uint8_t res = procRicIndication(pdu, id);
  }
	else {
		std::cout << "process_ric_indication, retval.code " << retval.code << std::endl;
	}
}

/**
 * Handle RIC indication
 * TODO doxy
 */
uint8_t procRicIndication(E2AP_PDU_t *e2apMsg, transaction_identifier gnb_id)
{
   uint8_t idx;
   uint8_t ied;
   uint8_t ret = RC_OK;
   uint32_t recvBufLen;
   RICindication_t *ricIndication;
   RICaction_ToBeSetup_ItemIEs_t *actionItem;

    // final payload to send to the ML agent
    std::string payload;

   printf("\nE2AP : RIC Indication received");
   ricIndication = &e2apMsg->choice.initiatingMessage->value.choice.RICindication;

   printf("protocolIEs elements %d\n", ricIndication->protocolIEs.list.count);

   for (idx = 0; idx < ricIndication->protocolIEs.list.count; idx++)
   {
      switch(ricIndication->protocolIEs.list.array[idx]->id)
      {
				case 28:  // RIC indication type
				{
					long ricindicationType = ricIndication->protocolIEs.list.array[idx]->value.choice.RICindicationType;
					printf("RIC IndicationType %ld\n", ricindicationType);
					break;
				}
				case 26:  // RIC indication message
				{
                     int payload_size = ricIndication->protocolIEs.list.array[idx]->value.choice.RICindicationMessage.size;
                     char* payload = (char*) calloc(payload_size, sizeof(char));
                     memcpy(payload, ricIndication->protocolIEs.list.array[idx]->value.choice.RICindicationMessage.buf, payload_size);

                    auto *e2SmIndicationMessage = (E2SM_KPM_IndicationMessage_t *) calloc(1, sizeof(E2SM_KPM_IndicationMessage_t));
                    ASN_STRUCT_RESET(asn_DEF_E2SM_KPM_IndicationMessage, e2SmIndicationMessage);

                    asn_decode(nullptr, ATS_ALIGNED_BASIC_PER, &asn_DEF_E2SM_KPM_IndicationMessage,
                               (void **) &e2SmIndicationMessage, ricIndication->protocolIEs.list.array[idx]->value.choice.RICindicationMessage.buf,
                               ricIndication->protocolIEs.list.array[idx]->value.choice.RICindicationMessage.size);
//                    xer_fprint(stderr, &asn_DEF_E2SM_KPM_IndicationMessage, e2SmIndicationMessage);

                    if (e2SmIndicationMessage->present == E2SM_KPM_IndicationMessage_PR_indicationMessage_Format1) {
                        mdclog_write(MDCLOG_INFO, "Format 1 present\n");

                        E2SM_KPM_IndicationMessage_Format1_t *e2SmIndicationMessageFormat1 = e2SmIndicationMessage->choice.indicationMessage_Format1;

                        // extract RAN container, which is where we put our payload
                        vector<uint8_t*> serving_cell_payload_vec;
                        vector<uint8_t*> neighbor_cell_payload_vec;
                        for (size_t i = 0; i < e2SmIndicationMessageFormat1->pm_Containers.list.count; i++) {
                            PM_Containers_List *ranPmContainerList = e2SmIndicationMessageFormat1->pm_Containers.list.array[i];
                            RAN_Container *ranco = ranPmContainerList->theRANContainer;
                            xer_fprint(stderr, &asn_DEF_RAN_Container, ranco);

//                            CU_CP_Usage_Report_Per_UE_t *o_cu_cp_ue = &ranco->reportContainer.choice.oCU_CP_UE;
//
//                            for (int j = 0; j < o_cu_cp_ue->cellResourceReportList.list.count; j++) {
//                                CU_CP_Usage_Report_CellResourceReportItem_t *cellResourceReportItem = o_cu_cp_ue->cellResourceReportList.list.array[j];
//                                for (int z = 0; z < cellResourceReportItem->ueResourceReportList.list.count; z++) {
//                                    OCTET_STRING_t *serving_cell = cellResourceReportItem->ueResourceReportList.list.array[z]->serving_Cell_RF_Type;
//                                    OCTET_STRING_t *neighbor_cell = cellResourceReportItem->ueResourceReportList.list.array[z]->neighbor_Cell_RF;
//
//                                    mdclog_write(MDCLOG_INFO, "Received serving_Cell_RF_Type: %s, neighbor_Cell_RF: %s\n", serving_cell->buf, neighbor_cell->buf);
//                                    serving_cell_payload_vec.push_back(serving_cell->buf);
//                                    neighbor_cell_payload_vec.push_back(neighbor_cell->buf);
//                                }
//                            }
                        }

//                        // combine content of vectors, there should be a single entry in the vector anyway
//                        std::ostringstream serving_cell_payload_oss;
//                        std::copy(serving_cell_payload_vec.begin(), serving_cell_payload_vec.end() - 1, std::ostream_iterator<uint8_t*>(serving_cell_payload_oss, ", "));
//                        serving_cell_payload_oss << serving_cell_payload_vec.back();
//                        std::string serving_cell_payload = serving_cell_payload_oss.str();
//
//                        std::ostringstream neighbor_cell_payload_oss;
//                        std::copy(neighbor_cell_payload_vec.begin(), neighbor_cell_payload_vec.end() - 1, std::ostream_iterator<uint8_t*>(neighbor_cell_payload_oss, ", "));
//                        neighbor_cell_payload_oss << neighbor_cell_payload_vec.back();
//                        std::string neighbor_cell_payload = neighbor_cell_payload_oss.str();
//
//                        mdclog_write(MDCLOG_INFO, "String conversion: serving_Cell_RF_Type %s, neighbor_Cell_RF: %s\n", serving_cell_payload.c_str(), neighbor_cell_payload.c_str());
//
//                        // assemble final payload
//                        if (serving_cell_payload.length() > 0 && neighbor_cell_payload.length() > 0) {
//                            payload = serving_cell_payload + ", " + neighbor_cell_payload;
//                        }
//                        else if (serving_cell_payload.length() > 0 && neighbor_cell_payload.length() == 0) {
//                            payload = serving_cell_payload;
//                        }
//                        else if (serving_cell_payload.length() == 0 && neighbor_cell_payload.length() > 0) {
//                            payload = neighbor_cell_payload;
//                        }
//                        else {
//                            payload = "";
//                        }
//
//                        mdclog_write(MDCLOG_INFO, "Payload from RIC Indication message: %s\n", payload.c_str());
                    }
                    else {
                        mdclog_write(MDCLOG_INFO, "No payload received in RIC Indication message (or was unable to decode received payload\n");
                    }

//                    add_gnb_to_vector_unique(gnb_id);
//
//                    if (payload.length() > 0) {
//                        // add gnb id to payload
//                        payload += "\n{\"gnb_id\": \"" + std::string(reinterpret_cast<char const*>(gnb_id)) + "\"}";
//
//                        mdclog_write(MDCLOG_INFO, "Sending RIC Indication message to agent\n");
//                        send_socket(payload.c_str());
//                    }
//                    else if (payload.length() <= 0) {
//                        mdclog_write(MDCLOG_INFO, "Received empty payload\n");
//                    }
//                    else {
//                        mdclog_write(MDCLOG_INFO, "Returned empty agent IP\n");
//                    }
					break;
				}
      }
   }
   return ret; // TODO update ret value in case of errors
}
