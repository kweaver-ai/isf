#ifndef __NC_DB_LOCK_MANAGER_MOCK_H
#define __NC_DB_LOCK_MANAGER_MOCK_H

#include <gmock/gmock.h>
#include <acsdb/public/ncIDBLockManager.h>

class ncDBLockManagerMock: public ncIDBLockManager
{
    XPCOM_OBJECT_MOCK (ncDBLockManagerMock);

public:
    MOCK_METHOD1(Add, void(const ncDBLockInfo&));
    MOCK_METHOD2(GetAutoExpired, void(int64, vector<ncDBLockInfo>&));
    MOCK_METHOD2(GetAppointedExpired, void(int64, vector<ncDBLockInfo>&));
    MOCK_METHOD2(Get, bool(const String&, String&));
    MOCK_METHOD1(Delete, void(const String &));
    MOCK_METHOD1(DeleteSubs, void(const String &));
    MOCK_METHOD1(DeleteByUserId, void(const String &));
    MOCK_METHOD1(DeleteByDocId, void(const String &));
    MOCK_METHOD2(GetAllLocked, void(const String &, vector<ncDBLockInfo>&));
    MOCK_METHOD4(SearchFileLockInfos, void(const String&, const vector<String>&, const String&, map<String, ncDBLockFileInfo>&));
};

#endif // End __NC_DB_LOCK_MANAGER_MOCK_H
