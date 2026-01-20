#ifndef __NC_DB_TOKEN_MANAGER_H
#define __NC_DB_TOKEN_MANAGER_H

#include <dboperator/public/ncIDBOperator.h>
#include <tedbcm/public/ncITEDBCManager.h>

#include "./public/ncIDBTokenManager.h"

/* Header file */
class ncDBTokenManager : public ncIDBTokenManager
{
    AB_DECLARE_THREADSAFE_SINGLETON (ncDBTokenManager)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_NCIDBTOKENMANAGER

    ncDBTokenManager();
    ~ncDBTokenManager();

protected:
    // 获取数据库连接并初始化表
    virtual ncIDBOperator* GetDBOperator ();
};

#endif // __NC_DB_TOKEN_MANAGER_H
