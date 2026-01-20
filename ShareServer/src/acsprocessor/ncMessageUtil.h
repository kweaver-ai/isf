#ifndef __NC_FORMAT_MESSAGE_H__
#define __NC_FORMAT_MESSAGE_H__

#include <dataapi/ncJson.h>
#include <acssharemgnt/public/ncIACSShareMgnt.h>
#include <acsprocessor/public/ncIACSOwnerManager.h>
#include <acsprocessor/public/ncIACSMessageManager.h>
#include <drivenadapter/public/userManagementInterface.h>

class ncMessageUtil
{
    AB_DECLARE_THREADSAFE_SINGLETON (ncMessageUtil)

public:
    ncMessageUtil ();
    ~ncMessageUtil ();

public:
    void GetAllUsers (const String & departId, vector<String> & userlist);
    String AccessorTypeToStr (int accessorType);
    void RemoveDuplicateStrs (vector<String>& strs);
    void CalcMsgReceivers (shared_ptr<acsMessageInfo> msg, std::vector<String>& receivers);
private:
    String getUserNameById (const String& userId);

private:
    nsCOMPtr<ncIACSShareMgnt>           _acsShareMgnt;
    nsCOMPtr<ncIACSOwnerManager>        _acsOwnerManager;
    nsCOMPtr<userManagementInterface>   _userManager;
};

#endif
