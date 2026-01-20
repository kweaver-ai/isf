/***************************************************************************************************
acsServiceAccessConfig.h:
    Copyright (c) Eisoo Software Inc. (2009 - 2013), All rights reserved.

Purpose:
    配置文件信息读取


Creating Time:
    2021-04-28
***************************************************************************************************/
#ifndef __NC_ACS_SERVICE_ACCESS_CONFIG_H__
#define __NC_ACS_SERVICE_ACCESS_CONFIG_H__
#include <abprec.h>

///////////////////////////////////////////////////////////////////////////////////////////////////
// 公共类型

class AcsServiceAccessConfig {

    AB_DECLARE_THREADSAFE_SINGLETON(AcsServiceAccessConfig)
private:
    AcsServiceAccessConfig();
    ~AcsServiceAccessConfig();
public:
    String   sharemgntHost;
    int      sharemgntPort;
    String   eacpThriftHost;
    int      eacpThriftPort;
    String   policyMgntHost;
    int      policyMgntPort;
};

#endif
