/***************************************************************************************************
serviceAccessConfig.h:
    Copyright (c) Eisoo Software Inc. (2009 - 2013), All rights reserved.

Purpose:
    配置文件信息读取


Creating Time:
    2021-04-28
***************************************************************************************************/
#ifndef __NC_SERVICE_ACCESS_CONFIG_H__
#define __NC_SERVICE_ACCESS_CONFIG_H__
#include <abprec.h>

///////////////////////////////////////////////////////////////////////////////////////////////////
// 公共类型

class ServiceAccessConfig {

    AB_DECLARE_THREADSAFE_SINGLETON(ServiceAccessConfig)
private:
    ServiceAccessConfig();
    ~ServiceAccessConfig();
public:
    String   deployHost;
    int      deployPost;
    String   hydraAdminHost;
    int      hydraAdminPort;
    String   userManagePrivateHost;
    int      userManagePrivatePort;
    String   policyEngineHost;
    int      policyEnginePort;
    String   policyMgntHost;
    int      policyMgntPort;
    String   ossgatewayPrivateProtocol;
    String   ossgatewayPrivateHost;
    int      ossgatewayPrivatePort;
    String   authenticationPrivateHost;
    int      authenticationPrivatePort;
};

#endif
