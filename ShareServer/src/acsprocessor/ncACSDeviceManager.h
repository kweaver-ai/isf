#ifndef __NC_ACS_DEVICE_MANAGER_H
#define __NC_ACS_DEVICE_MANAGER_H

#include <acsdb/public/ncIDBDeviceManager.h>
#include <acsdb/public/ncIDBTokenManager.h>
#include <acsprocessor/public/ncIACSTokenManager.h>
#include <acssharemgnt/public/ncIACSShareMgnt.h>
#include <drivenadapter/public/hydraInterface.h>
#include "./public/ncIACSDeviceManager.h"

class ncACSDeviceManager : public ncIACSDeviceManager
{
    AB_DECLARE_THREADSAFE_SINGLETON (ncACSDeviceManager)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_NCIACSDEVICEMANAGER

    ncACSDeviceManager();

private:
    ~ncACSDeviceManager();

    void checkDeviceInfo(const ncDeviceBaseInfo & baseInfo);

protected:
    nsCOMPtr<ncIDBDeviceManager>    _dbDeviceManager;
    nsCOMPtr<ncIDBTokenManager>     _dbTokenManager;
    nsCOMPtr<ncIACSShareMgnt>       _acsShareMgnt;
    nsCOMPtr<ncIACSTokenManager>    _acsTokenManager;
    nsCOMPtr<hydraInterface>        _hydra;
};

#endif // __NC_ACS_DEVICE_MANAGER_H
