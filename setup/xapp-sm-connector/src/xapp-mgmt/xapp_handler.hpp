/*
 * xapp_handler.hpp
 *
 *  Created on: Mar 16, 2020
 *  Author: Shraboni Jana
 */

#ifndef SRC_XAPP_MGMT_XAPP_HANDLER_HPP_
#define SRC_XAPP_MGMT_XAPP_HANDLER_HPP_

class XappHandler{
	XappHandler *xhandler;
public:
	virtual ~XappHandler(){delete xhandler;};
	virtual void register_handler(XappHandler *xhandler) = 0;
	virtual XappHandler* get_handler() = 0;
};



#endif /* SRC_XAPP_MGMT_XAPP_HANDLER_HPP_ */
