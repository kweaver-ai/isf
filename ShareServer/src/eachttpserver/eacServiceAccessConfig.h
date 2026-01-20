/***************************************************************************************************
eacServiceAccessConfig.h:
    Copyright (c) Eisoo Software Inc. (2009 - 2013), All rights reserved.

Purpose:
    控制管理的公共头文件


Creating Time:
    2021-04-28
***************************************************************************************************/
#ifndef __NC_EAC_SERVICE_ACCESS_CONFIG_H__
#define __NC_EAC_SERVICE_ACCESS_CONFIG_H__
#include <abprec.h>

///////////////////////////////////////////////////////////////////////////////////////////////////
// 公共类型

class EacServiceAccessConfig {

    AB_DECLARE_THREADSAFE_SINGLETON(EacServiceAccessConfig)
private:
    EacServiceAccessConfig();
    ~EacServiceAccessConfig();
public:
    String   sharemgntHost;
    int      sharemgntPort;
    String   eacpThriftHost;
    int      eacpThriftPort;
};

#endif
