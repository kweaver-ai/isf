#include <abprec.h>
#include <license.h>

#include "ncACSLicenseManager.h"

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1(ncACSLicenseManager, ncIACSLicenseManager)

// protected
NS_IMETHODIMP_(nsrefcnt) ncACSLicenseManager::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) ncACSLicenseManager::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (ncACSLicenseManager)

ncACSLicenseManager::ncACSLicenseManager()
    : _mutex ()
{
    /* member initializers and constructor code */
}

ncACSLicenseManager::~ncACSLicenseManager()
{
    /* destructor code */
}

/* [notxpcom] int GetLicenseInfo ([const] in StringRef license, in StringRef info); */
NS_IMETHODIMP_(int) ncACSLicenseManager::GetLicenseInfo(const String & license, String & info)
{
    AutoLock<ThreadMutexLock> lock (&_mutex);
    string tmpInfo;
    // 由于授权码体系变更，旧的授权码不允许被激活
    // 先用ncGetLicenseInfo解析，如果不等于ok，再调用ncGetLicenseInfo_old解析
    int ret = ncGetLicenseInfo (toSTLString (license), tmpInfo);
    if (ret == LR_OK) {
        info = toCFLString (tmpInfo);
    }
    else {
        ret = ncGetLicenseInfo_Old (toSTLString (license), tmpInfo);
        if (ret == LR_OK) {
            info = toCFLString (tmpInfo);
        }
    }

    return ret;
}

/* [notxpcom] int VerifyActiveCode ([const] in StringRef license, [const] in StringRef machineCode, [const] in StringRef activeCode); */
NS_IMETHODIMP_(int) ncACSLicenseManager::VerifyActiveCode(const String & license, const String & machineCode, const String & activeCode)
{
    return ncVerifyActiveCode (toSTLString (license), toSTLString (machineCode), toSTLString (activeCode));
}
