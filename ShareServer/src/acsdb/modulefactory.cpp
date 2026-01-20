#include <abprec.h>
#include <dboperatormanager/public/ncIDBOperatorManager.h>

#include "acsdb.h"
#include "ncDBOwnerManager.h"
#include "ncDBTokenManager.h"
#include "ncDBPermManager.h"
#include "ncDBLockManager.h"
#include "ncDBConfManager.h"
#include "ncDBDeviceManager.h"
#include "ncDBMessageManager.h"
#include "ncDBOutboxManager.h"
#include "common/util.h"

NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (ncDBOwnerManager, ncDBOwnerManager::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (ncDBTokenManager, ncDBTokenManager::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (ncDBPermManager, ncDBPermManager::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (ncDBLockManager, ncDBLockManager::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (ncDBConfManager, ncDBConfManager::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (ncDBDeviceManager, ncDBDeviceManager::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (ncDBMessageManager, ncDBMessageManager::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (ncDBOutboxManager, ncDBOutboxManager::getInstance)

// 定义要注册的组件的列表
static const nsModuleComponentInfo components[] = {
    {
        "ncDBOwnerManager",
        NC_DB_OWNER_MANAGER_CID,
        NC_DB_OWNER_MANAGER_CONTRACTID,
        ncDBOwnerManagerConstructor
    },
    {
        "ncDBTokenManager",
        NC_DB_TOKEN_MANAGER_CID,
        NC_DB_TOKEN_MANAGER_CONTRACTID,
        ncDBTokenManagerConstructor
    },
    {
        "ncDBPermManager",
        NC_DB_PERM_MANAGER_CID,
        NC_DB_PERM_MANAGER_CONTRACTID,
        ncDBPermManagerConstructor
    },
    {
        "ncDBLockManager",
        NC_DB_LOCK_MANAGER_CID,
        NC_DB_LOCK_MANAGER_CONTRACTID,
        ncDBLockManagerConstructor
    },
    {
        "ncDBConfManager",
        NC_DB_CONF_MANAGER_CID,
        NC_DB_CONF_MANAGER_CONTRACTID,
        ncDBConfManagerConstructor
    },
    {
        "ncDBDeviceManager",
        NC_DB_DEVICE_MANAGER_CID,
        NC_DB_DEVICE_MANAGER_CONTRACTID,
        ncDBDeviceManagerConstructor
    },
    {
        "ncDBMessageManager",
        NC_DB_MESSAGE_MANAGER_CID,
        NC_DB_MESSAGE_MANAGER_CONTRACTID,
        ncDBMessageManagerConstructor
    },
    {
        "ncDBOutboxManager",
        NC_DB_OUTBOX_MANAGER_CID,
        NC_DB_OUTBOX_MANAGER_CONTRACTID,
        ncDBOutboxManagerConstructor
    }
};

// 定义resLoader
IResourceLoader* ncACSDBResLoader = NULL;

/*
 * ncACSDBLibrary
 */
class ncACSDBLibrary : public ISharedLibrary
{
public:
    virtual void onInitLibrary (const AppSettings *appSettings,
                                const AppContext *appCtx)
    {
        try {
            ncACSDBResLoader = new MoResourceLoader (::getResourceFileName (ACS_DB,
                appSettings,
                appCtx,
                AB_RESOURCE_MO_EXT_NAME));
        }
        catch (Exception &e) {
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                            _T("Unabled to load resource acsprocessor.po. [Error: %s]"),
                                            e.toString ().getCStr ());
        }
    }

    virtual void onCloseLibrary (void) NO_THROW
    {
    }

    virtual void onInstall (const AppSettings* appSettings,
                            const AppContext* appCtx)
    {
    }

    virtual void onUninstall (void) NO_THROW
    {
    }

    virtual const tchar_t* getLibName (void) const
    {
        return ACS_DB;
    }

}; // class ncACSDBLibrary

AB_IMPL_NSGETMODULE_WITH_LIB(ACS_DB, components, ncACSDBLibrary)

ncIDBOperator* ncACSDBGetDBOperator (int timeout /* = 0*/)
{
    static bool _bInitPool = false;
    static ThreadMutexLock _initPoolLock;
    static nsCOMPtr<ncIDBOperatorPool> _pool;

    if (_bInitPool == false) {

        AutoLock<ThreadMutexLock> lock (&_initPoolLock);
        if (_bInitPool == false) {
            nsresult ret;
            nsCOMPtr<ncIDBOperatorPoolFactory> poolFactory = do_CreateInstance (NC_DB_OPERATOR_POOL_FACTORY_CONTRACTID, &ret);
            if (NS_FAILED (ret)) {
                THROW_E (ACS_DB, FAILED_TO_CREATE_DB_OPERATOR_POOL_FACTORY,
                    _T("Faild to create db operator pool factory: 0x%x"), ret);
            }
            String dbName = Util::getDBName(ANYSHARE_DB_NAME);
            _pool = getter_AddRefs (poolFactory->GetDBOperatorPool (dbName.getCStr(), CT_RW_NODE_A));

            _bInitPool = true;
        }
    }

    nsCOMPtr<ncIDBOperator> dbOper;
    if (timeout == 0) {
        dbOper = getter_AddRefs (_pool->GetDBOperator ());
    }
    else {
        dbOper = getter_AddRefs (_pool->GetDBOperatorWithTimeout (timeout));
    }
    NS_ADDREF (dbOper.get ());

    return dbOper;
}
