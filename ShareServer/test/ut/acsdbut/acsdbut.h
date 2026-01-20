#ifndef __NC_ACS_DB_UT_H__
#define __NC_ACS_DB_UT_H__

struct dbCAuthInfo;
struct ncDbFinderInfo;
struct dbOwnerInfo;
struct dbCustomPermInfo;
struct dbTokenInfo;

class String;
class ncIDBOperator;

//////////////////////////////////////////////////////////////////////////
//
// ncDBUTApi
//
class ncDBUTApi
{
public:
    #define NC_TB_T_ACS_DOC                _T("t_acs_doc")
    #define NC_TB_T_ACS_OWNER            _T("t_acs_owner")
    #define NC_TB_T_ACS_CUSTOM_PERM        _T("t_acs_custom_perm")
    #define NC_TB_T_ACS_ACCESS_TOKEN    _T("t_acs_access_token")
    #define NC_TB_T_ACS_FINDER            _T("t_acs_finder")

    #define NC_GENERATE_ANYSHARE_TB(anyshareOperator, tbName)                                    \
        if (!anyshareOperator.get()){                                                        \
            anyshareOperator = getter_AddRefs (ncDBUTApi::GenerateAnyshareTB (tbName));        \
        }

public:
    static ncIDBOperator* GetDBOperator (void);
    static ncIDBOperator* GenerateAnyshareTB (const String& tbName);
    static void CleanTable (ncIDBOperator* dbOperator, const String& tableName);
    static String FindCreateTBSQL (const String& tb);
};

//////////////////////////////////////////////////////////////////////////
//
// ncDBDocUTApi
//
class ncDBDocUTApi
{
public:
    // 前提条件
    #define INVALID_DOCINFO_1(info)                                                \
            ASSERT_TRUE (true == ncOIDUtil::IsObjectID (info.objectId));        \
            ASSERT_TRUE (true == ncOIDUtil::IsDocID (info.docId));                \
            ASSERT_TRUE (true == ncOIDUtil::IsGUID (info.createrId));

    #define INVALID_DOCINFO_2(info1, info2)                                        \
            INVALID_DOCINFO_1(info1);                                            \
            INVALID_DOCINFO_1(info2);

    #define INVALID_DOCINFO_3(info1, info2, info3)                                \
            INVALID_DOCINFO_1 (info1);                                            \
            INVALID_DOCINFO_2 (info2, info3);

public:
    static const String GetCreateTBSQL (void);
};

//////////////////////////////////////////////////////////////////////////
//
// ncDBFinderUTApi
//
class ncDBFinderUTApi
{
public:
    static const String GetCreateTBSQL (void);
    static bool EqualDBFinderInfo (const ncDbFinderInfo& info1, const ncDbFinderInfo& info2);
    static void InsertFinderInfo1 (ncIDBOperator* anyshareOperator, ncDbFinderInfo& info1);
    static void InsertFinderInfo2 (ncIDBOperator* anyshareOperator, ncDbFinderInfo& info1, ncDbFinderInfo& info2);
};

//////////////////////////////////////////////////////////////////////////
//
// ncDBOwnerUTApi
//
class ncDBOwnerUTApi
{
public:
    static const String GetCreateTBSQL (void);
    static bool EqualDBOwnerInfo (const dbOwnerInfo& info1, const dbOwnerInfo& info2);
    static void InsertOwnerInfo1 (ncIDBOperator* anyshareOperator, dbOwnerInfo& info1);
    static void InsertOwnerInfo2 (ncIDBOperator* anyshareOperator, dbOwnerInfo& info1, dbOwnerInfo& info2);
};


//////////////////////////////////////////////////////////////////////////
//
// ncDBPermUTApi
//
class ncDBPermUTApi
{
public:
    static const String GetCreateTBSQL (void);
    static bool EqualDBPermInfo (const dbCustomPermInfo& info1, const dbCustomPermInfo& info2);
    static void InsertPermInfo1 (ncIDBOperator* anyshareOperator, dbCustomPermInfo& info1);
    static void InsertPermInfo2 (ncIDBOperator* anyshareOperator, dbCustomPermInfo& info1, dbCustomPermInfo& info2);
};

//////////////////////////////////////////////////////////////////////////
//
// ncDBTokenUTApi
//
class ncDBTokenUTApi
{
public:
    static const String GetCreateTBSQL (void);
    static bool EqualDBTokenInfo (const dbTokenInfo& info1, const dbTokenInfo& info2);
    static void InsertTokenInfo1 (ncIDBOperator* anyshareOperator, dbTokenInfo& info1);
    static void InsertTokenInfo2 (ncIDBOperator* anyshareOperator, dbTokenInfo& info1, dbTokenInfo& info2);
};