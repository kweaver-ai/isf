#ifndef __NC_ACS_SHAREMGNT_H
#define __NC_ACS_SHAREMGNT_H

#include <dboperator/public/ncIDBOperator.h>
#include <dboperatormanager/public/ncIDBOperatorManager.h>
#include <tedbcm/public/ncITEDBCManager.h>

#include "./public/ncIACSShareMgnt.h"

#ifdef __WINDOWS__
#define TO_ATOMIC(value) reinterpret_cast<LONG volatile*>(value)
#else
#define TO_ATOMIC(value) reinterpret_cast<gint32 volatile*>(value)
#endif

/* Header file */
class ncACSShareMgnt : public ncIACSShareMgnt
{
    AB_DECLARE_THREADSAFE_SINGLETON (ncACSShareMgnt)

public:
    NS_DECL_ISUPPORTS
    NS_DECL_NCIACSSHAREMGNT

    ncACSShareMgnt();
    ~ncACSShareMgnt();

protected:
    // 获取数据库连接
    virtual ncIDBOperator* GetDBOperator ();

    // 向下遍历用户管理的部门的子部门id
    void GetAllManageDepIds (const String& userId, vector<String>& departIds);

    // [id1,id2,id3]-> "'id1', 'id2', 'id3'"
    String GenerateGroupStr (const vector<String>& ids);

    // [id1,id2,id3]-> "'id1', 'id2', 'id3'" id 不转译
    String GenerateGroupStrWithOutEscapeEx (const vector<String>& ids);

    // set(id1,id2,id3)-> "'id1', 'id2', 'id3'"
    String GenerateGroupStrBySet (const set<String>& ids);

    // 去掉重复的str
    void RemoveDuplicateStrs (vector<String>& strs);

    // 判断是否为同一网段
    bool IsSameNetworkSegment (const String& ip, const String& netMask, const String& accessIp);

    // 将获取到的信息按照权重值进行排序
    void SortObjIdsWithPriority (const String& groupStr, const size_t objType, vector<String>& objIds);

    // 根据配置获取搜索用户限制条件
    String GenerateSearchStr (const bool& exact_search_user, const int& searchRange, const String escKey, const String escLikeKey);

    // des加密
    string DESEncrypt(const string& plainText);

    // base加密
    string Base64Encode(const string& input);

private:
    set<string>                         _adminIds;     // 系统管理员id列表
};

#endif // __NC_ACS_SHAREMGNT_H
