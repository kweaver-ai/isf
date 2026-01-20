#ifndef __NC_EVFS_NAME_IOC_UT_H
#define __NC_EVFS_NAME_IOC_UT_H

#include <gtest/gtest.h>
#include <evfsioc/ncIEVFSNameIOC.h>

class ncEVFSNameIOCUT: public testing::Test
{
public:
    ncEVFSNameIOCUT (void);
    ~ncEVFSNameIOCUT (void);

    virtual void SetUp ();
    virtual void TearDown ();

protected:
    nsCOMPtr<ncIEVFSNameIOC>    _evfsNameIOC;
};

#endif // End __NC_EVFS_NAME_IOC_UT_H
