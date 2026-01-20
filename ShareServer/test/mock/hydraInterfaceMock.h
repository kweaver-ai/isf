/*************************************************************************************************
hydraInterfaceMock.h:
    Copyright (c) Eisoo Software, Inc.(2012 - 2020), All rights reserved.

Purpose:
    drivenadapters hydraInterfaceMock mock

Author:
    Sunshine.tang@aishu.cn

Creating Time:
    2021-03-16
***************************************************************************************************/
#ifndef __HYDRA_MOCK_H
#define __HYDRA_MOCK_H

#include <gmock/gmock.h>
#include <drivenadapter/public/hydraInterface.h>

class hydraMock: public hydraInterface
{
    XPCOM_OBJECT_MOCK (hydraMock);

public:
    MOCK_METHOD2(IntrospectToken, void(const String&, ncTokenIntrospectInfo&));
    MOCK_METHOD2(DeleteConsentAndLogin, void(const String&, const String&));
    MOCK_METHOD2(GetConsentInfo, void(const String&, vector<ncTokenIntrospectInfo>&));
};

#endif // End __HYDRA_MOCK_H
