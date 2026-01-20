#include <abprec.h>

#include "acsprocessor.h"
#include "ncACSTokenManager.h"

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (ncACSTokenManager, ncIACSTokenManager)

// protected
NS_IMETHODIMP_(nsrefcnt) ncACSTokenManager::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) ncACSTokenManager::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (ncACSTokenManager)

ncACSTokenManager::ncACSTokenManager()
    : _dbTokenManager ()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);

    nsresult ret;
    _dbTokenManager = do_CreateInstance (NC_DB_TOKEN_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, ACS_PROCESSOR_INTERNAL_ERR,
            _T("Failed to create db token manager: 0x%x"), ret);
    }

    _hydra = do_CreateInstance (HYDRA_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, ACS_PROCESSOR_INTERNAL_ERR,
            _T("Failed to create hydra adapter instance: 0x%x"), ret);
    }
}

ncACSTokenManager::ncACSTokenManager(ncIDBTokenManager *dbTokenManager, hydraInterface *hydra)
    : _dbTokenManager(dbTokenManager),
      _hydra(hydra)
{

}

ncACSTokenManager::~ncACSTokenManager()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void DeleteTokenByUserId ([const] in StringRef userId); */
NS_IMETHODIMP_(void) ncACSTokenManager::DeleteTokenByUserId(const String & userId)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p, userId: %s begin"),
        this, userId.getCStr ());

    _hydra->DeleteConsentAndLogin (String::EMPTY, userId);

    NC_ACS_PROCESSOR_TRACE (_T("this: %p, userId: %s end"),
        this, userId.getCStr ());
}

NS_IMETHODIMP_(bool) ncACSTokenManager::HasTokenByUDID(const String & userId, const String & udid)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p, userId: %s, udid: %s begin"), this, userId.getCStr (), udid.getCStr());

    if(userId.isEmpty () || udid.isEmpty ())
    {
        return false;
    }

    vector<ncTokenIntrospectInfo> tokenInfos;
    _hydra->GetConsentInfo(userId, tokenInfos);

    bool ret = false;
    for(size_t i = 0; i < tokenInfos.size(); ++i) {
        if(tokenInfos[i].udid == udid){
            ret = true;
            break;
        }
    }

    NC_ACS_PROCESSOR_TRACE (_T("this: %p, userId: %s, udid: %s end, ret: %s"), this, userId.getCStr (), udid.getCStr(), String::toString (ret).getCStr ());
    return ret;
}
