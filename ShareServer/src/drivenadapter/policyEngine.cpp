/***************************************************************************************************
policyEngine.cpp:
    Copyright (c) Eisoo Software Inc. (2020), All rights reserved.

Purpose:
    policyEngine 实现

Author:
    Will.lv@aishu.cn

Creating Time:
    2020-11-10
***************************************************************************************************/
#include <abprec.h>
#include "policyEngine.h"
#include "drivenadapter.h"
#include "serviceAccessConfig.h"
#include <dataapi/ncJson.h>
#include <ossclient/public/ncIOSSClient.h>
#include <ossclient/ossclient.h>

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (policyEngine, policyEngineInterface)

// protected
NS_IMETHODIMP_(nsrefcnt) policyEngine::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) policyEngine::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (policyEngine)

policyEngine::policyEngine (void): _ossClientPtr(0)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("this: %p"), this);
}

policyEngine::~policyEngine (void)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("this: %p"), this);
}

void policyEngine::createOSSClient ()
{
    if (nullptr == _ossClientPtr.get ()) {
        nsresult ret;
        nsCOMPtr<ncIOSSClient> ossClient = do_CreateInstance (NC_OSS_CLIENT_CONTRACTID, &ret);
        if (NS_FAILED (ret)){
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR,
                     _T("Failed Create OSSClient: 0x%x"), ret);
        }
        _ossClientPtr.reset (new nsCOMPtr<ncIOSSClient>(ossClient));
    }
}


NS_IMETHODIMP_(bool) policyEngine::Audit_ClientRestriction (const String& clientType)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("[BEGIN] this: %p"), this);
    createOSSClient ();
    ncOSSResponse res;
    vector<string> inHeaders;
    String url;
    url.format(_T("http://%s:%d/api/policy-management/v1/sign_in_policy/client_restriction?client_type=%s"),
               ServiceAccessConfig::getInstance()->policyMgntHost.getCStr(), ServiceAccessConfig::getInstance()->policyMgntPort, clientType.getCStr());
    (*_ossClientPtr)->Get (url.getCStr (), inHeaders, 30, res);
    if (res.code != 200){
        if (res.code == 0) {
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Could not connect to server"));
        }
        else{
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Code:%d. Cause:%s."), res.code, res.body.c_str ());
        }
    }

    JSON::Value response;
    JSON::Reader::read (response, res.body.c_str (), res.body.size ());
    bool needAudit = response["result"].b ();
    NC_DRIVEN_ADAPTER_TRACE (_T("[END] this: %p"), this);
    return needAudit;
}
