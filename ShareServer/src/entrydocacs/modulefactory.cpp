#include <abprec.h>

#include "entrydocacs.h"
#include "ncEntryDocIOC.h"

NS_GENERIC_FACTORY_CONSTRUCTOR (ncEntryDocIOC)

// 定义要注册的组件的列表
static const nsModuleComponentInfo components[] = {
    {
        "ncEntryDocIOC",
        NC_ENTRY_DOC_IOC_CID,
        NC_ENTRY_DOC_IOC_CONTRACTID,
        ncEntryDocIOCConstructor
    },
};

/*
 * ncEntryDocACSLibrary
 */
class ncEntryDocACSLibrary : public ISharedLibrary
{
public:
    virtual void onInitLibrary (const AppSettings *appSettings,
                                const AppContext *appCtx)
    {
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
        return ENTRY_DOC_ACS;
    }

}; // class ncEFSPACSLibrary

AB_IMPL_NSGETMODULE_WITH_LIB(ENTRY_DOC_ACS, components, ncEntryDocACSLibrary)
