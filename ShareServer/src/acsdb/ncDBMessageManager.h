#ifndef __NC_DB_MESSAGE_MANAGER_H
#define __NC_DB_MESSAGE_MANAGER_H

#include <acsdb/public/ncIDBMessageManager.h>

class ncIDBOperator;
/* Header file */
class ncDBMessageManager : public ncIDBMessageManager
{
    AB_DECLARE_THREADSAFE_SINGLETON (ncDBMessageManager)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_NCIDBMESSAGEMANAGER

    ncDBMessageManager();

private:
    ~ncDBMessageManager();

protected:
    // 获取数据库连接并初始化表
    virtual ncIDBOperator* GetDBOperator ();
};

#endif // __NC_DB_MESSAGE_MANAGER_H
