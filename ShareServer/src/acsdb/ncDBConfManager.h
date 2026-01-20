#ifndef __NC_DB_CONF_MANAGER_H
#define __NC_DB_CONF_MANAGER_H

#include <dboperator/public/ncIDBOperator.h>
#include <tedbcm/public/ncITEDBCManager.h>

#include "./public/ncIDBConfManager.h"

/* Header file */
class ncDBConfManager : public ncIDBConfManager
{
    AB_DECLARE_THREADSAFE_SINGLETON (ncDBConfManager)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_NCIDBCONFMANAGER

    ncDBConfManager();
    ~ncDBConfManager();

protected:
    // 获取数据库连接并初始化表
    virtual ncIDBOperator* GetDBOperator ();

    // 初始化数据
    void Init (ncIDBOperator* dbOper);
};

#endif // __NC_DB_CONF_MANAGER_H
