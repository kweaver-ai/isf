#include <abprec.h>
#include <gmock/gmock.h>
#include "ncDBTokenManagerUT.h"
#include <boost/date_time/posix_time/posix_time.hpp>

using namespace testing;

ncDBTokenManagerUT::ncDBTokenManagerUT ()
    : _anyshareOperator (0)
    , _manager (0)
{
    _anyshareOperator = ncDBUTApi::GetDBOperator();
}

ncDBTokenManagerUT::~ncDBTokenManagerUT ()
{
}

void ncDBTokenManagerUT::SetUp ()
{
    _manager = new ncDBTokenManagerTest;
}

void ncDBTokenManagerUT::TearDown ()
{

}

TEST_F (ncDBTokenManagerUT, DoCreateInstance)
{
    nsresult ret;
    nsCOMPtr<ncIDBTokenManager> manager = do_CreateInstance (NC_DB_TOKEN_MANAGER_CONTRACTID, &ret);

    if (NS_FAILED (ret)) {
        printMessage2 (_T("Failed to create instance: 0x%x"), ret);
        ASSERT_EQ (1, 0);
    }
}