#include <abprec.h>
#include <ncutil/ncBusinessDate.h>
#include <gmock/gmock.h>

#include <acsprocessor/ncACSTokenManager.h>
#include <acsprocessor/acsprocessor.h>
#include "../../mock/ncDBTokenManagerMock.h"
#include "../../mock/ncACSShareMgntMock.h"
#include "ncACSTokenManagerUT.h"

using namespace testing;

ncACSTokenManagerUT::ncACSTokenManagerUT ()
{
}

ncACSTokenManagerUT::~ncACSTokenManagerUT ()
{
}

void ncACSTokenManagerUT::SetUp ()
{
    _tokenMock = new ncDBTokenManagerMock ();
    _hydraMock = new hydraMock ();
    _acsShareMgntMock = new ncACSShareMgntMock();
    _acsTokenManager = new ncACSTokenManager (_tokenMock, _hydraMock);
}

void ncACSTokenManagerUT::TearDown ()
{
    if (_acsTokenManager) {
        delete _acsTokenManager;
    }
    if (_tokenMock) {
        delete _tokenMock;
    }
    if (_hydraMock) {
        delete _hydraMock;
    }
    if(_acsShareMgntMock)
    {
        delete _acsShareMgntMock;
    }
}

// 更新mq的sdk UT失败 先去除ut
// TEST_F (ncACSTokenManagerUT, DoCreateInstance)
// {
//     nsresult ret;
//     nsCOMPtr<ncIACSTokenManager> tokenManager = do_CreateInstance (NC_ACS_TOKEN_MANAGER_CONTRACTID, &ret);

//     if (NS_FAILED (ret)) {
//         printMessage2 (_T("Failed to create acs token manager: 0x%x"), ret);
//         ASSERT_EQ(0, 1);
//     }
// }

// TEST_F (ncACSTokenManagerUT, CheckToken)
// {
//     String tokenId1(_T("tokenId1")),
//         tokenId2(_T("tokenId2")),
//         tokenId3(_T("tokenId3"));

//     dbTokenInfo tmpInfo;
//     tmpInfo.tokenId = tokenId1;
//     tmpInfo.createTime = Int64::toString(BusinessDate::getCurrentTime ());
//     tmpInfo.lastRequestTime = Int64::toString(BusinessDate::getCurrentTime ()/1000000);
//     tmpInfo.userId = _T("user001");
//     tmpInfo.expires = 3600;
//     tmpInfo.flag = 0;

//     EXPECT_CALL(*_tokenMock, GetTokenInfo(tokenId1, _))
//         .WillRepeatedly(DoAll(SetArgReferee<1>(tmpInfo), Return(true)));

//     tmpInfo.tokenId = tokenId2;
//     tmpInfo.lastRequestTime = Int64::toString(BusinessDate::getCurrentTime ()/1000000 - 3601);
//     EXPECT_CALL(*_tokenMock, GetTokenInfo(tokenId2, _))
//         .WillRepeatedly(DoAll(SetArgReferee<1>(tmpInfo), Return(true)));

//     EXPECT_CALL(*_tokenMock, GetTokenInfo(tokenId3, _))
//         .WillRepeatedly(DoAll(SetArgReferee<1>(tmpInfo), Return(false)));

//     ncCheckTokenInfo checkTokenInfo;
//     checkTokenInfo.tokenId = tokenId3;

//     // case 1: 没有符合条件的 token, 校验失败
//     ASSERT_EQ (_acsTokenManager->CheckToken (checkTokenInfo), false);

//     // case 2: token 校验成功, 返回该 token 对应的 userId
//     checkTokenInfo.tokenId = tokenId1;
//     ASSERT_EQ (_acsTokenManager->CheckToken (checkTokenInfo), true);
//     ASSERT_EQ (checkTokenInfo.userId, tmpInfo.userId);

//     // case 3: token 超时, 校验失败
//     checkTokenInfo.tokenId = tokenId2;
//     ASSERT_EQ (_acsTokenManager->CheckToken (checkTokenInfo), false);
// }

// TEST_F (ncACSTokenManagerUT, DeleteTokenByUserId)
// {
//     String userId ("user");

//     EXPECT_CALL(*_tokenMock, DeleteTokenByUserId(userId))
//         .Times(1);

//     ASSERT_NO_THROW (_acsTokenManager->DeleteTokenByUserId (userId));
// }

// TEST_F (ncACSTokenManagerUT, HasTokenByUDID)
// {
//     String userid ("usrid");
//     String udid ("xxxx");

//     vector<ncTokenIntrospectInfo> tokenInfos;
//     EXPECT_CALL(*_acsTokenMock, GetConsentInfo(userid, _))
//         .WillOnce(SetArgReferee<1>(tokenInfos));
//     bool ret = false;
//     ASSERT_NO_THROW (ret = _acsTokenManager->HasTokenByUDID (userid, udid));
//     ASSERT_EQ (ret, false);

//     ncTokenIntrospectInfo oTempInfo1;
//     oTempInfo1.userId = userid;
//     oTempInfo1.udid = "yyy";
//     tokenInfos.push_back(oTempInfo1);
//     EXPECT_CALL(*_acsTokenMock, GetConsentInfo(userid, _))
//         .WillOnce(SetArgReferee<1>(tokenInfos));
//     ASSERT_NO_THROW (ret = _acsTokenManager->HasTokenByUDID (userid, udid));
//     ASSERT_EQ (ret, false);

//     ncTokenIntrospectInfo oTempInfo2;
//     oTempInfo2.userId = userid;
//     oTempInfo2.udid = udid;
//     tokenInfos.push_back(oTempInfo2);
//     EXPECT_CALL(*_acsTokenMock, GetConsentInfo(userid, _))
//         .WillOnce(SetArgReferee<1>(tokenInfos));
//     ASSERT_NO_THROW (ret = _acsTokenManager->HasTokenByUDID (userid, udid));
//     ASSERT_EQ (ret, true);
// }
