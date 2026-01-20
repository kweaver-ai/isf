#ifndef __NC_ACS_LICENSE_MANAGER_UT_H
#define __NC_ACS_LICENSE_MANAGER_UT_H

#include <gtest/gtest.h>

class ncACSLicenseManagerUT: public testing::Test
{
public:
    ncACSLicenseManagerUT (void);
    ~ncACSLicenseManagerUT (void);

    virtual void SetUp ();
    virtual void TearDown ();
};

#endif // End __NC_ACS_LICENSE_MANAGER_UT_H
