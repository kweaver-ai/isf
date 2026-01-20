#include <abprec.h>
#include <thread>
#include "ncACSDataStore.h"
#include <ossclient/ossclient.h>
#include "acsdatastore.h"

#include "gen-cpp/ncTEVFS.h"
#include "gen-cpp/EVFS_constants.h"
#include <ethriftutil/ncThriftClient.h>
#define RETRY_COUNT 100

#define REQUEST_HTTP_CONTINUE                 (100)
#define REQUEST_HTTP_INTERNAL_SERVER_ERROR    (500)
#define REQUEST_HTTP_RETRY_TIMES              (10)

#define NC_ACSDS_REQUEST_TRY                                                                         \
    try {                                                                                            \

#define NC_ACSDS_REQUEST_CATCH(method, objectId, sn)                                                 \
    }                                                                                                \
    catch (Exception& e) {                                                                           \
        if (e.getErrorProviderName () == OSSCLIENT) {                                                \
            THROW_E (ACS_DATA_STORE, ACS_DATA_STORE_INTERNAL_ERR,                                    \
                     _T("Request (%s) to storage server error (object: %s, sn: %d, retcode: %d, msg: %s)."), \
                     method, objectId.getCStr (), sn, e.getErrorId (), e.toFullString ().getCStr ());        \
        }                                                                                            \
        else {                                                                                       \
            throw;                                                                                   \
        }                                                                                            \
    }                                                                                                \


#define NC_ACSDS_REQUEST_RETRY_BEGIN                                                                 \
        int i = 0;                                                                                   \
        while ((++i) <= REQUEST_HTTP_RETRY_TIMES) {                                                  \
            try {                                                                                    \

#define NC_ACSDS_REQUEST_RETRY_END                                                                   \
            }                                                                                        \
            catch (Exception& e) {                                                                   \
                if ((e.getErrorProviderName () == OSSCLIENT) && (i != REQUEST_HTTP_RETRY_TIMES)) {   \
                    res.code = 0;                                                                    \
                    res.body.clear ();                                                               \
                    res.headers.clear ();                                                            \
                    std::this_thread::sleep_for (std::chrono::milliseconds (i * 500));               \
                    continue;                                                                        \
                }                                                                                    \
                else {                                                                               \
                    throw;                                                                           \
                }                                                                                    \
            }                                                                                        \
            if ((res.code < REQUEST_HTTP_CONTINUE || res.code >= REQUEST_HTTP_INTERNAL_SERVER_ERROR) && (i != REQUEST_HTTP_RETRY_TIMES)) {    \
                res.code = 0;                                                                        \
                res.body.clear ();                                                                   \
                res.headers.clear ();                                                                \
                std::this_thread::sleep_for (std::chrono::milliseconds (i * 500));                   \
                continue;                                                                            \
            }                                                                                        \
            break;                                                                                   \
        }                                                                                            \


NS_IMPL_ISUPPORTS1(ncDataCopier, ncIDataCopier)

/* [notxpcom] void OnCopyData (in ucharPtr buf, in uint len); */
NS_IMETHODIMP_(void)
ncDataCopier::OnCopyData (unsigned char * buf, unsigned int len)
{
    _data.append (reinterpret_cast<char*>(buf), len);
    _length += len;
}

NS_IMPL_QUERY_INTERFACE1 (ncACSDataStore, ncIACSDataStore)

NS_IMETHODIMP_(nsrefcnt) ncACSDataStore::AddRef (void)
{
    return 1;
}

NS_IMETHODIMP_(nsrefcnt) ncACSDataStore::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (ncACSDataStore)

ncACSDataStore::ncACSDataStore ()
            : _ossClientPtr (0)
{
    NC_ACS_DATA_STORE_TRACE (_T("this: %p"), this);

    nsresult ret;
    _ossGateway = do_CreateInstance (OSS_GATEWAY_CONTRACTID, &ret);
    if (NS_FAILED (ret)) {
        THROW_E (ACS_DATA_STORE, ACS_DATA_STORE_INTERNAL_ERR,
            _T("Failed to create ossGateway instance: 0x%x"), ret);
    }
}

ncACSDataStore::~ncACSDataStore()
{
    NC_ACS_DATA_STORE_TRACE (_T("this: %p"), this);
}

void ncACSDataStore::createOSSClient ()
{
    if (nullptr == _ossClientPtr.get ()) {
        nsresult ret;
        nsCOMPtr<ncIOSSClient> ossClient = do_CreateInstance (NC_OSS_CLIENT_CONTRACTID, &ret);
        if (NS_FAILED (ret)){
            THROW_E (ACS_DATA_STORE, ACS_DATA_STORE_INTERNAL_ERR,
                     _T("Failed to create OSSClient: %s(0x%x)"),
                     String::toString ((int64)ret).getCStr (), (int64)ret);
        }
        _ossClientPtr.reset (new nsCOMPtr<ncIOSSClient>(ossClient));
    }
}

/* [notxpcom] String CreateAccountId (); */
NS_IMETHODIMP_(String) ncACSDataStore::CreateAccountId ()
{
    NC_ACS_DATA_STORE_TRACE (_T("this: %p begin"), this);

    ncObjectID id;
    id.GenerateObjectID ();

    NC_ACS_DATA_STORE_TRACE (_T("this: %p id: %s end"), this, id.GetString ().getCStr ());
    return String (id.GetString (), acsDataStorePoolAllocator);
}

/*[notxpcom] String GenerateObjectId ();*/
NS_IMETHODIMP_(String) ncACSDataStore::GenerateObjectId ()
{
    NC_ACS_DATA_STORE_TRACE (_T("this: %p begin"), this);

    ncObjectID id;
    id.GenerateObjectID ();

    NC_ACS_DATA_STORE_TRACE (_T("this: %p id: %s end"), this, id.GetString ().getCStr ());
    return String (id.GetString (), acsDataStorePoolAllocator);
}

/*[notxpcom] String GetAvailableOSSID ();*/
NS_IMETHODIMP_(String) ncACSDataStore::GetAvailableOSSID ()
{
    NC_ACS_DATA_STORE_TRACE (_T("this: %p begin"), this);

    String ossId;
    vector<OSSInfo> infos;
    _ossGateway->GetLocalOSSInfo(infos);
    for (int i = 0; i < infos.size(); i++)
    {
        if (i == 0){
            ossId = infos[i].id;
        }

        if (infos[i].bDefault == true)
        {
            ossId = infos[i].id;
            break;
        }
    }

    NC_ACS_DATA_STORE_TRACE (_T("this: %p end"), this);
    return String (ossId, acsDataStorePoolAllocator);
}

void CommonException (const String& method, const String& objectName, const ncOSSResponse& res)
{
    if (res.code == -1) {
        THROW_E (ACS_DATA_STORE, ACS_DATA_STORE_INTERNAL_ERR,
                 _T("Cannot connect to storage server (method:%s, object: %s)."),
                 method.getCStr (), objectName.getCStr ());
    }
    else if (res.code == 400) {
        THROW_E (ACS_DATA_STORE, ACS_DATA_STORE_INTERNAL_ERR,
                 _T("Request to storage server error, bad request (method:%s, object: %s, retcode: %d, msg: %s)."),
                 method.getCStr (), objectName.getCStr (), res.code, res.body.c_str ());
    }
    else if (res.code == 403) {
        THROW_E (ACS_DATA_STORE, ACS_DATA_STORE_INTERNAL_ERR,
                 _T("Request to storage server error, wrong authentication (method:%s, object: %s, retcode: %d, msg: %s)."),
                 method.getCStr (), objectName.getCStr (), res.code, res.body.c_str ());
    }
    else {
        THROW_E (ACS_DATA_STORE, ACS_DATA_STORE_INTERNAL_ERR,
                 _T("Request to storage server error (method:%s, object: %s, retcode: %d, msg: %s)."),
                 method.getCStr (), objectName.getCStr (), res.code, res.body.c_str ());
    }
}

/* NS_IMETHOD_(void) InitUpload(const String &prefix, const String & accountId, const String & objId, const String & ossId, ncUploadInfo & info) */
NS_IMETHODIMP_(void) ncACSDataStore::InitUpload (const String &prefix, const String& accountId, const String& objId, const String& ossId, ncUploadInfo& info)
{
    NC_ACS_DATA_STORE_TRACE (_T("this: %p accountId: %s, objId: %s, ossId: %s begin."),
                             this, accountId.getCStr (), objId.getCStr (), ossId.getCStr ());

    String objectName = prefix;
    if (prefix != "")
    {
        objectName.append ("/");
    }
    objectName.append (accountId);
    objectName.append ("/");
    objectName.append (objId);

    OSSUploadInfo tempInfo;
    _ossGateway->UploadInfo(objectName, ossId, tempInfo);
    info.partSize = tempInfo.partSize;
    info.uploadId = tempInfo.uploadId;

    NC_ACS_DATA_STORE_TRACE (_T("this: %p accountId: %s, objId: %s, ossId: %s , uploadId: %s end."),
                             this, accountId.getCStr (), objId.getCStr (), ossId.getCStr (), info.uploadId.getCStr ());
}

/* [notxpcom] void UploadBlock ([const] in StringRef prefix, [const] in StringRef accountId, [const] in StringRef objId, [const] in StringRef ossId, [const] in StringRef uploadId, [const] in StringRef content, in int64 offest, in int sn, in partInfoRef partInfo); */
NS_IMETHODIMP_(void) ncACSDataStore::UploadBlock(const String& prefix,
                                                 const String& accountId,
                                                 const String& objId,
                                                 const String& ossId,
                                                 const String& uploadId,
                                                 const String& content,
                                                 int64 offest,
                                                 int sn,
                                                 ncUploadPartInfo& partInfo)
{
    NC_ACS_DATA_STORE_TRACE (_T("this: %p accountId: %s, objId: %s, ossId: %s, uploadId: %s, offest: %lld, sn: %d begin"),
        this, accountId.getCStr (), objId.getCStr (), ossId.getCStr (), uploadId.getCStr (), offest, sn);

    createOSSClient ();

    // 上传到开放存储，数据块号从1开始；上传第一块之前，初始化分块上传
    ++sn;
    String objectName = prefix;
    if (prefix != "")
    {
        objectName.append ("/");
    }
    objectName.append (accountId);
    objectName.append ("/");
    objectName.append (objId);

    RequestInfo tempRes;
    _ossGateway->UploadPart(objectName, ossId, uploadId, sn, tempRes);

    // 根据返回值上传数据
    ncOSSResponse res;
    NC_ACSDS_REQUEST_TRY
        NC_ACSDS_REQUEST_RETRY_BEGIN
                (*_ossClientPtr)->Put (tempRes.url, toSTLString (content), tempRes.headers, 300L, res);
        NC_ACSDS_REQUEST_RETRY_END
        if (res.code != 200 && res.code != 201) {
            CommonException ("PUT", objId, res);
        }
    NC_ACSDS_REQUEST_CATCH ("PUT", objId, sn)

    for (map<string, string>::const_iterator iter = res.headers.begin (); iter != res.headers.end (); ++iter) {
        String key = toCFLString (iter->first);
        if (0 == key.compareIgnoreCase ("etag")) {
             partInfo.etag = iter->second;
             break;
        }
    }
    partInfo.size = content.getLength ();

    NC_ACS_DATA_STORE_TRACE (_T("this: %p accountId: %s, objId: %s, uploadId:%s, offest: %lld, sn: %d end"),
        this, accountId.getCStr (), objId.getCStr (), uploadId.getCStr (), offest, sn);
}

/* [notxpcom] void CompleteUpload ([const] in StringRef prefix, [const] in StringRef accountId, [const] in StringRef objId, [const] in StringRef ossId, [const] in StringRef uploadId, [const] in partInfoMapRef partInfos); */
NS_IMETHODIMP_(void) ncACSDataStore::CompleteUpload (const String& prefix,
                                                     const String& accountId,
                                                     const String& objId,
                                                     const String& ossId,
                                                     const String& uploadId,
                                                     const map<int, ncUploadPartInfo>& partInfos)
{
    NC_ACS_DATA_STORE_TRACE (_T("this: %p accountId: %s, objId: %s, ossId: %s, uploadId: %s begin."),
                             this, accountId.getCStr (), objId.getCStr (), ossId.getCStr (), uploadId.getCStr ());


    map<int, UploadPartInfo> OSSPartInfo;
    for (map<int, ncUploadPartInfo>::const_iterator iter = partInfos.begin (); iter != partInfos.end (); ++iter) {
        UploadPartInfo partInfo((iter->second).etag, (iter->second).size);
        OSSPartInfo.insert (make_pair ((iter->first), partInfo));
    }

    String objectName = prefix;
    if (prefix != "")
    {
        objectName.append ("/");
    }
    objectName.append (accountId);
    objectName.append ("/");
    objectName.append (objId);

    RequestInfo tempRes;
    _ossGateway->CompleteUpload(objectName, ossId, uploadId, OSSPartInfo, tempRes);

    createOSSClient ();
    ncOSSResponse res;
    NC_ACSDS_REQUEST_TRY
        NC_ACSDS_REQUEST_RETRY_BEGIN
                if (0 == tempRes.method.compare ("POST")) {
                    (*_ossClientPtr)->Post (tempRes.url, tempRes.body, tempRes.headers, 900L, res);
                }
                else {
                    (*_ossClientPtr)->Put (tempRes.url, tempRes.body, tempRes.headers, 300L, res);
                }
        NC_ACSDS_REQUEST_RETRY_END
        if (res.code != 200 && res.code != 201) {
            CommonException (tempRes.method.c_str (), objId, res);
        }
    NC_ACSDS_REQUEST_CATCH (tempRes.method.c_str (), objId, 1) // POST/PUT 1 代表上传索引文件

    NC_ACS_DATA_STORE_TRACE (_T("this: %p accountId: %s, objId: %s, uploadId: %s end."),
                             this, accountId.getCStr (), objId.getCStr (), uploadId.getCStr ());
}

/* [notxpcom] void ReadByOffest ([const] in StringRef prefix, [const] in StringRef fileId, in int64 offset, in int length, in int64 fileSize, in stlStringRef content); */
NS_IMETHODIMP_(void) ncACSDataStore::ReadByOffest (const String& prefix,
                                                   const String& fileId,
                                                   const String& ossId,
                                                   int64 offset,
                                                   int length,
                                                   int64 fileSize,
                                                   string& content)
{
    NC_ACS_DATA_STORE_TRACE (_T("this: %p fileId: %s, ossId: %s, offset: %lld, length: %d, fileSize: %lld begin"),
                             this, fileId.getCStr (), ossId.getCStr (), offset, length, fileSize);

    String accountId, objId;
    getAccountObjId (fileId, accountId, objId);

    String objectName = prefix;
    if (prefix != "")
    {
        objectName.append ("/");
    }
    objectName.append (accountId);
    objectName.append ("/");
    objectName.append (objId);

    RequestInfo tempRes;
    _ossGateway->GetDownLoadInfo(objectName, ossId, "", 3600, tempRes);

    offset = (offset < 0) ? 0 : offset;     // offset 从0开始
    length = (length < -1) ? -1 : length;   // length 支持 -1, 读取到文件结尾
    // offset 大于等于数据实际大小，返回; 获取的 length 为0，返回
    if ((offset >= fileSize) || (length == 0)) {
        return;
    }
    int64 rangeEnd = offset + length - 1;
    if ((rangeEnd + 1 > fileSize) || (length == -1)) {
        rangeEnd = fileSize - 1;
    }
    String rangeHeader;
    rangeHeader.format ("Range: bytes=%lld-%lld", offset, rangeEnd);
    tempRes.headers.push_back (toSTLString (rangeHeader));

    createOSSClient ();
    ncOSSResponse res;
    NC_ACSDS_REQUEST_TRY
        NC_ACSDS_REQUEST_RETRY_BEGIN
                (*_ossClientPtr)->Get (tempRes.url, tempRes.headers, 300L, res);
        NC_ACSDS_REQUEST_RETRY_END
        if (res.code != 200 && res.code != 206) {
            if (res.code == 404) {
                THROW_E (ACS_DATA_STORE, ACS_DATA_STORE_INTERNAL_ERR,
                         _T("Version object does not exist in storage server (%s)."),
                         fileId.getCStr ());
            }
            else if (res.code == 416) {
                String objectMsg;
                objectMsg.format ("fileId: %s, rangeBegin:%lld, rangeEnd:%lld",
                                  fileId.getCStr (), offset, rangeEnd);
                THROW_E (ACS_DATA_STORE, ACS_DATA_STORE_INTERNAL_ERR,
                         _T("Version object does not exist in storage server (%s)."),
                         objectMsg.getCStr ());
            }
            else {
                CommonException ("GET", fileId, res);
            }
        }
    NC_ACSDS_REQUEST_CATCH ("GET", objId, -1) // -1 代表上传索引文件
    content = res.body;

    NC_ACS_DATA_STORE_TRACE (_T("this: %p fileId: %s, offset: %lld, length: %d, fileSize:%lld end"),
                             this, fileId.getCStr (), offset, length, fileSize);
}

void ncACSDataStore::getAccountObjId (const String& id, String& accountId, String& objId)
{
    size_t index = id.findFirstOf ('/');
    if (index != String::NO_POSITION) {
        accountId = id.subString (0, index);
    }
    else {
        THROW_E (ACS_DATA_STORE, ACS_DATA_STORE_INVALID_ARGUMENT,
                 _T("invalid fileId : %s"), id.getCStr ());
    }

    index = id.findLastOf ('/');
    if (index != String::NO_POSITION) {
        objId = id.subString (index + 1, id.getLength () - index);
    }
    else {
        THROW_E (ACS_DATA_STORE, ACS_DATA_STORE_INVALID_ARGUMENT,
                 _T("invalid fileId : %s"), id.getCStr ());
    }
}
