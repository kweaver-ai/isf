#ifndef __NC_EAC_DEVICE_HANDLER_H__
#define __NC_EAC_DEVICE_HANDLER_H__

#include <acsprocessor/public/ncIACSDeviceManager.h>

class ncEACDeviceHandler
{
public:
    ncEACDeviceHandler (ncIACSDeviceManager* deviceManager);
    ~ncEACDeviceHandler (void);

    void doDeviceRequestHandler (brpc::Controller* cntl);

protected:
    /***
     * 获取用户的设备信息
     */
    void onList (brpc::Controller* cntl, const String& userId);

    /***
     * 禁用设备
     */
    void onDisable (brpc::Controller* cntl, const String& userId);

    /***
     * 启用设备
     */
    void onEnable (brpc::Controller* cntl, const String& userId);

    /***
     * 擦除设备
     */
    void onErase (brpc::Controller* cntl, const String& userId);

    /***
     * 获取设备状态
     */
    void onGetStatus (brpc::Controller* cntl, const String& userId);

    /***
     * 收到擦除成功消息
     */
    void onEraseSuc (brpc::Controller* cntl, const String& userId);

private:
    typedef void (ncEACDeviceHandler::*ncMethodFunc) (brpc::Controller*, const String&);
    map<String, ncMethodFunc>            _methodFuncs;

private:
    nsCOMPtr<ncIACSDeviceManager>       _acsDeviceManager;  // device管理
};

#endif  // __NC_EAC_DEVICE_HANDLER_H__
