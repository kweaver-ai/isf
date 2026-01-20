/***************************************************************************************************
ossGateway.h:
    Copyright (c) Eisoo Software Inc. (2020), All rights reserved.

Purpose:
    ossGateway

Author:
    Young.yu@aishu.cn

Creating Time:
    2023-05-30
***************************************************************************************************/
#ifndef __OSS_GATEWAY_H
#define __OSS_GATEWAY_H

#include "public/ossGatewayInterface.h"
#include <ossclient/public/ncIOSSClient.h>
#include <boost/thread/tss.hpp>

/* Header file */
class ossGateway : public ossGatewayInterface
{
    AB_DECLARE_THREADSAFE_SINGLETON (ossGateway)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_OSSGATEWAYINTERFACE

    ossGateway ();
    ~ossGateway ();

private:

    boost::thread_specific_ptr<nsCOMPtr<ncIOSSClient>> _ossClientPtr;
    String _getLocalStorageInfoUrl;
    String _uploadInfoUrl;
    String _uploadPartUrl;
    String _completeUploadUrl;
    String _getDownloadInfoUrl;

    void createOSSClient ();
    bool auditApprovalStatus(const String& user);
    String UrlEncode3986 (const String &input);
};

#endif // __OSS_GATEWAY_H

