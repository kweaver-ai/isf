#include <abprec.h>
#include <ncutil/ncBusinessDate.h>
#include <boost/progress.hpp>
#include <boost/regex.hpp>
#include <boost/date_time/posix_time/posix_time.hpp>
#include <boost/date_time/local_time_adjustor.hpp>
#include <boost/date_time/c_local_time_adjustor.hpp>

#include <openssl/rsa.h>
#include <openssl/pem.h>
#include <openssl/err.h>
#include <openssl/evp.h>
#include <openssl/des.h>
#include <openssl/rand.h>

#include <mysql/mysql.h>

#include <ncutil/ncutil.h>
#include <ncutil/ncPerformanceProfilerPrec.h>
#include <ethriftutil/ncThriftClient.h>
#include <dboperator/public/ncIDBOperator.h>

#include "./gen-cpp/ncTEVFS.h"
#include "./gen-cpp/EVFS_constants.h"

#include "./gen-cpp/ncTEACP.h"
#include "./gen-cpp/EACP_constants.h"

#include "ncACSUtil.h"

ncACSUtil::ncACSUtil ()
{
    AddCommandFun (_T("ptime"), (ncCommandFunc)&ncACSUtil::ptime);
    AddCommandFun (_T("ticks2String"), (ncCommandFunc)&ncACSUtil::ticks2String);
    AddCommandFun (_T("md5_encrypt"), (ncCommandFunc)&ncACSUtil::md5_encrypt);
    AddCommandFun (_T("base64_encrypt"), (ncCommandFunc)&ncACSUtil::base64_encrypt);
    AddCommandFun (_T("base64_decrypt"), (ncCommandFunc)&ncACSUtil::base64_decrypt);
    AddCommandFun (_T("rsa_encrypt"), (ncCommandFunc)&ncACSUtil::rsa_encrypt);
    AddCommandFun (_T("rsa_decrypt"), (ncCommandFunc)&ncACSUtil::rsa_decrypt);
    AddCommandFun (_T("des_encrypt"), (ncCommandFunc)&ncACSUtil::des_encrypt);
    AddCommandFun (_T("des_decrypt"), (ncCommandFunc)&ncACSUtil::des_decrypt);
    AddCommandFun (_T("mtdecrypt"), (ncCommandFunc)&ncACSUtil::mtdecrypt);
    AddCommandFun (_T("getEndOfDayTicks"), (ncCommandFunc)&ncACSUtil::getEndOfDayTicks);
    AddCommandFun (_T("testString"), (ncCommandFunc)&ncACSUtil::testString);
    AddCommandFun (_T("getParentString"), (ncCommandFunc)&ncACSUtil::getParentString);
    AddCommandFun (_T("getUTF8StringSize"), (ncCommandFunc)&ncACSUtil::getUTF8StringSize);
    AddCommandFun (_T("updateDB"), (ncCommandFunc)&ncACSUtil::updateDB);
    AddCommandFun (_T("testdb"), (ncCommandFunc)&ncACSUtil::testdb);
    AddCommandFun (_T("conndb"), (ncCommandFunc)&ncACSUtil::conndb);
    AddCommandFun (_T("isValidURL"), (ncCommandFunc)&ncACSUtil::isValidURL);

    AddCommandFun (_T("createuserdoc"), (ncCommandFunc)&ncACSUtil::createuserdoc);
    AddCommandFun (_T("setquota"), (ncCommandFunc)&ncACSUtil::setquota);
}

ncACSUtil::~ncACSUtil ()
{
}

void ncACSUtil::ticks2String()
{
    int64 ticks = BusinessDate::getCurrentTime ();
    Date dateValue (ticks);
    printMessage2 (_T("%lld: %s"), ticks, dateValue.toString (FD_GENERAL_FULL).getCStr ());

    using boost::date_time::c_local_adjustor;
    using boost::posix_time::from_time_t;
    using boost::posix_time::ptime;

    ptime pt = c_local_adjustor<ptime>::utc_to_local(from_time_t((ticks / Date::ticksPerMillisecond / 1000)));
    printMessage2 (_T("%lld: %s"), ticks, boost::posix_time::to_simple_string(pt).c_str ());

    boost::posix_time::ptime ptCurrent(boost::posix_time::second_clock::local_time());
    printMessage2 (_T("%s"), boost::posix_time::to_simple_string(ptCurrent).c_str ());
}

void ncACSUtil::ptime ()
{
    String timeValue = GetString (_T("value"));

    int64 iValue = Int64::getValue (timeValue);
    Date dateValue (iValue);

    printMessage2 (_T("%s: %s"), timeValue.getCStr (),
        dateValue.toString (FD_GENERAL_FULL).getCStr ());
}

void ncACSUtil::md5_encrypt ()
{
    String value = GetString (_T("value"));

    String str1 = genMD5String (value);
    String str2 = genMD5String2 (value);

    printMessage2 (_T("value: %s, md5str: %s,%s"), value.getCStr (),
        str1.getCStr (), str2.getCStr ());
}

void ncACSUtil::base64_encrypt ()
{
    String value = GetString (_T("value"));
    Base64Encode(toSTLString(value));
}

void ncACSUtil::base64_decrypt ()
{
    String value = GetString (_T("value"));
    string tmp = toSTLString(value);
    boost::replace_all(tmp, "\\\\", "\\");
    boost::replace_all(tmp, "\\t",  "\t");
    boost::replace_all(tmp, "\\n",  "\n");

    Base64Decode(tmp);
}

void ncACSUtil::rsa_encrypt ()
{
    String value = GetString (_T("value"));
    string plainText = toSTLString(value);

    string cipherText (RSAEncrypt(plainText, "test_pub.key"));
    string encodeText (Base64Encode (cipherText));

    // escape 回车符
    String printstr;
    for (size_t i = 0; i < encodeText.length (); ++i) {
        if (encodeText[i] == _T('\n')) {
            printstr.append (_T("\\n"));
        }
        else {
            printstr.append (encodeText[i], 1);
        }
    }

    printMessage2 (_T("base64(rsa(%s)): %s"), value.getCStr (), printstr.getCStr ());
}

void ncACSUtil::rsa_decrypt ()
{
    String value = GetString (_T("value"));

    string tmp = toSTLString(value);
    boost::replace_all(tmp, "\\\\", "\\");
    boost::replace_all(tmp, "\\t",  "\t");
    boost::replace_all(tmp, "\\n",  "\n");
    boost::replace_all(tmp, "\\r",  "\r");

    string cipherText(Base64Decode(tmp));
    string plainText(RSADecrypt(cipherText, "test.key"));

    // escape 回车符
    String printstr;
    for (size_t i = 0; i < plainText.length (); ++i) {
        if (plainText[i] == _T('\n')) {
            printstr.append (_T("\\n"));
        }
        else {
            printstr.append (plainText[i], 1);
        }
    }

    printMessage2 (_T("rsa'(base64'(%s)): %s"), value.getCStr (), printstr.getCStr ());
}

void ncACSUtil::des_encrypt()
{
    String value = GetString("value");
    string plainText = toSTLString(value);

    Base64Encode(DESEncrypt(plainText));
}

void ncACSUtil::des_decrypt ()
{
    String value = GetString("value");
    string encodeText = toSTLString(value);

    DESDecrypt(Base64Decode(encodeText));
}

/*
 * RSA解密线程
 */
class ncDecryptThread : public Thread
{
public:
    ncDecryptThread ()
    {
    }

    ~ncDecryptThread ()
    {
    }

public:
    virtual void run ();
};

void ncDecryptThread::run ()
{
    while (1) {
/*        string password = "WFOLzdG0DN605Ebq7XOKd4MZnsy2jC50ygEzp6Bf8HpVfiD74HD9ufguK3wMR6SP\nkGkzzfrsJYP5+2EnPU3KM3TnBaQVw0fL8O79zQ4hrNc20rPv7mC1+ZBVwsM6h/nt\npeiRiEMflQ4tBFpZg12ktnBnxHSRkLjHauiTTQ7FDs0=\n";
        string decodePwd (Base64Decode (password));
        string originPassword (RSADecrypt(decodePwd.c_str (), "test.key"));
        string noPaddingPwd (originPassword.data ());*/

        //printMessage2(_T("%s"), noPaddingPwd.c_str());
    }
}

void ncACSUtil::mtdecrypt ()
{
    vector<ncDecryptThread*> threads;
    for(int i = 0; i < 10; ++i) {
        ncDecryptThread* decryptThread = new ncDecryptThread ();
        decryptThread->start();

        threads.push_back(decryptThread);
    }

    for(int i = 0; i < 10; ++i) {
        threads[i]->join();
    }
}

void ncACSUtil::getEndOfDayTicks()
{
    struct tm * timeinfo;

    int64 ticks1 = BusinessDate::getCurrentTime ();

    time_t rt = ticks1 / 1000000;
    timeinfo = localtime ( &rt );

    printMessage2 (_T("CurTime: %s, %s"), String::toString (ticks1).getCStr (), asctime (timeinfo));

    timeinfo->tm_hour = 23;
    timeinfo->tm_min = 59;
    timeinfo->tm_sec = 59;

    time_t endtt = mktime (timeinfo);
    int64 ticks2 = endtt * (int64)1000000;

    printMessage2 (_T("EndOfDay: %s, %s"), String::toString (ticks2).getCStr (), asctime (timeinfo));
}

void ncACSUtil::printByte (const String& str)
{
    printf ("%s = 0x", str.getCStr ());

    unsigned char* pstart = (unsigned char*)str.getCStr ();
    for (size_t i = 0; i < str.getLength (); ++i) {
        printf ("%02x", *pstart);
        pstart++;
    }
    printf ("\n");
}

void ncACSUtil::testString ()
{
    // 存储char
    char charStr[128] = "hello";
    String str (charStr);
    printByte (str);

    // 存储unicode
    char unicodeStr[128] = "\xa4\x7f\xc4\x7e\x30\x00\x31\x00";
    str = unicodeStr;
    printByte (str);

    // 存储utf-8
    char utf8Str[128] = "\xe7\xbe\xa4\xe7\xbb\x84\x30\x31";
    str = utf8Str;
    printByte (str);
}

void ncACSUtil::getParentString ()
{
    vector<String> strs;

    // 获得 gns://abc
    strs.push_back (_T("gns://abc"));
    strs.push_back (_T("gns://abc/def"));
    strs.push_back (_T("gns://abc/def/ghi"));

    // 获取 gns://测试
    strs.push_back (_T("gns://测试"));
    strs.push_back (_T("gns://测试/01"));
    strs.push_back (_T("gns://测试/哈哈"));
    strs.push_back (_T("gns://测试/哈哈/01"));

    // 获取gns://test1/test2
    strs.push_back (_T("gns://test1/test2/测试"));
    strs.push_back (_T("gns://test1/test2"));

    // 先排序
    sort (strs.begin (), strs.end ());

    for (size_t i = 0 ; i < strs.size (); ++i) {
        printMessage2 (strs[i].getCStr ());
    }

    // 在找出所有的parent目录
    vector<String> resultStrs;
    for (size_t i = 0; i < strs.size (); ++i) {
        if (i == 0) {
            resultStrs.push_back (strs[i]);
        }
        else {
            size_t index = resultStrs.size () - 1;
            String lastStr = resultStrs[index];

            printMessage2 (_T("%s - %s"), lastStr.getCStr (), strs[i].getCStr ());
            if (strs[i].find (lastStr) == String::NO_POSITION) {
                resultStrs.push_back (strs[i]);
            }
        }
    }

    for (size_t i = 0 ; i < resultStrs.size (); ++i) {
        printMessage2 (resultStrs[i].getCStr ());
    }
}

void ncACSUtil::getUTF8StringSize ()
{
    String utf8String = GetString (_T("str"));

    char* str = (char*)utf8String.getCStr ();

    size_t index = 0;
    size_t count = 0;

    // 默认传进来的是大端编码，低地址存的高位的值
    while (str[index]) {
        if ((str[index] & 0xc0) != 0x80) {
            count++;
        }

        index++;
    }

    printMessage2 (_T("count = %d"), count);
}

// 最终插入的表信息
struct ncDocInfo {
    String        docId;
    int            docType;
    String        typeName;
    String        objId;
    String        name;
    String        createrId;
};

void ncACSUtil::updateDB ()
{
    String ip = GetString (_T("ip"));
    int anysharePort = GetInt (_T("anysharePort"));
    int sharemgntPort = GetInt (_T("sharemgntPort"));

    if (anysharePort < 0) {
        throw Exception (_T("Invalid anyshare port: %d"), anysharePort);
    }
    if (sharemgntPort < 0) {
        throw Exception (_T("Invalid sharemgnt port: %d"), sharemgntPort);
    }

    // 连接anyshare数据库
    ncDBConnectionInfo connInfo;

    printf (_T("connect anyshare db:"));
    nsresult ret;
    nsCOMPtr<ncIDBOperator> anyshareDBOperator = do_CreateInstance (NC_MYSQL_OPERATOR_CID, &ret);
    if (NS_FAILED (ret)) {
        String error;
        error.format (_T("Failed to create eacp db operator: 0x%x"), ret);

        throw Exception (error);
    }

    connInfo.ip = ip;
    connInfo.port = anysharePort;
    connInfo.user = _T("DataEngine");
    connInfo.password = _T("DataEngineMysqlClient");
    connInfo.db = _T("anyshare");
    anyshareDBOperator->Connect (connInfo);

    printf (_T("success\n"));

    // 连接sharemgnt数据库
    printf (_T("connect sharemgnt db:"));
    nsCOMPtr<ncIDBOperator> sharemgntDBOperator = do_CreateInstance (NC_MYSQL_OPERATOR_CID, &ret);
    if (NS_FAILED (ret)) {
        String error;
        error.format (_T("Failed to create sharemgnt db operator: 0x%x"), ret);

        throw Exception (error);
    }
    connInfo.ip = ip;
    connInfo.port = sharemgntPort;
    connInfo.user = _T("DataEngine");
    connInfo.password = _T("DataEngineMysqlClient");
    connInfo.db = _T("sharemgnt_db");
    sharemgntDBOperator->Connect (connInfo);

    printf (_T("success\n"));

    // 备份anyshare数据库
    printf (_T("backup anyshare db:"));
    String backupCmd;
    String fileName;
    fileName.format (_T("anyshare_bak_%s.sql"), String::toString (BusinessDate::getCurrentTime ()).getCStr ());
    backupCmd.format (_T("mysqldump -h%s -P%d -uDataEngine -pDataEngineMysqlClient anyshare > %s"),
        ip.getCStr (), anysharePort, fileName.getCStr ());

    if (system (backupCmd.getCStr ()) != 0) {
        // 删除文件
        File tmpFile (fileName);
        tmpFile.remove ();

        throw Exception (_T("back up anyshare database error"));
    }

    printf (_T("file: %s, success\n"), fileName.getCStr ());

    // 获取所有文档信息
    printf (_T("get old doc info:"));

    String strSql;
    strSql.format (_T("select f_doc_id,f_obj_id,f_obj_type from t_acs_doc"));

    ncDBRecords records;
    anyshareDBOperator->Select (strSql, records);

    printf (_T("success\n"));

    ncDocInfo tmpInfo;
    vector<ncDocInfo> docInfos;

    // 处理用户文档
    printf (_T("process old user doc info:"));
    for (size_t i = 0; i < records.size (); ++i) {
        if (records[i][2] == _T("1")) {

            tmpInfo.docId = records[i][0];
            tmpInfo.docType = 1;
            tmpInfo.typeName = _T("ACS_USER_DOC");
            tmpInfo.objId = records[i][1];
            tmpInfo.name = _T("");
            tmpInfo.createrId = records[i][1];

            docInfos.push_back (tmpInfo);
        }
    }

    printf (_T("success\n"));

    // 获取admin的id
    printf (_T("process old department doc info:"));
    strSql.format (_T("select f_user_id from t_user where f_login_name = 'admin';"));

    ncDBRecords tmpRecords1;
    sharemgntDBOperator->Select (strSql, tmpRecords1);

    if (tmpRecords1.size () != 1) {
        throw Exception (_T("Get admin id error"));
    }

    String adminId = tmpRecords1[0][0];

    for (size_t i = 0; i < records.size (); ++i) {
        if (records[i][2] == _T("3")) {

            tmpInfo.docId = records[i][0];
            tmpInfo.docType = 3;
            tmpInfo.typeName = _T("部门文档");
            tmpInfo.objId = records[i][1];

            strSql.format (_T("select f_name from t_department where f_department_id = '%s';"),
                records[i][1].getCStr ());

            ncDBRecords tmpRecords;
            sharemgntDBOperator->Select (strSql, tmpRecords);

            if (tmpRecords.size () != 1) {
                String error;
                error.format (_T("Get department %s error"), records[i][1].getCStr ());
                throw Exception (error);
            }

            tmpInfo.name = tmpRecords[0][0];
            tmpInfo.createrId = adminId;

            docInfos.push_back (tmpInfo);
        }
    }

    printf (_T("success\n"));

    printf (_T("rename t_acs_doc->t_acs_doc_obsolete:"));
    // 重命名旧的t_acs_doc->t_acs_doc_obsolete
    strSql.format (_T("alter table t_acs_doc rename to t_acs_doc_obselete;"));
    anyshareDBOperator->Execute (strSql);
    printf (_T("success\n"));

    printf (_T("rename t_acs_group_share->t_acs_group_share_obsolete:"));
    // 重命名旧的t_acs_group_share->t_acs_group_share_obsolete
    strSql.format (_T("alter table t_acs_group_share rename to t_acs_group_share_obsolete;"));
    anyshareDBOperator->Execute (strSql);
    printf (_T("success\n"));

    // 新建新的表t_acs_doc
    printf (_T("create new t_acs_doc table:"));
    strSql = _T("CREATE TABLE IF NOT EXISTS `t_acs_doc` (") \
        _T("`f_doc_id` char(40) NOT NULL,")    \
        _T("`f_doc_type` tinyint(4) NOT NULL,") \
        _T("`f_type_name` char(128) NOT NULL,") \
        _T("`f_obj_id` char(40) NOT NULL,") \
        _T("`f_name` char(128) NOT NULL,") \
        _T("`f_create_time` bigint(20) NOT NULL DEFAULT '0',") \
        _T("`f_creater_id` char(40) NOT NULL,") \
        _T("PRIMARY KEY (`f_doc_id`),") \
        _T("KEY `t_doc_f_doc_id_index` (`f_doc_id`) USING BTREE,") \
        _T("KEY `t_doc_f_doc_type_index` (`f_doc_type`) USING BTREE,") \
        _T("KEY `t_doc_f_obj_id_index` (`f_obj_id`) USING BTREE,") \
        _T("KEY `t_doc_f_name_index` (`f_name`) USING BTREE") \
        _T(") ENGINE=InnoDB DEFAULT CHARSET=utf8;");
    anyshareDBOperator->Execute (strSql);
    printf (_T("success\n"));

    // 写入最终的文档信息
    printMessage2 (_T("write new doc info begin"));
    for (size_t i = 0; i < docInfos.size (); ++i) {
        printMessage2 (_T("%s_%d_%s_%s_%s_%s"),
            docInfos[i].docId.getCStr (), docInfos[i].docType, docInfos[i].typeName.getCStr (),
            docInfos[i].objId.getCStr (), docInfos[i].name.getCStr (), docInfos[i].createrId.getCStr ());

        strSql.format (_T("insert into t_acs_doc (f_doc_type, f_type_name, f_doc_id, f_obj_id, f_name, f_creater_id) values (%d, '%s', '%s', '%s', '%s', '%s')"),
            docInfos[i].docType, docInfos[i].typeName.getCStr (), docInfos[i].docId.getCStr (),
            docInfos[i].objId.getCStr (), docInfos[i].name.getCStr (), docInfos[i].createrId.getCStr ());

        anyshareDBOperator->Execute (strSql);
    }
    printMessage2 (_T("write new doc info success"));
}

/*
 * 事务线程
 */
class ncTransactionThread : public Thread
{
public:
    ncTransactionThread (ncIDBOperator* dbOperator)
        : _dbOperator (dbOperator)
    {
    }

    ~ncTransactionThread ()
    {
    }

public:
    virtual void run ();

private:
    nsCOMPtr<ncIDBOperator>    _dbOperator;
};

void ncTransactionThread::run ()
{
    while (1) {

        printMessage2 (_T("trancation begin"));

        _dbOperator->StartTransaction ();

        String strSql;
        for (size_t i = 0; i < 10000; ++i) {
            strSql.format (_T("update t_acs_access_token set f_last_request_time = 123 where f_token_id = 'user001';"));

            _dbOperator->Execute (strSql);
        }

        _dbOperator->Commit ();

        printMessage2 (_T("trancation end"));

        sleep (1*1000);
    }
}

/*
 * 普通查询线程
 */
class ncQueryThread : public Thread
{
public:
    ncQueryThread (ncIDBOperator* writeOperator, ncIDBOperator* readOperator)
        : _writeOperator (writeOperator),
        _readOperator (readOperator)
    {
    }

    ~ncQueryThread ()
    {
    }

public:
    virtual void run ();

private:
    nsCOMPtr<ncIDBOperator>    _writeOperator;
    nsCOMPtr<ncIDBOperator>    _readOperator;
};

void ncQueryThread::run ()
{
    while (1) {

        Guid guid;
        String strGUID = guid.toString ();
        int64 createTime = BusinessDate::getCurrentTime ();

        String escUserId;
        _writeOperator->Escape (_T("user01"), escUserId);

        String strSql;
        strSql.format (_T("insert into t_acs_access_token (f_token_id, f_user_id, f_create_time, f_last_request_time) values ('%s', '%s', %s, %s)"),
            strGUID.getCStr (), escUserId.getCStr (),
            String::toString (createTime).getCStr (),
            String::toString (createTime).getCStr ());

        _writeOperator->Execute (strSql);

        ncDBRecords tmpRecords;
        strSql.format (_T("select * from t_acs_access_token where f_token_id ='%s';"), strGUID.getCStr ());
        _readOperator->Select (strSql, tmpRecords);

        if (tmpRecords.size () == 1) {
            printf (_T("%s_%s_%s\n"),
                tmpRecords[0][0].getCStr (), tmpRecords[0][1].getCStr (), tmpRecords[0][2].getCStr ());
        }
        else {
            printf (_T("error\n"));
        }

        sleep (1000);
    }
}

void ncACSUtil::testdb ()
{
    String ip = GetString (_T("ip"));
    int port = GetInt (_T("port"));

    // 连接anyshare数据库
    ncDBConnectionInfo connInfo;
    nsresult ret;
    nsCOMPtr<ncIDBOperator> writeOperator = do_CreateInstance (NC_MYSQL_OPERATOR_CID, &ret);
    if (NS_FAILED (ret)) {
        String error;
        error.format (_T("Failed to create eacp db operator: 0x%x"), ret);

        throw Exception (error);
    }

    nsCOMPtr<ncIDBOperator> readOperator = do_CreateInstance (NC_MYSQL_OPERATOR_CID, &ret);
    if (NS_FAILED (ret)) {
        String error;
        error.format (_T("Failed to create eacp db operator: 0x%x"), ret);

        throw Exception (error);
    }

    connInfo.ip = ip;
    connInfo.port = port;
    connInfo.user = _T("DataEngine");
    connInfo.password = _T("DataEngineMysqlClient");
    connInfo.db = _T("anyshare");
    writeOperator->Connect (connInfo);
    readOperator->Connect (connInfo);

    ncQueryThread* pQueryThread = new ncQueryThread (writeOperator, readOperator);
    pQueryThread->start ();

    ncTransactionThread* pTrancThread = new ncTransactionThread (writeOperator);
    pTrancThread->start ();

    pQueryThread->join ();
    pTrancThread->join ();
}

void ncACSUtil::conndb ()
{
    String ip = GetString (_T("ip"));
    int port = GetInt (_T("port"));

    MYSQL* _mysql = NULL;

    // 初始化mysql环境
    if ((_mysql = ::mysql_init (NULL)) == NULL) {
        THROW_E ("conndb", 1, _T("%s(%d)"), mysql_error(_mysql), mysql_errno (_mysql));
    }

    unsigned int time1 = 30, time2 = 60, time3 = 60;

    // 设置连接超时时间.
    mysql_options (_mysql, MYSQL_OPT_CONNECT_TIMEOUT, (const char *)&time1);

    // 设置查询数据库(select)超时时间
    mysql_options (_mysql, MYSQL_OPT_READ_TIMEOUT, (const char *)&(time2));

    // 设置写数据库(update,delect,insert,replace等)的超时时间。
    mysql_options (_mysql, MYSQL_OPT_WRITE_TIMEOUT, (const char *)&time3);

    // 设置 utf8 字符集
    //::mysql_set_character_set (_mysql, "utf8");

    if (::mysql_real_connect (_mysql, ip.getCStr(),
                                "AnyShareAdmin", "Any1Share2Mysql3Client4",
                                "anyshare", port, NULL, 0) == NULL) {
        THROW_E ("conndb", 1, _T("%s(%d)"), mysql_error(_mysql), mysql_errno (_mysql));
    }

    printMessage2 (_T("connect %s:%d success"), ip.getCStr(), port);
}

void ncACSUtil::isValidURL ()
{
    validateURL ("http://www.abc.com/a1");
}

void ncACSUtil::validateURL (const String& str)
{
    string regularStr = "(http|https):\\/\\/(\\w+\\.)*(\\w*)\\/([\\w\\d]+\\/{0,1})+";
    printMessage2 (regularStr.c_str ());

    boost::regex re (regularStr);
    if (boost::regex_match (toSTLString (str), re)) {
        printMessage2 (_T("%s: valid url"), str.getCStr ());
    }
    else {
        printMessage2 (_T("%s: invalid url"), str.getCStr ());
    }
}

void ncACSUtil::createuserdoc ()
{
    printMessage2 (_T("createuserdoc"));

    // 读取参数
    String ip = GetString (_T("ip"));
    int anysharePort = GetInt (_T("anysharePort"));
    int sharemgntPort = GetInt (_T("sharemgntPort"));

    if (ip.isEmpty () == true) {
        throw Exception (_T("Invalid null ip."));
    }

    if (anysharePort < 0) {
        throw Exception (_T("Invalid anyshare port: %d"), anysharePort);
    }
    if (sharemgntPort < 0) {
        throw Exception (_T("Invalid sharemgnt port: %d"), anysharePort);
    }

    // 连接anyshare数据库
    ncDBConnectionInfo connInfo;

    printf (_T("connect anyshare db:"));
    nsresult ret;
    nsCOMPtr<ncIDBOperator> anyshareDBOperator = do_CreateInstance (NC_MYSQL_OPERATOR_CID, &ret);
    if (NS_FAILED (ret)) {
        String error;
        error.format (_T("Failed to create eacp db operator: 0x%x"), ret);

        throw Exception (error);
    }

    connInfo.ip = ip;
    connInfo.port = anysharePort;
    connInfo.user = _T("DataEngine");
    connInfo.password = _T("DataEngineMysqlClient");
    connInfo.db = _T("anyshare");
    anyshareDBOperator->Connect (connInfo);

    printf (_T("success\n"));

    // 连接sharemgnt数据库
    printf (_T("connect sharemgnt db:"));
    nsCOMPtr<ncIDBOperator> sharemgntDBOperator = do_CreateInstance (NC_MYSQL_OPERATOR_CID, &ret);
    if (NS_FAILED (ret)) {
        String error;
        error.format (_T("Failed to create sharemgnt db operator: 0x%x"), ret);

        throw Exception (error);
    }

    connInfo.ip = ip;
    connInfo.port = sharemgntPort;
    connInfo.user = _T("DataEngine");
    connInfo.password = _T("DataEngineMysqlClient");
    connInfo.db = _T("sharemgnt_db");
    sharemgntDBOperator->Connect (connInfo);

    printf (_T("success\n"));

    // 读取t_user
    printf (_T("get t_user info:"));

    String strSql;
    strSql.format (_T("select f_user_id,f_login_name from t_user;"));

    ncDBRecords records;
    sharemgntDBOperator->Select (strSql, records);

    printf (_T("success\n"));

    // 连接eacp.thrift进行创建个人文档
    try {
        // 遍历所有用户，如果没有创建个人文档，就创建个人文档
        int64 defaultQuotaBytes = 5368709120;
        for (size_t i = 0; i < records.size (); ++i) {
            printMessage2 (_T("%s_%s"),
                records[i][0].getCStr (), records[i][1].getCStr ());

            String userId = records[i][0];
            String loginName = records[i][1];

            // 创建个人文档
            strSql.format (_T("select f_doc_id from t_acs_doc where f_obj_id = '%s';"), records[i][0].getCStr ());

            ncDBRecords existRecords;;
            anyshareDBOperator->Select (strSql, existRecords);
            if (existRecords.size () == 0) {
                printMessage2 (_T("not exists, need to create user doc."));
            }
            else {
                printMessage2 (_T("exist."));
            }
        }
    }
    catch (ncTException & e) {
        String error;
        error.format (_T("Failed to set quota: %s"), e.expMsg.c_str ());
        throw Exception (error);
    }
    catch (TException & e) {
        String error;
        error.format (_T("Failed to set quota: %s"), e.what ());
        throw Exception (error);
    }
}

void ncACSUtil::setquota ()
{
    printMessage2 (_T("setquota"));

    // 读取参数
    String ip = GetString (_T("ip"));
    int anysharePort = GetInt (_T("anysharePort"));

    if (anysharePort < 0) {
        throw Exception (_T("Invalid anyshare port: %d"), anysharePort);
    }

    // 连接anyshare数据库
    ncDBConnectionInfo connInfo;

    printf (_T("connect anyshare db:"));
    nsresult ret;
    nsCOMPtr<ncIDBOperator> anyshareDBOperator = do_CreateInstance (NC_MYSQL_OPERATOR_CID, &ret);
    if (NS_FAILED (ret)) {
        String error;
        error.format (_T("Failed to create eacp db operator: 0x%x"), ret);

        throw Exception (error);
    }

    connInfo.ip = ip;
    connInfo.port = anysharePort;
    connInfo.user = _T("DataEngine");
    connInfo.password = _T("DataEngineMysqlClient");
    connInfo.db = _T("anyshare");
    anyshareDBOperator->Connect (connInfo);

    printf (_T("success\n"));

    // 读取t_acs_doc
    printf (_T("get t_acs_doc info:"));

    String strSql;
    strSql.format (_T("select f_doc_id,f_doc_type,f_creater_id from t_acs_doc where f_doc_type != 3 order by f_creater_id;"));

    ncDBRecords records;
    anyshareDBOperator->Select (strSql, records);

    printf (_T("success\n"));

    vector<ncQuotaGroup> quotaGroups;

    String lastCreaterId;
    ncQuotaGroup lastQuotaGroup;
    for (size_t i = 0; i < records.size (); ++i) {
        String curCreaterId = records[i][2];

        if (curCreaterId != lastCreaterId) {
            if (lastCreaterId != String::EMPTY) {
                quotaGroups.push_back (lastQuotaGroup);

                lastQuotaGroup.reset ();
            }
            lastCreaterId = curCreaterId;

            printMessage2 (_T("-----------------------------------"));
        }

        if (records[i][1] == _T("1")) {
            lastQuotaGroup.userDocId = records[i][0];
        }
        else {
            String error;
            error.format (_T("invalid doc type: %s"), records[i][1].getCStr ());
            throw Exception (_T("invalid doc type"));
        }

        printMessage2 (_T("%s_%s_%s"), records[i][0].getCStr (), records[i][1].getCStr (), records[i][2].getCStr ());
    }

    executeQuotaProcess (quotaGroups, ip);
}

void ncACSUtil::executeQuotaProcess (const vector<ncQuotaGroup> quotaGroups, const String& ip)
{
    for (size_t i = 0; i < quotaGroups.size (); ++i) {
        printMessage2 (_T("userdocid: %s"), quotaGroups[i].userDocId.getCStr ());

        if (i != quotaGroups.size () - 1) {
            printMessage2 (_T("-------------------------"));
        }
    }

    int64 userDocQuotaBytes = 5368709120;    // 个人文档配额设置为5G

    try {
        ncThriftClient<ncTEVFSClient> evfsClient (ip, g_EVFS_constants.NCT_EVFS_PORT);

        for (size_t i = 0; i < quotaGroups.size (); ++i) {

            // 设置用户文档配额
            evfsClient->SetQuotaInfo (toSTLString (quotaGroups[i].userDocId), userDocQuotaBytes);

            printMessage2 (_T("set user quta success: %s"), quotaGroups[i].userDocId.getCStr ());
        }
    }
    catch (ncTException & e) {
        String error;
        error.format (_T("Failed to set quota: %s"), e.expMsg.c_str ());
        throw Exception (error);
    }
    catch (TException & e) {
        String error;
        error.format (_T("Failed to set quota: %s"), e.what ());
        throw Exception (error);
    }
    catch (Exception & e) {
        String error;
        error.format (_T("Failed to set quota: %s"), e.toString ().getCStr ());
        throw Exception (error);
    }
}

void ncACSUtil::getUsage ()
{
    printMessage2 (_T("Usage:"));
    printMessage2 (_T("./acs_util \"cmd=ptime&value=$value\""));
    printMessage2 (_T("./acs_util \"cmd=ticks2String\""));
    printMessage2 (_T("./acs_util \"cmd=getmd5&value=$value\""));
    printMessage2 (_T("./acs_util \"cmd=encrypt&value=$value\""));
    printMessage2 (_T("./acs_util \"cmd=base64_encrypt&value=$value\""));
    printMessage2 (_T("./acs_util \"cmd=base64_decrypt&value=$value\""));
    printMessage2 (_T("./acs_util \"cmd=rsa_encrypt&value=$value\""));
    printMessage2 (_T("./acs_util \"cmd=rsa_decrypt&value=$value\""));
    printMessage2 (_T("./acs_util \"cmd=des_encrypt&value=$value\""));
    printMessage2 (_T("./acs_util \"cmd=des_decrypt&value=$value\""));
    printMessage2 (_T("./acs_util \"cmd=testString\""));
    printMessage2 (_T("./acs_util \"cmd=getParentString\""));
    printMessage2 (_T("./acs_util \"cmd=getUTF8StringSize&str=$str\""));
    printMessage2 (_T("./acs_util \"cmd=updateDB&ip=$ip&anysharePort=$port&sharemgntPort=$port\""));
    printMessage2 (_T("./acs_util \"cmd=testdb\""));
    printMessage2 (_T("./acs_util \"cmd=createuserdoc&ip=$ip&anysharePort=$port&sharemgntPort=$port\""));
    printMessage2 (_T("./acs_util \"cmd=setquota&ip=$ip&anysharePort=$port\""));
}

string ncACSUtil::RSAEncrypt (const string& plainText, const char *path_key)
{
    string cipherText;
    RSA *p_rsa;
    FILE *file;
    int rsa_len = 0;

    if ((file = fopen (path_key,"r")) == NULL) {
        throw Exception (_T("open key file error"));
    }
    if ((p_rsa = PEM_read_RSA_PUBKEY (file, NULL, NULL, NULL)) == NULL) {
        throw Exception (_T("PEM_read_RSA_PUBKEY error"));
    }

    rsa_len = RSA_size (p_rsa);
    cipherText.resize (rsa_len, 0);

    printf("Before rsa encrypt:\n");
    printBytes(plainText);

    int num = RSA_public_encrypt (plainText.length(), (unsigned char *)plainText.c_str(), (unsigned char*)cipherText.c_str(), p_rsa, RSA_PKCS1_PADDING);
    if (num != rsa_len){
        RSA_free (p_rsa);
        fclose (file);
        throw Exception (_T("RSA_public_encrypt error"));
    }

    printf("After rsa decrypt:\n");
    printBytes(cipherText);

    RSA_free (p_rsa);
    fclose (file);

    return cipherText;
}

string ncACSUtil::RSADecrypt (const string& cipherText, const char *path_key)
{
    static ThreadMutexLock sLock;
    static RSA* p_rsa = NULL;
    AutoLock<ThreadMutexLock> lock (&sLock);
    if(p_rsa == NULL) {
        FILE *file;
        if ((file = fopen (path_key,"r")) == NULL) {
            throw Exception(_T("fopen error"));
        }
        if ((p_rsa = PEM_read_RSAPrivateKey (file, NULL, NULL, NULL)) == NULL) {
            fclose (file);
            throw Exception(_T("PEM_read_RSAPrivateKey error"));
        }
        fclose (file);
    }

    string plainText;
    int rsa_len = RSA_size (p_rsa);
    plainText.resize (rsa_len, 0);

    printf("Before rsa encrypt:\n");
    printBytes(cipherText);

    if (RSA_private_decrypt (cipherText.length(), (unsigned char *)cipherText.c_str(), (unsigned char*)plainText.c_str(), p_rsa, RSA_PKCS1_PADDING)<0) {
        throw Exception(_T("RSA_private_decrypt error"));
    }

    string ret(plainText.c_str());
    printf("After rsa decrypt:\n");
    printBytes(ret);

    return ret;
}

string ncACSUtil::Base64Encode(const string& input)
{
    printf("Before Base64Encode:\n");
    printBytes(input);

    BIO * bmem = NULL;
    BIO * b64 = NULL;
    BUF_MEM * bptr = NULL;

    b64 = BIO_new(BIO_f_base64());
    bmem = BIO_new(BIO_s_mem());
    b64 = BIO_push(b64, bmem);
    BIO_write(b64, (char*)input.c_str(), input.length());
    BIO_flush(b64);
    BIO_get_mem_ptr(b64, &bptr);

    string buffer;
    buffer.assign(bptr->data, bptr->length);

    BIO_free_all(b64);

    printf("After Base64Encode: \n%s", buffer.c_str());
    printBytes(buffer);

    return buffer;
}

size_t ncACSUtil::calcDecodeLength(const string& tmp)
{
    // Calculates the length of a decoded string
    size_t len = tmp.length();
    size_t padding = 0;

    if (tmp[len-1] == '=' && tmp[len-2] == '=') //last two chars are =
        padding = 2;
    else if (tmp[len-1] == '=') //last char is =
        padding = 1;

    return (len*3)/4 - padding;
}

string ncACSUtil::Base64Decode(const string& input)
{
    printf("Before Base64Decode:\n");
    printBytes(input);

    // 去除所有的\n
    string tmp = input;
    boost::replace_all(tmp, "\n", "");
    boost::replace_all(tmp, "\r", "");

    string buffer;
    buffer.resize(calcDecodeLength(tmp));

    BIO * b64 = NULL;
    BIO * bmem = NULL;

    // 使用没有\n的密文进行解码
    b64 = BIO_new(BIO_f_base64());
    BIO_set_flags(b64, BIO_FLAGS_BASE64_NO_NL);
    bmem = BIO_new_mem_buf((char*)tmp.c_str(), tmp.length());
    bmem = BIO_push(b64, bmem);
    BIO_read(bmem, (char*)buffer.c_str(), tmp.length());
    BIO_free_all(bmem);

    printf("After Base64Decode:\n");
    printBytes(buffer);

    return buffer;
}

string ncACSUtil::DESEncrypt(const string& plainText)
{
    string cipherText;
    int cipherTextLength = 0;
    if (plainText.length() % 8 == 0) {
        cipherTextLength = plainText.length();
    }
    else {
        cipherTextLength = plainText.length() + (8 - plainText.length() % 8);
    }
    cipherText.assign(cipherTextLength, '*');

    DES_cblock key = {'E', 't', 'm', 'B', '8', '?', 's', 'f'};
    DES_cblock ivec = {'E', 't', 'm', 'B', '8', '?', 's', 'f'};
    DES_key_schedule keysched;

    DES_set_odd_parity(&key);
    if (DES_set_key_checked((C_Block *)key, &keysched)) {
        throw Exception("Unable to set key schedule");
    }

    printf("Before des encrypt:\n");
    printBytes(plainText);

    DES_ncbc_encrypt((unsigned char*)plainText.c_str(), (unsigned char*)cipherText.c_str(), plainText.length(), &keysched, &ivec, DES_ENCRYPT);

    printf("After des encrypt:\n");
    printBytes(cipherText);

    return cipherText;
}

string ncACSUtil::DESDecrypt(const string& cipherText)
{
    string plainText;
    plainText.assign(cipherText.length(), '*');

    DES_cblock key = {'E', 't', 'm', 'B', '8', '?', 's', 'f'};
    DES_cblock ivec = {'E', 't', 'm', 'B', '8', '?', 's', 'f'};
    DES_key_schedule keysched;

    DES_set_odd_parity(&key);
    if (DES_set_key_checked((C_Block *)key, &keysched)) {
        throw Exception("Unable to set key schedule");
    }

    printf("Before des decrypt:\n");
    printBytes(cipherText);

    DES_ncbc_encrypt((unsigned char*)cipherText.c_str(), (unsigned char*)plainText.c_str(), cipherText.length(), &keysched, &ivec, DES_DECRYPT);

    // 去除末尾的\0
    string ret(plainText.c_str());
    printf("After des decrypt:\n");
    printBytes(ret);

    return ret;
}

void ncACSUtil::printBytes(const string& str)
{
    for(int i = 0; i < str.length(); ++i) {
        if((unsigned char)str[i] == 0) {
            printf("[*00*] ");
        }
        else {
            printf("[%02x] ", (unsigned char)str[i]);
        }
    }
    printf("\n");
}
