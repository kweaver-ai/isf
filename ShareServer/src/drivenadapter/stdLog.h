/***************************************************************************************************
stdLog.h:
    Copyright (c) Eisoo Software Inc. (2021), All rights reserved.

Purpose:
    stdLog manager

Author:
    xu.zhi@aishu.cn

Creating Time:
    2022-08-02
***************************************************************************************************/
#ifndef __STD_LOG_H
#define __STD_LOG_H

#include "public/stdLogInterface.h"
#include <dataapi/ncJson.h>
#include <ncMQClient.h>

#define TELEMETRY_LOG_USER_OPERATION       _T("as.audit_log.operation_log.user_login")

/* Header file */
class stdLog : public stdLogInterface
{
    AB_DECLARE_THREADSAFE_SINGLETON (stdLog)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_STDLOGINTERFACE

    stdLog ();
    ~stdLog ();

private:
    String getLogDescription (const String& actorName, ncClientType clientType);
    void formatLog (const ncOperator& actor, JSON::Value& eventJson);

private:
    map<ncClientType, String>         _clientTypeIntToStr;
    boost::shared_ptr<ncMQClient>     _mqClient;
};

#endif
