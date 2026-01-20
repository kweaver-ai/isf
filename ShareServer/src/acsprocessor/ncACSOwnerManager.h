#ifndef __NC_ACS_OWNER_MANAGER_H
#define __NC_ACS_OWNER_MANAGER_H

#include <acsdb/public/ncIDBOwnerManager.h>
#include <acssharemgnt/public/ncIACSShareMgnt.h>
#include <acsprocessor/public/ncIACSOwnerManager.h>
#include "drivenadapter/public/userManagementInterface.h"

#include "ncACSProcessorUtil.h"

/* Header file */
class ncACSOwnerManager : public ncIACSOwnerManager
{
    AB_DECLARE_THREADSAFE_SINGLETON (ncACSOwnerManager)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_NCIACSOWNERMANAGER

    ncACSOwnerManager();

    //
    // 做ut测试传桩对象使用。
    //
    ncACSOwnerManager(ncIACSProcessorUtil* acsProcessorUtil,
        ncIDBOwnerManager* dbOwnerManager,
        ncIACSShareMgnt* acsShareMgnt,
        userManagementInterface* userManagement)
        : _acsProcessorUtil (acsProcessorUtil),
          _dbOwnerManager (dbOwnerManager),
          _acsShareMgnt (acsShareMgnt),
          _userManager (userManagement)
    {}

private:
    ~ncACSOwnerManager();

    // 获取所有者名称信息
    void getOwnerNameByIDs(const vector<String> &ids, map<String, String> &names);

protected:
    ncIACSProcessorUtil*               _acsProcessorUtil;
    nsCOMPtr<ncIDBOwnerManager>        _dbOwnerManager;
    nsCOMPtr<ncIACSShareMgnt>          _acsShareMgnt;
    nsCOMPtr<userManagementInterface>   _userManager;
};

#endif // __NC_ACS_OWNER_MANAGER_H
