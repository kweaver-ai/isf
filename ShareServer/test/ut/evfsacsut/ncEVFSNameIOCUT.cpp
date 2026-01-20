#include <abprec.h>
#include <gmock/gmock.h>

#include "ncEVFSNameIOCUT.h"

using namespace testing;

ncEVFSNameIOCUT::ncEVFSNameIOCUT ()
    : _evfsNameIOC (NULL)
{

}

ncEVFSNameIOCUT::~ncEVFSNameIOCUT ()
{

}

void ncEVFSNameIOCUT::SetUp ()
{
}

void ncEVFSNameIOCUT::TearDown ()
{
}

TEST_F (ncEVFSNameIOCUT, DoCreateInstance)
{
    nsresult ret;
    nsCOMPtr<ncIEVFSNameIOC> evsfUserNameIOC = do_CreateInstance (NC_EVFS_NAME_IOC_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        printMessage2 (_T("Failed to create evfs name ioc: 0x%x"), ret);
        ASSERT_EQ (0, 1);
    }
}
