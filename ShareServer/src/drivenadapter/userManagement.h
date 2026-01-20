/***************************************************************************************************
usermanagement.h:
    Copyright (c) Eisoo Software, Inc.(2009 - 2020), All rights reserved

Purpose:
    usermanagement access control

Author:
    Young.yu@aishu.cn

Creating Time:
    2020-11-17
***************************************************************************************************/

#ifndef __USER_MANAGEMENT_H
#define __USER_MANAGEMENT_H

#include <abprec.h>

#include "public/userManagementInterface.h"
#include <ossclient/public/ncIOSSClient.h>
#include <boost/thread/tss.hpp>

/* Header file */
class userManagement : public userManagementInterface
{
    AB_DECLARE_THREADSAFE_SINGLETON (userManagement)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_USERMANAGEMENTINTERFACE

    userManagement();
    ~userManagement();

private:
    String                              _getAccessorIDsByDepartIDUrl;
    String                              _getAccessorIDsByUserIDUrl;
    String                              _getOrgNamesByIDsUrl;
    String                              _deleteDepartUrl;
    boost::thread_specific_ptr<nsCOMPtr<ncIOSSClient>>  _ossClientPtr;

    // 枚举映射 map 表
    map<string, ncUserRoleType>         _userRoleTypeMap;

private:
    void createOSSClient ();
};


#endif // __USER_MANAGEMENT_H
