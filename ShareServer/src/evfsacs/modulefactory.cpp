#include <abprec.h>

#include "evfsacs.h"
#include "ncEVFSAccessControlIOC.h"
#include "ncEVFSNameIOC.h"


NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (ncEVFSAccessControlIOC, ncEVFSAccessControlIOC::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (ncEVFSNameIOC, ncEVFSNameIOC::getInstance)

// 定义要注册的组件的列表
static const nsModuleComponentInfo components[] = {
    {
        "ncEVFSAccessControlIOC",
        NC_EVFS_ACCESS_CONTROL_IOC_CID,
        NC_EVFS_ACCESS_CONTROL_IOC_CONTRACTID,
        ncEVFSAccessControlIOCConstructor
    },
    {
        "ncEVFSNameIOC",
        NC_EVFS_NAME_IOC_CID,
        NC_EVFS_NAME_IOC_CONTRACTID,
        ncEVFSNameIOCConstructor
    }
};

IResourceLoader* ncEVFSACSResLoader = NULL;

/*
 * ncEVFSACSLibrary
 */
class ncEVFSACSLibrary : public ISharedLibrary
{
public:
    ncEVFSACSLibrary (void)
    {
    }

    virtual ~ncEVFSACSLibrary (void)
    {
        if (0 != ncEVFSACSResLoader) {
            delete ncEVFSACSResLoader;
            ncEVFSACSResLoader = 0;
        }
    }
public:
    virtual void onInitLibrary (const AppSettings *appSettings,
                                const AppContext *appCtx)
    {
        if (ncEVFSACSResLoader == 0) {
            try {
                ncEVFSACSResLoader = new MoResourceLoader (::getResourceFileName (EVFS_ACS,
                                                                                    appSettings,
                                                                                    appCtx,
                                                                                    AB_RESOURCE_MO_EXT_NAME));

                AppSettings* appSetting = AppSettings::getCFLAppSettings ();
                try {
                    appSetting->load ();
                }
                catch (Exception&) {
                    CFLASSERT (false);
                }
            }
            catch (Exception &e) {
                throw SharedLibraryException (e.getMessage (),
                                              e.getErrorId (),
                                              e.getErrorProvider ());
            }

        }
    }

    virtual void onCloseLibrary (void) NO_THROW
    {
        if (0 != ncEVFSACSResLoader) {
            delete ncEVFSACSResLoader;
            ncEVFSACSResLoader = 0;
        }
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
        return EVFS_ACS;
    }

}; // class ncEVFSACSLibrary

AB_IMPL_NSGETMODULE_WITH_LIB(EVFS_ACS, components, ncEVFSACSLibrary)
