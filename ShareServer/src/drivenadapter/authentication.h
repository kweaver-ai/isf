/***************************************************************************************************
authentication.h:
    Copyright (c) Eisoo Software Inc. (2020), All rights reserved.

Purpose:
    authentication

Author:
    Yuanbin.yan@aishu.cn

Creating Time:
    2023-08-28
***************************************************************************************************/
#ifndef __AUTHENTICATION_H
#define __AUTHENTICATION_H

#include "public/authenticationInterface.h"
#include <ossclient/public/ncIOSSClient.h>
#include <boost/thread/tss.hpp>

#include "gen-cpp/ncTShareMgnt.h"


/* Header file */
class authentication : public authenticationInterface
{
    AB_DECLARE_THREADSAFE_SINGLETON (authentication)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_AUTHENTICATIONINTERFACE

    authentication ();
    ~authentication ();

private:
    // hydra 微服务信息
    String                                              _loginUrl;
    String                                              _auditLogUrl;
    boost::thread_specific_ptr<nsCOMPtr<ncIOSSClient>>  _ossClientPtr;

private:
    void createOSSClient ();

    void handleMessage(const String& userId, ncTokenVisitorType typ, ncTLogLevel level, int opType, const String& msg,
                        const String& exmsg, const String& ip, const String& macAddress, const String& userAgent, JSON::Value & logMsgs);
};

#endif // __AUTHENTICATION_H
