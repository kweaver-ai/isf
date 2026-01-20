#ifndef __NC_DB_PERM_MANAGER_H
#define __NC_DB_PERM_MANAGER_H

#include <dboperator/public/ncIDBOperator.h>
#include <tedbcm/public/ncITEDBCManager.h>

#include "./public/ncIDBPermManager.h"

/* Header file */
class ncDBPermManager : public ncIDBPermManager
{
    AB_DECLARE_THREADSAFE_SINGLETON (ncDBPermManager)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_NCIDBPERMMANAGER

    ncDBPermManager();
    ~ncDBPermManager();

protected:
    // 获取数据库连接并初始化表
    virtual ncIDBOperator* GetDBOperator ();
};

#endif // __NC_DB_DOC_MANAGER_H
