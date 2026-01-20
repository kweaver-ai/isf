/***************************************************************************************************
modulefactory.cpp:
    Copyright (c) Eisoo Software, Inc.(2020), All rights reserved.

Purpose:
    注册组件

Author:
    Young.yu

Creating Time:
    2020-11-17
***************************************************************************************************/
#include <abprec.h>

#include "drivenadapter.h"
#include "userManagement.h"
#include "deploy.h"
#include "hydra.h"
#include "nsq.h"
#include "policyEngine.h"
#include "pluginMessage.h"
#include "stdLog.h"
#include "ossGateway.h"
#include "authentication.h"

NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (nsq, nsq::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (userManagement, userManagement::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (deploy, deploy::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (hydra, hydra::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (policyEngine, policyEngine::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (pluginMessage, pluginMessage::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (stdLog, stdLog::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (ossGateway, ossGateway::getInstance)
NS_GENERIC_FACTORY_SINGLETON_CONSTRUCTOR (authentication, authentication::getInstance)

// 定义要注册的组件的列表
static const nsModuleComponentInfo components[] = {
    {
        "userManagement",
        USER_MANAGEMENT_CID,
        USER_MANAGEMENT_CONTRACTID,
        userManagementConstructor
    },
    {
        "hydra",
        HYDRA_CID,
        HYDRA_CONTRACTID,
        hydraConstructor
    },
    {
        "deploy",
        DEPLOY_CID,
        DEPLOY_CONTRACTID,
        deployConstructor
    },
    {
        "nsq",
        NSQ_CID,
        NSQ_CONTRACTID,
        nsqConstructor
    },
    {
        "policyEngine",
        POLICY_ENGINE_CID,
        POLICY_ENGINE_CONTRACTID,
        policyEngineConstructor
    },
    {
        "pluginMessage",
        PLUGIN_MESSAGE_CID,
        PLUGIN_MESSAGE_CONTRACTID,
        pluginMessageConstructor
    },
    {
        "stdLog",
        STD_LOG_CID,
        STD_LOG_CONTRACTID,
        stdLogConstructor
    },
    {
        "ossGateway",
        OSS_GATEWAY_CID,
        OSS_GATEWAY_CONTRACTID,
        ossGatewayConstructor
    },
    {
        "authentication",
        AUTHENTICATION_CID,
        AUTHENTICATION_CONTRACTID,
        authenticationConstructor
    },
};

// 定义resLoader
IResourceLoader* ncDrivenAdapterLoader = NULL;


/*
 * ncDrivenAdapterLibrary
 */
class ncDrivenAdapterLibrary : public ISharedLibrary
{
public:
    virtual void onInitLibrary (const AppSettings *appSettings,
                                const AppContext *appCtx)
    {
        try {
            ncDrivenAdapterLoader = new MoResourceLoader (::getResourceFileName (ACS_DRIVEN_ADAPTER,
                appSettings,
                appCtx,
                AB_RESOURCE_MO_EXT_NAME));
        }
        catch (Exception &e) {
            SystemLog::getInstance ()->log (__FILE__, __LINE__, ERROR_LOG_TYPE,
                                            _T("Unabled to load resource drivenadapter.po. [Error: %s]"),
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
        return ACS_DRIVEN_ADAPTER;
    }

}; // class ncDrivenAdapterLibrary

AB_IMPL_NSGETMODULE_WITH_LIB(ACS_DRIVEN_ADAPTER, components, ncDrivenAdapterLibrary)
