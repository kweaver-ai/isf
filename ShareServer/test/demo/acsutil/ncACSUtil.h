#ifndef __NC_ACS_UTIL_H
#define __NC_ACS_UTIL_H

// 遍历每个createrId 分配配额（用户文档=107374182400）
struct ncQuotaGroup {
    String userDocId;

    void reset ()
    {
        userDocId = String::EMPTY;
    }
};

class ncACSUtil: public ncBaseCommandHandler
{
public:
    ncACSUtil ();
    ~ncACSUtil ();

public:
    void ptime ();
    void ticks2String(void);
    void md5_encrypt ();
    void base64_encrypt ();
    void base64_decrypt ();
    void rsa_encrypt ();
    void rsa_decrypt ();
    void des_encrypt ();
    void des_decrypt ();
    void mtdecrypt ();
    void getEndOfDayTicks ();
    void testString ();
    void getParentString ();
    void getUTF8StringSize ();
    void insertPermRecords ();
    void checkToken ();
    void updateDB ();
    void testdb ();
    void conndb ();
    void isValidURL ();

    // 为没有个人文档的用户创建文档
    void createuserdoc ();

    // 升级用户的配额空间
    void setquota ();

    void getUsage ();

private:
    String removeBlankAndDot (const String& str);
    void printByte (const String& str);
    void validateURL (const String& str);

    void executeQuotaProcess (const vector<ncQuotaGroup> quotaGroups, const String& ip);
    string RSAEncrypt (const string& plainText, const char *path_key);
    string RSADecrypt (const string& cipherText, const char *path_key);

    size_t calcDecodeLength(const string& tmp);
    string Base64Encode (const string& plainText);
    string Base64Decode (const string& cipherText);

    string DESEncrypt(const string& plainText);
    string DESDecrypt(const string& cipherText);

    void printBytes(const string& str);
};

#endif // __NC_ACS_UTIL_H
