#ifndef __NC_ACS_OWNER_MANAGER_MOCK_H
#define __NC_ACS_OWNER_MANAGER_MOCK_H

#include <gmock/gmock.h>
#include <acsprocessor/public/ncIACSOwnerManager.h>

class ncACSOwnerManagerMock: public ncIACSOwnerManager
{
    XPCOM_OBJECT_MOCK (ncACSOwnerManagerMock)

public:
    MOCK_METHOD2(SetAllOwner, void(const String&, vector<dbOwnerInfo>&));
    MOCK_METHOD2(IsOwner, bool(const String&, const String&));
    MOCK_METHOD2(CheckIsOwner, void(const String&, const String&));
    MOCK_METHOD2(GetOwnerIds, void(const String&, vector<String>&));
    MOCK_METHOD2(GetOwnerInfos, void(const String&, vector<ncOwnerInfo>&));
    MOCK_METHOD2(GetOwnerDocIds, void(const String&, set<String>&));
};

#endif // End __NC_ACS_OWNER_MANAGER_MOCK_H
