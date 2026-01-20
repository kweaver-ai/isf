#include <abprec.h>
#include <gmock/gmock.h>
#include "ncDBConfManagerUT.h"

using namespace testing;

ncDBConfManagerUT::ncDBConfManagerUT ()
    : _anyshareOperator (0)
    , _manager (0)
{
}

ncDBConfManagerUT::~ncDBConfManagerUT ()
{
}

void ncDBConfManagerUT::SetUp ()
{
    _manager = new ncDBConfManagerTest;
}

void ncDBConfManagerUT::TearDown ()
{
}

TEST_F (ncDBConfManagerUT, DoCreateInstance)
{
    nsresult ret;
    nsCOMPtr<ncIDBConfManager> manager = do_CreateInstance (NC_DB_CONF_MANAGER_CONTRACTID, &ret);

    if (NS_FAILED (ret)) {
        printMessage2 (_T("Failed to create instance: 0x%x"), ret);
        ASSERT_EQ (1, 0);
    }
}

TEST_F (ncDBConfManagerUT, SetConfig)
{
    _manager->SetConfig("abc", "def");
    ASSERT_EQ(_manager->GetConfig("abc"), "def");

    _manager->SetConfig("abc", "def1");
    ASSERT_EQ(_manager->GetConfig("abc"), "def1");
}

TEST_F (ncDBConfManagerUT, GetConfig)
{
    ASSERT_EQ(_manager->GetConfig("non-exist"), "");
}
