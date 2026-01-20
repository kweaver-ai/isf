#ifndef __NC_ACS_TOKEN_MANAGER_UT_H
#define __NC_ACS_TOKEN_MANAGER_UT_H

#include <gtest/gtest.h>
#include <acssharemgnt/public/ncIACSShareMgnt.h>
#include "../../mock/hydraInterfaceMock.h"

class ncACSTokenManagerUT: public testing::Test
{
public:
    ncACSTokenManagerUT (void);
    ~ncACSTokenManagerUT (void);

    virtual void SetUp ();
    virtual void TearDown ();

protected:
    ncDBTokenManagerMock*                    _tokenMock;
    hydraMock*                              _hydraMock;
    ncACSShareMgntMock*                    _acsShareMgntMock;
    ncIACSTokenManager*                    _acsTokenManager;
};

#endif // End __NC_ACS_TOKEN_MANAGER_UT_H
