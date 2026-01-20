#ifndef __NC_DATA_STORE_MOCK_H
#define __NC_DATA_STORE_MOCK_H

#include <gmock/gmock.h>
#include <acsdatastore/public/ncIACSDataStore.h>

class ncACSDataStoreMock: public ncIACSDataStore
{
    XPCOM_OBJECT_MOCK (ncACSDataStoreMock);

public:
    MOCK_METHOD0 (CreateAccountId, String(void));
    MOCK_METHOD0 (GenerateObjectId, String(void));
    MOCK_METHOD0 (GetAvailableOSSID, String(void));
    MOCK_METHOD5 (InitUpload, void (const String&, const String&, const String&, const String&, ncUploadInfo&));
    MOCK_METHOD9 (UploadBlock, void(const String&, const String&, const String&, const String&, const String&, const String&, int64, int, ncUploadPartInfo&));
    MOCK_METHOD6 (CompleteUpload, void(const String&, const String&, const String&, const String&, const String&, const map<int, ncUploadPartInfo>&));
    MOCK_METHOD7 (ReadByOffest, void(const String&, const String&, const String&, int64, int, int64, string&));
};

#endif // End __NC_DATA_STORE_MOCK_H
