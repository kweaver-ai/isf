#ifndef __NC_ACS_CONF_MANAGER_H
#define __NC_ACS_CONF_MANAGER_H

#include <acsdb/public/ncIDBConfManager.h>
#include <acssharemgnt/public/ncIACSShareMgnt.h>
#include "./public/ncIACSConfManager.h"
#include "ncACSConfManager.h"

#ifdef __WINDOWS__
#define TO_ATOMIC(value) reinterpret_cast<LONG volatile*>(value)
#else
#define TO_ATOMIC(value) reinterpret_cast<gint32 volatile*>(value)
#endif

class ncACSConfManager;

class ncACSConfThread : public cpp::system::mt::Thread
{
public:
    ncACSConfThread (ncACSConfManager* configManager);
    ~ncACSConfThread ();

private:
    void run(void);

    bool                    _threadToRun;
    ncACSConfManager*       _configManager;
};


class ncACSConfManager : public ncIACSConfManager
{
    AB_DECLARE_THREADSAFE_SINGLETON (ncACSConfManager)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_NCIACSCONFMANAGER

    ncACSConfManager();
    ncACSConfManager(ncIDBConfManager* confMgr):_confManager(confMgr){}
    ~ncACSConfManager();

    void InitConfig ();
    void UpdateValue ();
    bool GetDocWatermarkStatus ();
    int GetDocWatermarkTypeByDocID (const String docID);
    int GetDocWatermarkTypeByDocType (const int docType);
    void InitWatermarkSupportType ();

protected:
    nsCOMPtr<ncIDBConfManager>      _confManager;
    nsCOMPtr<ncIACSShareMgnt>       _acsShareMgnt;

private:
    bool                            _initConfig;
    ncACSConfThread*                _configThread;
    ThreadMutexLock                 _initConfiglock;
    AtomicValue                     _enable_doc_watermark;                  // 水印配置策略开关
    AtomicValue                     _enableMessagePlugin;                   // 开启消息推送插件
    map<String, int>                _watermarkDocsMap;                      // 文档库水印配置
    ThreadMutexLock                 _watermarkDocsMapLock;
    map<int, int>                   _watermarkDocTypesMap;                  // 不同文档库类型的水印设置
    ThreadMutexLock                 _watermarkDocTypesMapLock;
    vector<String>                  _needwatermarkType;                     // 支持的水印文件类型
    bool                            _enableNetDocsLimit;                    // 网段文档库限制
    AtomicValue                     _fileCrawlStatus;                       // 文档抓取策略开关
    String                          _vcodeServerStatus;                     // 验证码发送服务器开关状态配置信息
    ThreadMutexLock                 _vcodeServerStatusLock;
    AtomicValue                     _passwordErrCnt;                        // 密码错误次数
    String                          _productName;                           // OEM中Product信息
    ThreadMutexLock                 _productNameLock;
    String                          _authVcodeServerStatus;                 //双因子验证码发送服务器开关状态配置信息
    ThreadMutexLock                 _authVcodeServerStatusLock;
};

#endif // __NC_ACS_CONF_MANAGER_H
