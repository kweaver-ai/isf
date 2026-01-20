#ifndef __NC_DB_PERM_MANAGER_UT_H
#define __NC_DB_PERM_MANAGER_UT_H

#include <gtest/gtest.h>
#include <acsdb/public/ncIDBPermManager.h>
#include <acsdb/ncDBPermManager.h>
#include "acsdbut.h"

class ncDBPermManagerTest: public ncDBPermManager
{
protected:
    // 获取数据库连接，通过虚函数注入
    virtual ncIDBOperator* GetDBOperator ()
    {
        return ncDBUTApi::GetDBOperator ();
    }
};

class ncDBPermManagerUT: public testing::Test
{
public:
    ncDBPermManagerUT (void);
    ~ncDBPermManagerUT (void);

    virtual void SetUp ();
    virtual void TearDown ();

protected:
    nsCOMPtr<ncIDBPermManager>    _manager;
    nsCOMPtr<ncIDBOperator> _anyshareOperator;
};

#endif // __NC_EACP_DOC_OPERATION_UT_H__
