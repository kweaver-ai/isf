/***************************************************************************************************
deploy.h:
    Copyright (c) Eisoo Software, Inc.(2009 - 2020), All rights reserved

Purpose:
    deploy

Author:
    xu.zhi@aishu.cn

Creating Time:
    2021-04-29
***************************************************************************************************/

#ifndef __DEPLOY_H
#define __DEPLOY_H

#include <abprec.h>

#include "public/deployInterface.h"
#include <ossclient/public/ncIOSSClient.h>
#include <boost/thread/tss.hpp>

/* Header file */
class deploy : public deployInterface
{
    AB_DECLARE_THREADSAFE_SINGLETON (deploy)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_DEPLOYINTERFACE

    deploy();
    ~deploy();

private:
    String                              _getVipUrl;
    boost::thread_specific_ptr<nsCOMPtr<ncIOSSClient>> _ossClientPtr;

private:
    void createOSSClient ();
};


#endif // __DEPLOY_H
