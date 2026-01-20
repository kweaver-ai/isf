#include <abprec.h>

#include "acsdb.h"
#include "ncDBConfManager.h"
#include "common/util.h"

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (ncDBConfManager, ncIDBConfManager)

// protected
NS_IMETHODIMP_(nsrefcnt) ncDBConfManager::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) ncDBConfManager::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (ncDBConfManager)

ncDBConfManager::ncDBConfManager()
{
    NC_ACS_DB_TRACE (_T("this: %p"), this);
}

ncDBConfManager::~ncDBConfManager()
{
    NC_ACS_DB_TRACE (_T("this: %p"), this);
}

/* [notxpcom] void SetConfig ([const] in StringRef key, [const] in StringRef value); */
NS_IMETHODIMP_(void) ncDBConfManager::SetConfig(const String & key, const String & value)
{
    NC_ACS_DB_TRACE (_T("this: %p begin"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String escKey = dbOper->EscapeEx(key);

    String escValue = dbOper->EscapeEx(value);

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("select f_key from %s.t_conf where f_key = '%s'"),
                    dbName.getCStr(), escKey.getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    if(results.size() == 0) {
        strSql.format (_T("insert into %s.t_conf (f_key,f_value) values ('%s','%s')"),
                            dbName.getCStr(), escKey.getCStr (), escValue.getCStr ());
    }
    else {
        strSql.format (_T("update %s.t_conf set f_value = '%s' where f_key = '%s'"),
                            dbName.getCStr(), escValue.getCStr (), escKey.getCStr ());
    }

    dbOper->Execute (strSql);

    NC_ACS_DB_TRACE (_T("this: %p end"), this);
}

/* [notxpcom] String GetConfig ([const] in StringRef key); */
NS_IMETHODIMP_(String) ncDBConfManager::GetConfig(const String & key)
{
    NC_ACS_DB_TRACE (_T("this: %p begin"), this);

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String escKey = dbOper->EscapeEx(key);

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("select f_value from %s.t_conf where f_key = '%s'"),
                    dbName.getCStr(), escKey.getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    String value;

    if(results.size() == 1) {
        value = results[0][0];
    }

    NC_ACS_DB_TRACE (_T("this: %p end"), this);

    return value;
}

/* [notxpcom] String BatchGetConfig (in VectorStringRef keys, in StringMapRef kvMap); */
NS_IMETHODIMP_(void) ncDBConfManager::BatchGetConfig(vector<String>& keys, map<String, String>& kvMap)
{
    NC_ACS_DB_TRACE (_T("this: %p begin"), this);

    kvMap.clear();

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String groupStr;
    for (size_t i = 0; i < keys.size (); ++i) {
        groupStr.append ("\'", 1);
        groupStr.append (dbOper->EscapeEx(keys[i]));
        groupStr.append ("\'", 1);

        if (i != (keys.size () -1)) {
            groupStr.append (",", 1);
        }
    }

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("select f_key, f_value from %s.t_conf where f_key in (%s)"),
                    dbName.getCStr(), groupStr.getCStr ());

    ncDBRecords results;
    dbOper->Select (strSql, results);

    for (size_t i = 0; i < results.size (); ++i) {
        kvMap[results[i][0]] = results[i][1];
    }

    NC_ACS_DB_TRACE (_T("this: %p end"), this);
}

ncIDBOperator* ncDBConfManager::GetDBOperator ()
{
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (ncACSDBGetDBOperator ());
    NS_ADDREF (dbOper.get ());

    Init (dbOper);

    return dbOper;
}

void ncDBConfManager::Init (ncIDBOperator* dbOper)
{
    static bool sIsTableCreated = false;
    static ThreadMutexLock sInitLock;
    // 已连接直接返回
    if (sIsTableCreated == true) {
        return;
    }

    AutoLock<ThreadMutexLock> lock (&sInitLock);
    if (sIsTableCreated == true) {
        return;
    }

    NC_ACS_DB_TRACE (_T("%p"), this);

    // 默认数据
    map<String, String> defaultValues;

    // 是否开启文件自动锁
    defaultValues["auto_lock"] = "true";

    // 是否开启文件自动锁提醒
    defaultValues["auto_lock_remind"] = "true";

    // webclient的http端口
    defaultValues["web_client_http_port"] = "80";

    // 内外网一致域名地址
    defaultValues["web_client_host"] = "";

    // webclient的https端口
    defaultValues["web_client_port"] = "443";

    // 客户端是否可以配置无限期的权限
    defaultValues["oem_indefinite_perm"] = "true";

    // 客户端是否可以配置所有者
    defaultValues["oem_allow_owner"] = "true";

    // 客户端是否可以记住密码
    defaultValues["oem_remember_pass"] = "true";

    // webconsole最大的密码过期天数，-1表示无限制
    defaultValues["oem_max_pass_expired_days"] = "-1";

    // 是否可以共享高密级的文件给低密级的用户
    defaultValues["oem_allow_auth_low_csf_user"] = "true";

    // 客户端超时退出分钟，-1表示无限制
    defaultValues["oem_client_logout_time"] = "-1";

    // 客户端配置权限时，默认的有效天数，-1表示无限期
    defaultValues["oem_default_perm_expired_days"] = "-1";

    // 是否开启文件传输限制功能
    defaultValues["oem_enable_file_transfer_limit"] = "false";

    // 入口文档视图模式
    defaultValues["entrydoc_view_config"] = "1";

    // 是否启用onedrive跳转
    defaultValues["oem_enable_onedrive"] = "false";

    // 是否开启内外网数据交换
    defaultValues["enable_exchange_file"] = "false";

    // 是否启用外部程序超级表格
    defaultValues["enable_chaojibiaoge"] = "false";

    // 是否启用秦淮电教馆功能定制按钮， 供web界面显示
    defaultValues["enable_qhdj"] = "false";

    // 是否允许AnyShare客户端手动登录。
    defaultValues["client_manual_login"] = "true";

    // 是否开启消息推送至第三方应用程序，"null"为关闭，否则为第三方应用Id，如广联达OA系统为"glodon"
    defaultValues["push_message_to_third_party"] = "null";

    // 内链地址的前缀
    defaultValues["internal_link_prefix"] = "AnyShare://";

    // 是否开启邮件通知分享功能，"true" 为开启， "false" 为关闭
    defaultValues["enable_send_share_mail"] = "false";

    // 入口文档视图模式
    defaultValues["show_knowledge_page"] = "0";

    // 是否启用消息通知
    defaultValues["enable_message_notify"] = "true";

    // 浩辰CAD使用大图插件的临界值
    defaultValues["cad_plugin_threshold"] = "10485760";

    // 是否启用修改密码签名认证
    defaultValues["enable_eacp_check_sign"] = "false";

    // appid登录，是否进行设备绑定检查
    defaultValues["check_appid_login_device_bind"] = "true";

    String dbName = Util::getDBName("anyshare");
    for(map<String, String>::iterator iter = defaultValues.begin(); iter != defaultValues.end(); ++iter) {
        String strSql;
        strSql.format (_T("select f_key from %s.t_conf where f_key = '%s'"), dbName.getCStr(), iter->first.getCStr());

        ncDBRecords results;
        dbOper->Select (strSql, results);

        if(results.size() == 0) {
            strSql.format (_T("insert into %s.t_conf (f_key,f_value) values ('%s','%s')"),
                dbName.getCStr(), iter->first.getCStr(), iter->second.getCStr());
            dbOper->Execute (strSql);
        }
    }

    sIsTableCreated = true;
}
