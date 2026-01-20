#ifndef __NC_ACS_TOKEN_MANAGER_H
#define __NC_ACS_TOKEN_MANAGER_H

#include <acsdb/public/ncIDBTokenManager.h>
#include "./public/ncIACSTokenManager.h"
#include <acssharemgnt/public/ncIACSShareMgnt.h>
#include <drivenadapter/public/hydraInterface.h>

/* Header file */
class ncACSTokenManager : public ncIACSTokenManager
{
    AB_DECLARE_THREADSAFE_SINGLETON (ncACSTokenManager)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_NCIACSTOKENMANAGER

    ncACSTokenManager();
    ncACSTokenManager (ncIDBTokenManager *dbTokenManager, hydraInterface *hydra);

private:
    ~ncACSTokenManager();

protected:
    nsCOMPtr<ncIDBTokenManager>        _dbTokenManager;
    nsCOMPtr<hydraInterface>           _hydra;

    // sharemgnt接口
    nsCOMPtr<ncIACSShareMgnt> _acsShareMgnt;

};

#endif // __NC_ACS_TOKEN_MANAGER_H
