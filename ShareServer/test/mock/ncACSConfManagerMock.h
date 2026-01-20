#ifndef __NC_ACS_CONF_MANAGER_MOCK_H
#define __NC_ACS_CONF_MANAGER_MOCK_H

#include <gmock/gmock.h>
#include <acsprocessor/public/ncIACSConfManager.h>

class ncACSConfManagerMock: public ncIACSConfManager
{
    XPCOM_OBJECT_MOCK (ncACSConfManagerMock);

public:
    MOCK_METHOD2(SetConfig, void(const String&, const String&));
    MOCK_METHOD1(GetConfig, String(const String&));
    MOCK_METHOD2(BatchGetConfig, void(vector<String>& keys, map<String, String>& kvMap));
    MOCK_METHOD4(IsDownloadWatermarkDoc, bool(const String&, const int, const int64, const String&));
    MOCK_METHOD0(GetNetDocsLimitStatus, bool ());
    MOCK_METHOD0(GetFileCrawlStatus, bool ());
    MOCK_METHOD0(GetMessagePluginStatus, bool ());
    MOCK_METHOD0(GetVcodeServerStatus, String ());
    MOCK_METHOD0(GetProductName, String ());
    MOCK_METHOD0(GetAuthVcodeServerStatus, String ());
    MOCK_METHOD0(GetPasswordErrCnt, int ());
};

#endif // End __NC_ACS_CONF_MANAGER_MOCK_H
