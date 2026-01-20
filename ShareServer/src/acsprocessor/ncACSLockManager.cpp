#include <abprec.h>
#include <ncutil/ncBusinessDate.h>
#include <dataapi/ncGNSUtil.h>

#include "acsprocessor.h"
#include "ncACSLockManager.h"

#define AUTO_LOCK_CONFIG_KEY            "auto_lock"
#define AUTO_LOCK_EXPIRED_INTERVAL      "auto_lock_expired_interval"
#define DEFAULT_EXPIRED_INTERVAL        180

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (ncACSLockManager, ncIACSLockManager)

// protected
NS_IMETHODIMP_(nsrefcnt) ncACSLockManager::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) ncACSLockManager::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (ncACSLockManager)

ncACSLockManager::ncACSLockManager()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);

    nsresult ret;
    _dbLockManager = do_CreateInstance (NC_DB_LOCK_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_DB_LOCK_MANANGER,
            _T("Failed to create db perm manager: 0x%x"), ret);
    }

    _dbConfManager = do_CreateInstance (NC_DB_CONF_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_DB_CONF_MANANGER,
            _T("Failed to create db conf manager: 0x%x"), ret);
    }

}

ncACSLockManager::~ncACSLockManager()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);
}


NS_IMETHODIMP_(void) ncACSLockManager::SetAutolockConfig(bool isEnable, int64 expiredInterval)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p begin"), this);

    if (isEnable) {
        _dbConfManager->SetConfig(AUTO_LOCK_CONFIG_KEY, "true");
        _dbConfManager->SetConfig(AUTO_LOCK_EXPIRED_INTERVAL, String::toString (expiredInterval).getCStr());
    }
    else {
        _dbConfManager->SetConfig(AUTO_LOCK_CONFIG_KEY, "false");
    }

    NC_ACS_PROCESSOR_TRACE (_T("this: %p end"), this);
}

/* [notxpcom] bool IsAutolockEnabled (); */
NS_IMETHODIMP_(bool) ncACSLockManager::IsAutolockEnabled()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p begin"), this);

    String val = _dbConfManager->GetConfig(AUTO_LOCK_CONFIG_KEY);

    bool ret = false;

    if(val == "") {
        _dbConfManager->SetConfig(AUTO_LOCK_CONFIG_KEY, "true");
        ret = true;
    }
    else {
        ret = val.compareIgnoreCase("true") == 0;
    }

    NC_ACS_PROCESSOR_TRACE (_T("this: %p end"), this);

    return ret;
}

/* [notxpcom] int64 GetExpiredInterval (int64 & expiredInterval); */
NS_IMETHODIMP_(int64) ncACSLockManager::GetExpiredInterval()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p begin"), this);

    int64 expiredInterval;
    String tmpExpiredInterval;
    tmpExpiredInterval = _dbConfManager->GetConfig (AUTO_LOCK_EXPIRED_INTERVAL);
    if (!tmpExpiredInterval.isEmpty()) {
        expiredInterval = Int64::getValue(tmpExpiredInterval);
    }
    else {
        expiredInterval = DEFAULT_EXPIRED_INTERVAL;
        // 默认三分钟
        _dbConfManager->SetConfig(AUTO_LOCK_EXPIRED_INTERVAL, String::toString (expiredInterval).getCStr());
    }

    NC_ACS_PROCESSOR_TRACE (_T("this: %p end"), this);

    return expiredInterval;
}

/*[notxpcom] void Delete ([const] in StringRef fileId);*/
NS_IMETHODIMP_(void) ncACSLockManager::Delete (const String& fileId)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p, fileId: %s begin"), this, fileId.getCStr ());

    _dbLockManager->Delete (fileId);

    NC_ACS_PROCESSOR_TRACE (_T("his: %p, fileId: %s end"), this, fileId.getCStr ());
}

/* [notxpcom] void DeleteSubs ([const] in StringRef dirId); */
NS_IMETHODIMP_(void) ncACSLockManager::DeleteSubs(const String & dirId)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p, dirId: %s begin"),
        this, dirId.getCStr ());

    _dbLockManager->DeleteSubs(dirId);

    NC_ACS_PROCESSOR_TRACE (_T("this: %p, dirId: %s end"),
        this, dirId.getCStr ());
}

/* [notxpcom] void DeleteByUserId ([const] in StringRef userId); */
NS_IMETHODIMP_(void) ncACSLockManager::DeleteByUserId(const String & userId)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p, userId: %s begin"),
        this, userId.getCStr ());

    _dbLockManager->DeleteByUserId(userId);

    NC_ACS_PROCESSOR_TRACE (_T("this: %p, userId: %s end"),
        this, userId.getCStr ());
}
