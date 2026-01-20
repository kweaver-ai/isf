#include <abprec.h>

#include "acsdatastore.h"
#include "ncACSDataStore.h"

NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (ncACSDataStore, ncACSDataStore::getInstance)

// 定义要注册的组件的列表
static const nsModuleComponentInfo components[] = {
    {
        "ncACSDataStore",
        NC_ACS_DATA_STORE_CID,
        NC_ACS_DATA_STORE_CONTRACTID,
        ncACSDataStoreConstructor
    }
};

NC_DEFINE_UMM_ALLOCATOR (acsDataStorePoolAllocator);

// 定义资源加载器
//IResourceLoader* ncACSDataStoreResLoader = NULL;

/*
 * ncACSDataStoreLibrary
 */
class ncACSDataStoreLibrary : public ISharedLibrary
{
public:
    virtual void onInitLibrary (const AppSettings *appSettings,
                                const AppContext *appCtx)
    {
        try {
            //ncACSDataStoreResLoader = new MoResourceLoader (::getResourceFileName (ACS_DATA_STORE,
            //    appSettings,
            //    appCtx,
            //    AB_RESOURCE_MO_EXT_NAME));

            NC_CREATE_UMM_ALLOACTOR (ncModulePool::getInstance (),
                                    acsDataStorePoolAllocator,
                                    UMM_ERROR);
        }
        catch (Exception &e) {
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                            _T("Unabled to load resource acsdatastore.po. [Error: %s]"),
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
        return ACS_DATA_STORE;
    }

}; // class ncACSDataStoreLibrary

AB_IMPL_NSGETMODULE_WITH_LIB(ACS_DATA_STORE, components, ncACSDataStoreLibrary)
