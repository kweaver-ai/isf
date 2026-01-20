#ifndef __NC_ACS_DB_OWNER_SHARE_MOCK_H
#define __NC_ACS_DB_OWNER_SHARE_MOCK_H

#include <gmock/gmock.h>
#include <acsdb/public/ncIDBOwnerManager.h>

class ncDBOwnerManagerMock: public ncIDBOwnerManager
{
    XPCOM_OBJECT_MOCK (ncDBOwnerManagerMock);

public:
    MOCK_METHOD1(AddOwner, void(const dbOwnerInfo&));
    MOCK_METHOD1(DeleteOwner, void(const dbOwnerInfo&));
    MOCK_METHOD1(DeleteOwnerInfosByDocId, void(const String&));
    MOCK_METHOD2(GetOwnerInfosByDocId, void(const String&, vector<dbOwnerInfo>&));
    MOCK_METHOD3(GetInheritOwnerInfosByDocId, void(const String&, vector<dbOwnerInfo>&, bool));
    MOCK_METHOD2(GetOwnerInfosByUserId, void(const String&, vector<dbOwnerInfo>&));
    MOCK_METHOD1(DeleteOwnerByFileId, void(const String&));
    MOCK_METHOD1(DeleteOwnerByDirId, void(const String&));
    MOCK_METHOD1(DeleteOwnerByUserId, void(const String&));
};

#endif // End __NC_ACS_DB_MOCK_H
