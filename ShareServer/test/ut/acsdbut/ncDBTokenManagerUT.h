#ifndef __NC_TOKEN_MANAGER_UT_H
#define __NC_TOKEN_MANAGER_UT_H

#include <gtest/gtest.h>
#include <acsdb/public/ncIDBTokenManager.h>
#include <acsdb/ncDBTokenManager.h>
#include "acsdbut.h"

class ncDBTokenManagerTest: public ncDBTokenManager
{
protected:
    // 获取数据库连接，通过虚函数注入
    virtual ncIDBOperator* GetDBOperator ()
    {
        return ncDBUTApi::GetDBOperator ();
    }
};

class ncDBTokenManagerUT: public testing::Test
{
public:
    ncDBTokenManagerUT (void);
    ~ncDBTokenManagerUT (void);

    virtual void SetUp ();
    virtual void TearDown ();

protected:
    nsCOMPtr<ncIDBTokenManager>        _manager;
    nsCOMPtr<ncIDBOperator> _anyshareOperator;
};

#endif // __NC_TOKEN_MANAGER_UT_H
