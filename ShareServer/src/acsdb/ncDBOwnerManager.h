#ifndef __NC_DB_OWNER_MANAGER_H
#define __NC_DB_OWNER_MANAGER_H

#include <dboperator/public/ncIDBOperator.h>
#include <tedbcm/public/ncITEDBCManager.h>

#include "./public/ncIDBOwnerManager.h"

/* Header file */
class ncDBOwnerManager : public ncIDBOwnerManager
{
    AB_DECLARE_THREADSAFE_SINGLETON (ncDBOwnerManager)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_NCIDBOWNERMANAGER

    ncDBOwnerManager();
    ~ncDBOwnerManager();

protected:
    // 获取数据库连接并初始化表
    virtual ncIDBOperator* GetDBOperator ();

    String GenerateGroupStr (const vector<String>& strs);
};

#endif // __NC_DB_OWNER_MANAGER_H
