#include <abprec.h>

#include "acsprocessor.h"
#include "ncACSPermManager.h"
#include "ncACSTokenManager.h"
#include "ncACSOwnerManager.h"
#include "ncACSLockManager.h"
#include "ncACSLicenseManager.h"
#include "ncACSDeviceManager.h"
#include "ncACSConfManager.h"
#include "ncACSMessageManager.h"
#include "ncACSPolicyManager.h"
#include "ncACSOutboxManager.h"

NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (ncACSPermManager, ncACSPermManager::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (ncACSTokenManager, ncACSTokenManager::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (ncACSOwnerManager, ncACSOwnerManager::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (ncACSLockManager, ncACSLockManager::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (ncACSLicenseManager, ncACSLicenseManager::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (ncACSDeviceManager, ncACSDeviceManager::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (ncACSConfManager, ncACSConfManager::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (ncACSMessageManager, ncACSMessageManager::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (ncACSPolicyManager, ncACSPolicyManager::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (ncACSOutboxManager, ncACSOutboxManager::getInstance)

// 组件列表
static const nsModuleComponentInfo components[] = {
    {
        "ncACSPermManager",
        NC_ACS_PERM_MANAGER_CID,
        NC_ACS_PERM_MANAGER_CONTRACTID,
        ncACSPermManagerConstructor
    },
    {
        "ncACSTokenManager",
        NC_ACS_TOKEN_MANAGER_CID,
        NC_ACS_TOKEN_MANAGER_CONTRACTID,
        ncACSTokenManagerConstructor
    },
    {
        "ncACSOwnerManager",
        NC_ACS_OWNER_MANAGER_CID,
        NC_ACS_OWNER_MANAGER_CONTRACTID,
        ncACSOwnerManagerConstructor
    },
    {
        "ncACSLockManager",
        NC_ACS_LOCK_MANAGER_CID,
        NC_ACS_LOCK_MANAGER_CONTRACTID,
        ncACSLockManagerConstructor
    },
    {
        "ncACSLicenseManager",
        NC_ACS_LICENSE_MANAGER_CID,
        NC_ACS_LICENSE_MANAGER_CONTRACTID,
        ncACSLicenseManagerConstructor
    },
    {
        "ncACSDeviceManager",
        NC_ACS_DEVICE_MANAGER_CID,
        NC_ACS_DEVICE_MANAGER_CONTRACTID,
        ncACSDeviceManagerConstructor
    },
    {
        "ncACSConfManager",
        NC_ACS_CONF_MANAGER_CID,
        NC_ACS_CONF_MANAGER_CONTRACTID,
        ncACSConfManagerConstructor
    },
    {
        "ncACSMessageManager",
        NC_ACS_MESSAGE_MANAGER_CID,
        NC_ACS_MESSAGE_MANAGER_CONTRACTID,
        ncACSMessageManagerConstructor
    },
    {
        "ncACSPolicyManager",
        NC_ACS_POLICY_MANAGER_CID,
        NC_ACS_POLICY_MANAGER_CONTRACTID,
        ncACSPolicyManagerConstructor
    },
    {
        "ncACSOutboxManager",
        NC_ACS_OUTBOX_MANAGER_CID,
        NC_ACS_OUTBOX_MANAGER_CONTRACTID,
        ncACSOutboxManagerConstructor
    }
};

// 资源加载器
IResourceLoader* ncACSProcessorResLoader = NULL;

/*
 * ncACSProcessorLibrary
 */
class ncACSProcessorLibrary : public ISharedLibrary
{
public:
    virtual void onInitLibrary (const AppSettings *appSettings,
                                const AppContext *appCtx)
    {
        try {
            ncACSProcessorResLoader = new MoResourceLoader (::getResourceFileName (ACS_PROCESSOR,
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
        return ACS_PROCESSOR;
    }

}; // class ncACSProcessorLibrary

AB_IMPL_NSGETMODULE_WITH_LIB(ACS_PROCESSOR, components, ncACSProcessorLibrary)
