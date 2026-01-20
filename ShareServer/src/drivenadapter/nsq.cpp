/***************************************************************************************************
nsq.cpp:
    Copyright (c) Eisoo Software Inc. (2020), All rights reserved.

Purpose:
    NSQ manager

Author:
    Will.lv@aishu.cn

Creating Time:
    2020-08-26
***************************************************************************************************/
#include <abprec.h>
#include "nsq.h"
#include "drivenadapter.h"
#include <dataapi/ncJson.h>
#include <ossclient/public/ncIOSSClient.h>
#include <ossclient/ossclient.h>
#include "serviceAccessConfig.h"
#include <dataapi/dataapi.h>
#include <ncutil/ncBusinessDate.h>


// 原子权限
enum PermValue {
    ACS_AP_DISPLAY      =       0x00000001,       // 显示
    ACS_AP_PREVIEW      =       0x00000002,       // 预览
    ACS_AP_READ         =       0x00000004,       // 下载 枚举值名称未改
    ACS_AP_CREATE       =       0x00000008,       // 新建
    ACS_AP_EDIT         =       0x00000010,       // 修改
    ACS_AP_DELETE       =       0x00000020,       // 删除

    ACS_CP_MIN          =       0x00000001,       // 权限配置最小值
    ACS_CP_MAX          =       0x0000003F,       // 权限配置最大值
};

// 访问者类型
enum AccessorType {
    ACS_USER            =       0x00000001,         // 用户
    ACS_DEPARTMENT      =       0x00000002,         // 部门
    ACS_CONTACTOR       =       0x00000003,         // 联系人
    ACS_ANONYMOUS_USER  =       0x00000004,         // 匿名用户
    ACS_GROUP           =       0x00000005,         // 用户组
};

enum PermPair {
    OPER_NONE,
    OPER_ADD,
    OPER_UPDATE,
    OPER_DELETE
};

enum DocType {
    ACS_USER_DOC = 1,
    ACS_DEPARTMENT_DOC = 2,
    ACS_CUSTOM_DOC = 3,
    ACS_SHARE_DOC = 4,
    ACS_ARCHIVE_DOC = 5,
    ACS_KNOWLEDGE_DOC = 6,
};

/* Implementation file */
NS_IMPL_QUERY_INTERFACE1 (nsq, nsqInterface)

// protected
NS_IMETHODIMP_(nsrefcnt) nsq::AddRef (void)
{
    return 1;
}

// protected
NS_IMETHODIMP_(nsrefcnt) nsq::Release (void)
{
    return 1;
}

AB_DEFINE_THREADSAFE_SINGLETON_NO_POOL (nsq)

nsq::nsq (void)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("this: %p"), this);

    _eventTopicMap[ncNSQEventType::NSQ_DOC_SHARE_REALNAME_APPLY] = _T("core.audit.share.realname.apply");   //审核申请 实名
    _eventTopicMap[ncNSQEventType::NSQ_DOC_SHARE_ANONYMOUS_APPLY] = _T("core.audit.share.anonymous.apply");  //审核申请 匿名
    _eventTopicMap[ncNSQEventType::NSQ_DOC_SHARE_CANCEL] = _T("core.audit.share.cancel");                    // 取消审核申请 实名 匿名
    _eventTopicMap[ncNSQEventType::NSQ_SHARE_PERM_CHANGE] = _T("core.share.perm.change");                    // 权限变更
    _eventTopicMap[ncNSQEventType::NSQ_CORE_USER_ANONYMOUS_OUT_OF_SCOPE] = _T("core.user.anonymous.out_of_scope");

    _mqClient = ncMQClient::GetConnectorFromFile("/sysvol/conf/service_conf/mq_config.yaml");
}

nsq::~nsq (void)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("this: %p"), this);
}

// 格式化NSQ消息
String nsq::FormatNSQMsg (const ncNSQEventType& nsqType, const NSQMsg& nsqMsg)
{
    JSON::Value contentJson;
    // 实名、匿名 共享
    if (nsqType == ncNSQEventType::NSQ_DOC_SHARE_REALNAME_APPLY
        || nsqType == ncNSQEventType::NSQ_DOC_SHARE_ANONYMOUS_APPLY) {

        contentJson["apply_id"] = nsqMsg.applyId.getCStr ();
        if (!nsqMsg.conflictApplyId.isEmpty ()) {
            contentJson["conflict_apply_id"] = nsqMsg.conflictApplyId.getCStr ();
        }
        contentJson["user_id"] = nsqMsg.userId.getCStr ();
        // docInfo
        JSON::Object& docInfo = contentJson["doc"].o ();
        docInfo["id"] = nsqMsg.docId.getCStr ();
        docInfo["path"] = nsqMsg.docPath.getCStr ();
        docInfo["doc_lib_type"] = docLibType2String(nsqMsg.docLibType).getCStr();
        docInfo["type"] = nsqMsg.isFile? "file":"folder";
        docInfo[nsqMsg.isFile? "csf_level":"max_csf_level"] = nsqMsg.docCsfLevel;

        if (nsqMsg.applyType == ncNSQApplyType::NSQ_APPLY_PERM) {
            contentJson["type"] = applyType2String(nsqMsg.applyType).getCStr ();
            // accessor 访问者
            JSON::Object& accessor = contentJson["accessor"].o ();
            accessor["id"] = nsqMsg.accessorId.getCStr ();
            accessor["name"] = nsqMsg.accessorName.getCStr ();
            accessor["type"] = accessorType2String(nsqMsg.accessorType).getCStr ();
            // 权限操作
            contentJson["operation"] = opType2String(nsqMsg.operation).getCStr ();
            if (nsqMsg.expiresAt == -1 ){
                contentJson["expires_at"] = "-1";
            }else {
                contentJson["expires_at"] = ticks2String(nsqMsg.expiresAt).getCStr ();
            }
            // 权限信息
            JSON::Object& perm = contentJson["perm"].o ();
            JSON::Array& allowArray = perm["allow"].a ();
            JSON::Array& denyArray = perm["deny"].a ();
            intToPermArray(nsqMsg.allowValue, allowArray);
            intToPermArray(nsqMsg.denyValue, denyArray);
        }else if (nsqMsg.applyType == ncNSQApplyType::NSQ_APPLY_INHERIT){
            contentJson["type"] = applyType2String(nsqMsg.applyType).getCStr ();
            contentJson["inherit"] = nsqMsg.inherit;
        }else if (nsqMsg.applyType == ncNSQApplyType::NSQ_APPLY_OWNER){
            contentJson["type"] = applyType2String(nsqMsg.applyType).getCStr ();
            // accessor 访问者
            JSON::Object& accessor = contentJson["accessor"].o ();
            accessor["id"] = nsqMsg.accessorId.getCStr ();
            accessor["name"] = nsqMsg.accessorName.getCStr ();
            accessor["type"] = accessorType2String(nsqMsg.accessorType).getCStr ();
            // 操作
            contentJson["operation"] = opType2String(nsqMsg.operation).getCStr ();
            if (nsqMsg.expiresAt == -1 ){
                contentJson["expires_at"] = "-1";
            }else {
                contentJson["expires_at"] = ticks2String(nsqMsg.expiresAt).getCStr ();
            }
        }else if (nsqMsg.applyType == ncNSQApplyType::NSQ_APPLY_ANONYMOUS){
            // 操作
            contentJson["operation"] = opType2String(nsqMsg.operation).getCStr ();
            // linkId
            contentJson["link_id"] = nsqMsg.linkId.getCStr ();
            contentJson["title"] = nsqMsg.title.getCStr ();
            contentJson["password"] = nsqMsg.password.getCStr ();
            contentJson["access_limit"] = nsqMsg.accessLimit;
            if (nsqMsg.expiresAt == -1 ){
                contentJson["expires_at"] = "-1";
            }else {
                contentJson["expires_at"] = ticks2String(nsqMsg.expiresAt).getCStr ();
            }
            JSON::Array& perm = contentJson["perm"].a ();
            intToPermArray(nsqMsg.allowValue, perm);
        }
    }else if (nsqType == ncNSQEventType::NSQ_DOC_SHARE_CANCEL){
        JSON::Array& cancelApplyArry = contentJson["apply_ids"].a ();
        // 取消审核申请
        for (auto it = nsqMsg.cancelApplyIds.begin(); it != nsqMsg.cancelApplyIds.end(); it++){
            cancelApplyArry.push_back (it->getCStr ());
        }
    }else if (nsqType == ncNSQEventType::NSQ_SHARE_PERM_CHANGE){
        contentJson["doc_id"] = nsqMsg.docId.getCStr ();
    }else if (nsqType == ncNSQEventType::NSQ_CORE_USER_ANONYMOUS_OUT_OF_SCOPE){
        contentJson["id"] = nsqMsg.userId.getCStr();
    }else {
        THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR, _T("MessageType is wrong"));
    }
    string content;
    JSON::Writer::write(contentJson.o (), content);
    return toCFLString (content);
}

// 发送一个MQ消息
// private
void nsq::Publish (const ncNSQEventSPtr eventSPtr)
{
     // 过滤空消息
    if (eventSPtr->nsqMsg.isEmpty ()) {
        return;
    }

    if (_eventTopicMap.count(eventSPtr->eventType)){
        try {
            _mqClient->Pub(_eventTopicMap[eventSPtr->eventType], eventSPtr->nsqMsg);
        } catch (Exception& e) {
            THROW_E (ACS_DRIVEN_ADAPTER, DRIVEN_HTTP_SERVER_INTERNAL_ERR,
                _T("MQ pub message failed, topic %s, errorId %d, err %s"), _eventTopicMap[eventSPtr->eventType].getCStr (),
                e.getErrorId (), e.toFullString ().getCStr ());
        }
    }
}

/* [notxpcom] void PublishNSQMessage ([const] in ncNSQEventTypeRef nsqType, [const] in NSQMsgRef nsqMsg); */
NS_IMETHODIMP_(void) nsq::PublishNSQMessage (const ncNSQEventType& nsqType, const NSQMsg& nsqMsg)
{
    NC_DRIVEN_ADAPTER_TRACE (_T("[BEGIN] this: %p"), this);

    String nsqMsgStr = FormatNSQMsg (nsqType, nsqMsg);
    auto eventSPtr = std::make_shared<ncNSQEvent> (nsqType, nsqMsgStr);
    Publish (eventSPtr);

    NC_DRIVEN_ADAPTER_TRACE (_T("[END] this: %p"), this);
}

void nsq::intToPermArray (int perm, JSON::Array& permArray)
{
    if (perm & PermValue::ACS_AP_DISPLAY) {
        permArray.push_back (_T("display"));
    }
    if (perm & PermValue::ACS_AP_PREVIEW) {
        permArray.push_back (_T("preview"));
    }
    if (perm & PermValue::ACS_AP_READ) {
        permArray.push_back (_T("download"));
    }
    if (perm & PermValue::ACS_AP_CREATE) {
        permArray.push_back (_T("create"));
    }
    if (perm & PermValue::ACS_AP_EDIT) {
        permArray.push_back (_T("modify"));
    }
    if (perm & PermValue::ACS_AP_DELETE) {
        permArray.push_back (_T("delete"));
    }
}

String nsq::docLibType2String(int docType)
{
    String result;
    if (docType == DocType::ACS_USER_DOC) {
        result = toCFLString("user_doc_lib");
    }else if (docType == DocType::ACS_DEPARTMENT_DOC) {
        result = toCFLString("department_doc_lib");
    }else if (docType == DocType::ACS_CUSTOM_DOC){
        result = toCFLString("custom_doc_lib");
    }else if (docType == DocType::ACS_KNOWLEDGE_DOC){
        result = toCFLString("knowledge_doc_lib");
    }
    return result;
}

String nsq::accessorType2String(int accessorType)
{
    String result;
    if(accessorType == AccessorType::ACS_USER){
        result =  toCFLString("user");
    } else if (accessorType == AccessorType::ACS_DEPARTMENT){
        result =  toCFLString("department");
    } else if (accessorType == AccessorType::ACS_CONTACTOR){
        result =  toCFLString("contactor");
    } else if (accessorType == AccessorType::ACS_GROUP){
        result =  toCFLString("group");
    }
    return result;
}

String nsq::opType2String(int opType)
{
    String result;
    if (opType == PermPair::OPER_ADD){
        result = toCFLString("create");
    }else if (opType == PermPair::OPER_UPDATE) {
        result = toCFLString("modify");
    }else if (opType == PermPair::OPER_DELETE){
        result = toCFLString("delete");
    }
    return result;
}

String nsq::applyType2String(ncNSQApplyType applyType)
{
    String result;
    if (applyType == ncNSQApplyType::NSQ_APPLY_PERM){
        result = toCFLString("perm");
    }else if (applyType == ncNSQApplyType::NSQ_APPLY_OWNER) {
        result = toCFLString("owner");
    }else if (applyType == ncNSQApplyType::NSQ_APPLY_INHERIT){
        result = toCFLString("inherit");
    }
    return result;
}

String nsq::ticks2String(int64 ticks)
{
    Date date (ticks);
    Date localCur (date.getLocalTime ());

    String ptStr;
    ptStr.format (_T("%04d-%02d-%02d %02d:%02d:%02d"),
        (int)localCur.getYear (), (int)localCur.getMonth (), (int)localCur.getDay (),
        (int)localCur.getHours (), (int)localCur.getMinutes (), (int)localCur.getSeconds ());

    return ptStr;
}
