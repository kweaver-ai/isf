#ifndef __NC_EVFS_ACCESS_CONTROL_IOC_H
#define __NC_EVFS_ACCESS_CONTROL_IOC_H

#include <acssharemgnt/public/ncIACSShareMgnt.h>
#include <acsprocessor/public/ncIACSPermManager.h>
#include <acsprocessor/public/ncIACSOwnerManager.h>
#include <acsprocessor/public/ncIACSLockManager.h>
#include <acsprocessor/public/ncIACSConfManager.h>
#include <acsprocessor/public/ncIACSCommon.h>
#include <evfsioc/ncIEVFSAccessControlIOC.h>

/* Header file */
class ncEVFSAccessControlIOC : public ncIEVFSAccessControlIOC
{
    AB_DECLARE_THREADSAFE_SINGLETON (ncEVFSAccessControlIOC)

    typedef void (ncEVFSAccessControlIOC::*ncAccessControlFunc) (const ncACSSubjectAttr&, const ncACSObjectAttr&);
    typedef void (ncEVFSAccessControlIOC::*ncAccessFinishFunc) (const String& gnsPath, const String& userId);

public:
    NS_DECL_ISUPPORTS
    NS_DECL_NCIEVFSACCESSCONTROLIOC

    ncEVFSAccessControlIOC ();
    ncEVFSAccessControlIOC (ncIACSPermManager* acsPermManager,
                            ncIACSOwnerManager* acsOwnerManager,
                            ncIACSLockManager* acsLockManager,
                            ncIACSConfManager* acsConfManager,
                            ncIACSShareMgnt* acsShareMgnt);
    ~ncEVFSAccessControlIOC();

private:
    // 文件操作开始时的处理（用以检查权限）
    void onListFileVersion (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr);
    void onListDir (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr);
    void onRecycleFile (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr);
    void onRecycleDir (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr);
    void onCopyFile (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr);
    void onGetFile (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr);
    void onListRecycleBinDir (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr);
    void onSetRecyclePolicy (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr);
    void onDeleteFile (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr);
    void onDeleteDir (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr);
    void onPreviewFile (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr);
    void onGetFileMeta (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr);
    void onGetAttr (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr);
    void onGetDir (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr);
    void onSetCSFLevel (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr);
    void onSetTag (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr);
    void onQuarantineAppeal (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr);
    void onSetDocDue (const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr);

    // 文件操作结束后的处理（用以删除gns对象的权限信息）
    void onRecycleFileEnd (const String& gnsPath, const String& userId);
    void onRecycleDirEnd (const String& gnsPath, const String& userId);
    void onDeleteFileEnd (const String& gnsPath, const String& userId);
    void onDeleteDirEnd (const String& gnsPath, const String& userId);

public:
    void initHandlers();
    void checkPermHelper(const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr, int permValue, const String &who);
    ncCheckPermCode checkPermission(const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr, int permValue);
    void checkIsOwner(const ncACSSubjectAttr& subjectAttr, const ncACSObjectAttr& objAttr, const String &who);
    bool isFileCrawlOperation(const ncACSSubjectAttr & subjectAttr, const ncACSObjectAttr & objAttr);

protected:
    vector<int>                                         _allPerms;
    map<int, int>                                       _allowAttrMap;
    map<int, int>                                       _denyAttrMap;
    map<ncEVFSAccessType, ncAccessControlFunc>          _beginHandlers;
    map<ncEVFSAccessType, ncAccessFinishFunc>           _endHandlers;

    nsCOMPtr<ncIACSPermManager>                         _acsPermManager;
    nsCOMPtr<ncIACSOwnerManager>                        _acsOwnerManager;
    nsCOMPtr<ncIACSLockManager>                         _acsLockManager;
    nsCOMPtr<ncIACSShareMgnt>                           _acsShareMgnt;
    nsCOMPtr<ncIACSConfManager>                         _acsConfManager;
};

#endif // __NC_EVFS_ACCESS_CONTROL_IOC_H
