/***************************************************************************************************
pluginMessage.cpp:
    Copyright (c) Eisoo Software Inc. (2022), All rights reserved.

Purpose:
    pluginMessage 服务接口调用

Author:
    Dylan.gao@aishu.cn

Creating Time:
    2022-03-25
***************************************************************************************************/
#include <abprec.h>
#include "pluginMessage.h"

#include "drivenadapter.h"
#include <dataapi/ncJson.h>
#include <ossclient/public/ncIOSSClient.h>
#include <ossclient/ossclient.h>

#include "serviceAccessConfig.h"

#define NC_LINE_BREAK                   _T("/r/n")
#define NSQ_THIRD_PLUGIN_MESSAGE_PUSH   _T("thirdparty_message_plugin.message.push")

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (pluginMessage, pluginMessageInterface)

// protected
NS_IMETHODIMP_(nsrefcnt) pluginMessage::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) pluginMessage::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (pluginMessage)

pluginMessage::pluginMessage ()
{
    NC_DRIVEN_ADAPTER_TRACE (_T("this: %p"), this);

    _mqClient = ncMQClient::GetConnectorFromFile("/sysvol/conf/service_conf/mq_config.yaml");

    nsresult ret;
}

pluginMessage::~pluginMessage (void)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void SendPluginMessage ([const] in AcsMessageRef msg); */
NS_IMETHODIMP_(void) pluginMessage::SendPluginMessage(const acsMessage& msg)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("[BEGIN] this: %p"), this);

    // 封装headers
    JSON::Object headers;
    headers["msg_id"] = msg.msgId.getCStr ();
    JSON::Array& receiversArr = headers["receivers"].a ();
    for (auto receiverIt = msg.receivers.begin (); receiverIt != msg.receivers.end (); ++receiverIt) {
        receiversArr.push_back (JSON::OBJECT);
        JSON::Object& receiver = receiversArr.back ().o ();
        receiver["id"] = receiverIt->id.getCStr ();
        receiver["account"] = receiverIt->account.getCStr ();
        receiver["name"] = receiverIt->name.getCStr ();
        receiver["email"] = receiverIt->email.getCStr ();
        receiver["telephone"] = receiverIt->telephone.getCStr ();
        receiver["third_attr"] = receiverIt->thirdAttr.getCStr ();
        receiver["third_id"] = receiverIt->thirdId.getCStr ();
    }
    std::string strHeaders;
    JSON::Writer::write (headers, strHeaders);

    SystemLog::getInstance ()->log (__FILE__, __LINE__, INFORMATION_LOG_TYPE,
                                            _T("Info : send Plugin Message: channel is %s\n; headers is %s\n, payload is %s\n"),
                                            msg.channel.getCStr (), strHeaders.c_str (), msg.content.getCStr ());

    String message;
    message.append (msg.channel + NC_LINE_BREAK);
    message.append (toCFLString (strHeaders) + NC_LINE_BREAK);
    message.append (msg.content);

    try {
        _mqClient->Pub (NSQ_THIRD_PLUGIN_MESSAGE_PUSH, message);
    } catch (Exception& e) {
        THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR,
            _T("MQ pub message failed, topic %s, errorId %d, err %s"), msg.channel.getCStr (),
            e.getErrorId (), e.toFullString ().getCStr ());
    }

    NC_DRIVEN_ADAPTER_TRACE (_T("[END] this: %p"), this);
    return;
}