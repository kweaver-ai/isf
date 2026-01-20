/***************************************************************************************************
ncACSOutboxManager.h:
    Copyright (c) Eisoo Software Inc. (2021), All rights reserved.

Purpose:
    acs outbox manager 接口

Author:
    Dylan.gao@aishu.cn

Creating Time:
    2021-06-18
***************************************************************************************************/
#ifndef __NC_ACS_OUTBOX_MANAGER_H
#define __NC_ACS_OUTBOX_MANAGER_H

#include <acsprocessor/public/ncIACSOutboxManager.h>
#include <acsdb/public/ncIDBOutboxManager.h>
#include <thread>

/* Header file */
class ncACSOutboxManager : public ncIACSOutboxManager
{
    AB_DECLARE_THREADSAFE_SINGLETON (ncACSOutboxManager)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_NCIACSOUTBOXMANAGER

    ncACSOutboxManager();

private:
    ~ncACSOutboxManager();
    void pushOutbox (void);

private:
    std::unique_ptr<std::thread>       _pushOutboxThread;
    nsCOMPtr<ncIDBOutboxManager>       _dbOutboxManager;
    Semaphore                          _notiSem;
    bool                               _needNotify;
    ThreadMutexLock                    _needNotifyLock;
};

#endif
