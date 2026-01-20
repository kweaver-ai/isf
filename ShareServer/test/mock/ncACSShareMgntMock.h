#ifndef __NC_ACS_SHAREMGNT_MOCK_H
#define __NC_ACS_SHAREMGNT_MOCK_H

#include <gmock/gmock.h>
#include <acssharemgnt/public/ncIACSShareMgnt.h>

class ncACSShareMgntMock: public ncIACSShareMgnt
{
    XPCOM_OBJECT_MOCK (ncACSShareMgntMock)

public:
    MOCK_METHOD2(GetNameByAccessorId, String(const String&, ncIOCAccesorType));
    MOCK_METHOD2(GetDirectBelongDepartmentIds, void(const String&, vector<String>&));
    MOCK_METHOD2(GetDirectBelongOrganIds, void(const String&, vector<String>&));
    MOCK_METHOD2(GetAllBelongDepartmentIds, void(const String &, vector<String> &));
    MOCK_METHOD1(IsAdminId, bool(const String&));
    MOCK_METHOD1(IsCustomDocManager, bool(const String&));
    MOCK_METHOD1(GetUserType, ncACSUserType(const String&));
    MOCK_METHOD2(GetUserSiteId, void(const String&, String&));
    MOCK_METHOD2(GetUserOSSId, void(const String&, String&));
    MOCK_METHOD2(GetDepartmentOSSId, void(const String&, String&));
    MOCK_METHOD2(GetAllGroups, void(const String&, vector<ncGroupInfo>&));
    MOCK_METHOD2(GetContactors, void(const String&, vector<ncACSUserInfo>&));
    MOCK_METHOD2(GetAllBelongGroups, void (const String&, vector<String>&));
    MOCK_METHOD2(GetUserDisplayName, void (const String&, String&));
    MOCK_METHOD3(GetUserName, void (const String&, String&, String&));
    MOCK_METHOD1(CheckDisplayNameIsExist, bool (const String&));
    MOCK_METHOD2(GetUserInfoById, bool (const String&, ncACSUserInfo&));
    MOCK_METHOD2(GetAppInfoById, bool (const String&, ncACSAppInfo&));
    MOCK_METHOD3(GetUserInfoByAccount, bool (const String&, ncACSUserInfo&, int&));
    MOCK_METHOD1(GetAccountType, int (const String&));
    MOCK_METHOD1(IsUserEnabled, bool (const String&));
    MOCK_METHOD2(GetUserInfoByThirdId, bool (const String&, ncACSUserInfo&));
    MOCK_METHOD2(GetDepartInfoById, bool (const String&, ncACSDepartInfo&));
    MOCK_METHOD2(GetUserInfoByIdBatch, void (const vector<String>&, map<String, ncACSUserInfo>&));
    MOCK_METHOD2(GetDepartInfoByIdBatch, void (const vector<String>&, map<String, ncACSDepartInfo>&));
    MOCK_METHOD2(GetGroupInfoByIdBatch, void (const vector<String>&, map<String, ncGroupInfo>&));
    MOCK_METHOD2(GetParentDepartIds, void (const String&, set<String>&));
    MOCK_METHOD2(GetSubDeps, void (const String&, vector<ncACSDepartInfo>&));
    MOCK_METHOD2(GetSubUsers, void (const String&, vector<ncACSUserInfo>&));
    MOCK_METHOD1(GetAllUser, void(vector<ncACSUserInfo> &));
    MOCK_METHOD6(SearchGivenOrganization, void (const vector<String>&, const String&, int, int, vector<ncACSUserInfo>&, vector<ncACSDepartInfo>&));
    MOCK_METHOD3(SearchGivenOrganizationCount, int (const vector<String>&, const String&, bool));
    MOCK_METHOD6(SearchAllOrganization, void (const String&, const String&, int, int, vector<ncACSUserInfo>&, vector<ncACSDepartInfo>&));
    MOCK_METHOD2(SearchAllOrganizationCount, int (const String&, const String&));
    MOCK_METHOD3(SearchDepartment, void (const String&, const String&, vector<ncACSUserInfo>&));
    MOCK_METHOD6(SearchContactGroup, void (const String&, const String&, int, int, vector<ncACSUserInfo>&, vector<ncGroupInfo>&));
    MOCK_METHOD3(SearchContactGroupCount, int (const String&, const String&, bool));
    MOCK_METHOD2(GetBelongDepartByIdBatch, void (const vector<String>&, map<String, ncACSDepartInfo>&));
    MOCK_METHOD3(GetBelongGroupByIdBatch, void (const String&, const vector<String>&, map<String, ncGroupInfo>&));
    MOCK_METHOD5(ExtLogin, String (const String&, const String&, const String&, const vector<String>&, int&));
    MOCK_METHOD2(GetUserCSFLevel, int (const String&, const ncVisitorType));
    MOCK_METHOD0(GetPermShareLimitStatus, bool ());
    MOCK_METHOD0(GetFindShareLimitStatus, bool ());
    MOCK_METHOD0(GetLinkShareLimitStatus, bool ());
    MOCK_METHOD2(GetOrgInfoByDeptId, void (const String&, vector<ncOrganizationInfo>&));
    MOCK_METHOD2(GetUserPermScopeInfos, void (const String&, map<String, ncPermScopeObjInfo>&));
    MOCK_METHOD2(GetOrgIdsByScopeInfo, void (const map<String, ncPermScopeObjInfo>&, vector<ncOrganizationInfo>&));
    MOCK_METHOD2(GetDeptRootPath, void (const String&, vector<String>&));
    MOCK_METHOD2(GetParentDeptRootPathName, void (const String&, String&));
    MOCK_METHOD2(GetUserRootPath, void (const String&, vector<vector<String> >&));
    MOCK_METHOD2(GetScopeOrgInfo, void (const String&, vector<pair<ncOrganizationInfo,bool> >&));
    MOCK_METHOD3(GetScopeSubDeps, void (const String&, const String&, vector<pair<ncACSDepartInfo,bool> >&));
    MOCK_METHOD3(GetScopeSubUsers, void (const String&, const String&, vector<ncACSUserInfo>&));
    MOCK_METHOD6(SearchScopeOrganization, void (const String&, const String&, int, int, vector<ncACSUserInfo>&, vector<ncACSDepartInfo>&));
    MOCK_METHOD3(SearchScopeOrganizationCount, int (const String&, const String&, bool));
    MOCK_METHOD2(CheckUsrInPermScope, bool (const String&, const String&));
    MOCK_METHOD2(CheckDeptInPermScope, bool (const String&, const String&));
    MOCK_METHOD1(IsUserFindEnabled, bool (const String&));
    MOCK_METHOD1(IsUserLinkEnabled, bool (const String&));
    MOCK_METHOD1(GetAllUserContactIds, void (map<String, vector<String> >&));
    MOCK_METHOD2(DeleteContactsByPatch, void (const String&, const vector<String>&));
    MOCK_METHOD3(GetUsrIdsOutOfPermScope, void (const String&, const vector<String>&, vector<String>&));
    MOCK_METHOD1(IsUndistirbutedUser, bool (const String&));
    MOCK_METHOD2(GetAllOrgInfo, void(const String&, vector<ncOrganizationInfo>&));
    MOCK_METHOD0(GetLeakProofStatus, bool());
    MOCK_METHOD1(GetLeakProofPerm, int(const String&));
    MOCK_METHOD0(GetClearCacheInterval, int ());
    MOCK_METHOD0(GetClearCacheSize, int64 ());
    MOCK_METHOD1(GetPwdControl, int (const String&));
    MOCK_METHOD1(GetUserAuthType, int (const String&));
    MOCK_METHOD0(GetClearClientCacheStatus, bool ());
    MOCK_METHOD0(GetHideClientCacheStatus, bool ());
    MOCK_METHOD0(GetMutiTenantStatus, bool ());
    MOCK_METHOD2(GetUserIdByDisplayName, void(const String&, vector<String>&));
    MOCK_METHOD2(GetOrgIdByUserId, void(const String&, vector<String>&));
    MOCK_METHOD2(GetManageDepIds, void(const String&, vector<String>&));
    MOCK_METHOD1(GetShareMgntConfig, String(const String&));
    MOCK_METHOD0(GetNetDocsLimitStatus, bool());
    MOCK_METHOD2(CheckNetDocLimit, bool(const String&, const String&));
    MOCK_METHOD2(FilterByNetDocLimit, void(set<String>&, const String&));
    MOCK_METHOD0(GetDefaulStrategySuperimStatus, bool());
    MOCK_METHOD1(UpdateUserLastRequestTime, void(const String&));
    MOCK_METHOD2(BatchGetConfig, void(vector<String>& keys, map<String, String>& kvMap));
    MOCK_METHOD2(GetCustomConfigOfString, bool (const String &, String &));
    MOCK_METHOD0(GetFreezeStatus, bool ());
    MOCK_METHOD1(IsUserFreeze, bool (const String&));
    MOCK_METHOD1(AgreedToTermsOfUse, void (const String&));
    MOCK_METHOD2(UpdateUserDocumentReadStatus, void (const String&, int));
    MOCK_METHOD4(GetMailAddress, void (const String& reveiverId, const String& tbName, const String& fieldName, String& email));
    MOCK_METHOD0(GetWaterMarkStrategy, bool ());
    MOCK_METHOD1(GetDownloadWatermarkDocs, void (map<String, int> &));
    MOCK_METHOD1(GetDownloadWatermarkDocTypes, void (map<int, int> &));
    MOCK_METHOD0(GetRealNameAuthStatus, bool ());
    MOCK_METHOD1(IsUserRealNameAuth, bool (const String&));
    MOCK_METHOD0(GetFileCrawlStatus, bool ());
    MOCK_METHOD2(IsFileCrawlStrategy, bool (const String&, const String&));
    MOCK_METHOD1(IsDepartmentExist, bool (const String&));
    MOCK_METHOD2(GetDeepestSubDepartIds, void (vector<String>&, const vector<String>&));
    MOCK_METHOD2(GetSubDepartIds, void (vector<String>&, const vector<String>&));
    MOCK_METHOD2(GetAllSubDepartIds, void (vector<String>&, const vector<String>&));
    MOCK_METHOD3(GetSubDepartIdsByLevel, void (vector<String>&, const vector<String>&, const int));
    MOCK_METHOD4(GetRelateDepartIds, void (vector<String>&, const vector<String>&, const bool, const int));
    MOCK_METHOD3(GetDocInfoOfDeparts, void (vector<ncDepDocInfo>&, const vector<String>&, const vector<dbOwnerInfo>&));
    MOCK_METHOD2(GetParentDeptPath, void (const String&, vector<String>&));
    MOCK_METHOD0(GetMessagePluginStatus, bool ());
    MOCK_METHOD2(GetAuditSupervisoryUserIds, void (const String& userId, vector<String>& userIds));
    MOCK_METHOD2(GetUserRoleIds, void (const String& userId, vector<String>& roleIds));
    MOCK_METHOD2(GetUserIdsByRoleId, void (const String& roleId, vector<String>& userIds));
    MOCK_METHOD2(GetUserRole, void (const String& docId, vector<ncRoleInfo>& roleInfos));
    MOCK_METHOD1(IsAdminRole, bool (const String& userId));
    MOCK_METHOD1(IsSecuritRole, bool (const String& userId));
    MOCK_METHOD2(GetUserInfoByTelNumber, bool (const String &, ncACSUserInfo &));
    MOCK_METHOD2(GetUserInfoByEmail, bool (const String &, ncACSUserInfo &));
    MOCK_METHOD3(OEM_GetConfigByOption, bool (const String &, const String &, String &));
    MOCK_METHOD1(BatchUpdateUserLastRequestTime, void (const vector<ncACSRefreshInfo> &));
    MOCK_METHOD2(GetManagerQuota, bool (const String &, ncManagerQuota &));
    MOCK_METHOD2(GetUserIdsBydepIds, void (const vector<String>&, vector<String>&));
    MOCK_METHOD1(UpdateUserActivateStatus, void(const String&));
    MOCK_METHOD2(GetAppOrgPerm, int (const String&, const ncAppPermOrgType&));
};

#endif // End __NC_ACS_SHAREMGNT_MOCK_H
