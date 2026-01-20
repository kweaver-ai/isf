/***************************************************************************************************
hydra.h:
    Copyright (c) Eisoo Software Inc. (2020), All rights reserved.

Purpose:
    hydra

Author:
    Sunshine.tang@aishu.cn

Creating Time:
    2021-02-22
***************************************************************************************************/
#ifndef __HYDRA_H
#define __HYDRA_H

#include "public/hydraInterface.h"
#include <ossclient/public/ncIOSSClient.h>
#include <boost/thread/tss.hpp>

/* Header file */
class hydra : public hydraInterface
{
    AB_DECLARE_THREADSAFE_SINGLETON (hydra)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_HYDRAINTERFACE

    hydra ();
    ~hydra ();

private:
    String UrlEncode3986 (const String &input);

private:
    // hydra 微服务信息
    String                              consentSessionUrl;
    String                              loginSessionUrl;
    boost::thread_specific_ptr<nsCOMPtr<ncIOSSClient>> _ossClientPtr;

    // type 类型转换 string-->enum class
    map<string, ncTokenVisitorType>     _visitorTypeMap;
    map<string, ncAccountType>          _accountTypeMap;
    map<string, ncClientType>           _clientTypeMap;

private:
    void createOSSClient ();
};

#endif // __HYDRA_H
