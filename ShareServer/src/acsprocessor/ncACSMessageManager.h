#ifndef __NC_ACS_MESSAGE_MANAGER_H
#define __NC_ACS_MESSAGE_MANAGER_H

#include <acsdb/public/ncIDBMessageManager.h>
#include <acsdb/public/ncIDBConfManager.h>
#include <acsprocessor/public/ncIACSMessageManager.h>
#include <acsprocessor/public/ncIACSConfManager.h>
#include <drivenadapter/public/userManagementInterface.h>

class ncMessageThread2;
class ncPushMessageThread2;

/* Header file */
class ncACSMessageManager : public ncIACSMessageManager
{
    AB_DECLARE_THREADSAFE_SINGLETON (ncACSMessageManager)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_NCIACSMESSAGEMANAGER

    ncACSMessageManager();

private:
    ~ncACSMessageManager();

protected:
    nsCOMPtr<ncIDBMessageManager>       _dbMessageManager;
    nsCOMPtr<ncIDBConfManager>          _dbConfManager;
    nsCOMPtr<ncIACSConfManager>         _acsConfManager;
    nsCOMPtr<userManagementInterface>   _userManagement;

    ncMessageThread2*                   _messageThread2;
    ncPushMessageThread2*               _pushMessageThread2;
};

#endif
