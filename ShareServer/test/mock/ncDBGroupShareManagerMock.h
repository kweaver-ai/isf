#ifndef __NC_ACS_DB_GROUP_SHARE_MOCK_H
#define __NC_ACS_DB_GROUP_SHARE_MOCK_H

#include <gmock/gmock.h>
#include <acsdb/public/ncIDBGroupShareManager.h>

class ncDBGroupShareManagerMock: public ncIDBGroupShareManager
{
    XPCOM_OBJECT_MOCK (ncDBGroupShareManagerMock);

public:
    MOCK_METHOD1(AddGroupShare, void(const dbGroupShareInfo&));
    MOCK_METHOD2(RenameGroupShare, void(const String&, const String&));
    MOCK_METHOD1(DeleteGroupShare, void(const String&));
    MOCK_METHOD2(GetGroupShareById, bool(const String&, dbGroupShareInfo&));
    MOCK_METHOD3(GetGroupShareByCreater, bool(const String&, const String&, dbGroupShareInfo&));
};

#endif // End __NC_ACS_DB_MOCK_H
