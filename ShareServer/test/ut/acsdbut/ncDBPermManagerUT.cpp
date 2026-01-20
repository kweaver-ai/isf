#include <abprec.h>
#include <gmock/gmock.h>
#include "ncDBPermManagerUT.h"

using namespace testing;

ncDBPermManagerUT::ncDBPermManagerUT ()
    : _anyshareOperator (0)
    , _manager (0)
{
}

ncDBPermManagerUT::~ncDBPermManagerUT ()
{
}

void ncDBPermManagerUT::SetUp ()
{
    _manager = new ncDBPermManagerTest;
}

void ncDBPermManagerUT::TearDown ()
{
}

TEST_F (ncDBPermManagerUT, do_CreateInstance)
{
    nsresult ret;
    nsCOMPtr<ncIDBPermManager> manager = do_CreateInstance (NC_DB_PERM_MANAGER_CONTRACTID, &ret);

    if (NS_FAILED (ret)) {
        printMessage2 (_T("Failed to create instance: 0x%x"), ret);
        ASSERT_EQ (1, 0);
    }
}

TEST_F (ncDBPermManagerUT, GetCustomPermByDocIds)
{
    NC_GENERATE_ANYSHARE_TB(_anyshareOperator, NC_TB_T_ACS_CUSTOM_PERM);
    ncDBUTApi::CleanTable (_anyshareOperator.get (), NC_TB_T_ACS_CUSTOM_PERM);

    vector<String> docIds;
    vector<dbCustomPermInfo> infos;

    ASSERT_NO_THROW (_manager->GetCustomPermByDocIds (docIds, infos));

    ASSERT_EQ (infos.size (), 0);
}

TEST_F (ncDBPermManagerUT, GetExpirePermInfos)
{
    NC_GENERATE_ANYSHARE_TB(_anyshareOperator, NC_TB_T_ACS_CUSTOM_PERM);
    ncDBUTApi::CleanTable (_anyshareOperator.get (), NC_TB_T_ACS_CUSTOM_PERM);

    dbCustomPermInfo permInfo;
    permInfo.isAllowed = false;
    permInfo.docId = _T("gns://D1111111111111111111111111111111/D2222222222222222222222222222222");
    permInfo.accessorId = _T("B6511AE0-9D55-54EF-61F9-222222222222");
    permInfo.permValue = 1;
    permInfo.endTime = -1;
    ASSERT_NO_THROW (_manager->AddCustomPerm (permInfo));

    permInfo.isAllowed = true;
    permInfo.docId = _T("gns://D1111111111111111111111111111111");
    permInfo.accessorId = _T("B6511AE0-9D55-54EF-61F9-222222222222");
    permInfo.permValue = 5;
    permInfo.endTime = 1000;
    ASSERT_NO_THROW (_manager->AddCustomPerm (permInfo));

    vector<dbCustomPermInfo> infos;
    ASSERT_NO_THROW (_manager->GetExpirePermInfos (10000, infos));
    ASSERT_EQ (infos.size (), 1);
    ASSERT_EQ (infos[0].isAllowed, permInfo.isAllowed);
    ASSERT_EQ (infos[0].docId, permInfo.docId);
    ASSERT_EQ (infos[0].accessorId, permInfo.accessorId);
    ASSERT_EQ (infos[0].permValue, permInfo.permValue);
}

TEST_F (ncDBPermManagerUT, GetPermConfigsByAccessToken)
{
    NC_GENERATE_ANYSHARE_TB(_anyshareOperator, NC_TB_T_ACS_CUSTOM_PERM);
    ncDBUTApi::CleanTable (_anyshareOperator.get (), NC_TB_T_ACS_CUSTOM_PERM);

    String docId = "gns://";
    set<String> accessorIds;
    vector<dbPermConfig> infos;

    // docId深度为0
    ASSERT_NO_THROW (_manager->GetPermConfigsByAccessToken (docId, accessorIds, infos, true, false));

    ASSERT_EQ (infos.size (), 0);

    // accessToken不为空
    accessorIds.insert (_T("B6511AE0-9D55-54EF-61F9-111111111111"));
    ASSERT_NO_THROW (_manager->GetPermConfigsByAccessToken (docId, accessorIds, infos, true, false));

    ASSERT_EQ (infos.size (), 0);

    // 添加测试数据
    dbPermConfig data1;
    data1.denyValue = 0;
    data1.allowValue = 1;
    data1.accessorType = 2;
    data1.accessorId = _T("B6511AE0-9D55-54EF-61F9-111111111111");
    data1.docId = _T("gns://D1111111111111111111111111111111");
    data1.endTime = 1410504466643428;
    ASSERT_NO_THROW (_manager->AddPermConfig (data1));

    dbPermConfig data2;
    data2.denyValue = 0;
    data2.allowValue = 1;
    data2.accessorType = 2;
    data2.accessorId = _T("B6511AE0-9D55-54EF-61F9-222222222222");
    data2.docId = _T("gns://D1111111111111111111111111111111/D2222222222222222222222222222222");
    data2.endTime = 1410504466643428;
    ASSERT_NO_THROW (_manager->AddPermConfig (data2));

    // 添加测试数据
    dbPermConfig data3;
    data3.denyValue = 0;
    data3.allowValue = 1;
    data3.accessorType = 2;
    data3.accessorId = _T("B6511AE0-9D55-54EF-61F9-111111111111");
    data3.docId = _T("gns://D1111111111111111111111111111111/D2222222222222222222222222222222");
    data3.endTime = 1410504466643428;
    ASSERT_NO_THROW (_manager->AddPermConfig (data3));

    dbPermConfig data4;
    data4.denyValue = 0;
    data4.allowValue = 1;
    data4.accessorType = 2;
    data4.accessorId = _T("B6511AE0-9D55-54EF-61F9-222222222222");
    data4.docId = _T("gns://D1111111111111111111111111111111");
    data4.endTime = 1410504466643428;
    ASSERT_NO_THROW (_manager->AddPermConfig (data4));

    // docId 一层，accessorId 两个
    docId = "gns://D1111111111111111111111111111111";

    accessorIds.clear ();
    accessorIds.insert (_T("B6511AE0-9D55-54EF-61F9-111111111111"));
    accessorIds.insert (_T("B6511AE0-9D55-54EF-61F9-222222222222"));

    ASSERT_NO_THROW (_manager->GetPermConfigsByAccessToken (docId, accessorIds, infos, true, false));

    ASSERT_EQ (infos.size (), 2);

    // docId 两层，accessorId 1个
    docId = "gns://D1111111111111111111111111111111/D2222222222222222222222222222222";

    accessorIds.clear ();
    accessorIds.insert (_T("B6511AE0-9D55-54EF-61F9-111111111111"));

    ASSERT_NO_THROW (_manager->GetPermConfigsByAccessToken (docId, accessorIds, infos, true, false));

    ASSERT_EQ (infos.size (), 2);

    // docId 两个，accessorId 两个，获取继承权限
    docId = "gns://D1111111111111111111111111111111/D2222222222222222222222222222222";

    accessorIds.clear ();
    accessorIds.insert (_T("B6511AE0-9D55-54EF-61F9-111111111111"));
    accessorIds.insert (_T("B6511AE0-9D55-54EF-61F9-222222222222"));

    ASSERT_NO_THROW (_manager->GetPermConfigsByAccessToken (docId, accessorIds, infos, true, false));

    ASSERT_EQ (infos.size (), 4);

    // docid 两个，accessorId 两个，仅获取自身权限
    docId = "gns://D1111111111111111111111111111111/D2222222222222222222222222222222";

    accessorIds.clear ();
    accessorIds.insert (_T("B6511AE0-9D55-54EF-61F9-111111111111"));
    accessorIds.insert (_T("B6511AE0-9D55-54EF-61F9-222222222222"));

    ASSERT_NO_THROW (_manager->GetPermConfigsByAccessToken (docId, accessorIds, infos, false, false));

    ASSERT_EQ (infos.size (), 2);
}
