.. This work is licensed under a Creative Commons Attribution 4.0 International License.
.. SPDX-License-Identifier: CC-BY-4.0
.. Copyright (C) 2020 AT&T

============================================================================================ 
HelloWorld xAPP (C++) 
============================================================================================ 
-------------------------------------------------------------------------------------------- 
User's Guide 
-------------------------------------------------------------------------------------------- 
 
Introduction 
============================================================================================ 

The RIC platform provides set of functions that the xAPPs can use to accomplish their tasks. 
The HW xAPP is envisioned to provide xAPP developers, examples of implementing these sets of functions. 
Note, HW xAPP does not address/implement any RIC Usecases. 

HelloWorld xAPP Features 
============================================================================================ 

RIC Platform provides many APIs and libraries to aid the development of xAPPs. All xAPPs will have some custom 
processing functional logic core to the xApp and some additional non-functional platform related processing using 
these APIs and libraries. HW xAPP attempts to show the usage of such additional platform processing using RIC platform APIs and libraries.


The Hello World xApp demonstrates how an xApp uses the A1, and E2 interfaces and persistent database read-write operations.
The following paragraphs cover the various steps involved to create an HelloWorld xApp instance, setting its configuration,
retrieving R-NIB data, sending subscription, connecting SDL, RME & A1 Healthcheck and usage of "Hello World SM"

HelloWorld Creation
============================================================================================ 
 
The creation of the xApp instance is as simple as invoking 
the object's constructor with two required parameters: 


HW xAPP, may choose to create following objects for obtaining desired set of functionalities provided under xapp-utils:

XappRmr
-------------------------------------------------------------------------------------------- 
An xAPP can have the capability of receiving and sending rmr messages. This is achieved by creating an XappRmr object. The constructor of xAPPRMR object requires xAPP developer to provide  
xAPP's listening port and developer configurable number of attempts need to be made to send the message. The key functionalities of the class being :
        
1. Setting RMR initial context: ...xapp_rmr_init(...)
        
2. Sending RMR message: ...xapp_rmr_send(xapp_rmr_header, void*)
        
3. Receiving RMR message: ...xapp_rmr_receive(msghandler,...)

The RMR Header can be defined using xapp_rmr_header :
::   

        typedef struct{
			struct timespec ts;
			int32_t message_type; //mandatory
			int32_t state;
			int32_t payload_length; //mandatory
			unsigned char sid[RMR_MAX_SID]; 
			unsigned char src[RMR_MAX_SRC]; 
			unsigned char meid[RMR_MAX_MEID];

		}  xapp_rmr_header;

Except for message type and payload length, its developers prerogative to use remaining header information. 
The XappMsgHandler (msghandler) instance in xapp_rmr_receive function handles received messages. The handling of messages is based on
the usecase catered by a xAPP. Hence, XappMsgHandler class used in HW xAPP is not very comprehensive and addresses only Healthcheck Messages.

XappSettings
------------------------------------------------------------------------------------------- 
An xAPP has the capability to use environment variables or xapp-descriptor information as its configuration settings 
creating XappSettings object, whose key functions being :
        
1. Loading Default Settings: ...loadDefaultSettings()
        
2. Loading Environment Variables: ...loadEnvVarSettings()
        
3. Loading Command Line Settings: ...loadCmdlineSettings(argc, argv)

XappSDL
-------------------------------------------------------------------------------------------- 
An xAPP can have the capability to read and write into a persistent storage (key-value store) creating XappSDL object. 
for a namespace. The key functionalities available currently are: 


1. Getting Data from SDL: ... get_data(...);

2. Setting Data to SDL:	... set_data(...);	

The HW xAPP can be instantiationed as following:
::

	HW_Xapp = Xapp(XappRmr object, XappSettings object,...);


HelloWorld E2 and A1 Message Handling
============================================================================================ 
Helper Objects
-------------------------------------------------------------------------------------------- 
HW xAPP creates wrapper datastructures mirroring ASN and JSON messages. These datastructures facilitate processing of 
E2 and A1 messages in the xAPP. A sample  helper object for A1 Health Check message being:
::

	struct a1_policy_helper{
		std::string operation;
		std::string policy_type_id;
		std::string policy_instance_id;
		std::string handler_id;
		std::string status;
	};

And a sample E2AP Control datastructure:
::

	struct ric_control_helper{
  		ric_control_helper(void):req_id(1), req_seq_no(1), func_id(0), action_id(1), control_ack(-1), cause(0), sub_cause(0), control_status(1), control_msg(0), control_msg_size(0), control_header(0), control_header_size(0), call_process_id(0), call_process_id_size(0){};
  		long int req_id, req_seq_no, func_id, action_id,  control_ack, cause, sub_cause, control_status;
  
  		unsigned char* control_msg;
  		size_t control_msg_size;
  
  		unsigned char* control_header;
  		size_t control_header_size;
  
  		unsigned char *call_process_id;
  		size_t call_process_id_size;
  
	};

As mentioned, these datastructures are very much tied to the message specifications.



ASN Encoding/Decoding
-------------------------------------------------------------------------------------------- 
RIC platform provided ASN1C (modified) library is used for processing ASN1 messages. HW xAPP, for each 
ASN message type, uses a class which is responsible for handling a particular message type.
The class encapsulates, the APIs and datastructures used in ASN1C using helper objects. For example:
::

	class ric_control_response{
		...
		bool encode_e2ap_control_response(..., ric_control_helper &);
		bool set_fields(..., ric_control_helper &);
		bool get_fields(..., ric_control_helper &);
		...
	}

Note, the helper objects and message type processing classes can be found under xapp-asn subdirectories.

E2AP Subscription
-------------------------------------------------------------------------------------------- 
In HW xAPP, we consider sunny-side scenario, in which for a E2AP subscription request sent, it is assumed, 
that HW xAPP will be receiving E2AP subscription response. Handling advanced subscription (class SubscriptionHandler) flows is out of the 
scope of HW xAPP. Current form of class SubscriptionHandler has following key functionalities:

1. manage_subscription_request(...)

2. manage_subscription_response(...)


The manage_subscription_request function waits for the response for a specified time for subscription response 
and if no response is received within a specified time, gives a time out error message. A subscription message 
is created using ASN Encodong/Decoding and Helper classes. (Refer test_sub.h). HW xAPP sends the subscriptions based 
on the gNodeB IDs received from RNIB. Please refer following function in xapp.* for RNIB transactions: set_rnib_gnblist(...) 


E2SM Subscription, Indication, Control
-------------------------------------------------------------------------------------------- 
HellowWorld E2SM (e2sm-HelloWorld-v001.asn) is an example E2SM available in the docs directory. The Helper and 
encoding/decoding classes are in xapp-asn/e2sm. Sample code for control message E2SM:
::

	//ControlHeader 
	unsigned char header_buf[128];
	size_t header_buf_len = 128;

	//ControlMessage
	unsigned char msg_buf[128];
	size_t msg_buf_len = 128;

	bool res;
	
	e2sm_control_helper e2sm_cntrldata; //helper object
	e2sm_control e2sm_cntrl; //encoding/decoding object

	unsigned char msg[20] = "HelloWorld";

	e2sm_cntrldata.header = 1001;
	e2sm_cntrldata.message = msg;
	e2sm_cntrldata.message_len = strlen((const char*)e2sm_cntrldata.message);


	// Encode the control header
	res = e2sm_cntrl.encode_control_header(&header_buf[0], &header_buf_len, e2sm_cntrldata);
	if(!res)
		std::cout << e2sm_cntrl.get_error() << std::endl;
	
	// Encode the control message
	res = e2sm_cntrl.encode_control_message(&msg_buf[0], &msg_buf_len, e2sm_cntrldata);
	if(!res)
		std::cout << e2sm_cntrl.get_error() << std::endl;


RMR and A1 Healtcheck
-------------------------------------------------------------------------------------------- 
On receiving health check request message types (A1_HEALTHCHECK_REQ, RMR_HEALTHCHECK_REQ), HW xAPP sends RMR
response (A1_HEALTHCHECK_RES, RMR_HEALTHCHECK_RESP) adding appropriate responses using RMR Return to Sender
functionality.  
