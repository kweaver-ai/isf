/***************************************************************************************************
deploy.cpp:
    Copyright (c) Eisoo Software Inc. (2020), All rights reserved.

Purpose:
    deploy 服务接口调用

Author:
    xu.zhi@aishu.cn

Creating Time:
    2021-04-29
***************************************************************************************************/
#include <abprec.h>
#include "deploy.h"

#include "drivenadapter.h"
#include <dataapi/ncJson.h>
#include <ossclient/public/ncIOSSClient.h>
#include <ossclient/ossclient.h>

#include "serviceAccessConfig.h"

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (deploy, deployInterface)

// protected
NS_IMETHODIMP_(nsrefcnt) deploy::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) deploy::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (deploy)

deploy::deploy ()
    : _getVipUrl()
    , _ossClientPtr(0)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("this: %p"), this);

    _getVipUrl.format(_T("http://%s:%d/api/deploy-manager/v1/access-addr/app"),
        ServiceAccessConfig::getInstance()->deployHost.getCStr(), ServiceAccessConfig::getInstance()->deployPost);
}

deploy::~deploy (void)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("this: %p"), this);
}

void deploy::createOSSClient ()
{
    if (nullptr == _ossClientPtr.get ()) {
        nsresult ret;
        nsCOMPtr<ncIOSSClient> ossClient = do_CreateInstance (NC_OSS_CLIENT_CONTRACTID, &ret);
        if (NS_FAILED (ret)){
            THROW_E (ACS_DRIVEN_ADAPTER, FAILED_TO_CREATE_XPCOM_INSTANCE,
                     _T("Failed to create OSSClient: %s(0x%x)"),
                     String::toString ((int64)ret).getCStr (), (int64)ret);
        }
        _ossClientPtr.reset (new nsCOMPtr<ncIOSSClient>(ossClient));
    }
}

/* [notxpcom] void GetAccessAddr (in StringRef host, in StringRef port); */
NS_IMETHODIMP_(void) deploy::GetAccessAddr(String & host, String & port)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("[BEGIN] this: %p"), this);
    createOSSClient ();
    ncOSSResponse res;
    vector<string> inHeaders;

    (*_ossClientPtr)->Get (_getVipUrl.getCStr (), inHeaders, 30, res);
    if (res.code != 200)
    {
        if (res.code == 0) {
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("Could not connect to server"));
        }
        else{
            JSON::Value errorJson;
            errorJson["code"] = res.code;
            errorJson["body"] = res.body;
            throw errorJson;
        }
    }

    JSON::Value accessAddr;
    JSON::Reader::read (accessAddr, res.body.c_str (), res.body.length ());

    host = toCFLString (accessAddr["host"].s ().c_str ());
    port = toCFLString (accessAddr["port"].s ().c_str ());

    NC_DRIVEN_ADAPTER_TRACE (_T("[END] this: %p"), this);
}
