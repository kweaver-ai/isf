#ifndef __NC_EVFS_NAME_IOC_H
#define __NC_EVFS_NAME_IOC_H

#include <acssharemgnt/public/ncIACSShareMgnt.h>
#include <evfsioc/ncIEVFSNameIOC.h>

/* Header file */
class ncEVFSNameIOC : public ncIEVFSNameIOC
{
    AB_DECLARE_THREADSAFE_SINGLETON (ncEVFSNameIOC)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_NCIEVFSNAMEIOC

    ncEVFSNameIOC();

private:
    ~ncEVFSNameIOC();

protected:
    nsCOMPtr<ncIACSShareMgnt>        _acsShareMgnt;
};

#endif // __NC_EVFS_NAME_IOC_H
