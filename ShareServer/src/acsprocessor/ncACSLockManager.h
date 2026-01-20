#ifndef __NC_ACS_LOCK_MANAGER_H
#define __NC_ACS_LOCK_MANAGER_H

#include <acsdb/public/ncIDBLockManager.h>
#include <acsdb/public/ncIDBConfManager.h>

#include <acsprocessor/public/ncIACSLockManager.h>

#include "ncACSProcessorUtil.h"

/* Header file */
class ncACSLockManager : public ncIACSLockManager
{
    AB_DECLARE_THREADSAFE_SINGLETON (ncACSLockManager)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_NCIACSLOCKMANAGER

    ncACSLockManager();

    //
    // ut使用
    //
    ncACSLockManager (ncIDBLockManager* dbLockManager,
                    ncIDBConfManager* dbConfManager)
                    : _dbLockManager (dbLockManager),
                      _dbConfManager (dbConfManager)
    {}

    ~ncACSLockManager();

protected:
    nsCOMPtr<ncIDBLockManager>      _dbLockManager;
    nsCOMPtr<ncIDBConfManager>      _dbConfManager;
};

#endif // __NC_ACS_LOCK_MANAGER_H
