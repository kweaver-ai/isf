#include <abprec.h>
#include <gmock/gmock.h>
#include "ncACSShareMgntUT.h"
#include <dboperator/public/ncIDBOperator.h>
#include <boost/date_time/posix_time/posix_time.hpp>

using namespace testing;

static ncIDBOperator* GetShareMgntDBOperator ()
{
    // 创建新的数据库连接
    nsresult ret;
    nsCOMPtr<ncIDBOperator> dbOper = do_CreateInstance (NC_MYSQL_OPERATOR_CID, &ret);
    if (NS_FAILED (ret)) {
        String error;
        error.format (_T("Failed to create db operator: 0x%x"), ret);
        throw Exception (error);
    }

    ncDBConnectionInfo connInfo;
    connInfo.ip = _T("127.0.0.1");
    connInfo.port = 3306;
    connInfo.user = _T("root");
    connInfo.password = _T("eisoo.com");
    connInfo.db = _T("sharemgnt_db");

    dbOper->Connect (connInfo);

    NS_ADDREF (dbOper.get ());
    return dbOper;
}

ncIDBOperator* ncACSShareMgntTested::GetDBOperator ()
{
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    NS_ADDREF (dbOper.get ());
    return dbOper;
}

ncACSShareMgntUT::ncACSShareMgntUT ()
{
}

ncACSShareMgntUT::~ncACSShareMgntUT ()
{
}

void ncACSShareMgntUT::SetUp ()
{
    _manager = new ncACSShareMgntTested ();
}

void ncACSShareMgntUT::TearDown ()
{
}

TEST_F (ncACSShareMgntUT, do_CreateInstance)
{
    nsresult ret;
    nsCOMPtr<ncIACSShareMgnt> manager = do_CreateInstance (NC_ACS_SHARE_MGNT_CONTRACTID, &ret);

    if (NS_FAILED (ret)) {
        printMessage2 (_T("Failed to create instance: 0x%x"), ret);
        ASSERT_EQ (1, 0);
    }
}

TEST_F (ncACSShareMgntUT, GetUserInfoByIdBatch)
{
    // 插入2条用户信息，再去获取
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_user(f_user_id,f_login_name,f_display_name,f_mail_address) values ")
                    _T("('userid1', 'account1', 'name1', 'address1'),")
                    _T("('userid2', 'account2', 'name2', 'address2');"));

    dbOper->Execute (strSql);

    vector<String> userIds;
    userIds.push_back (_T("userid1"));
    userIds.push_back (_T("userid2"));
    map<String, ncACSUserInfo> userInfoMap;

    ASSERT_NO_THROW (_manager->GetUserInfoByIdBatch (userIds, userInfoMap));

    map<String, ncACSUserInfo>::iterator iter1 = userInfoMap.find (_T("userid1"));
    ASSERT_NE (iter1, userInfoMap.end ());
    ASSERT_EQ (iter1->second.id, _T("userid1"));
    ASSERT_EQ (iter1->second.account, _T("account1"));
    ASSERT_EQ (iter1->second.visionName, _T("name1"));
    ASSERT_EQ (iter1->second.email, _T("address1"));

    map<String, ncACSUserInfo>::iterator iter2 = userInfoMap.find (_T("userid2"));
    ASSERT_NE (iter2, userInfoMap.end ());
    ASSERT_EQ (iter2->second.id, _T("userid2"));
    ASSERT_EQ (iter2->second.account, _T("account2"));
    ASSERT_EQ (iter2->second.visionName, _T("name2"));
    ASSERT_EQ (iter2->second.email, _T("address2"));

    // 清除掉插入的数据
    strSql.format (_T("delete from t_user;"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, GetDepartInfoByIdBatch)
{
    // 插入2条部门信息，再去获取
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_department(f_department_id,f_name) values ('departid1', 'name1'), ('departid2', 'name2');"));
    dbOper->Execute (strSql);

    vector<String> departIds;
    departIds.push_back (_T("departid1"));
    departIds.push_back (_T("departid2"));
    map<String, ncACSDepartInfo> departInfoMap;

    ASSERT_NO_THROW (_manager->GetDepartInfoByIdBatch (departIds, departInfoMap));

    map<String, ncACSDepartInfo>::iterator iter1 = departInfoMap.find (_T("departid1"));
    ASSERT_NE (iter1, departInfoMap.end ());
    ASSERT_EQ (iter1->second.name, _T("name1"));

    map<String, ncACSDepartInfo>::iterator iter2 = departInfoMap.find (_T("departid2"));
    ASSERT_NE (iter2, departInfoMap.end ());
    ASSERT_EQ (iter2->second.name, _T("name2"));

    // 清除掉插入的数据
    strSql.format (_T("delete from t_department;"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, GetGroupInfoByIdBatch)
{
    // 插入2条联系人组信息，再去获取
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_person_group(f_group_id,f_user_id,f_group_name,f_person_count) values ")
                        _T("('groupid1', 'userid1', 'name1', 1),")
                        _T("('groupid2', 'userid2', 'name2', 2);"));
    dbOper->Execute (strSql);

    vector<String> groupIds;
    groupIds.push_back (_T("groupid1"));
    groupIds.push_back (_T("groupid2"));
    map<String, ncGroupInfo> groupInfoMap;

    ASSERT_NO_THROW (_manager->GetGroupInfoByIdBatch (groupIds, groupInfoMap));

    map<String, ncGroupInfo>::iterator iter1 = groupInfoMap.find (_T("groupid1"));
    ASSERT_NE (iter1, groupInfoMap.end ());
    ASSERT_EQ (iter1->second.id, _T("groupid1"));
    ASSERT_EQ (iter1->second.createrId, _T("userid1"));
    ASSERT_EQ (iter1->second.groupName, _T("name1"));
    ASSERT_EQ (iter1->second.count, 1);

    map<String, ncGroupInfo>::iterator iter2 = groupInfoMap.find (_T("groupid2"));
    ASSERT_NE (iter2, groupInfoMap.end ());
    ASSERT_EQ (iter2->second.id, _T("groupid2"));
    ASSERT_EQ (iter2->second.createrId, _T("userid2"));
    ASSERT_EQ (iter2->second.groupName, _T("name2"));
    ASSERT_EQ (iter2->second.count, 2);

    // 清除掉插入的数据
    strSql.format (_T("delete from t_person_group;"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, GetParentDepartIds)
{
    // 插入2条联系人组信息，再去获取
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_department_relation (f_department_id, f_parent_department_id) values")
        _T("('departid1', 'dep'),")
        _T("('departid2', 'departid1'),")
        _T("('departid3', 'departid2'),")
        _T("('departid5', 'dep'),")
        _T("('departid7', 'dep');"));

    dbOper->Execute (strSql);

    set<String> parentDepIds;
    ASSERT_NO_THROW (_manager->GetParentDepartIds ("departid3", parentDepIds));

    ASSERT_EQ(parentDepIds.size(), 3);
    ASSERT_EQ(parentDepIds.count("departid2"), 1);
    ASSERT_EQ(parentDepIds.count("departid1"), 1);
    ASSERT_EQ(parentDepIds.count("dep"), 1);

    ASSERT_NO_THROW (_manager->GetParentDepartIds ("departid5", parentDepIds));
    ASSERT_EQ(parentDepIds.count("dep"), 1);
    ASSERT_NE(find(parentDepIds.begin(), parentDepIds.end(), "dep"), parentDepIds.end());

    ASSERT_NO_THROW (_manager->GetParentDepartIds ("dep", parentDepIds));
    ASSERT_EQ(parentDepIds.size(), 0);

    // 清除掉插入的数据
    strSql.format (_T("delete from t_department_relation;"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, GetNameByAccessorId)
{
    // 插入2条用户信息，再去获取
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_user(f_user_id,f_login_name,f_display_name,f_mail_address) values ")
        _T("('userid1', 'account1', '用户01', 'address1'),")
        _T("('userid2', 'account2', '用户02', 'address2');"));

    dbOper->Execute (strSql);

    // 插入2条部门信息，再去获取
    strSql.format (_T("insert into t_department(f_department_id,f_name) values ('departid1', 'depname1'), ('departid2', 'depname2');"));
    dbOper->Execute (strSql);

    // 插入2条联系人组信息，再去获取
    strSql.format (_T("insert into t_person_group(f_group_id,f_user_id,f_group_name,f_person_count) values ")
        _T("('groupid1', 'userid1', 'groupname1', 1),")
        _T("('groupid2', 'userid2', 'groupname2', 2);"));
    dbOper->Execute (strSql);

    // 获取用户
    ASSERT_EQ (_manager->GetNameByAccessorId (_T("userid1"), IOC_USER), _T("用户01"));
    ASSERT_EQ (_manager->GetNameByAccessorId (_T("userid2"), IOC_USER), _T("用户02"));

    // 获取部门
    ASSERT_EQ (_manager->GetNameByAccessorId (_T("departid1"), IOC_DEPARTMENT), _T("depname1"));

    // 获取联系人组
    ASSERT_EQ (_manager->GetNameByAccessorId (_T("groupid1"), IOC_USER_GROUP), _T("groupname1"));

    // 其他
    ASSERT_EQ (_manager->GetNameByAccessorId (_T("groupid2"), (ncIOCAccesorType)6), _T("groupid2"));

    // 清除掉插入的数据
    strSql.format (_T("delete from t_user;"));
    dbOper->Execute (strSql);

    strSql.format (_T("delete from t_person_group;"));
    dbOper->Execute (strSql);

    strSql.format (_T("delete from t_department;"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, GetDirectBelongDepartmentIds)
{
    // 插入2条用户信息，再去获取
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_user_department_relation(f_user_id,f_department_id) values ")
        _T("('userid1', 'departid1'),")
        _T("('userid1', 'departid2'),")
        _T("('userid1', 'departid3');"));

    dbOper->Execute (strSql);

    // 获取父部门
    vector <String> depIds;
    ASSERT_NO_THROW (_manager->GetDirectBelongDepartmentIds (_T("userid1"), depIds));

    ASSERT_EQ (3, depIds.size ());
    ASSERT_EQ (_T("departid1"), depIds[0]);
    ASSERT_EQ (_T("departid2"), depIds[1]);
    ASSERT_EQ (_T("departid3"), depIds[2]);

    // 清除掉插入的数据
    strSql.format (_T("delete from t_user_department_relation;"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, GetAllBelongDepartmentIds)
{
    /* 插入用户及部门信息
     *        dep
     *       /   \
     *     dep1  dep2    dep3
     *      |     |       |
     *    user1  user1   user1
     */
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_user_department_relation(f_user_id,f_department_id) values ")
        _T("('userid1', 'departid1'),")
        _T("('userid1', 'departid2'),")
        _T("('userid1', 'departid3');"));

    dbOper->Execute (strSql);

    strSql.format (_T("insert into t_department_relation (f_department_id, f_parent_department_id) values")
        _T("('departid1', 'dep'),")
        _T("('departid2', 'dep');"));

    dbOper->Execute (strSql);

    vector<String> departIds;
    ASSERT_NO_THROW (_manager->GetAllBelongDepartmentIds (_T("userid1"), departIds));

    ASSERT_EQ (4, departIds.size ());

    // 清除掉插入的数据
    strSql.format (_T("delete from t_user_department_relation;"));
    dbOper->Execute (strSql);

    strSql.format (_T("delete from t_department_relation;"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, IsCustomDocManager)
{
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_department_responsible_person(f_user_id,f_department_id) values ")
        _T("('userid1', 'departid1'),")
        _T("('userid1', 'departid2'),")
        _T("('userid1', 'departid3');"));

    dbOper->Execute (strSql);

    // 指定userid1为组织管理员
    strSql.format (_T("insert into t_user_role_relation(f_user_id,f_role_id) values ")
        _T("('userid1', 'e63e1c88-ad03-11e8-aa06-000c29358ad6');"));

    dbOper->Execute (strSql);

    bool isManager;
    ASSERT_NO_THROW (isManager = _manager->IsCustomDocManager (_T("userid1")));
    ASSERT_EQ (true, isManager);

    // 清除掉插入的数据
    strSql.format (_T("delete from t_user_department_relation;"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, GetUserType)
{
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_user(f_user_id,f_login_name,f_display_name,f_auth_type) values ")
        _T("('userid0', 'account0', 'name0', 0),")
        _T("('userid1', 'account1', 'name1', 1),")
        _T("('userid2', 'account2', 'name2', 2),")
        _T("('userid3', 'account3', 'name3', 3),")
        _T("('userid4', 'account4', 'name4', 4);"));

    dbOper->Execute (strSql);

    ASSERT_EQ (_manager->GetUserType (_T("notexist")), USER_TYPE_NONE);
    ASSERT_EQ (_manager->GetUserType (_T("userid0")), USER_TYPE_NONE);
    ASSERT_EQ (_manager->GetUserType (_T("userid1")), USER_TYPE_LOCAL);
    ASSERT_EQ (_manager->GetUserType (_T("userid2")), USER_TYPE_DOMAIN);
    ASSERT_EQ (_manager->GetUserType (_T("userid3")), USER_TYPE_THIRD);
    ASSERT_EQ (_manager->GetUserType (_T("userid4")), USER_TYPE_NONE);

    // 清除掉插入的数据
    strSql.format (_T("delete from t_user;"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, GetAllGroups)
{
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_person_group(f_group_id,f_user_id,f_group_name,f_person_count) values ")
                    _T("('group1', 'userid1', '才', 10),")
                    _T("('group2', 'userid1', '啊', 9),")
                    _T("('group3', 'userid1', '吧', 5),")
                    _T("('group4', 'userid1', '临时联系人', 5),")
                    _T("('group5', 'userid2', '的', 5)"));

    dbOper->Execute (strSql);

    vector<ncGroupInfo> groups;
    ASSERT_NO_THROW (_manager->GetAllGroups (_T("userid1"), groups));

    ASSERT_EQ (groups.size (), 4);

    // 结果是按照f_group_name(GBK)排序的
    // ASSERT_EQ (groups[0].id, _T("group2"));
    // ASSERT_EQ (groups[1].id, _T("group3"));
    // ASSERT_EQ (groups[2].id, _T("group1"));
    // ASSERT_EQ (groups[3].id, _T("group4"));

    // 清除掉插入的数据
    strSql.format (_T("delete from t_person_group;"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, GetContactors)
{
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_user(f_user_id,f_login_name,f_display_name,f_mail_address) values ")
        _T("('userid0', 'account0', '才', 'mail0'),")
        _T("('userid1', 'account1', '啊', 'mail1'),")
        _T("('userid2', 'account2', '吧', 'mail2'),")
        _T("('userid3', 'account3', '的', 'mail3');"));
    dbOper->Execute (strSql);

    strSql.format (_T("insert into t_contact_person(f_group_id,f_user_id) values ")
        _T("('group1', 'userid0'),")
        _T("('group1', 'userid1'),")
        _T("('group1', 'userid2');"));

    dbOper->Execute (strSql);

    vector<ncACSUserInfo> userInfos;
    ASSERT_NO_THROW (_manager->GetContactors (_T("group1"), userInfos));

    ASSERT_EQ (3, userInfos.size ());
    // ASSERT_EQ (userInfos[0].id, _T("userid1"));
    // ASSERT_EQ (userInfos[1].id, _T("userid2"));
    // ASSERT_EQ (userInfos[2].id, _T("userid0"));

    // 清除掉插入的数据
    strSql.format (_T("delete from t_contact_person;"));
    dbOper->Execute (strSql);

    strSql.format (_T("delete from t_user;"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, GetAllBelongGroups)
{
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_contact_person(f_group_id,f_user_id) values ")
        _T("('group1', 'userid1'),")
        _T("('group2', 'userid1'),")
        _T("('group3', 'userid1');"));

    dbOper->Execute (strSql);

    vector<String> groupIds;
    ASSERT_NO_THROW (_manager->GetAllBelongGroups (_T("userid1"), groupIds));
    ASSERT_EQ (3, groupIds.size ());
    // ASSERT_EQ (_T("group1"), groupIds[0]);
    // ASSERT_EQ (_T("group2"), groupIds[1]);
    // ASSERT_EQ (_T("group3"), groupIds[2]);

    // 清除掉插入的数据
    strSql.format (_T("delete from t_contact_person;"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, GetUserDisplayName)
{
    // 插入2条用户信息，再去获取
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_user(f_user_id,f_login_name,f_display_name,f_mail_address) values ")
        _T("('userid1', 'account1', 'username1', 'address1'),")
        _T("('administrator', 'admin', '管理员', 'address2');"));

    dbOper->Execute (strSql);

    String name;
    ASSERT_NO_THROW (_manager->GetUserDisplayName (_T("userid1"), name));
    ASSERT_EQ (name, _T("username1"));

    // 清除掉插入的数据
    strSql.format (_T("delete from t_user;"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, GetUserInfoById)
{
    // 插入2条用户信息，再去获取
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_user(f_user_id,f_login_name,f_display_name,f_mail_address,f_tel_number) values ")
        _T("('userid1', 'account1', 'username1', 'address1', '13000000000'),")
        _T("('administrator', 'admin', '管理员', 'address2', '13300000000');"));

    dbOper->Execute (strSql);

    ncACSUserInfo userInfo;
    ASSERT_NO_THROW (_manager->GetUserInfoById (_T("userid1"), userInfo));

    ASSERT_EQ (userInfo.account, _T("account1"));
    ASSERT_EQ (userInfo.visionName, _T("username1"));
    ASSERT_EQ (userInfo.email, _T("address1"));
    ASSERT_EQ (userInfo.telNumber, _T("13000000000"));

    // 清除掉插入的数据
    strSql.format (_T("delete from t_user;"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, GetUserInfoByAccount)
{
    // 插入2条用户信息，再去获取
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_user(f_user_id,f_login_name,f_display_name,f_mail_address) values ")
        _T("('userid1', 'account1', 'username1', 'address1'),")
        _T("('administrator', 'admin', '管理员', 'address2');"));

    dbOper->Execute (strSql);

    ncACSUserInfo userInfo;
    int accountType = 0;
    _manager->GetUserInfoByAccount (_T("account1"), userInfo, accountType);

    ASSERT_EQ (userInfo.account, _T("account1"));
    ASSERT_EQ (userInfo.visionName, _T("username1"));
    ASSERT_EQ (userInfo.email, _T("address1"));

    ASSERT_EQ(_manager->GetUserInfoByAccount (_T("abcc"), userInfo, accountType), false);

    // 清除掉插入的数据
    strSql.format (_T("delete from t_user;"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, IsUserEnabled)
{
    // 插入2条用户信息，再去获取
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_user(f_user_id,f_login_name,f_display_name,f_mail_address,f_status) values ")
        _T("('userid1', 'account1', 'username1', 'address1', 0),")
        _T("('administrator', 'admin', '管理员', 'address2', 1);"));

    dbOper->Execute (strSql);

    // 运行
    ASSERT_EQ (_manager->IsUserEnabled (_T("userid1")), true);
    ASSERT_EQ (_manager->IsUserEnabled (_T("admin")), false);

    // 清除掉插入的数据
    strSql.format (_T("delete from t_user;"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, GetUserInfoByThirdId)
{
    // 插入2条用户信息，再去获取
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_user(f_user_id,f_login_name,f_display_name,f_mail_address,f_third_party_id) values ")
        _T("('userid1', 'account1', 'username1', 'address1', 'thirdid1'),")
        _T("('administrator', 'admin', '管理员', 'address2', 'thirdid2');"));

    dbOper->Execute (strSql);

    ncACSUserInfo userInfo;
    ASSERT_NO_THROW (_manager->GetUserInfoByThirdId (_T("thirdid1"), userInfo));

    ASSERT_EQ (userInfo.id, _T("userid1"));
    ASSERT_EQ (userInfo.account, _T("account1"));
    ASSERT_EQ (userInfo.visionName, _T("username1"));
    ASSERT_EQ (userInfo.email, _T("address1"));

    // 清除掉插入的数据
    strSql.format (_T("delete from t_user;"));
    dbOper->Execute (strSql);
}


TEST_F (ncACSShareMgntUT, GetDepartInfoById)
{
    // 插入2条部门信息，再去获取
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_department(f_department_id,f_name) values")
        _T("('departid1', 'name1'),")
        _T("('departid2', 'name2');"));
    dbOper->Execute (strSql);

    bool ret = false;
    ncACSDepartInfo departInfo;
    ASSERT_NO_THROW (ret = _manager->GetDepartInfoById (_T("departid1"), departInfo));

    ASSERT_EQ (ret, true);
    ASSERT_EQ (departInfo.id, _T("departid1"));
    ASSERT_EQ (departInfo.name, _T("name1"));

    ASSERT_NO_THROW (ret = _manager->GetDepartInfoById (_T("departid2"), departInfo));

    ASSERT_EQ (ret, true);
    ASSERT_EQ (departInfo.id, _T("departid2"));
    ASSERT_EQ (departInfo.name, _T("name2"));

    ASSERT_NO_THROW (ret = _manager->GetDepartInfoById (_T("none"), departInfo));
    ASSERT_EQ (ret, false);

    // 清除掉插入的数据
    strSql.format (_T("delete from t_department;"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, GetSubDeps)
{
    // 准备数据
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_department_relation (f_department_id, f_parent_department_id) values")
        _T("('departid1', 'dep'),")
        _T("('departid2', 'dep'),")
        _T("('departid3', 'dep');"));

    dbOper->Execute (strSql);

    // 插入2条部门信息，再去获取
    strSql.format (_T("insert into t_department(f_department_id,f_name) values ")
                    _T("('departid1', '吧'), ")
                    _T("('departid2', '啊'),")
                    _T("('departid3', '才'),")
                    _T("('dep', '组织');"));
    dbOper->Execute (strSql);

    vector<ncACSDepartInfo> depInfos;
    ASSERT_NO_THROW (_manager->GetSubDeps (_T("dep"), depInfos));
    ASSERT_EQ (depInfos.size (), 3);
    // ASSERT_EQ (depInfos[0].id, _T("departid2"));
    // ASSERT_EQ (depInfos[1].id, _T("departid1"));
    // ASSERT_EQ (depInfos[2].id, _T("departid3"));

    strSql.format (_T("delete from t_department_relation;"));
    dbOper->Execute (strSql);

    strSql.format (_T("delete from t_department;"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, GetSubUsers)
{
    // 准备数据
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_user(f_user_id,f_login_name,f_display_name,f_mail_address) values ")
        _T("('userid1', 'account1', '吧', 'address1'),")
        _T("('userid2', 'account2', '啊', 'address2'),")
        _T("('userid3', 'account3', '才', 'address3');"));

    dbOper->Execute (strSql);

    strSql.format (_T("insert into t_user_department_relation(f_user_id,f_department_id) values ")
        _T("('userid1', 'departid1'),")
        _T("('userid2', 'departid1'),")
        _T("('userid3', 'departid1');"));

    dbOper->Execute (strSql);

    vector<ncACSUserInfo> userInfos;
    ASSERT_NO_THROW (_manager->GetSubUsers (_T("departid1"), userInfos));
    ASSERT_EQ (userInfos.size (), 3);

    // ASSERT_EQ (userInfos[0].id, _T("userid2"));
    // ASSERT_EQ (userInfos[1].id, _T("userid1"));
    // ASSERT_EQ (userInfos[2].id, _T("userid3"));

    strSql.format (_T("delete from t_user_department_relation"));
    dbOper->Execute (strSql);

    strSql.format (_T("delete from t_user"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, GetAllUser)
{
    // 准备数据
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_user(f_user_id,f_login_name,f_display_name,f_mail_address) values ")
        _T("('userid1', 'account1', 'username1', 'address1'),")
        _T("('userid2', 'admin', 'enginsl', 'address2'),")
        _T("('userid3', 'xiaosan', 'sage', 'address3');"));

    dbOper->Execute (strSql);

    vector<ncACSUserInfo> userInfos;
    ASSERT_NO_THROW (_manager->GetAllUser (userInfos));
    ASSERT_EQ (userInfos.size (), 3);

    for (size_t i = 0; i < userInfos.size (); ++i) {
        if (userInfos[i].id == _T("userid1")) {
            ASSERT_EQ (userInfos[i].account, _T("account1"));
            ASSERT_EQ (userInfos[i].visionName, _T("username1"));
            ASSERT_EQ (userInfos[i].email, _T("address1"));
        }
        else if (userInfos[i].id == _T("userid2")) {
            ASSERT_EQ (userInfos[i].account, _T("admin"));
            ASSERT_EQ (userInfos[i].visionName, _T("enginsl"));
            ASSERT_EQ (userInfos[i].email, _T("address2"));
        }
        else if (userInfos[i].id == _T("userid3")) {
            ASSERT_EQ (userInfos[i].account, _T("xiaosan"));
            ASSERT_EQ (userInfos[i].visionName, _T("sage"));
            ASSERT_EQ (userInfos[i].email, _T("address3"));
        }
        else
            ASSERT_TRUE (false);
    }

    strSql.format (_T("delete from t_user"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, SearchDepartment)
{
    // 准备数据
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    // 存在用户userid1,userid2,userid3
    String strSql;
    strSql.format (_T("insert into t_user(f_user_id,f_login_name,f_display_name,f_mail_address) values ")
        _T("('266c6a42-6131-4d62-8f39-853e7093701c', 'admin', '0key0', 'address0'),")
        _T("('userid', 'user00', 'name0', 'address00'),")
        _T("('userid1', 'user01', '1key1', 'address1'),")
        _T("('userid2', '2key2', 'name2', 'address2'),")
        _T("('userid3', 'user03', 'name3', 'address3');"));

    dbOper->Execute (strSql);

    // admin 能够搜索所有的用户
    vector<ncACSUserInfo> userInfos;
    ASSERT_NO_THROW (_manager->SearchDepartment (_T("266c6a42-6131-4d62-8f39-853e7093701c"), _T("key"), userInfos));
    ASSERT_EQ (userInfos.size (), 2);

    // 空的key时，返回的信息为空
    ASSERT_NO_THROW (_manager->SearchDepartment (_T("266c6a42-6131-4d62-8f39-853e7093701c"), _T(""), userInfos));
    ASSERT_EQ (userInfos.size (), 0);

    // 不管理任何部门的用户无法搜索到用户信息
    ASSERT_NO_THROW (_manager->SearchDepartment (_T("userid"), _T("key"), userInfos));
    ASSERT_EQ (userInfos.size (), 0);

    /*
     *        dep1(userid)        dep3(userid)
              /                     |
            dep2               userid1
           /
      userid2   userid3
    */
    // 存在部门dep1,dep2,dep3
    strSql.format (_T("insert into t_department(f_department_id,f_name) values ")
        _T("('departid1', 'dep1'),")
        _T("('departid2', 'dep2'),")
        _T("('departid3', 'dep3');"));
    dbOper->Execute (strSql);

    // 部门关系
    strSql.format (_T("insert into t_department_relation(f_department_id,f_parent_department_id) values ")
        _T("('departid2', 'departid1');"));
    dbOper->Execute (strSql);

    // 用户-部门关系
    strSql.format (_T("insert into t_user_department_relation(f_user_id,f_department_id) values ")
        _T("('userid', 'departid1'),")
        _T("('userid2', 'departid2'),")
        _T("('userid3', 'departid2'),")
        _T("('userid', 'departid3'),")
        _T("('userid1', 'departid3');"));
    dbOper->Execute (strSql);

    ASSERT_NO_THROW (_manager->SearchDepartment (_T("userid"), _T("key"), userInfos));
    ASSERT_EQ (userInfos.size (), 0);

    // 清理数据
    strSql.format (_T("delete from t_user;"));
    dbOper->Execute (strSql);

    strSql.format (_T("delete from t_department;"));
    dbOper->Execute (strSql);

    strSql.format (_T("delete from t_department_relation;"));
    dbOper->Execute (strSql);

    strSql.format (_T("delete from t_user_department_relation;"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, SearchContactGroup)
{
    // 准备数据
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    // 存在用户userid1,userid2,userid3
    String strSql;
    // 清除掉插入的数据
    strSql.format (_T("delete from t_user;"));
    dbOper->Execute (strSql);
    strSql.format (_T("delete from t_contact_person"));
    dbOper->Execute (strSql);
    strSql.format (_T("delete from t_person_group"));
    dbOper->Execute (strSql);
    strSql.format (_T("delete from t_sharemgnt_config"));
    dbOper->Execute (strSql);

    // 添加搜索配置信息
    strSql.format (_T("insert into t_sharemgnt_config(f_key, f_value) values ")
        _T("('search_user_config', '{\"exactSearch\":false, \"searchRange\":3, \"searchResults\":2}');"));
    dbOper->Execute (strSql);
    strSql.format (_T("insert into t_sharemgnt_config(f_key, f_value) values ")
        _T("('only_share_to_user', 0);"));
    dbOper->Execute (strSql);

    strSql.format (_T("insert into t_user(f_user_id,f_login_name,f_display_name,f_mail_address) values ")
        _T("('userid1', 'user01', 'name1', 'address1'),")
        _T("('userid2', 'akey1', 'name2', 'address2'),")
        _T("('userid3', 'user03', 'akey2', 'address3');"));

    dbOper->Execute (strSql);

    // userid1创建了2个联系人组
    strSql.format (_T("insert into t_person_group(f_user_id,f_group_id,f_group_name) values ") \
        _T("('userid1', 'groupid1', 'group_key1'),") \
        _T("('userid1', 'groupid2', 'name2');"));
    dbOper->Execute (strSql);

    // 组groupid1里面有userdid2,userid3，组groupid2里面有userid2
    strSql.format (_T("insert into t_contact_person(f_group_id,f_user_id) values ")
        _T("('groupid1', 'userid2'),")
        _T("('groupid1', 'userid3'),")
        _T("('groupid2', 'userid2');"));

    dbOper->Execute (strSql);

    // 执行
    vector<ncACSUserInfo> userInfos;
    vector<ncGroupInfo> groupInfos;

    // 空的key时，返回的信息为空
    ASSERT_NO_THROW (_manager->SearchContactGroup (_T("userid1"), _T(""), 0, 10, userInfos, groupInfos));
    ASSERT_EQ (userInfos.size (), 0);
    ASSERT_EQ (groupInfos.size (), 0);

    // 输入非法的key时，不抛出异常
    ASSERT_NO_THROW (_manager->SearchContactGroup (_T("userid1"), _T("1' or delete from t_user"), 0, 10, userInfos, groupInfos));

    // 输入%，不能获取任何信息
    ASSERT_NO_THROW (_manager->SearchContactGroup (_T("userid1"), _T("%%"), 0, 10, userInfos, groupInfos));
    ASSERT_EQ (userInfos.size (), 0);
    ASSERT_EQ (groupInfos.size (), 0);

    // userid1去搜索_，能搜索到userid1,userid2,departid1,departid2
    ASSERT_NO_THROW (_manager->SearchContactGroup (_T("userid1"), _T("_"), 0, 10, userInfos, groupInfos));
    ASSERT_EQ (userInfos.size (), 0);
    ASSERT_EQ (groupInfos.size (), 1);

    ASSERT_NO_THROW (_manager->SearchContactGroup (_T("userid1"), _T("key"), 0, 10, userInfos, groupInfos));

    // 检查结果
    ASSERT_EQ (userInfos.size (), 2);

    int count = 0;
    for (size_t i = 0; i < userInfos.size (); ++i) {
        if (userInfos[i].id == _T("userid2")) {
            ASSERT_EQ (userInfos[i].account, _T("akey1"));
            ASSERT_EQ (userInfos[i].visionName, _T("name2"));
            ASSERT_EQ (userInfos[i].email, _T("address2"));
            ++count;
        }
        else if (userInfos[i].id == _T("userid3")) {
            ASSERT_EQ (userInfos[i].account, _T("user03"));
            ASSERT_EQ (userInfos[i].visionName, _T("akey2"));
            ASSERT_EQ (userInfos[i].email, _T("address3"));
            ++count;
        }
        else
            ASSERT_TRUE (false);
    }

    ASSERT_EQ (count, 2);

    count = 0;
    for (size_t i = 0; i < groupInfos.size (); ++i) {
        if (groupInfos[i].id == _T("groupid1")) {
            ASSERT_EQ (groupInfos[i].createrId, _T("userid1"));
            ASSERT_EQ (groupInfos[i].groupName, _T("group_key1"));
            ++count;
        }
        else
            ASSERT_TRUE (false);
    }

    ASSERT_EQ (count, 1);

    // 清理数据
    strSql.format (_T("delete from t_contact_person"));
    dbOper->Execute (strSql);

    strSql.format (_T("delete from t_person_group"));
    dbOper->Execute (strSql);

    strSql.format (_T("delete from t_user"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, GetBelongDepartByIdBatch)
{
    // 准备数据
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());
    String strSql;
    // 清除掉插入的数据
    strSql.format (_T("delete from t_user;"));
    dbOper->Execute (strSql);
    strSql.format (_T("insert into t_user(f_user_id,f_login_name,f_display_name,f_mail_address) values ")
        _T("('userid1', 'account1', 'username1', 'address1'),")
        _T("('userid2', 'admin', 'enginsl', 'address2'),")
        _T("('userid3', 'xiaosan', 'sage', 'address3');"));

    dbOper->Execute (strSql);

    strSql.format (_T("insert into t_user_department_relation(f_user_id,f_department_id) values ")
        _T("('userid1', 'departid1'),")
        _T("('userid2', 'departid1'),")
        _T("('userid2', 'departid2'),")
        _T("('userid3', 'departid2');"));

    dbOper->Execute (strSql);

    strSql.format (_T("insert into t_department(f_department_id,f_name) values ") \
        _T("('departid1', 'bumen1'),") \
        _T("('departid2', 'bumen2')"));
    dbOper->Execute (strSql);

    // 执行
    map<String, ncACSDepartInfo> infoMap;
    vector<String> userIds;
    userIds.push_back (_T("userid1"));
    userIds.push_back (_T("userid2"));
    userIds.push_back (_T("userid3"));

    ASSERT_NO_THROW (_manager->GetBelongDepartByIdBatch (userIds, infoMap));

    // 检查结果
    ASSERT_EQ (infoMap.size (), 3);

    ncACSDepartInfo info = infoMap[_T("userid1")];
    ASSERT_EQ (info.id, _T("departid1"));
    ASSERT_EQ (info.name, _T("bumen1"));

    info = infoMap[_T("userid2")];
    ASSERT_EQ (info.id, _T("departid1"));
    ASSERT_EQ (info.name, _T("bumen1"));

    info = infoMap[_T("userid3")];
    ASSERT_EQ (info.id, _T("departid2"));
    ASSERT_EQ (info.name, _T("bumen2"));

    strSql.format (_T("delete from t_user_department_relation"));
    dbOper->Execute (strSql);

    strSql.format (_T("delete from t_user"));
    dbOper->Execute (strSql);

    strSql.format (_T("delete from t_department"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, GetBelongGroupByIdBatch)
{
    // 准备数据
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    // 存在用户userid1,userid2,userid3
    String strSql;
    // 清除掉插入的数据
    strSql.format (_T("delete from t_user;"));
    dbOper->Execute (strSql);
    strSql.format (_T("delete from t_contact_person"));
    dbOper->Execute (strSql);
    strSql.format (_T("delete from t_person_group"));
    dbOper->Execute (strSql);

    strSql.format (_T("insert into t_user(f_user_id,f_login_name,f_display_name,f_mail_address) values ")
        _T("('userid1', 'user01', 'name1', 'address1'),")
        _T("('userid2', 'akey1', 'name2', 'address2'),")
        _T("('userid3', 'user03', 'akey2', 'address3');"));

    dbOper->Execute (strSql);

    // userid1创建了2个联系人组
    strSql.format (_T("insert into t_person_group(f_user_id,f_group_id,f_group_name) values ") \
        _T("('userid1', 'groupid1', 'group_key1'),") \
        _T("('userid1', 'groupid2', 'name2'),")
        _T("('userid2', 'groupid3', 'name3'),")
        _T("('userid3', 'groupid4', 'name5');"));
    dbOper->Execute (strSql);

    // 组groupid1里面有userdid2,userid3，组groupid2里面有userid2
    strSql.format (_T("insert into t_contact_person(f_group_id,f_user_id) values ")
        _T("('groupid1', 'userid2'),")
        _T("('groupid1', 'userid3'),")
        _T("('groupid2', 'userid2'),")
        _T("('groupid3', 'userid1'),")
        _T("('groupid3', 'userid2'),")
        _T("('groupid4', 'userid3');"));

    dbOper->Execute (strSql);

    // 执行
    vector<String> userIds;
    userIds.push_back (_T("userid2"));
    userIds.push_back (_T("userid3"));

    map<String, ncGroupInfo> infoMap;
    ASSERT_NO_THROW (_manager->GetBelongGroupByIdBatch (_T("userid1"), userIds, infoMap));

    ASSERT_EQ (infoMap.size (), 2);

    // 检查结果
    ncGroupInfo& info1 = infoMap[_T("userid2")];
    ASSERT_EQ (info1.id, _T("groupid1"));
    ASSERT_EQ (info1.groupName, _T("group_key1"));

    ncGroupInfo& info2 = infoMap[_T("userid3")];
    ASSERT_EQ (info2.id, _T("groupid1"));
    ASSERT_EQ (info2.groupName, _T("group_key1"));

    // 清理数据
    strSql.format (_T("delete from t_contact_person"));
    dbOper->Execute (strSql);

    strSql.format (_T("delete from t_person_group"));
    dbOper->Execute (strSql);

    strSql.format (_T("delete from t_user"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, ExtLogin)
{
    // 准备数据
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    // 存在用户userid1,userid2
    // 清除掉插入的数据
    String strSql;
    strSql.format (_T("delete from t_user;"));
    dbOper->Execute (strSql);

    strSql.format (_T("insert into t_user(f_user_id,f_login_name,f_display_name,f_mail_address) values ")
        _T("('userid1', 'user01', 'name1', 'address1'),")
        _T("('userid2', 'user02', 'name2', 'address2'),")
        _T("('userid3', 'admin', 'name2', 'address2');"));

    dbOper->Execute (strSql);

    //受信任的第三方app
    strSql.format (_T("insert into t_third_auth_info(f_app_id,f_app_key,f_enabled) values ")
        _T("('appid', 'app_secret0', '0'),")
        _T("('chuxiong', 'app_secret1', '1'),")
        _T("('ppsuc', 'app_secret2', '0');"));

    dbOper->Execute (strSql);

    int accountType = 0;
    // 测试用例 通用
    {
        // appid 不存在，返回String::EMPTY
        vector<String> params;
        ASSERT_EQ (_manager->ExtLogin (_T("notexistappid"), _T("user01"), _T("key"), params, accountType), String::EMPTY);
    }

    {
        // appid 存在但未开启，返回String::EMPTY
        vector<String> params;
        ASSERT_EQ (_manager->ExtLogin (_T("appid"), _T("user01"), _T("key"), params, accountType), String::EMPTY);
    }

    // 测试用例 楚雄供电局
    {
        // account不存在，返回String::EMPTY
        vector<String> params;
        ASSERT_EQ (_manager->ExtLogin (_T("chuxiong"), _T("notexistuser"), _T("key"), params, accountType), String::EMPTY);
    }

    {
        // account存在，key错误，返回String::EMPTY
        vector<String> params;
        ASSERT_EQ (_manager->ExtLogin (_T("chuxiong"), _T("user01"), _T("key"), params, accountType), String::EMPTY);
    }

    {
        // account存在，account为管理员账号，返回String::EMPTY
        vector<String> params;
        ASSERT_EQ (_manager->ExtLogin (_T("chuxiong"), _T("admin"), _T("key"), params, accountType), String::EMPTY);
    }

    {
        // account存在，account用户账号，key错误，返回String::EMPTY
        vector<String> params;
        ASSERT_EQ (_manager->ExtLogin (_T("chuxiong"), _T("user01"), _T("key"), params, accountType), String::EMPTY);
    }

    {
        // account存在，account为用户账号，key正确，校验成功，返回userid1
        String key = genMD5String2 (_T("chuxiongapp_secret1user01"));
        key.toLower ();
        vector<String> params;
        ASSERT_EQ (_manager->ExtLogin (_T("chuxiong"), _T("user01"), key, params, accountType), _T("userid1"));
    }

    //修改数据库，对不同受信任第三方做登录测试
    strSql.format (_T("delete from t_third_auth_info;"));
    dbOper->Execute (strSql);

    strSql.format (_T("insert into t_third_auth_info(f_app_id,f_app_key,f_enabled) values ")
        _T("('appid', 'app_secret0', '0'),")
        _T("('chuxiong', 'app_secret1', '0'),")
        _T("('ppsuc', 'app_secret2', '1');"));

    dbOper->Execute (strSql);

    //测试用例 中国人民公安大学
    {
        // account 不存在，返回String::EMPTY
        vector<String> params;
        params.push_back("time");
        ASSERT_EQ (_manager->ExtLogin (_T("ppsuc"), _T("notexistuser"), _T("key"), params, accountType), String::EMPTY);
    }

    {
        // account存在，account为管理员账号，返回String::EMPTY
        vector<String> params;
        params.push_back("time");
        ASSERT_EQ (_manager->ExtLogin (_T("ppsuc"), _T("admin"), _T("key"), params, accountType), String::EMPTY);
    }

    {
        // account 存在，key 错误，返回String::EMPTY
        vector<String> params;
        params.push_back("time");
        ASSERT_EQ (_manager->ExtLogin (_T("ppsuc"), _T("user01"), _T("key"), params, accountType), String::EMPTY);
    }

    {
        // account 存在，key 正确，返回userid1
        vector<String> params;
        String time = _T("2015-10-29 15:04:32");
        params.push_back(time);
        String key = genMD5String2 (_T("user02app_secret22015-10-29 15:04:32"));
        key.toLower ();
        ASSERT_EQ (_manager->ExtLogin (_T("ppsuc"), _T("user02"), key, params, accountType), _T("userid2"));
    }

    // 清理数据
    strSql.format (_T("delete from t_user;"));
    dbOper->Execute (strSql);

    strSql.format (_T("delete from t_third_auth_info;"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, GetUserCSFLevel)
{
    // 准备数据
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    // 存在用户userid1

    // 清除掉插入的数据
    String strSql;
    strSql.format (_T("delete from t_user;"));
    dbOper->Execute (strSql);

    strSql.format (_T("insert into t_user(f_user_id,f_login_name,f_display_name,f_mail_address, f_csf_level) values ")
        _T("('userid1', 'user01', 'name1', 'address1', 1),")
        _T("('userid2', 'user02', 'name2', 'address2', 2),")
        _T("('userid3', 'user03', 'name3', 'address3', 3);"));

    dbOper->Execute (strSql);

    {
        // 用户不存在，抛出异常
        ASSERT_ANY_THROW(_manager->GetUserCSFLevel("notexists", ncVisitorType::REALNAME));
    }

    {
        // 用户存在，获取密级
        ASSERT_EQ(_manager->GetUserCSFLevel("userid1", ncVisitorType::REALNAME), 1);
        ASSERT_EQ(_manager->GetUserCSFLevel("userid2", ncVisitorType::REALNAME), 2);
        ASSERT_EQ(_manager->GetUserCSFLevel("userid3", ncVisitorType::REALNAME), 3);
    }
}

// GetMailAddress(const String & receiverId, const String & tbName, const String & fieldName, String & email)
TEST_F (ncACSShareMgntUT, GetMailAddress)
{
    // 准备数据
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    // 清除掉插入的数据
    String strSql;
    strSql.format (_T("delete from t_department;"));
    dbOper->Execute (strSql);

    strSql.format (_T("insert into t_department(f_department_id, f_auth_type, f_name, f_is_enterprise, f_mail_address) values ")
                    _T("('departmentId', 1, 'depatName', 1, '122@eisoo.com');"));
    dbOper->Execute (strSql);

    // case 1: 组织部门不存在
    String email;
    _manager->GetMailAddress ("notexists", "t_department", "f_department_id", email);
    ASSERT_TRUE (email.isEmpty ());

    // case 2: 组织/部门存在，获取邮箱地址
    _manager->GetMailAddress ("departmentId", "t_department", "f_department_id", email);
    ASSERT_EQ (email, "122@eisoo.com");
}

// GetUserInfoByTelNumber
TEST_F (ncACSShareMgntUT, GetUserInfoByTelNumber)
{
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_user(f_user_id,f_login_name,f_display_name,f_tel_number,f_mail_address) values ")
        _T("('useridaaaaa', 'aa', 'aa', '15554334344', 'address1');"));

    dbOper->Execute (strSql);

    ncACSUserInfo userInfo;

    // 判断手机号是否存在，返回userid
    String telNumber = "15554334344";
    ASSERT_EQ(_manager->GetUserInfoByTelNumber (telNumber, userInfo), true);

    ASSERT_EQ (userInfo.id, _T("useridaaaaa"));

    // 清除掉插入的数据
    strSql.format (_T("delete from t_user;"));
    dbOper->Execute (strSql);
}
// GetUserInfoByEmail
TEST_F (ncACSShareMgntUT, GetUserInfoByEmail)
{
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_user(f_user_id,f_login_name,f_display_name,f_tel_number,f_mail_address) values ")
        _T("('useridbbbbb', 'bb', 'bb', '', 'bb@eisoo.com');"));

    dbOper->Execute (strSql);
    ncACSUserInfo userInfo;


    String email = "bb@eisoo.com";
    ASSERT_EQ(_manager->GetUserInfoByEmail (email, userInfo), true);

    ASSERT_EQ (userInfo.id, _T("useridbbbbb"));

    // 清除掉插入的数据
    strSql.format (_T("delete from t_user;"));
    dbOper->Execute (strSql);
}
// GetCustomConfigOfString
TEST_F (ncACSShareMgntUT, GetCustomConfigOfString)
{
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_sharemgnt_config(f_key, f_value) values")
    _T( "('vcode_server_status','{send_vcode_by_sms : false, send_vcode_by_email:false}');"));
    dbOper->Execute (strSql);

    ncACSUserInfo userInfo;

    String key = "vcode_server_status";
    String testvalue;
    ASSERT_EQ(_manager->GetCustomConfigOfString (key, testvalue),true);

    // 清除掉插入的数据
    strSql.format (_T("delete from t_sharemgnt_config;"));
    dbOper->Execute (strSql);
}

// GetCustomConfigOfString
TEST_F (ncACSShareMgntUT, OEM_GetConfigByOption)
{
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String strSql;
    strSql.format (_T("insert into t_oem_config(f_section, f_option, f_value) values")
    _T( "('shareweb_en-us', 'product', 'Anyshare');"));
    dbOper->Execute (strSql);

    ncACSUserInfo userInfo;

    String key = "vcode_server_status";
    String section = "shareweb_en-us";
    String option = "product";
    String value;
    ASSERT_EQ(_manager->OEM_GetConfigByOption (section, option, value),true);

    // 清除掉插入的数据
    strSql.format (_T("delete from t_sharemgnt_config;"));
    dbOper->Execute (strSql);
}

TEST_F (ncACSShareMgntUT, BatchUpdateUserLastRequestTime)
{
    nsCOMPtr<ncIDBOperator> dbOper = getter_AddRefs (GetShareMgntDBOperator ());

    String userId = "94752844-BDD0-4B9E-8927-111111111";
    String nTime = "2020-05-02 13:41:15";

    String strSql;
    strSql.format (_T("insert into t_user(f_user_id, f_login_name, f_display_name, f_last_request_time) values")
    _T( "('%s','xxxx','xxxx','%s');"),
    userId.getCStr(), nTime.getCStr() );

    dbOper->Execute (strSql);

    vector<ncACSRefreshInfo> infos;
    ncACSRefreshInfo oInfo;
    oInfo.userId = "94752844-BDD0-4B9E-8927-111111111";
    oInfo.lastRequestTime = "2020-05-08 13:41:15";
    oInfo.bUpdateClientTime = true;
    infos.push_back(oInfo);
    ASSERT_NO_THROW(_manager->BatchUpdateUserLastRequestTime (infos));


    strSql.format (_T(" select f_last_request_time, f_last_client_request_time from t_user where f_user_id = '%s'; "), userId.getCStr());
    ncDBRecords results;
    dbOper->Select (strSql, results);

    ASSERT_EQ(results.size (), 1);
    ASSERT_EQ(results[0][0], oInfo.lastRequestTime);
    ASSERT_EQ(results[0][1], oInfo.lastRequestTime);


    vector<ncACSRefreshInfo> infos1;
    ncACSRefreshInfo oInfo1;
    oInfo1.userId = "94752844-BDD0-4B9E-8927-111111111";
    oInfo1.lastRequestTime = "2020-05-09 13:41:15";
    oInfo1.bUpdateClientTime = false;
    infos1.push_back(oInfo1);
    ASSERT_NO_THROW(_manager->BatchUpdateUserLastRequestTime (infos1));


    strSql.format (_T(" select f_last_request_time, f_last_client_request_time from t_user where f_user_id = '%s'; "), userId.getCStr());
    ncDBRecords results1;
    dbOper->Select (strSql, results1);

    ASSERT_EQ(results1.size (), 1);
    ASSERT_EQ(results1[0][0], oInfo1.lastRequestTime);
    ASSERT_EQ(results1[0][1], oInfo.lastRequestTime);


    // 清除掉插入的数据
    strSql.format (_T("delete from t_user;"));
    dbOper->Execute (strSql);
}
