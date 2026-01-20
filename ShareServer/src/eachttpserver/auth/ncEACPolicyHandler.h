/****************************************************************************************************
ncEACPolicyHandler.h
     Copyright (c) Eisoo Software, Inc.(2009 - 2015), All rights reserved.

Purpose:
    ncEACPolicyHandler

Author:
    Sunshine.tang@aishu.cn

Created Time:
    2021-02-03
****************************************************************************************************/
#ifndef __NC_EAC_POLICY_HANDLER_H__
#define __NC_EAC_POLICY_HANDLER_H__

#include <acsprocessor/public/ncIACSPolicyManager.h>

class ncEACPolicyHandler
{
public:
    ncEACPolicyHandler (ncIACSPolicyManager* acsPolicyManager);
    ~ncEACPolicyHandler (void);

    /***
     * 检查可执行性
     */
    void CheckPolicy (brpc::Controller* cntl);

private:
    // 解析实名用户信息 json格式
    void parseUserInfo (brpc::Controller* cntl, ncPolicyCheckInfo& policyInfo, JSON::Value& requestJson);

private:
    nsCOMPtr<ncIACSPolicyManager>        _acsPolicyManager;   // policy管理

    //枚举映射map表
    map<string, ACSAccountType>          _accountTypeMap;
    map<string, ACSClientType>           _clientTypeMap;
};

#endif  // __NC_EAC_POLICY_HANDLER_H__
