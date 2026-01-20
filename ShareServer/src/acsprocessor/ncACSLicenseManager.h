#ifndef __NC_ACS_LICENSE_MANAGER_H
#define __NC_ACS_LICENSE_MANAGER_H

#include "./public/ncIACSLicenseManager.h"

/* Header file */
class ncACSLicenseManager : public ncIACSLicenseManager
{
    AB_DECLARE_THREADSAFE_SINGLETON (ncACSLicenseManager)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_NCIACSLICENSEMANAGER

    ncACSLicenseManager();

private:
    ~ncACSLicenseManager();

protected:
    ThreadMutexLock _mutex;
};

#endif // __NC_ACS_CAUTH_MANAGER_H
