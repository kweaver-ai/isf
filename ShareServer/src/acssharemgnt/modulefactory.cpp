#include <abprec.h>
#include "common/util.h"
#include "acssharemgnt.h"
#include "ncACSShareMgnt.h"

NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (ncACSShareMgnt, ncACSShareMgnt::getInstance)

// 定义要注册的组件的列表
static const nsModuleComponentInfo components[] = {
    {
        "ncACSShareMgnt",
        NC_ACS_SHARE_MGNT_CID,
        NC_ACS_SHARE_MGNT_CONTRACTID,
        ncACSShareMgntConstructor
    },
};

IResourceLoader* ncACSShareMgntResLoader = NULL;

/*
 * ncACSShareMgntLibrary
 */
class ncACSShareMgntLibrary : public ISharedLibrary
{
public:
    virtual void onInitLibrary (const AppSettings *appSettings,
                                const AppContext *appCtx)
    {
        try {
            ncACSShareMgntResLoader = new MoResourceLoader (::getResourceFileName (ACS_SHAREMGNT,
                appSettings,
                appCtx,
                AB_RESOURCE_MO_EXT_NAME));
        }
        catch (Exception &e) {
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                            _T("Unabled to load resource acssharemgnt.po. [Error: %s]"),
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
        return ACS_SHAREMGNT;
    }

}; // class ncACSShareMgntLibrary

AB_IMPL_NSGETMODULE_WITH_LIB(ACS_SHAREMGNT, components, ncACSShareMgntLibrary)

ncIDBOperator* ncACSShareMgntGetDBOperator ()
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
                THROW_E (ACS_SHAREMGNT, FAILED_TO_CREATE_DB_OPERATOR_POOL_FACTORY,
                    _T("Faild to create db operator pool factory: 0x%x"), ret);
            }
            String dbName = Util::getDBName(SHAREMGNT_DB_NAME);
            _pool = getter_AddRefs (poolFactory->GetDBOperatorPool (dbName, CT_RW_NODE_A));

            _bInitPool = true;
        }
    }

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (_pool->GetDBOperator ());
    NS_ADDREF (dbOper.get ());

    return dbOper;
}
