#ifndef __NC_DB_CONF_MANAGER_MOCK_H
#define __NC_DB_CONF_MANAGER_MOCK_H

#include <gmock/gmock.h>
#include <acsdb/public/ncIDBConfManager.h>

class ncDBConfManagerMock: public ncIDBConfManager
{
    XPCOM_OBJECT_MOCK (ncDBConfManagerMock);

public:
    MOCK_METHOD2(SetConfig, void(const String&, const String&));
    MOCK_METHOD1(GetConfig, String(const String&));
    MOCK_METHOD2(BatchGetConfig, void(vector<String>& keys, map<String, String>& kvMap));
};

#endif // End __NC_DB_CONF_MANAGER_MOCK_H
