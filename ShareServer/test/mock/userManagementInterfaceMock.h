/*************************************************************************************************
userManagementInterfaceMock.h:
    Copyright (c) Eisoo Software, Inc.(2012 - 2020), All rights reserved.

Purpose:
    drivenadapters userManagementInterfaceMock mock

Author:
    Young.yu@aishu.cn

Creating Time:
    2020-12-04
***************************************************************************************************/
#ifndef __USER_MANAGEMENT_MOCK_H
#define __USER_MANAGEMENT_MOCK_H

#include <gmock/gmock.h>
#include <drivenadapter/public/userManagementInterface.h>

class userManagementMock: public userManagementInterface
{
    XPCOM_OBJECT_MOCK (userManagementMock);

public:
    MOCK_METHOD2(GetAccessorIDsByDepartID, void(const String&, set<String>&));
    MOCK_METHOD2(GetAccessorIDsByUserID, void(const String&, set<String>&));
    MOCK_METHOD2(GetOrgNameIDInfo, void(const ncOrgIDInfo&, ncOrgNameIDInfo&));
    MOCK_METHOD2(GetUserInfo, void(const String&, UserInfo&));
    MOCK_METHOD2(GetAppInfo, void(const String&, AppInfo&));
    MOCK_METHOD2(BatchGetUserInfo, void(const vector<String>& userIds, vector<UserInfo>& userInfos));
    MOCK_METHOD1(DeleteDepart, void(const String& departID));
};

#endif // End __USER_MANAGEMENT_MOCK_H
