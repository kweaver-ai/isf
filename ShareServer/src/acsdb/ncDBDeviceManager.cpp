#include <abprec.h>

#include "acsdb.h"
#include "ncDBDeviceManager.h"
#include "gen-cpp/ShareMgnt_constants.h"
#include "common/util.h"

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (ncDBDeviceManager, ncIDBDeviceManager)

// protected
NS_IMETHODIMP_(nsrefcnt) ncDBDeviceManager::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) ncDBDeviceManager::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (ncDBDeviceManager)

ncDBDeviceManager::ncDBDeviceManager()
{
    /* member initializers and constructor code */
}

ncDBDeviceManager::~ncDBDeviceManager()
{
    /* destructor code */
}

/* [notxpcom] void GetDevicesByUserIdAndUdid ([const] in StringRef userId, [const] in StringRef udid, in dbDeviceInfoVecRef deviceInfos, in int start, in int limit); */
NS_IMETHODIMP_(void) ncDBDeviceManager::GetDevicesByUserIdAndUdid(const String & userId, const String & udid, vector<dbDeviceInfo> & deviceInfos, int start, int limit)
{
    NC_ACS_DB_TRACE (_T("userId: %s begin"), userId.getCStr());
    deviceInfos.clear();

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    limit = limit < 0 ? Int::MAX_VALUE : limit;

    String dbName = Util::getDBName("anyshare");
    String strSql;
    if (udid.isEmpty()) {
        strSql.format (_T("select f_udid,f_name,f_os_type,f_device_type,f_last_login_ip,f_last_login_time,f_erase_flag,f_last_erase_time,f_disable_flag,f_bind_flag ")
                        _T("from %s.t_device where f_user_id = '%s'")
                        _T("order by f_last_login_time desc ")
                        _T("limit %d, %d"),
                        dbName.getCStr(), dbOper->EscapeEx(userId).getCStr(),
                        start, limit);
    }
    else if (userId.isEmpty()) {
        strSql.format (_T("select f_udid,f_name,f_os_type,f_device_type,f_last_login_ip,f_last_login_time,f_erase_flag,f_last_erase_time,f_disable_flag,f_bind_flag ")
                        _T("from %s.t_device where f_udid like '%%%s%%'")
                        _T("limit %d, %d"),
                        dbName.getCStr(), dbOper->EscapeEx(udid).getCStr(),
                        start, limit);
    }
    else {
        strSql.format (_T("select f_udid,f_name,f_os_type,f_device_type,f_last_login_ip,f_last_login_time,f_erase_flag,f_last_erase_time,f_disable_flag,f_bind_flag ")
                        _T("from %s.t_device where f_user_id = '%s' and (f_udid = '%s' or f_udid like '%s%%')")
                        _T("order by f_last_login_time desc ")
                        _T("limit %d, %d"),
                        dbName.getCStr(), dbOper->EscapeEx(userId).getCStr(), dbOper->EscapeEx(udid).getCStr(), dbOper->EscapeEx(udid).getCStr(),
                        start, limit);
    }

    ncDBRecords records;
    dbOper->Select(strSql, records);

    dbDeviceInfo tmpInfo;
    for(size_t i = 0; i < records.size(); ++i) {
        tmpInfo.baseInfo.udid = records[i][0];
        tmpInfo.baseInfo.name = records[i][1];
        tmpInfo.baseInfo.osType = Int::getValue(records[i][2]);
        tmpInfo.baseInfo.deviceType = records[i][3];
        tmpInfo.baseInfo.lastLoginIp = records[i][4];
        tmpInfo.baseInfo.lastLoginTime = Int64::getValue(records[i][5]);
        tmpInfo.eraseFlag = Int::getValue(records[i][6]);
        tmpInfo.lastEraseTime = Int64::getValue(records[i][7]);
        tmpInfo.disableFlag = Int::getValue(records[i][8]);
        tmpInfo.bindFlag = Int::getValue(records[i][9]);

        deviceInfos.push_back(tmpInfo);
    }

    NC_ACS_DB_TRACE (_T("userId: %s end"), userId.getCStr());
}

/* [notxpcom] void AddDevice ([const] in StringRef userId, [const] in dbDeviceBaseInfoRef deviceInfo); */
NS_IMETHODIMP_(void) ncDBDeviceManager::AddDevice(const String & userId, const dbDeviceBaseInfo & baseInfo)
{
    NC_ACS_DB_TRACE (_T("userId: %s, udid: %s begin"), userId.getCStr(), baseInfo.udid.getCStr());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String escUDID = dbOper->EscapeEx(baseInfo.udid);
    String escName = dbOper->EscapeEx(baseInfo.name);
    String escDeviceType = dbOper->EscapeEx(baseInfo.deviceType);

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format ("insert into %s.t_device(f_user_id, f_udid, f_name, f_os_type, f_device_type, f_last_login_ip, f_last_login_time) "
                    "values('%s','%s','%s', %d, '%s','%s',%lld)",
                    dbName.getCStr(),
                    dbOper->EscapeEx(userId).getCStr(),
                    escUDID.getCStr(),
                    escName.getCStr(),
                    baseInfo.osType,
                    escDeviceType.getCStr(),
                    baseInfo.lastLoginIp.getCStr(),
                    baseInfo.lastLoginTime);

    dbOper->Execute(strSql);

    NC_ACS_DB_TRACE (_T("userId: %s, udid: %s end"), userId.getCStr(), baseInfo.udid.getCStr());
}

/*[notxpcom] void AddDevices([const] in StringRef userId, [const] in dbDeviceBaseInfoVecRef baseInfos);*/
NS_IMETHODIMP_(void) ncDBDeviceManager::AddDevices(const String & userId, const vector<dbDeviceBaseInfo> & baseInfos)
{
    NC_ACS_DB_TRACE (_T("userId: %s begin"), userId.getCStr());

    if (baseInfos.size () == 0) {
        return;
    }

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String dbName = Util::getDBName("anyshare");
    // 获取 用户已有的设备信息
    String selectSql;
    selectSql.format ("select f_udid from %s.t_device where f_user_id = '%s'",
                        dbName.getCStr(), dbOper->EscapeEx(userId).getCStr());
    ncDBRecords records;
    dbOper->Select(selectSql, records);

    map<String, bool> mapExistData;
    for (size_t i = 0; i < records.size(); ++i) {
        mapExistData[records[i][0]] = true;
    }

    // 添加数据
    String sqlPrefix;
    sqlPrefix.format ("insert into %s.t_device(f_user_id, f_udid, f_name, f_os_type, f_device_type, f_last_login_ip, f_last_login_time) values",
                        dbName.getCStr());

    String valueString;
    int itemNum = 0;
    for (auto iter = baseInfos.begin (); iter != baseInfos.end (); iter++) {
        // 检查是否存在此数据
        auto temp = mapExistData.find(iter->udid);
        if (temp != mapExistData.end()) {
            continue;
        }

        String escUDID = dbOper->EscapeEx (iter->udid);
        String escName = dbOper->EscapeEx (iter->name);
        String escDeviceType = dbOper->EscapeEx (iter->deviceType);

        String valueItem;
        valueItem.format ("('%s','%s','%s', %d, '%s','%s',%lld)",
                          dbOper->EscapeEx(userId).getCStr(),
                          escUDID.getCStr(),
                          escName.getCStr(),
                          iter->osType,
                          escDeviceType.getCStr(),
                          iter->lastLoginIp.getCStr(),
                          iter->lastLoginTime);

        valueString.append(valueItem);
        ++itemNum;
        if (itemNum % 3000 == 0) { // 3000条记录执行一次SQL
            dbOper->Execute(sqlPrefix + valueString);
            valueString.clear ();
        }
        else {
            valueString.append(",");
        }
    }

    if (!valueString.isEmpty()) { // 插入剩余item记录
        String newValueString = valueString.subString(0, (valueString.getLength()-1)); // 去掉最后的逗号分隔符
        dbOper->Execute(sqlPrefix + newValueString);
    }

    NC_ACS_DB_TRACE (_T("userId: %s end"), userId.getCStr());
}

/*[notxpcom] void GetUserNameByUdid([const] in StringRef udid, in VecStringRef users, in int start, in int limit);*/
NS_IMETHODIMP_(void) ncDBDeviceManager::GetUserNameByUdid(const String& udid, vector<String>& users, int start, int limit)
{
    NC_ACS_DB_TRACE (_T("udid: %s begin"), udid.getCStr());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    users.clear();

    limit = limit < 0 ? Int::MAX_VALUE : limit;

    String escUDID = dbOper->EscapeEx(udid);

    String dbName = Util::getDBName("anyshare");
    String dbName2 = Util::getDBName("sharemgnt_db");
    String sql;
    sql.format(_T("select d.f_user_id,u.f_display_name as display_name from %s.t_device as d "
                  "left join %s.t_user as u on u.f_user_id = d.f_user_id "
                  "where d.f_udid = '%s' order by upper(display_name) limit %d,%d"),
               dbName.getCStr(),
               dbName2.getCStr(),
               escUDID.getCStr(), start, limit);

    ncDBRecords records;
    dbOper->Select(sql, records);
    for (size_t i = 0; i < records.size(); i++) {
        if (records[i][0] == toCFLString(g_ShareMgnt_constants.NCT_ALL_USER_GROUP)) {
            users.push_back(LOAD_STRING (_T("IDS_DB_ALL_USER_GROUP_NAME")));
        }
        else if (!records[i][1].isEmpty()) {
            users.push_back(records[i][1]);
        }
    }

    NC_ACS_DB_TRACE (_T("udid: %s end"), udid.getCStr());
}

/* [notxpcom] void DeleteDevices ([const] in StringRef userId, [const] in VecStringRef udid); */
NS_IMETHODIMP_(void) ncDBDeviceManager::DeleteDevices(const String & userId, const vector<String> & udids)
{
    NC_ACS_DB_TRACE("userId: %s begin", userId.getCStr());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String escUserId = dbOper->EscapeEx(userId);

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("delete from %s.t_device where f_user_id='%s'"),
                   dbName.getCStr(), escUserId.getCStr());

    // 设备ID子条件，如果 udids 为空，删除用户的所有设备
    String groupStr;
    for (size_t i = 0; i < udids.size (); ++i) {
        groupStr.append ("\'", 1);
        groupStr.append (dbOper->EscapeEx(udids[i]));
        groupStr.append ("\'", 1);

        if (i != (udids.size () -1)) {
            groupStr.append (",", 1);
        }
    }

    String udidClause;
    if (!groupStr.isEmpty()) {
        udidClause.format (" and f_udid in (%s)", groupStr.getCStr ());
    }
    strSql.append (udidClause);

    dbOper->Execute(strSql);

    NC_ACS_DB_TRACE("userId: %s end", userId.getCStr());
}

/* [notxpcom] bool GetDeviceByUDID ([const] in StringRef userId, [const] in StringRef udid, in dbDeviceInfoRef deviceInfo); */
NS_IMETHODIMP_(bool) ncDBDeviceManager::GetDeviceByUDID(const String & userId, const String & udid, dbDeviceInfo & deviceInfo)
{
    NC_ACS_DB_TRACE (_T("userId: %s, udid: %s begin"), userId.getCStr(), deviceInfo.baseInfo.udid.getCStr());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String escUserId = dbOper->EscapeEx(userId);
    String escUDID = dbOper->EscapeEx(udid);

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("select f_name,f_os_type,f_device_type,f_last_login_ip,f_last_login_time,f_erase_flag,f_last_erase_time,f_disable_flag,f_bind_flag from %s.t_device where f_user_id = '%s' and f_udid = '%s'"),
                    dbName.getCStr(), escUserId.getCStr(), escUDID.getCStr());

    ncDBRecords records;
    dbOper->Select(strSql, records);

    bool ret = false;
    if(records.size() >= 1) {
        deviceInfo.baseInfo.udid = udid;
        deviceInfo.baseInfo.name = records[0][0];
        deviceInfo.baseInfo.osType = Int::getValue(records[0][1]);
        deviceInfo.baseInfo.deviceType = records[0][2];
        deviceInfo.baseInfo.lastLoginIp = records[0][3];
        deviceInfo.baseInfo.lastLoginTime = Int64::getValue(records[0][4]);
        deviceInfo.eraseFlag = Int::getValue(records[0][5]);
        deviceInfo.lastEraseTime = Int64::getValue(records[0][6]);
        deviceInfo.disableFlag = Int::getValue(records[0][7]);
        deviceInfo.bindFlag = Int::getValue(records[0][8]);

        ret = true;
    }

    NC_ACS_DB_TRACE (_T("userId: %s, udid: %s end"), userId.getCStr(), deviceInfo.baseInfo.udid.getCStr());

    return ret;
}

/* [notxpcom] void UpdateDevice ([const] in StringRef userId, [const] in dbDeviceBaseInfoRef baseInfo); */
NS_IMETHODIMP_(void) ncDBDeviceManager::UpdateDevice(const String & userId, const dbDeviceBaseInfo & baseInfo)
{
    NC_ACS_DB_TRACE (_T("userId: %s, udid: %s begin"), userId.getCStr(), baseInfo.udid.getCStr());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String escUserId = dbOper->EscapeEx(userId);
    String escUDID = dbOper->EscapeEx(baseInfo.udid);
    String escName = dbOper->EscapeEx(baseInfo.name);
    String escDeviceType = dbOper->EscapeEx(baseInfo.deviceType);

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("update %s.t_device set f_name='%s',f_os_type=%d,f_device_type='%s',f_last_login_ip='%s',f_last_login_time=%lld where f_user_id='%s' and f_udid='%s'"),
                    dbName.getCStr(), escName.getCStr(), baseInfo.osType, escDeviceType.getCStr(), baseInfo.lastLoginIp.getCStr(), baseInfo.lastLoginTime, escUserId.getCStr(), escUDID.getCStr());

    dbOper->Execute(strSql);

    NC_ACS_DB_TRACE (_T("userId: %s, udid: %s end"), userId.getCStr(), baseInfo.udid.getCStr());
}

/* [notxpcom] void SetEraseStatus([const] in StringRef userId, [const] in StringRef udid); */
NS_IMETHODIMP_(void) ncDBDeviceManager::SetEraseStatus(const String & userId, const String & udid)
{
    NC_ACS_DB_TRACE (_T("userId: %s, udid: %s begin"), userId.getCStr(), udid.getCStr());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String escUserId = dbOper->EscapeEx(userId);
    String escUDID = dbOper->EscapeEx(udid);

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("update %s.t_device set f_erase_flag=1 where f_user_id='%s' and f_udid='%s'"),
                    dbName.getCStr(), escUserId.getCStr(), escUDID.getCStr());

    dbOper->Execute(strSql);

    NC_ACS_DB_TRACE (_T("userId: %s, udid: %s end"), userId.getCStr(), udid.getCStr());
}

/* [notxpcom] void SetEraseSucInfo([const] in StringRef userId, [const] in StringRef udid, in int64 date); */
NS_IMETHODIMP_(void) ncDBDeviceManager::SetEraseSucInfo(const String & userId, const String & udid, int64 date)
{
    NC_ACS_DB_TRACE("userId: %s, udid: %s begin", userId.getCStr(), udid.getCStr());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String escUserId = dbOper->EscapeEx(userId);
    String escUDID = dbOper->EscapeEx(udid);

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("update %s.t_device set f_erase_flag=0,f_last_erase_time=%lld where f_user_id='%s' and f_udid='%s'"),
                    dbName.getCStr(), date, escUserId.getCStr(), escUDID.getCStr());

    dbOper->Execute(strSql);

    NC_ACS_DB_TRACE("userId: %s, udid: %s begin", userId.getCStr(), udid.getCStr());
}

/* [notxpcom] void SetDisableStatus([const] in StringRef userId, [const] in StringRef udid); */
NS_IMETHODIMP_(void) ncDBDeviceManager::SetDisableStatus(const String& userId, const String & udid)
{
    NC_ACS_DB_TRACE (_T("userId: %s, udid: %s begin"), userId.getCStr(), udid.getCStr());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String escUserId = dbOper->EscapeEx(userId);
    String escUDID = dbOper->EscapeEx(udid);

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("update %s.t_device set f_disable_flag=1 where f_user_id='%s' and f_udid='%s'"),
                    dbName.getCStr(), escUserId.getCStr(), escUDID.getCStr());

    dbOper->Execute(strSql);

    NC_ACS_DB_TRACE (_T("userId: %s, udid: %s end"), userId.getCStr(), udid.getCStr());
}

/* [notxpcom] void SetEnableStatus([const] in StringRef userId, [const] in StringRef udid); */
NS_IMETHODIMP_(void) ncDBDeviceManager::SetEnableStatus(const String& userId, const String & udid)
{
    NC_ACS_DB_TRACE (_T("userId: %s, udid: %s begin"), userId.getCStr(), udid.getCStr());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String escUserId = dbOper->EscapeEx(userId);
    String escUDID = dbOper->EscapeEx(udid);

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("update %s.t_device set f_disable_flag=0 where f_user_id='%s' and f_udid='%s'"),
                    dbName.getCStr(), escUserId.getCStr(), escUDID.getCStr());

    dbOper->Execute(strSql);

    NC_ACS_DB_TRACE (_T("userId: %s, udid: %s end"), userId.getCStr(), udid.getCStr());
}

/* [notxpcom] void SetDisableStatus([const] in StringRef userId, [const] in VecStringRef udids); */
NS_IMETHODIMP_(void) ncDBDeviceManager::SetBindStatus(const String& userId, const vector<String> & udids)
{
    NC_ACS_DB_TRACE (_T("userId: %s begin"), userId.getCStr());

    if (udids.size () == 0) {
        return;
    }

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String groupStr;
    for (size_t i = 0; i < udids.size (); ++i) {
        groupStr.append ("\'", 1);
        groupStr.append (dbOper->EscapeEx(udids[i]));
        groupStr.append ("\'", 1);

        if (i != (udids.size () -1)) {
            groupStr.append (",", 1);
        }
    }

    String escUserId = dbOper->EscapeEx(userId);

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("update %s.t_device set f_bind_flag=1 where f_user_id='%s' and f_udid in(%s)"),
                    dbName.getCStr(), escUserId.getCStr(), groupStr.getCStr ());

    dbOper->Execute(strSql);

    NC_ACS_DB_TRACE (_T("userId: %s end"), userId.getCStr());
}

/* [notxpcom] void SetEnableStatus([const] in StringRef userId, [const] in StringRef udid); */
NS_IMETHODIMP_(void) ncDBDeviceManager::SetUnbindStatus(const String& userId, const String & udid)
{
    NC_ACS_DB_TRACE (_T("userId: %s, udid: %s begin"), userId.getCStr(), udid.getCStr());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());
    String escUserId = dbOper->EscapeEx(userId);
    String escUDID = dbOper->EscapeEx(udid);

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("update %s.t_device set f_bind_flag=0 where f_user_id='%s' and f_udid='%s'"),
                    dbName.getCStr(), escUserId.getCStr(), escUDID.getCStr());

    dbOper->Execute(strSql);

    NC_ACS_DB_TRACE (_T("userId: %s, udid: %s end"), userId.getCStr(), udid.getCStr());
}

/* [notxpcom] void GetBindDevices([const] in StringRef userId, in StringIntMapRef udidsMap); */
NS_IMETHODIMP_(void) ncDBDeviceManager::GetBindDevices(const String & userId, map<String, int>& udidsMap)
{
    NC_ACS_DB_TRACE (_T("userId: %s begin"), userId.getCStr());
    udidsMap.clear();

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String escUserId = dbOper->EscapeEx(userId);

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("select f_udid, f_disable_flag from %s.t_device where f_user_id='%s' and f_bind_flag=1"),
                    dbName.getCStr(), escUserId.getCStr());

    ncDBRecords records;
    dbOper->Select(strSql, records);

    for(size_t i = 0; i < records.size(); ++i) {
        records[i][0].toUpper();
        udidsMap[records[i][0]] = Int::getValue(records[i][1]);
    }

    NC_ACS_DB_TRACE (_T("userId: %s end"), userId.getCStr());
}

/* [notxpcom] void GetBindDevicesWithAllUser([const] in StringRef userIds, in StringVecMapRef udidsMap); */
NS_IMETHODIMP_(void) ncDBDeviceManager::GetBindDevicesWithAllUser(const String & userId, map<String, vector<String> >& udidsMap)
{
    NC_ACS_DB_TRACE (_T("userId: %s begin"), userId.getCStr());
    udidsMap.clear();

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String escUserId = dbOper->EscapeEx(userId);

    String userIdClause;
    if (userId != toCFLString(g_ShareMgnt_constants.NCT_ALL_USER_GROUP)) {
        userIdClause.format (" (f_user_id = '%s' or f_user_id = '%s') and ",
                            dbOper->EscapeEx(userId).getCStr (),
                            (g_ShareMgnt_constants.NCT_ALL_USER_GROUP).c_str ());
    }

    String dbName = Util::getDBName("anyshare");
    String strSql;
    strSql.format (_T("select f_user_id, f_udid from %s.t_device where %s f_bind_flag=1"),
                    dbName.getCStr(),
                    userIdClause.getCStr());

    ncDBRecords records;
    dbOper->Select(strSql, records);

    for(size_t i = 0; i < records.size(); ++i) {
        records[i][1].toUpper();
        udidsMap[records[i][0]].push_back(records[i][1]);
    }

    NC_ACS_DB_TRACE (_T("userId: %s end"), userId.getCStr());
}

/*[notxpcom] bool UserHasDeviceInfo([const] in StringRef userId, [const] in bool binded);*/
NS_IMETHODIMP_(bool) ncDBDeviceManager::UserHasDeviceInfo (const String & userId, bool binded)
{
    NC_ACS_DB_TRACE (_T("userId: %s begin"), userId.getCStr());

    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetDBOperator ());

    String escUserId = dbOper->EscapeEx(userId);

    String dbName = Util::getDBName("anyshare");
    String strSql;
    if (binded) {
        strSql.format (_T("select f_user_id as count from %s.t_device where f_user_id = '%s' and f_bind_flag = 1 limit 1"),
                        dbName.getCStr(), escUserId.getCStr ());
    }
    else {
        strSql.format (_T("select f_user_id as count from %s.t_device where f_user_id = '%s' limit 1"),
                        dbName.getCStr(), escUserId.getCStr ());
    }

    ncDBRecords records;
    dbOper->Select(strSql, records);

    NC_ACS_DB_TRACE (_T("userId: %s end"), userId.getCStr());

    return records.size() > 0;
}

ncIDBOperator* ncDBDeviceManager::GetDBOperator ()
{
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (ncACSDBGetDBOperator ());
    NS_ADDREF (dbOper.get ());
    return dbOper;
}
