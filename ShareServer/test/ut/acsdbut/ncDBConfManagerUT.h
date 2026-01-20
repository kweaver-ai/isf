#ifndef __NC_DB_CONF_MANAGER_UT_H
#define __NC_DB_CONF_MANAGER_UT_H

#include <gtest/gtest.h>
#include <acsdb/public/ncIDBConfManager.h>
#include <acsdb/ncDBConfManager.h>
#include "acsdbut.h"

class ncDBConfManagerTest: public ncDBConfManager
{
protected:
    // 获取数据库连接，通过虚函数注入
    virtual ncIDBOperator* GetDBOperator ()
    {
        return ncDBUTApi::GetDBOperator ();
    }
};

class ncDBConfManagerUT: public testing::Test
{
public:
    ncDBConfManagerUT (void);
    ~ncDBConfManagerUT (void);

    virtual void SetUp ();
    virtual void TearDown ();

protected:
    nsCOMPtr<ncIDBConfManager>          _manager;
    nsCOMPtr<ncIDBOperator>             _anyshareOperator;
};

#endif // __NC_DB_CONF_MANAGER_UT_H
