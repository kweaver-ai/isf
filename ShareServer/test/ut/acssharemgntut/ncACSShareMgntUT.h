#ifndef __NC_ACS_SHAREMGNT_UT_H
#define __NC_ACS_SHAREMGNT_UT_H

#include <gtest/gtest.h>
#include <acssharemgnt/public/ncIACSShareMgnt.h>
#include <acssharemgnt/ncACSShareMgnt.h>

class ncACSShareMgntTested: public ncACSShareMgnt
{
protected:
    // 获取数据库连接
    virtual ncIDBOperator* GetDBOperator ();
};

class ncACSShareMgntUT: public testing::Test
{
public:
    ncACSShareMgntUT (void);
    ~ncACSShareMgntUT (void);

    virtual void SetUp ();
    virtual void TearDown ();

protected:
    nsCOMPtr<ncIACSShareMgnt>        _manager;
};

#endif // __NC_ACS_SHAREMGNT_UT_H
