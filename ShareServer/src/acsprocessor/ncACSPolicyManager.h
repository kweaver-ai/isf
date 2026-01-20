/***************************************************************************************************
ncACSPolicyManager.cpp:
    Copyright (c) Eisoo Software Inc. (2009 - 2013), All rights reserved.

Purpose:
    acs policy manager 接口

Author:
    xu.zhi@aishu.cn

Creating Time:
    2020-7-30
***************************************************************************************************/
#ifndef __NC_ACS_POLICY_MANAGER_H
#define __NC_ACS_POLICY_MANAGER_H

#include <acsprocessor/public/ncIACSPolicyManager.h>
#include <ossclient/public/ncIOSSClient.h>
#include <acssharemgnt/public/ncIACSShareMgnt.h>
#include <acsprocessor/public/ncIACSDeviceManager.h>
#include "ncRefreshTokenThread.h"
#include "ncActiveRecordThread.h"
#include <boost/thread/tss.hpp>
#include "drivenadapter/public/policyEngineInterface.h"

/* Header file */
class ncACSPolicyManager : public ncIACSPolicyManager
{
    AB_DECLARE_THREADSAFE_SINGLETON (ncACSPolicyManager)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_NCIACSPOLICYMANAGER

    ncACSPolicyManager();

private:
    String checkIpUrl;
    String getClientByIdUrl;

private:
    ~ncACSPolicyManager();
    void createOSSClient ();


protected:
    boost::thread_specific_ptr<nsCOMPtr<ncIOSSClient>> _ossClientPtr;

    // 后台 token 刷新线程
    static ncRefreshTokenThread*       _srefreshTokenThread;

    // 后台活跃记录线程
    static ncActiveRecordThread*       _sactiveRecordThread;

    // sharemgnt接口
    nsCOMPtr<ncIACSShareMgnt> _acsShareMgnt;
    nsCOMPtr<policyEngineInterface> _policyEngine;
    map<ACSClientType, String> _clientStringTypeMap;
};

#endif // __NC_ACS_POLICY_MANAGER_H
