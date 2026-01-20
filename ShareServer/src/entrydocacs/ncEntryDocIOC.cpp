#include <abprec.h>
#include <ncutil/ncPerformanceProfilerPrec.h>

#include "entrydocacs.h"
#include "ncEntryDocIOC.h"

/* Implementation file */
NS_IMPL_THREADSAFE_ISUPPORTS1(ncEntryDocIOC, ncIEntryDocIOC)

ncEntryDocIOC::ncEntryDocIOC()
{
    NC_ENTRY_DOC_ACS_TRACE (_T("this: %p"), this);

    nsresult ret;
    _acsPermManager = do_CreateInstance (NC_ACS_PERM_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        String error;
        error.format (_T("Failed to create acs processor: 0x%x"), ret);
        THROW_ENTRY_DOC_ACS_ERROR (error, FAILED_TO_CREATE_ACS_PERM_MANAGER);
    }

    _acsShareMgnt = do_CreateInstance (NC_ACS_SHARE_MGNT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        String error;
        error.format (_T("Failed to create acs sharemgnt: 0x%x"), ret);
        THROW_ENTRY_DOC_ACS_ERROR (error, FAILED_TO_CREATE_ACS_SHAREMGNT);
    }

    _acsConfManager = do_CreateInstance (NC_ACS_CONF_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        String error;
        error.format (_T("Failed to create acs conf manager: 0x%x"), ret);
        THROW_ENTRY_DOC_ACS_ERROR (error, FAILED_TO_CREATE_ACS_CONF_MANANGER);
    }

    _acsOwnerManager = do_CreateInstance (NC_ACS_OWNER_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        String error;
        error.format (_T("Failed to create acs owner manager: 0x%x"), ret);
        THROW_ENTRY_DOC_ACS_ERROR (error, FAILED_TO_CREATE_ACS_OWNER_MANANGER);
    }

    _dbOwnerManager = do_CreateInstance (NC_DB_OWNER_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        String error;
        error.format (_T("Failed to create db owner manager: 0x%x"), ret);
        THROW_ENTRY_DOC_ACS_ERROR (error, FAILED_TO_CREATE_DB_OWNER_MANANGER);
    }
}

ncEntryDocIOC::~ncEntryDocIOC()
{
    NC_ENTRY_DOC_ACS_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void GetUserInfoByIdBatch (const vector<String> & userIds, map<String, ncACSUserInfo> & userInfoMap); */
NS_IMETHODIMP_(void) ncEntryDocIOC::GetUserInfoByIdBatch(const vector<String>& userIds, map<String, ncACSUserInfo>& userInfoMap)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    _acsShareMgnt->GetUserInfoByIdBatch(userIds, userInfoMap);
}

/* [notxpcom] void GetUserDisplayName (const String& userId, String& name);*/
NS_IMETHODIMP_(void) ncEntryDocIOC::GetUserDisplayName(const String& userId, String& name)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    _acsShareMgnt->GetUserDisplayName(userId, name);
}

/* [notxpcom] bool IsDownloadWatermarkDoc ([const] in StringRef docId, [const] in int docType, [const] in int64 size, [const] in StringRef path) */;
NS_IMETHODIMP_(bool) ncEntryDocIOC::IsDownloadWatermarkDoc (const String& docId, const int docType, const int64 size, const String& path)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    return _acsConfManager->IsDownloadWatermarkDoc(docId, docType, size, path);
}

/* [notxpcom] int CheckPermission ([const] in ncSubjectAttrRef subjectAttr, [const] in ncObjectAttrRef objectAttr, [const] in int permValue, in int checkFlags); */
NS_IMETHODIMP_(int) ncEntryDocIOC::CheckPermission(const ncSubjectAttr& subjectAttr, const ncObjectAttr& objectAttr, const int permValue, int checkFlags)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    ncOpsAttr opsAttr;
    opsAttr.permValue = permValue;

    return _acsPermManager->CheckPermission(subjectAttr, objectAttr, opsAttr, checkFlags);
}

/* [notxpcom] String GetNameByAccessorId ([const] in StringRef accessorId, [const] in int accessorType); */
NS_IMETHODIMP_(String) ncEntryDocIOC::GetNameByAccessorId(const String& accessorId, int accessorType)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    if (accessorType < IOC_USER || accessorType > IOC_USER_GROUP) {
        String error;
        error.format (_T("Invalid accessorType"));
        THROW_ENTRY_DOC_ACS_ERROR (error, INVALID_ACCESSOR_TYPE);
    }

    return _acsShareMgnt->GetNameByAccessorId(accessorId, ncIOCAccesorType(accessorType));
}

/* [notxpcom] void AddOwner ([const] in dbOwnerInfoRef ownerInfo); */
NS_IMETHODIMP_(void) ncEntryDocIOC::AddOwner(const dbOwnerInfo & ownerInfo)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    return _dbOwnerManager->AddOwner(ownerInfo);
}

/* [notxpcom] void GetRelateDepartIds(in StringVecRef relateDepartIds, [const] in StringVecRef tmpDepartIds, [const] in bool includeCurDepart, [const] in int level); */
NS_IMETHODIMP_(void) ncEntryDocIOC::GetRelateDepartIds(vector<String>& relateDepartIds, const vector<String>& tmpDepartIds, const bool includeCurDepart, const int level)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    _acsShareMgnt->GetRelateDepartIds(relateDepartIds, tmpDepartIds, includeCurDepart, level);
}

/* [notxpcom] void GetDocInfoOfDeparts(in ncDepDocInfoVecRef depDocInfos, [const] in StringVecRef departIds, [const] in dbOwnerInfoVecRef ownerInfoVec); */
NS_IMETHODIMP_(void) ncEntryDocIOC::GetDocInfoOfDeparts(vector<ncDepDocInfo>& depDocInfos, const vector<String>& departIds, const vector<dbOwnerInfo>& ownerInfoVec)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    _acsShareMgnt->GetDocInfoOfDeparts(depDocInfos, departIds, ownerInfoVec);
}

/* [notxpcom] bool CheckDisplayNameIsExist ([const] in StringRef displayName, in StringRef name);*/
NS_IMETHODIMP_(bool) ncEntryDocIOC::CheckDisplayNameIsExist(const String& displayName)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    return _acsShareMgnt->CheckDisplayNameIsExist (displayName);
}

/* [notxpcom] void GetParentDeptPath ([const] in StringRef deptId, in StringVecRef deptNames) */
NS_IMETHODIMP_(void) ncEntryDocIOC::GetParentDeptPath(const String& deptId, vector<String>& deptNames)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    _acsShareMgnt->GetParentDeptPath(deptId, deptNames);
}

/* [notxpcom] void SetAllOwner ([const] in StringRef docId, in dbOwnerInfoVecRef ownerInfos); */
NS_IMETHODIMP_(void) ncEntryDocIOC::SetAllOwner(const String & docId, vector<dbOwnerInfo> & ownerInfos)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    _acsOwnerManager->SetAllOwner (docId, ownerInfos);
}

/* [notxpcom] void AddCustomPermConfigsWithoutOwner ([const] in StringRef gnsPath, [const] in StringRef userId, [const] in ncCustomPermConfigVectorRef cpConfigs, in ncCustomPermConfigVectorRef addedConfigs); */
NS_IMETHODIMP_(void) ncEntryDocIOC::AddCustomPermConfigsWithoutOwner(const String & gnsPath, const vector<ncCustomPermConfig> & cpConfigs, vector<ncCustomPermConfig> & addedConfigs)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    _acsPermManager->AddCustomPermConfigsWithoutOwner (gnsPath, cpConfigs, addedConfigs);
}

/* [notxpcom] bool GetUserInfoById ([const] in StringRef userId, in ncACSUserInfoRef userInfo); */
NS_IMETHODIMP_(bool) ncEntryDocIOC::GetUserInfoById(const String & userId, ncACSUserInfo & userInfo)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    return _acsShareMgnt->GetUserInfoById(userId, userInfo);
}

/* [notxpcom] void GetMutiTenantStatus ();*/
NS_IMETHODIMP_(bool) ncEntryDocIOC::GetMutiTenantStatus()
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    return _acsShareMgnt->GetMutiTenantStatus ();
}

/* [notxpcom] bool IsAdminId ([const] in StringRef userId); */
NS_IMETHODIMP_(bool) ncEntryDocIOC::IsAdminId(const String & userId)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    return _acsShareMgnt->IsAdminId (userId);
}

/* [notxpcom] void GetOrgIdByUserId ([const] in StringRef userId, in StringVecRef orgIds);*/
NS_IMETHODIMP_(void) ncEntryDocIOC::GetOrgIdByUserId(const String & userId, vector<String>& orgIds)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    _acsShareMgnt->GetOrgIdByUserId (userId, orgIds);
}

/* [notxpcom] bool IsDepartmentExist([const] in StringRef departId); */
NS_IMETHODIMP_(bool) ncEntryDocIOC::IsDepartmentExist (const String& departId)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    return _acsShareMgnt->IsDepartmentExist (departId);
}

/* [notxpcom] void GetUserName ([const] in StringRef userId, in StringRef displayName, in StringRef account);*/
NS_IMETHODIMP_(void) ncEntryDocIOC::GetUserName(const String& userId, String& displayName, String& account)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    _acsShareMgnt->GetUserName (userId, displayName, account);
}

/* [notxpcom] bool CheckDisplayNameIsExist ([const] in StringRef displayName, in StringRef name);*/
NS_IMETHODIMP_(void) ncEntryDocIOC::GetUserIdByDisplayName(const String& displayName, vector<String>& userId)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    _acsShareMgnt->GetUserIdByDisplayName (displayName, userId);
}

/* [notxpcom] bool IsCustomDocManager ([const] in StringRef userId); */
NS_IMETHODIMP_(bool) ncEntryDocIOC::IsCustomDocManager(const String & userId)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    return _acsShareMgnt->IsCustomDocManager (userId);
}

/* [notxpcom] bool GetRealNameAuthStatus ();*/
NS_IMETHODIMP_(bool) ncEntryDocIOC::GetRealNameAuthStatus ()
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    return _acsShareMgnt->GetRealNameAuthStatus ();
}

/* [notxpcom] bool IsUserRealNameAuth ([const] in StringRef userId);*/
NS_IMETHODIMP_(bool) ncEntryDocIOC::IsUserRealNameAuth (const String& userId)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    return _acsShareMgnt->IsUserRealNameAuth (userId);
}

/* [notxpcom] bool IsUserEnabled ([const] in StringRef userId); */
NS_IMETHODIMP_(bool) ncEntryDocIOC::IsUserEnabled(const String & userId)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    return _acsShareMgnt->IsUserEnabled (userId);
}

/* [notxpcom] bool IsAdminRole ([const] in StringRef userId); */
NS_IMETHODIMP_(bool) ncEntryDocIOC::IsAdminRole(const String& userId)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    return _acsShareMgnt->IsAdminRole (userId);
}

/* [notxpcom] void DeleteCustomPermByDirId ([const] in StringRef dirId); */
NS_IMETHODIMP_(void) ncEntryDocIOC::DeleteCustomPermByDirId(const String& dirId)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    _acsPermManager->DeleteCustomPermByDirId (dirId);
}

/* [notxpcom] void DeleteByDirId ([const] in StringRef dirId); */
NS_IMETHODIMP_(void) ncEntryDocIOC::DeleteByDirId(const String& dirId)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);
}

/* [notxpcom] void DeleteOwnerByDirId ([const] in StringRef dirId); */
NS_IMETHODIMP_(void) ncEntryDocIOC::DeleteOwnerByDirId(const String & dirId)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    _dbOwnerManager->DeleteOwnerByDirId (dirId);
}

/* [notxpcom] void DeleteALLByUserId ([const] in StringRef userId, [const] in StringRef docId); */
NS_IMETHODIMP_(void) ncEntryDocIOC::DeleteALLByUserId (const String& userId, const String& docId)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);
}

NS_IMETHODIMP_(void) ncEntryDocIOC::GetManageDepIds (const String& userId, vector<String>& departIds)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    _acsShareMgnt->GetManageDepIds (userId, departIds);
}

/* [notxpcom] void GetSubUsers ([const] in StringRef depId, in ncACSUserInfoVecRef userInfos); */
NS_IMETHODIMP_(void) ncEntryDocIOC::GetSubUsers(const String& depId, vector<ncACSUserInfo>& userInfos)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    _acsShareMgnt->GetSubUsers (depId, userInfos);
}

/* [notxpcom] void GetOwnerInfos ([const] in StringRef docId, in dbOwnerInfoVecRef ownerInfos); */
NS_IMETHODIMP_(void) ncEntryDocIOC::GetOwnerInfos(const String & docId, vector<dbOwnerInfo> & ownerInfos)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    _dbOwnerManager->GetInheritOwnerInfosByDocId (docId, ownerInfos, false);
}

/* [notxpcom] bool IsSecuritRole ([const] in StringRef userId); */
NS_IMETHODIMP_(bool) ncEntryDocIOC::IsSecuritRole(const String& userId)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    return _acsShareMgnt->IsSecuritRole (userId);
}

/* [notxpcom] void GetUserRoleIds ([const] in StringRef userId, in StringVecRef roleIds); */
NS_IMETHODIMP_(void) ncEntryDocIOC::GetUserRoleIds(const String& userId, vector<String>& roleIds)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    _acsShareMgnt->GetUserRoleIds (userId, roleIds);
}

/* [notxpcom] void GetDepartInfoById ([const] in StringRef departId, in ncACSDepartInfoRef departInfo); */
NS_IMETHODIMP_(bool) ncEntryDocIOC::GetDepartInfoById(const String & departId, ncACSDepartInfo & departInfo)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    return _acsShareMgnt->GetDepartInfoById (departId, departInfo);
}

/* [notxpcom] String GetShareMgntConfig ([const] in StringRef key);*/
NS_IMETHODIMP_(String) ncEntryDocIOC::GetShareMgntConfig(const String &key)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    return _acsShareMgnt->GetShareMgntConfig (key);
}

/* [notxpcom] void ListEntryDocsWithLongPath ([const] in ncSubjectAttrRef subjectAttr, [const] in ncObjectAttrRef objectAttr, in StringNcAccessPermMapRef displayIdsMap); */
NS_IMETHODIMP_(void) ncEntryDocIOC::ListEntryDocsWithLongPath(const ncSubjectAttr & subjectAttr, const ncObjectAttr & objectAttr, map<String, ncAccessPerm> & displayIdsMap)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    _acsPermManager->ListEntryDocsWithLongPath (subjectAttr, objectAttr, displayIdsMap);
}

/* [notxpcom] bool GetManagerQuota ([const] in StringRef userId, in ncManagerQuotaRef managerQuota); */
NS_IMETHODIMP_(bool) ncEntryDocIOC::GetManagerQuota(const String & userId, ncManagerQuota & managerQuota)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    return _acsShareMgnt->GetManagerQuota (userId, managerQuota);
}

/* [notxpcom] String GetUserOSSId ([const] in StringRef userId); */
NS_IMETHODIMP_(String) ncEntryDocIOC::GetUserOSSId(const String & userId)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    String ossId;
    _acsShareMgnt->GetUserOSSId (userId, ossId);
    return ossId;
}

/* [notxpcom] String GetDepartmentOSSId ([const] in StringRef departmentId); */
NS_IMETHODIMP_(String) ncEntryDocIOC::GetDepartmentOSSId(const String & departmentId)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    String ossId;
    _acsShareMgnt->GetDepartmentOSSId (departmentId, ossId);
    return ossId;
}

/* [notxpcom] void GetUerIdsBydepIds ([const] in StringVecRef departIds, in StringVecRef userIds); */
NS_IMETHODIMP_(void) ncEntryDocIOC::GetUserIdsBydepIds(const vector<String>& departIds, vector<String>& userIds)
{
    NC_ENTRY_DOC_ACS_TRACE (_T("begin this: %p"), this);

    _acsShareMgnt->GetUserIdsBydepIds (departIds, userIds);
}
