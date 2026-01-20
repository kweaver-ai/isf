#include <abprec.h>
#include <dataapi/ncGNSUtil.h>

#include "acsprocessor.h"
#include "ncACSOwnerManager.h"

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (ncACSOwnerManager, ncIACSOwnerManager)

// protected
NS_IMETHODIMP_(nsrefcnt) ncACSOwnerManager::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) ncACSOwnerManager::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (ncACSOwnerManager)

ncACSOwnerManager::ncACSOwnerManager()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);

    nsresult ret;
    _dbOwnerManager = do_CreateInstance (NC_DB_OWNER_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_DB_OWNER_MANANGER,
            _T("Failed to create db owner manager: 0x%x"), ret);
    }

    _acsShareMgnt = do_CreateInstance (NC_ACS_SHARE_MGNT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E(ACS_PROCESSOR, FAILED_TO_CREATE_ACS_SHAREMGNT,
           _T("Failed to create sharemgnt : 0x%x"), ret);
    }

    _acsProcessorUtil = ncACSProcessorUtil::getInstance();

    _userManager = do_CreateInstance (USER_MANAGEMENT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_DRIVENADAPTER_MANANGER,
            _T("Failed to create usermanagement instance: 0x%x"), ret);
    }
}

ncACSOwnerManager::~ncACSOwnerManager()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);
}

void ncACSOwnerManager::getOwnerNameByIDs(const vector<String> &ids, map<String, String> &names)
{
    ncOrgNameIDInfo oNameIDInfo;
    ncOrgIDInfo oOrgIDInfo;
    oOrgIDInfo.vecUserIDs = ids;
    _userManager->GetOrgNameIDInfo (oOrgIDInfo, oNameIDInfo);
    names = oNameIDInfo.mapUserInfo;
}


/* [notxpcom] void SetAllOwner ([const] in StringRef docId, in dbOwnerInfoVecRef ownerInfos); */
NS_IMETHODIMP_(void) ncACSOwnerManager::SetAllOwner(const String & docId, vector<dbOwnerInfo> & ownerInfos)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p, docId: %s, ownerInfos size: %d begin"),
        this, docId.getCStr (), (int)ownerInfos.size ());

    // 删除掉旧的权限配置
    _dbOwnerManager->DeleteOwnerInfosByDocId (docId);

    for (size_t i = 0; i < ownerInfos.size (); ++i) {
        ownerInfos[i].docId = docId;
        _dbOwnerManager->AddOwner (ownerInfos[i]);
    }

    // 发送权限变更NSQ
    _acsProcessorUtil->SendPermChangeNSQ (docId);

    NC_ACS_PROCESSOR_TRACE (_T("this: %p, docId: %s, ownerInfos size: %d end"),
        this, docId.getCStr (), (int)ownerInfos.size ());
}

/* [notxpcom] bool IsOwner ([const] in StringRef docId, [const] in StringRef userId); */
NS_IMETHODIMP_(bool) ncACSOwnerManager::IsOwner(const String & docId, const String & userId)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p, docId: %s, userId: %s begin"),
        this, docId.getCStr (), userId.getCStr ());

    return _dbOwnerManager->IsOwner (docId, userId);
}

/* [notxpcom] void GetOwnerIds ([const] in StringRef docId, in VectorStringRef ownerIds); */
NS_IMETHODIMP_(void) ncACSOwnerManager::GetOwnerIds(const String & docId, vector<String> & ownerIds)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p, docId: %s begin"),
        this, docId.getCStr ());

    ownerIds.clear ();

    // 获取所有owner信息，包括继承下来的
    vector<dbOwnerInfo> dbOwnerInfos;
    _dbOwnerManager->GetInheritOwnerInfosByDocId (docId, dbOwnerInfos, true);

    set<String> tmpIds;
    for (size_t i = 0; i < dbOwnerInfos.size (); ++i) {
        tmpIds.insert (dbOwnerInfos[i].ownerId);
    }

    set<String>::iterator iter = tmpIds.begin ();
    for (; iter != tmpIds.end (); ++iter) {
        ownerIds.push_back (*iter);
    }

    NC_ACS_PROCESSOR_TRACE (_T("this: %p, docId: %s end, ret owernids size: %d"),
        this, docId.getCStr (), (int)ownerIds.size ());
}

/*[notxpcom] void GetOwnerDocIds ([const] in StringRef userId, in SetStringRef docIds);*/
NS_IMETHODIMP_(void) ncACSOwnerManager::GetOwnerDocIds (const String& userId, set<String>& docIds)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p, userId: %s begin"), this, userId.getCStr ());

    vector<dbOwnerInfo> infos;
    _dbOwnerManager->GetOwnerInfosByUserId (userId, infos);
    for (size_t i = 0; i < infos.size (); ++i) {
        docIds.insert (infos[i].docId);
    }

    NC_ACS_PROCESSOR_TRACE (_T("this: %p, userId: %s docIds: %d end"), this, userId.getCStr (), docIds.size ());
}

/* [notxpcom] void DeleteOwnerByFileId ([const] in StringRef fileId); */
NS_IMETHODIMP_(void) ncACSOwnerManager::DeleteOwnerByFileId(const String & fileId)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p, fileId: %s begin"),
        this, fileId.getCStr ());

    _dbOwnerManager->DeleteOwnerByFileId (fileId);

    // 发送权限变更NSQ
    _acsProcessorUtil->SendPermChangeNSQ (fileId);

    NC_ACS_PROCESSOR_TRACE (_T("this: %p, fileId: %s end"),
        this, fileId.getCStr ());
}

/* [notxpcom] void DeleteOwnerByDirId ([const] in StringRef dirId); */
NS_IMETHODIMP_(void) ncACSOwnerManager::DeleteOwnerByDirId(const String & dirId)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p, dirId: %s begin"),
        this, dirId.getCStr ());

    _dbOwnerManager->DeleteOwnerByDirId (dirId);

    // 发送权限变更NSQ
    _acsProcessorUtil->SendPermChangeNSQ (dirId);

    NC_ACS_PROCESSOR_TRACE (_T("this: %p, dirId: %s end"),
        this, dirId.getCStr ());
}

/* [notxpcom] void DeleteOwnerByUserId ([const] in StringRef userId); */
NS_IMETHODIMP_(void) ncACSOwnerManager::DeleteOwnerByUserId(const String & userId)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p, userId: %s begin"),
        this, userId.getCStr ());

    _dbOwnerManager->DeleteOwnerByUserId (userId);

    NC_ACS_PROCESSOR_TRACE (_T("this: %p, userId: %s end"),
        this, userId.getCStr ());
}
