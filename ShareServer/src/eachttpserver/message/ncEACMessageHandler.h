#ifndef __NC_EAC_MESSAGE_HANDLER_H__
#define __NC_EAC_MESSAGE_HANDLER_H__

#include <acsprocessor/public/ncIACSMessageManager.h>
#include <acsprocessor/public/ncIACSPermManager.h>
#include <acssharemgnt/public/ncIACSShareMgnt.h>

class ncEACMessageHandler
{
public:
    ncEACMessageHandler (ncIACSMessageManager* acsMessageManager, ncIACSPermManager* acsPermManager, ncIACSShareMgnt* acsShareMgnt);
    ~ncEACMessageHandler ();

    void doMessageRequestHandler (brpc::Controller* cntl);

    /***
     * 发送消息
     */
    void SendMessages (brpc::Controller* cntl);

    /***
     * 读取消息
     */
    void ReadMessage (brpc::Controller* cntl);

protected:
    /***
     * 获取未读消息
     */
    void Get (brpc::Controller* cntl, const String& userId);

    /***
     * 删除已读消息
     */
    void Read (brpc::Controller* cntl, const String& userId);

    /***
     * 设置消息已读
     */
    void Read2 (brpc::Controller* cntl, const String& userId);

    /***
     * 发送邮件接口
     */
    void SendMail (brpc::Controller* cntl, const String& userId);

    /***
     * 所有接收者读取消息
     */
    void ReadMessageForAllReceivers (brpc::Controller* cntl);

    /***
     * 指定接收者读取消息
     */
    void ReadMessageForSomeReceivers (brpc::Controller* cntl);

private:
    typedef void (ncEACMessageHandler::*ncMethodFunc) (brpc::Controller*, const String&);
    map<String, ncMethodFunc>            _methodFuncs;

private:
    nsCOMPtr<ncIACSMessageManager>          _acsMessageManager;      // 消息管理
    nsCOMPtr<ncIACSPermManager>             _acsPermManager;         // 权限管理
    nsCOMPtr<ncIACSShareMgnt>               _acsShareMgnt;           // acs sharemgnt

};

#endif  // __NC_EAC_MESSAGE_HANDLER_H__
