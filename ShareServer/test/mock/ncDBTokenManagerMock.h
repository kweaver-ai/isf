#ifndef __NC_ACS_DB_TOKEN_MOCK_H
#define __NC_ACS_DB_TOKEN_MOCK_H

#include <gmock/gmock.h>
#include <acsdb/public/ncIDBTokenManager.h>

class ncDBTokenManagerMock: public ncIDBTokenManager
{
    XPCOM_OBJECT_MOCK (ncDBTokenManagerMock);

public:
    MOCK_METHOD1(SaveActiveUser, void(map<String, String>&));
};

#endif // End __NC_ACS_DB_TOKEN_MOCK_H
