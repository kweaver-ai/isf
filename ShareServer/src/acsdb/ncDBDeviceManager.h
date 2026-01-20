#ifndef __NC_DB_DEVICE_MANAGER_H
#define __NC_DB_DEVICE_MANAGER_H

#include <dboperator/public/ncIDBOperator.h>
#include <tedbcm/public/ncITEDBCManager.h>

#include "./public/ncIDBDeviceManager.h"

/* Header file */
class ncDBDeviceManager : public ncIDBDeviceManager
{
    AB_DECLARE_THREADSAFE_SINGLETON (ncDBDeviceManager)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_NCIDBDEVICEMANAGER

    ncDBDeviceManager();
    ~ncDBDeviceManager();

protected:
    // 获取数据库连接并初始化表
    virtual ncIDBOperator* GetDBOperator ();
};

#endif // __NC_DB_DEVICE_MANAGER_H
