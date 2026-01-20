#ifndef __NC_ACS_TOKEN_MANAGER_MOCK_H
#define __NC_ACS_TOKEN_MANAGER_MOCK_H

#include <gmock/gmock.h>
#include <acsprocessor/public/ncIACSTokenManager.h>

class ncACSTokenManagerMock: public ncIACSTokenManager
{
    XPCOM_OBJECT_MOCK (ncACSTokenManagerMock)

public:
    MOCK_METHOD1(DeleteTokenByUserId, void(const String&));
    MOCK_METHOD2(HasTokenByUDID, bool(const String&, const String&));
};

#endif // End __NC_ACS_TOKEN_MANAGER_MOCK_H
