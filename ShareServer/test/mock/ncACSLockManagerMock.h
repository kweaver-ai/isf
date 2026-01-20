#ifndef __NC_ACS_LOCK_MANAGER_MOCK_H
#define __NC_ACS_LOCK_MANAGER_MOCK_H

#include <gmock/gmock.h>
#include <acsprocessor/public/ncIACSLockManager.h>

class ncACSLockManagerMock: public ncIACSLockManager
{
    XPCOM_OBJECT_MOCK (ncACSLockManagerMock);

public:
    MOCK_METHOD2(SetAutolockConfig, void(bool, int64));
    MOCK_METHOD0(IsAutolockEnabled, bool(void));
    MOCK_METHOD1(Delete, void(const String &));
    MOCK_METHOD1(DeleteSubs, void(const String &));
    MOCK_METHOD1(DeleteByUserId, void(const String &));
};

#endif // End __NC_ACS_LOCK_MANAGER_MOCK_H
