#include <abprec.h>

#include "teacpserver.h"
#include "ncTEACPServer.h"

NS_GENERIC_FACTORY_CONSTRUCTOR (ncTEACPServer)

// 定义要注册的组件的列表
static const nsModuleComponentInfo components[] = {
    {
        "ncTEACPServer",
        NC_THRIFT_EACP_SERVER_CID,
        NC_THRIFT_EACP_SERVER_CONTRACTID,
        ncTEACPServerConstructor
    },
};

// 定义资源加载器
IResourceLoader* teacpserverResLoader = NULL;

// 模块名称
#define LIB_NAME _T("teacpserver")

/*
 * ncTEACPServerLibrary
 */
class ncTEACPServerLibrary : public ISharedLibrary
{
public:
    virtual void onInitLibrary (const AppSettings *appSettings,
                                const AppContext *appCtx)
    {
        try {
            teacpserverResLoader = new MoResourceLoader (::getResourceFileName (LIB_NAME,
                appSettings,
                appCtx,
                AB_RESOURCE_MO_EXT_NAME));
        }
        catch (Exception &e) {
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                            _T("Unabled to load resource teacpserver.po. [Error: %s]"),
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
        return LIB_NAME;
    }

}; // class ncTEACPServerLibrary

AB_IMPL_NSGETMODULE_WITH_LIB(LIB_NAME, components, ncTEACPServerLibrary)
