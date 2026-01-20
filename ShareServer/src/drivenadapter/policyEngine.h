/***************************************************************************************************
policyEngine.h:
    Copyright (c) Eisoo Software Inc. (2020), All rights reserved.

Purpose:
    policyEngine

Author:
    Will.lv@aishu.cn

Creating Time:
    2020-11-10
***************************************************************************************************/
#ifndef __POLICY_ENGINE_H
#define __POLICY_ENGINE_H

#include "public/policyEngineInterface.h"
#include <ossclient/public/ncIOSSClient.h>
#include <boost/thread/tss.hpp>

/* Header file */
class policyEngine : public policyEngineInterface
{
    AB_DECLARE_THREADSAFE_SINGLETON (policyEngine)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_POLICYENGINEINTERFACE

    policyEngine ();
    ~policyEngine ();

private:
    // 策略引擎服务信息
    boost::thread_specific_ptr<nsCOMPtr<ncIOSSClient>> _ossClientPtr;

    void createOSSClient ();
    bool auditApprovalStatus(const String& user);
};

#endif // __POLICY_ENGINE_H

