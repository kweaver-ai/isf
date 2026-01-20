#ifndef __NC_EAC_CONFIG_HANDLER_H__
#define __NC_EAC_CONFIG_HANDLER_H__

#include <acssharemgnt/public/ncIACSShareMgnt.h>
#include "gen-cpp/ncTShareMgnt.h"
#include "gen-cpp/ShareMgnt_constants.h"

class ncEACConfigHandler
{
public:
    ncEACConfigHandler (ncIACSShareMgnt* acsShareMgnt);
    ~ncEACConfigHandler ();

    void doConfigRequestHandler (brpc::Controller* cntl);

protected:
    /***
     * 获取系统配置信息
     */
    void Get (brpc::Controller* cntl, const String& userId);

    /***
     * 获取OEM配置信息
     */
    void GetOEMConfigBySection (brpc::Controller* cntl, const String& fakeUserId);

    /***
    * 获取水印策略配置
    */
    void GetDocWatermarkConfig(brpc::Controller* cntl, const String& fakeUserId);

    /***
    * 获取文件抓取策略
    */
    void GetFileCrawlConfig (brpc::Controller* cntl, const String& userId);

    /***
    * 更新用户"快速入门"状态
    */
    void SetQuickStartStatus (brpc::Controller* cntl, const String& userId);

private:
    typedef void (ncEACConfigHandler::*ncMethodFunc) (brpc::Controller*, const String&);
    map<String, ncMethodFunc>            _methodFuncs;

private:
    nsCOMPtr<ncIACSShareMgnt>               _acsShareMgnt;      // sharemgnt管理
    set<string>                             _methodWhiteList;   // 不需要验证token的方法白名单
};

#endif  // __NC_EAC_CONFIG_HANDLER_H__
