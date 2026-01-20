#ifndef __NC_DB_LOCK_MANAGER_H
#define __NC_DB_LOCK_MANAGER_H

#include <dboperator/public/ncIDBOperator.h>
#include "./public/ncIDBLockManager.h"

/* Header file */
class ncDBLockManager : public ncIDBLockManager
{
    AB_DECLARE_THREADSAFE_SINGLETON (ncDBLockManager)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_NCIDBLOCKMANAGER

    ncDBLockManager();
    ~ncDBLockManager();

protected:
    // 获取数据库连接并初始化表
    virtual ncIDBOperator* GetDBOperator ();
};

#endif // __NC_DB_LOCK_MANAGER_H
