#include <abprec.h>
#include <ncutil/ncBusinessDate.h>
#include <dataapi/ncGNSUtil.h>
#include <dataapi/ncJson.h>
#include <ncutil/ncPerformanceProfilerPrec.h>
#include <ehttpserver/ehttpserver.h>

#include "evfsacs.h"
#include "ncEVFSAccessControlIOC.h"

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (ncEVFSAccessControlIOC, ncIEVFSAccessControlIOC)

// protected
NS_IMETHODIMP_(nsrefcnt) ncEVFSAccessControlIOC::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) ncEVFSAccessControlIOC::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (ncEVFSAccessControlIOC)

ncEVFSAccessControlIOC::ncEVFSAccessControlIOC()
    : _beginHandlers (),
     _endHandlers (),
     _acsPermManager (NULL),
     _acsOwnerManager (NULL),
     _acsLockManager (NULL),
     _acsShareMgnt (NULL),
     _acsConfManager (NULL)
{
    NC_EVFS_ACS_TRACE (_T("this: %p"), this);

    // 创建 acsprocessor
    nsresult ret;
    _acsPermManager = do_CreateInstance (NC_ACS_PERM_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        String error;
        error.format (_T("Failed to create acs perm manager: 0x%x"), ret);

        THROW_EVFS_ACS_ERROR (error, FAILED_TO_CREATE_ACS_PERM_MANAGER);
    }

    _acsOwnerManager = do_CreateInstance (NC_ACS_OWNER_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        String error;
        error.format (_T("Failed to create acs owner manager: 0x%x"), ret);

        THROW_EVFS_ACS_ERROR (error, FAILED_TO_CREATE_ACS_OWNER_MANANGER);
    }

    _acsLockManager = do_CreateInstance (NC_ACS_LOCK_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        String error;
        error.format (_T("Failed to create acs lock manager: 0x%x"), ret);

        THROW_EVFS_ACS_ERROR (error, FAILED_TO_CREATE_ACS_LOCK_MANAGER);
    }

    _acsShareMgnt = do_CreateInstance (NC_ACS_SHARE_MGNT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        String error;
        error.format (_T("Failed to create acs sharemgnt manager: 0x%x"), ret);

        THROW_EVFS_ACS_ERROR (error, FAILED_TO_CREATE_ACS_SHAREMGNT_MANAGER);
    }

    _acsConfManager = do_CreateInstance (NC_ACS_CONF_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        String error;
        error.format (_T("Failed to create acs conf manager: 0x%x"), ret);

        THROW_EVFS_ACS_ERROR (error, FAILED_TO_CREATE_ACS_CONF_MANAGER);
    }

    _allPerms.push_back(ACS_AP_DISPLAY);
    _allPerms.push_back(ACS_AP_PREVIEW);
    _allPerms.push_back(ACS_AP_READ);
    _allPerms.push_back(ACS_AP_CREATE);
    _allPerms.push_back(ACS_AP_EDIT);
    _allPerms.push_back(ACS_AP_DELETE);

    _allowAttrMap[ACS_AP_DISPLAY] = ACS_ATTR_ALLOW_DISPLAY;
    _allowAttrMap[ACS_AP_PREVIEW] = ACS_ATTR_ALLOW_PREVIEW;
    _allowAttrMap[ACS_AP_READ] = ACS_ATTR_ALLOW_READ;
    _allowAttrMap[ACS_AP_CREATE] = ACS_ATTR_ALLOW_CREATE;
    _allowAttrMap[ACS_AP_EDIT] = ACS_ATTR_ALLOW_EDIT;
    _allowAttrMap[ACS_AP_DELETE] = ACS_ATTR_ALLOW_DELETE;

    _denyAttrMap[ACS_AP_DISPLAY] = ACS_ATTR_DENY_DISPLAY;
    _denyAttrMap[ACS_AP_PREVIEW] = ACS_ATTR_DENY_PREVIEW;
    _denyAttrMap[ACS_AP_READ] = ACS_ATTR_DENY_READ;
    _denyAttrMap[ACS_AP_CREATE] = ACS_ATTR_DENY_CREATE;
    _denyAttrMap[ACS_AP_EDIT] = ACS_ATTR_DENY_EDIT;
    _denyAttrMap[ACS_AP_DELETE] = ACS_ATTR_DENY_DELETE;

    initHandlers();
}

ncEVFSAccessControlIOC::ncEVFSAccessControlIOC(ncIACSPermManager* acsPermManager,
                                               ncIACSOwnerManager* acsOwnerManager,
                                               ncIACSLockManager* acsLockManager,
                                               ncIACSConfManager* acsConfManager,
                                               ncIACSShareMgnt* acsShareMgnt)
    : _beginHandlers (),
     _endHandlers (),
     _acsPermManager (acsPermManager),
     _acsOwnerManager (acsOwnerManager),
     _acsLockManager (acsLockManager),
     _acsShareMgnt (acsShareMgnt),
     _acsConfManager (acsConfManager)
{
    _allPerms.push_back(ACS_AP_DISPLAY);
    _allPerms.push_back(ACS_AP_PREVIEW);
    _allPerms.push_back(ACS_AP_READ);
    _allPerms.push_back(ACS_AP_CREATE);
    _allPerms.push_back(ACS_AP_EDIT);
    _allPerms.push_back(ACS_AP_DELETE);

    _allowAttrMap[ACS_AP_DISPLAY] = ACS_ATTR_ALLOW_DISPLAY;
    _allowAttrMap[ACS_AP_PREVIEW] = ACS_ATTR_ALLOW_PREVIEW;
    _allowAttrMap[ACS_AP_READ] = ACS_ATTR_ALLOW_READ;
    _allowAttrMap[ACS_AP_CREATE] = ACS_ATTR_ALLOW_CREATE;
    _allowAttrMap[ACS_AP_EDIT] = ACS_ATTR_ALLOW_EDIT;
    _allowAttrMap[ACS_AP_DELETE] = ACS_ATTR_ALLOW_DELETE;

    _denyAttrMap[ACS_AP_DISPLAY] = ACS_ATTR_DENY_DISPLAY;
    _denyAttrMap[ACS_AP_PREVIEW] = ACS_ATTR_DENY_PREVIEW;
    _denyAttrMap[ACS_AP_READ] = ACS_ATTR_DENY_READ;
    _denyAttrMap[ACS_AP_CREATE] = ACS_ATTR_DENY_CREATE;
    _denyAttrMap[ACS_AP_EDIT] = ACS_ATTR_DENY_EDIT;
    _denyAttrMap[ACS_AP_DELETE] = ACS_ATTR_DENY_DELETE;

    initHandlers();
}

ncEVFSAccessControlIOC::~ncEVFSAccessControlIOC()
{
    NC_EVFS_ACS_TRACE (_T("this: %p"), this);
}

void ncEVFSAccessControlIOC::initHandlers()
{
    // 文件操作开始时的处理
    _beginHandlers.insert (pair<ncEVFSAccessType, ncAccessControlFunc>(ncEVFSAccessType::LIST_FILE_VERSION, &ncEVFSAccessControlIOC::onListFileVersion));
    _beginHandlers.insert (pair<ncEVFSAccessType, ncAccessControlFunc>(ncEVFSAccessType::LIST_DIR, &ncEVFSAccessControlIOC::onListDir));
    _beginHandlers.insert (pair<ncEVFSAccessType, ncAccessControlFunc>(ncEVFSAccessType::RECYCLE_FILE, &ncEVFSAccessControlIOC::onRecycleFile));
    _beginHandlers.insert (pair<ncEVFSAccessType, ncAccessControlFunc>(ncEVFSAccessType::RECYCLE_DIR, &ncEVFSAccessControlIOC::onRecycleDir));
    _beginHandlers.insert (pair<ncEVFSAccessType, ncAccessControlFunc>(ncEVFSAccessType::COPY_FILE, &ncEVFSAccessControlIOC::onCopyFile));
    _beginHandlers.insert (pair<ncEVFSAccessType, ncAccessControlFunc>(ncEVFSAccessType::GET_FILE, &ncEVFSAccessControlIOC::onGetFile));
    _beginHandlers.insert (pair<ncEVFSAccessType, ncAccessControlFunc>(ncEVFSAccessType::LIST_RECYCLE_BIN_DIR, &ncEVFSAccessControlIOC::onListRecycleBinDir));
    _beginHandlers.insert (pair<ncEVFSAccessType, ncAccessControlFunc>(ncEVFSAccessType::SET_RECYCLE_POLICY, &ncEVFSAccessControlIOC::onSetRecyclePolicy));
    _beginHandlers.insert (pair<ncEVFSAccessType, ncAccessControlFunc>(ncEVFSAccessType::DELETE_FILE, &ncEVFSAccessControlIOC::onDeleteFile));
    _beginHandlers.insert (pair<ncEVFSAccessType, ncAccessControlFunc>(ncEVFSAccessType::DELETE_DIR, &ncEVFSAccessControlIOC::onDeleteDir));
    _beginHandlers.insert (pair<ncEVFSAccessType, ncAccessControlFunc>(ncEVFSAccessType::PREVIEW_FILE, &ncEVFSAccessControlIOC::onPreviewFile));
    _beginHandlers.insert (pair<ncEVFSAccessType, ncAccessControlFunc>(ncEVFSAccessType::GET_FILE_META, &ncEVFSAccessControlIOC::onGetFileMeta));
    _beginHandlers.insert (pair<ncEVFSAccessType, ncAccessControlFunc>(ncEVFSAccessType::GET_ATTR, &ncEVFSAccessControlIOC::onGetAttr));
    _beginHandlers.insert (pair<ncEVFSAccessType, ncAccessControlFunc>(ncEVFSAccessType::GET_DIR, &ncEVFSAccessControlIOC::onGetDir));
    _beginHandlers.insert (pair<ncEVFSAccessType, ncAccessControlFunc>(ncEVFSAccessType::SET_CSFLEVEL, &ncEVFSAccessControlIOC::onSetCSFLevel));
    _beginHandlers.insert (pair<ncEVFSAccessType, ncAccessControlFunc>(ncEVFSAccessType::SET_TAG, &ncEVFSAccessControlIOC::onSetTag));
    _beginHandlers.insert (pair<ncEVFSAccessType, ncAccessControlFunc>(ncEVFSAccessType::QUARANTINE_APPEAL, &ncEVFSAccessControlIOC::onQuarantineAppeal));
    _beginHandlers.insert (pair<ncEVFSAccessType, ncAccessControlFunc>(ncEVFSAccessType::SET_DOC_DUE, &ncEVFSAccessControlIOC::onSetDocDue));

    // 文件操作结束后的处理
    _endHandlers.insert (pair<ncEVFSAccessType, ncAccessFinishFunc>(ncEVFSAccessType::RECYCLE_FILE, &ncEVFSAccessControlIOC::onRecycleFileEnd));
    _endHandlers.insert (pair<ncEVFSAccessType, ncAccessFinishFunc>(ncEVFSAccessType::RECYCLE_DIR, &ncEVFSAccessControlIOC::onRecycleDirEnd));
    _endHandlers.insert (pair<ncEVFSAccessType, ncAccessFinishFunc>(ncEVFSAccessType::DELETE_FILE, &ncEVFSAccessControlIOC::onDeleteFileEnd));
    _endHandlers.insert (pair<ncEVFSAccessType, ncAccessFinishFunc>(ncEVFSAccessType::DELETE_DIR, &ncEVFSAccessControlIOC::onDeleteDirEnd));

}

/*  [notxpcom] void CheckPermission ([const] in ncACSubjectAtrrRef subjectAttr, [const] in ncACSObjectAtrrRef objAttr, in ncEVFSAccessType accessType); */
NS_IMETHODIMP_(void) ncEVFSAccessControlIOC::CheckPermission(const ncACSSubjectAttr & subjectAttr, const ncACSObjectAttr & objAttr, ncEVFSAccessType accessType)
{
    NC_ADD_PPN_CODE_CLIP_BEGIN_EX_EE (acs_CheckPermission, _T("acs_CheckPermission"));
    if (subjectAttr.userId.compare (NC_EVFS_NAME_IOC_DATAEXCHANGE_ID) == 0) {
        return;
    }
    map<ncEVFSAccessType, ncAccessControlFunc>::iterator iter = _beginHandlers.find (accessType);
    if (iter == _beginHandlers.end ()) {
        NC_EVFS_ACS_TRACE (_T("this: %p, gns: %s, userId: %s, accessType: %d, no handler"),
            this, objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr (), accessType);
        return;
    }

    // handle
    ncAccessControlFunc fun = iter->second;
    (this->*fun) (subjectAttr, objAttr);
    NC_ADD_PPN_CODE_CLIP_END_EX_EE (acs_CheckPermission);
}

/*  [notxpcom] bool HasEditPermission ([const] in ncACSubjectAtrrRef subjectAttr, [const] in ncACSObjectAtrrRef objAttr); */
NS_IMETHODIMP_(bool) ncEVFSAccessControlIOC::HasEditPermission(const ncACSSubjectAttr & subjectAttr, const ncACSObjectAttr & objAttr)
{
    NC_ADD_PPN_CODE_CLIP_BEGIN_EX_EE (acs_HasEditPermission, _T("acs_HasEditPermission"));
    if (subjectAttr.userId.compare (NC_EVFS_NAME_IOC_DATAEXCHANGE_ID) == 0) {
        return true;
    }
    ncCheckPermCode code = checkPermission (subjectAttr, objAttr, ACS_AP_EDIT);

    if (code != CHECK_OK) {
        return false;
    }

    NC_ADD_PPN_CODE_CLIP_END_EX_EE (acs_HasEditPermission);
    return true;
}

/* [notxpcom] void OnAccessFinished ([const] in StringRef gns, [const] in StringRef userID, in ncEVFSAccessType accessType); */
NS_IMETHODIMP_(void) ncEVFSAccessControlIOC::OnAccessFinished(const String & gns, const String & userID, ncEVFSAccessType accessType)
{
    NC_ADD_PPN_CODE_CLIP_BEGIN_EX_EE (acs_OnAccessFinished, _T("acs_OnAccessFinished"));
    map<ncEVFSAccessType, ncAccessFinishFunc>::iterator iter = _endHandlers.find (accessType);
    if (iter == _endHandlers.end ()) {
        NC_EVFS_ACS_TRACE (_T("this: %p, gns: %s, userId: %s, accessType: %d, no handler"),
            this, gns.getCStr (), userID.getCStr (), accessType);
        return;
    }

    // handle
    ncAccessFinishFunc fun = iter->second;
    (this->*fun) (gns, userID);
    NC_ADD_PPN_CODE_CLIP_END_EX_EE (acs_OnAccessFinished);
}

/*[notxpcom] void FetchOwnerDocIdSet ([const] in StringRef userId, out StringSet ownerDocIdSet);*/
NS_IMETHODIMP_(void) ncEVFSAccessControlIOC::FetchOwnerDocIdSet (const String& userId, set<String>* ownerDocIdSet)
{
    if (ownerDocIdSet == NULL) {
        return;
    }

    _acsOwnerManager->GetOwnerDocIds (userId, *ownerDocIdSet);
}

/* [notxpcom] bool CheckFinderEnabled ([const] in StringRef gns); */
NS_IMETHODIMP_(bool) ncEVFSAccessControlIOC::CheckFinderEnabled(const String & gns)
{
    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s begin"), this, gns.getCStr ());

    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s end"), this, gns.getCStr ());
    return false;
}


NS_IMETHODIMP_(int) ncEVFSAccessControlIOC::GetUserCSFLevel(const String & userId, const ACSVisitorType visitorType)
{
    NC_EVFS_ACS_TRACE (_T("this: %p, userId: %s begin"), this, userId.getCStr ());

    int level = _acsShareMgnt->GetUserCSFLevel(userId, static_cast<ncVisitorType>(visitorType));

    NC_EVFS_ACS_TRACE (_T("this: %p, userId: %s end"), this, userId.getCStr ());

    return level;
}

void ncEVFSAccessControlIOC::onListFileVersion (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr)
{

    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s, userId: %s begin"),
        this, objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr ());

    if (isFileCrawlOperation(subjectAttr, objAttr)) {
        NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s, userId: %s end, has file crawl perm."),
            this, objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr ());
        return;
    }

    checkPermHelper(subjectAttr, objAttr, ACS_AP_DISPLAY, __FUNCTION__);
    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s, userId: %s success"),
        this, objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr ());
}

void ncEVFSAccessControlIOC::onListDir (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr)
{
}

void ncEVFSAccessControlIOC::onRecycleFile (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr)
{
    checkPermHelper(subjectAttr, objAttr, ACS_AP_DELETE, __FUNCTION__);
    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s, userId: %s success"),
        this, objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr ());
}

void ncEVFSAccessControlIOC::onRecycleDir (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr)
{
    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s, userId: %s success"),
        this, objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr ());
}

void ncEVFSAccessControlIOC::onCopyFile (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr)
{
    checkPermHelper(subjectAttr, objAttr, ACS_AP_DISPLAY | ACS_AP_READ, __FUNCTION__);
    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s, userId: %s success"),
        this, objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr ());
}

void ncEVFSAccessControlIOC::onGetFile (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr)
{

    if (checkPermission(subjectAttr, objAttr, ACS_AP_READ) != CHECK_OK) {
        String error;
        error.format (_T("onGetFile (%s, %s) failed, need permission: %d"),
            objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr (), ACS_AP_READ);

        THROW_EVFS_ACS_ERROR (error, CHECK_PERMISSION_FAILED);
    }

    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s, userId: %s success"),
        this, objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr ());
}

void ncEVFSAccessControlIOC::onListRecycleBinDir (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr)
{
    checkIsOwner(subjectAttr, objAttr, __FUNCTION__);
    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s, userId: %s success"),
        this, objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr ());
}

void ncEVFSAccessControlIOC::onSetRecyclePolicy (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr)
{
    checkIsOwner(subjectAttr, objAttr, __FUNCTION__);
    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s, userId: %s success"),
        this, objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr ());
}

void ncEVFSAccessControlIOC::onDeleteFile (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr)
{
    checkIsOwner(subjectAttr, objAttr, __FUNCTION__);
    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s, userId: %s success"),
        this, objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr ());
}

void ncEVFSAccessControlIOC::onDeleteDir (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr)
{
    checkIsOwner(subjectAttr, objAttr, __FUNCTION__);
    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s, userId: %s success"),
        this, objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr ());
}

// 文件操作结束后的处理
void ncEVFSAccessControlIOC::onRecycleFileEnd (const String& gnsPath, const String& userId)
{
    _acsPermManager->DeleteCustomPermByFileId (gnsPath);
    _acsOwnerManager->DeleteOwnerByFileId (gnsPath);
    _acsLockManager->Delete (gnsPath);

    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s, userId: %s success"),
        this, gnsPath.getCStr (), userId.getCStr ());
}

void ncEVFSAccessControlIOC::onRecycleDirEnd (const String& gnsPath, const String& userId)
{
    _acsPermManager->DeleteCustomPermByDirId (gnsPath);
    _acsOwnerManager->DeleteOwnerByDirId (gnsPath);
    _acsLockManager->DeleteSubs (gnsPath);

    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s, userId: %s success"),
        this, gnsPath.getCStr (), userId.getCStr ());
}

void ncEVFSAccessControlIOC::onDeleteFileEnd (const String& gnsPath, const String& userId)
{
    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s, userId: %s success"),
        this, gnsPath.getCStr (), userId.getCStr ());
}

void ncEVFSAccessControlIOC::onDeleteDirEnd (const String& gnsPath, const String& userId)
{
    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s, userId: %s success"),
        this, gnsPath.getCStr (), userId.getCStr ());
}

void ncEVFSAccessControlIOC::onPreviewFile (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr)
{
    if (checkPermission(subjectAttr, objAttr, ACS_AP_READ) != CHECK_OK) {
        String error;
        error.format (_T("onPreviewFile (%s, %s) failed, need permission: %d"),
            objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr (), ACS_AP_READ);

        THROW_EVFS_ACS_ERROR (error, CHECK_PERMISSION_FAILED);
    }

    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s, userId: %s success"),
        this, objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr ());
}

void ncEVFSAccessControlIOC::onGetFileMeta (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr)
{
    if (checkPermission(subjectAttr, objAttr, ACS_AP_DISPLAY) != CHECK_OK) {
        String error;
        error.format (_T("onGetFileMeta (%s, %s) failed, need permission: %d"),
            objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr (), ACS_AP_DISPLAY);
        THROW_EVFS_ACS_ERROR (error, CHECK_PERMISSION_FAILED);
    }
    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s, userId: %s success"),
        this, objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr ());
}

void ncEVFSAccessControlIOC::onGetAttr (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr)
{
    checkPermHelper(subjectAttr, objAttr, ACS_AP_DISPLAY, __FUNCTION__);
    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s, userId: %s success"),
        this, objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr ());
}

void ncEVFSAccessControlIOC::onGetDir (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr)
{
    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s, userId: %s success"),
        this, objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr ());
}

void ncEVFSAccessControlIOC::onSetCSFLevel (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr)
{
    // 配置文件密级时，检查是否是所有者
    checkIsOwner(subjectAttr, objAttr, __FUNCTION__);
    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s, userId: %s success"),
        this, objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr ());
}

void ncEVFSAccessControlIOC::onSetTag (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr)
{
    checkPermHelper(subjectAttr, objAttr, ACS_AP_EDIT, __FUNCTION__);
    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s, userId: %s success"),
        this, objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr ());
}

void ncEVFSAccessControlIOC::onQuarantineAppeal (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr)
{
    checkIsOwner(subjectAttr, objAttr, __FUNCTION__);
    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s, userId: %s success"),
        this, objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr ());
}

void ncEVFSAccessControlIOC::onSetDocDue (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr)
{
    checkIsOwner(subjectAttr, objAttr, __FUNCTION__);
    NC_EVFS_ACS_TRACE (_T("this: %p, gnsPath: %s, userId: %s success"),
        this, objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr ());
}

/* [notxpcom] int64 GetLatestShareTime ([const] in StringRef userId, [const] in StringRef docId, [const] in ACSVisitorType visitorType, in SetStringRef subShareTimeSet); */
NS_IMETHODIMP_(int64) ncEVFSAccessControlIOC::GetLatestShareTime(const String & userId, const String & docId, const ACSVisitorType visitorType, set<String>& subShareTimeSet)
{
    return _acsPermManager->GetLatestShareTime(userId, docId, static_cast<ncVisitorType>(visitorType), subShareTimeSet);
}

void ncEVFSAccessControlIOC::checkPermHelper(const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr, int permValue, const String &who)
{
    // 检查用户权限
    ncCheckPermCode code = checkPermission (subjectAttr, objAttr, permValue);

    if (code != CHECK_OK) {
        String error;
        error.format (_T("%s (%s, %s) failed, need permission: %d"), who.getCStr(), objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr (), permValue);

        THROW_EVFS_ACS_ERROR (error, CHECK_PERMISSION_FAILED);
    }
}

ncCheckPermCode ncEVFSAccessControlIOC::checkPermission(const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr, int permValue)
{
    // 检查用户权限
    ncSubjectAttr ncSubAttr;
    ncSubAttr.userId = subjectAttr.userId;
    ncSubAttr.accessIp = subjectAttr.accessIp;
    ncSubAttr.visitorType = static_cast<ncVisitorType>(subjectAttr.visitorType);

    ncObjectAttr ncObjAttr;
    ncObjAttr.gnsPath = objAttr.gnsPath;

    ncOpsAttr ncOpsAttr;
    ncOpsAttr.permValue = permValue;

    ncCheckPermCode code = _acsPermManager->CheckPermission (ncSubAttr, ncObjAttr, ncOpsAttr, ncCheckFlags::NO_CHECK);
    return code;
}

void ncEVFSAccessControlIOC::checkIsOwner(const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr, const String &who)
{
    // 检查是否是所有者
    if (_acsOwnerManager->IsOwner (objAttr.gnsPath, subjectAttr.userId) == false) {
        String error;
        error.format (_T("%s (%s, %s) failed, need be owner"),
            who.getCStr(), objAttr.gnsPath.getCStr (), subjectAttr.userId.getCStr ());

        THROW_EVFS_ACS_ERROR (error, CHECK_PERMISSION_FAILED);
    }
}

bool ncEVFSAccessControlIOC::isFileCrawlOperation(const ncACSSubjectAttr & subjectAttr, const ncACSObjectAttr & objAttr)
{
    // 非匿名用户文档抓取策略是否开启，下列动作放权：新建、列举、编辑
    if (ACSVisitorType::REALNAME == subjectAttr.visitorType && _acsConfManager->GetFileCrawlStatus ()) {
        return _acsShareMgnt->IsFileCrawlStrategy (subjectAttr.userId, objAttr.gnsPath);
    }

    return false;
}
