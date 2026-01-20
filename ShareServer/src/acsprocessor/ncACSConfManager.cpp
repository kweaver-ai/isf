#include <abprec.h>
#include <dataapi/ncGNSUtil.h>
#include <dataapi/ncJson.h>
#include <ethriftutil/ncThriftClient.h>

#include "gen-cpp/ncTShareMgnt.h"
#include "gen-cpp/ShareMgnt_constants.h"
#include "acsprocessor.h"
#include "ncACSConfManager.h"
#include "acsServiceAccessConfig.h"

#define NC_WATERMARK_CONFIG_FILENAME              _T("/sysvol/apphome/lib64/evfs_watermark.config")

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (ncACSConfManager, ncACSConfManager)

// protected
NS_IMETHODIMP_(nsrefcnt) ncACSConfManager::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) ncACSConfManager::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (ncACSConfManager)

ncACSConfManager::ncACSConfManager()
    : _initConfig (false)
    , _configThread()
    , _enable_doc_watermark (0)
    , _enableNetDocsLimit(false)
    , _fileCrawlStatus(false)
    , _vcodeServerStatus()
    , _productName()
    , _authVcodeServerStatus()
    , _passwordErrCnt (5)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);
    nsresult ret;
    _confManager = do_CreateInstance (NC_DB_CONF_MANAGER_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_DB_CONF_MANANGER,
            _T("Failed to create db conf manager: 0x%x"), ret);
    }

    _acsShareMgnt = do_CreateInstance (NC_ACS_SHARE_MGNT_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_PROCESSOR, FAILED_TO_CREATE_ACS_SHAREMGNT,
            _T("Failed to create acs sharemgnt: 0x%x"), ret);
    }

}

ncACSConfManager::~ncACSConfManager()
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void SetConfig ([const] in StringRef key, [const] in StringRef value); */
NS_IMETHODIMP_(void) ncACSConfManager::SetConfig(const String & key, const String & value)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p begin"), this);

    _confManager->SetConfig(key, value);

    NC_ACS_PROCESSOR_TRACE (_T("this: %p end"), this);
}

/* [notxpcom] String GetConfig ([const] in StringRef key); */
NS_IMETHODIMP_(String) ncACSConfManager::GetConfig(const String & key)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p begin"), this);

    return _confManager->GetConfig(key);
}

/* [notxpcom] String BatchGetConfig (in VectorStringRef keys, in StringMapRef kvMap); */
NS_IMETHODIMP_(void) ncACSConfManager::BatchGetConfig(vector<String>& keys, map<String, String>& kvMap)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p begin"), this);
    _confManager->BatchGetConfig(keys, kvMap);
    NC_ACS_PROCESSOR_TRACE (_T("this: %p end"), this);
}

/* [notxpcom] bool IsDownloadWatermarkDoc ([const] in StringRef docId, [const] in int docType, [const] in int64 size, [const] in StringRef path) */;
NS_IMETHODIMP_(bool) ncACSConfManager::IsDownloadWatermarkDoc (const String & docId, const int docType, const int64 size, const String & path)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p begin"), this);

    InitConfig();

    // 若文件为不支持加水印的格式直接返回false，文件夹按照cid处理
    if (size >= 0)
    {
        String extensionName = Path::getExtensionName (path);
        extensionName.toLower ();
        if (_needwatermarkType.end () == find (_needwatermarkType.begin (), _needwatermarkType.end (), extensionName)) {
            return false;
        }
    }

    if (GetDocWatermarkStatus())
    {
        int watermarkType = GetDocWatermarkTypeByDocID (ncGNSUtil::GetCIDPath (docId));
        if (2 == watermarkType || 3 == watermarkType)
        {
            return true;
        }
        else if (-1 != watermarkType)
        {
            return false;
        }

        // 没有对该CID配置水印策略时，查看全局策略
        watermarkType = GetDocWatermarkTypeByDocType (docType);
        if (2 == watermarkType || 3 == watermarkType)
        {
            return true;
        }
    }

    NC_ACS_PROCESSOR_TRACE (_T("this: %p end"), this);

    return false;
}

/* [notxpcom] bool GetFileCrawlStatus ();*/
NS_IMETHODIMP_(bool) ncACSConfManager::GetFileCrawlStatus ()
{
    return InterlockedExchangeAdd (TO_ATOMIC (&_fileCrawlStatus), 0);
}

/* [notxpcom] String GetVcodeServerStatus (); */
NS_IMETHODIMP_(String) ncACSConfManager::GetVcodeServerStatus()
{
    InitConfig();
    AutoLock<ThreadMutexLock> lock (&_vcodeServerStatusLock);
    return _vcodeServerStatus;
}

/* [notxpcom] String GetAuthVcodeServerStatus (); */
NS_IMETHODIMP_(String) ncACSConfManager::GetAuthVcodeServerStatus()
{
    InitConfig();
    AutoLock<ThreadMutexLock> lock (&_authVcodeServerStatusLock);
    return _authVcodeServerStatus;
}

/* [notxpcom] String GetProductName (); */
NS_IMETHODIMP_(String) ncACSConfManager::GetProductName()
{
    InitConfig();
    AutoLock<ThreadMutexLock> lock (&_productNameLock);
    return _productName;
}

/*[notxpcom] bool GetMessagePluginStatus ();*/
NS_IMETHODIMP_(bool) ncACSConfManager::GetMessagePluginStatus (void)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p begin"), this);
    InitConfig();

    bool ret = InterlockedExchangeAdd (TO_ATOMIC (&_enableMessagePlugin), 0);

    NC_ACS_PROCESSOR_TRACE (_T("this: %p end"), this);
    return ret;
}

/* [notxpcom] int GetPasswordErrCnt (); */
NS_IMETHODIMP_(int) ncACSConfManager::GetPasswordErrCnt()
{
    return InterlockedExchangeAdd (TO_ATOMIC (&_passwordErrCnt), 0);
}


// public
void
ncACSConfManager::InitConfig (void)
{
    if (!_initConfig) {
        AutoLock<ThreadMutexLock> lock (&_initConfiglock);
        if (!_initConfig) {
            UpdateValue ();

            // 初始化水印文件支持类型
            InitWatermarkSupportType();

            // 启动配置更新线程
            _configThread = new ncACSConfThread(this);
            _configThread->start ();
            _initConfig = true;
        }
    }
}

bool
ncACSConfManager::GetDocWatermarkStatus (void)
{
    return InterlockedExchangeAdd (TO_ATOMIC (&_enable_doc_watermark), 0);
}

int
ncACSConfManager::GetDocWatermarkTypeByDocID (const String docID)
{
    AutoLock<ThreadMutexLock> lock (&_watermarkDocsMapLock);
    map<String, int>::iterator iter = _watermarkDocsMap.find(docID);
    if (iter != _watermarkDocsMap.end ())
    {
        return iter->second;
    }
    return -1;
}

int
ncACSConfManager::GetDocWatermarkTypeByDocType (const int docType)
{
    AutoLock<ThreadMutexLock> lock (&_watermarkDocTypesMapLock);
    map<int, int>::iterator iter = _watermarkDocTypesMap.find (docType);
    if (iter != _watermarkDocTypesMap.end ())
    {
        return iter->second;
    }
    return -1;
}

void
ncACSConfManager::InitWatermarkSupportType ()
{
    try {
        FileInputStream fin (NC_WATERMARK_CONFIG_FILENAME, READ_WRITE_SHARE);
        int64 flen = fin.available ();
        string json (flen, '\0');
        fin.read (reinterpret_cast<unsigned char*> (&json[0]), flen);

        JSON::Value v;
        JSON::Reader::read (v, json.c_str (), json.length ());

        String needWatermarkStr = v["needwatermark"].s ().c_str ();
        needWatermarkStr.split (' ', _needwatermarkType);
    }
    catch (Exception& e) {
        SystemLog::getInstance ()->log (__FILE__, __LINE__, WARNING_LOG_TYPE,
                                        _T("EACP get config throw msg: %s"), e.toFullString ().getCStr ());
    }
    catch (...) {
        SystemLog::getInstance ()->log (__FILE__, __LINE__, WARNING_LOG_TYPE,
                                        _T("Get third party config throw msg: UNKnown"));
    }

    // 配置项不完整
    if (_needwatermarkType.empty ()) {
        SystemLog::getInstance ()->log (__FILE__, __LINE__, WARNING_LOG_TYPE,
                                        _T("Get third party config throw msg: Download watermark configuration information is incomplete or corrupted"));
    }
}

bool
ncACSConfManager::GetNetDocsLimitStatus (void)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p begin"), this);
    InitConfig();

    bool ret = InterlockedExchangeAdd (TO_ATOMIC (&_enableNetDocsLimit), 0);

    NC_ACS_PROCESSOR_TRACE (_T("this: %p end"), this);
    return ret;
}

void
ncACSConfManager::UpdateValue (void)
{
    // 水印策略开关
    // bool enable = _acsShareMgnt->GetWaterMarkStrategy();
    // InterlockedExchange (TO_ATOMIC (&_enable_doc_watermark), (int)enable);
    bool enable = false;

    // 文档库水印配置
    {
        AutoLock<ThreadMutexLock> lock (&_watermarkDocsMapLock);
        _acsShareMgnt->GetDownloadWatermarkDocs (_watermarkDocsMap);
    }

    // 不同文档库类型的水印设置
    {
        AutoLock<ThreadMutexLock> lock (&_watermarkDocTypesMapLock);
        _acsShareMgnt->GetDownloadWatermarkDocTypes (_watermarkDocTypesMap);
    }

    // 网段文档库限制
    enable = _acsShareMgnt->GetNetDocsLimitStatus();
    InterlockedExchange (TO_ATOMIC (&_enableNetDocsLimit), (int)enable);

    // 文档抓取策略
    enable = _acsShareMgnt->GetFileCrawlStatus();
    InterlockedExchange (TO_ATOMIC (&_fileCrawlStatus), (int)enable);

    // 消息推送插件
    enable = _acsShareMgnt->GetMessagePluginStatus();
    InterlockedExchange (TO_ATOMIC (&_enableMessagePlugin), (int)enable);


    vector<String> sharemgntKeys;
    map<String, String> sharemgntKvMap;
    sharemgntKeys.push_back("vcode_server_status");             // 验证码发送服务器开关状态信息
    sharemgntKeys.push_back("dualfactor_auth_server_status");   // 双因子验证码服务器开关状态信息
    sharemgntKeys.push_back("pwd_err_cnt");                     // 密码错误次数
    _acsShareMgnt->BatchGetConfig(sharemgntKeys, sharemgntKvMap);

    // 验证码发送服务器开关状态信息
    {
        AutoLock<ThreadMutexLock> lock (&_vcodeServerStatusLock);
        if (sharemgntKvMap.find("vcode_server_status") != sharemgntKvMap.end()) {
            _vcodeServerStatus = sharemgntKvMap["vcode_server_status"];
        }
    }

    // 双因子验证码服务器开关状态信息
    {
        AutoLock<ThreadMutexLock> lock (&_authVcodeServerStatusLock);
        if (sharemgntKvMap.find("dualfactor_auth_server_status") != sharemgntKvMap.end()) {
            _authVcodeServerStatus = sharemgntKvMap["dualfactor_auth_server_status"];
        }
    }

    // 密码错误次数
    if (sharemgntKvMap.find("pwd_err_cnt") != sharemgntKvMap.end()) {
        InterlockedExchange (TO_ATOMIC (&_passwordErrCnt), Int::getValue (sharemgntKvMap["pwd_err_cnt"]));
    }

    // OEM中Product信息
    {
        AutoLock<ThreadMutexLock> lock (&_productNameLock);
        String section = "shareweb_en-us";
        String option = "product";
        _acsShareMgnt->OEM_GetConfigByOption(section, option, _productName);
    }
}

ncACSConfThread::ncACSConfThread(ncACSConfManager* configManager)
    : _threadToRun (true)
    , _configManager (configManager)
{
}

ncACSConfThread::~ncACSConfThread (void)
{
    NC_ACS_PROCESSOR_TRACE (_T("this: %p"), this);

    _threadToRun = false;
    join ();
}

void
ncACSConfThread::run (void)
{
    while (_threadToRun) {
        try {
            _configManager->UpdateValue ();
        }
        catch (Exception& e) {
            SystemLog::getInstance ()->log (__FILE__, __LINE__, WARNING_LOG_TYPE,
                                            _T("Update eacp config value thread throw msg:%s"), e.toFullString ().getCStr ());
        }
        catch (...) {
            SystemLog::getInstance ()->log (__FILE__, __LINE__, WARNING_LOG_TYPE,
                                            _T("Update eacp config value thread throw msg: UNKnown"));
        }
        Thread::sleep (10 * 1000);
    }
}
