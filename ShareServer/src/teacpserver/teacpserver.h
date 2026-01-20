#ifndef __T_EACP_SERVER_H
#define __T_EACP_SERVER_H

#include <abprec.h>
#include <ncutil/ncutil.h>

// 语言资源加载器
extern IResourceLoader* teacpserverResLoader;
#define LOAD_STRING(strID)                        \
    teacpserverResLoader->loadString (strID)

// 错误提供者名称
#define T_EACP_SERVER_ERR_PROVIDER_NAME            _T("eacp thrift server")

// 内部错误
#define FAILED_TO_CREATE_ACS_PROCESSOR              0x00001001L
#define FAILED_TO_CREATE_TUSERM                     0x00001002L
#define FAILED_TO_CREATE_ACS_DOC_MANANGER           0x00001003L
#define FAILED_TO_CREATE_ACS_PERM_MANANGER          0x00001004L
#define FAILED_TO_CREATE_ACS_TOKEN_MANANGER         0x00001005L
#define FAILED_TO_CREATE_ACS_OWNER_MANANGER         0x00001006L
#define FAILED_TO_CREATE_ACS_CONN_MANANGER          0x00001007L
#define FAILED_TO_CREATE_ACS_SHAREMGNT              0x00001008L
#define FAILED_TO_CREATE_DB_DOC_MANAGER             0x0000100AL
#define FAILED_TO_CREATE_ACS_CAUTH_MANAGER          0x0000100BL
#define CONVERT_PATH_ERROR                          0x0000100EL
#define FAILED_TO_CREATE_ACS_OAUTH_MANAGER          0x0000100FL
#define FAILED_TO_CREATE_DB_LOCK_MANAGER            0x00001010L
#define FAILED_TO_CREATE_ACS_LOCK_MANAGER           0x00001011L
#define FAILED_TO_CREATE_ACS_DEVICE_MANAGER         0x00001012L
#define FAILED_TO_CREATE_ACS_CSF_MANAGER            0x00001013L
#define FAILED_TO_CREATE_ACS_AUDIT_MANAGER          0x00001015L
#define FAILED_TO_CREATE_DB_CONF_MANANGER           0x00001018L
#define FAILED_TO_CREATE_ACS_MESSAGE_MANANGER       0x00001019L
#define FAILED_TO_CREATE_EMD_METADATA_MANAGER       0x0000101BL
#define FAILED_TO_CREATE_ACS_ENTRY_DOC_IOC          0x0000101CL
#define FAILED_TO_CREATE_OSS_CLIENT                 0x0000101DL

// 应用逻辑错误
#define FROBID_QUOTA_SIZE_DECREASED                 0x00002004L
#define EXCEED_MAX_CUSTOM_DOC_TYPE_COUN             0x00002005L
#define INVALID_USER_TOTAL_QUOTA_BYTES              0x00002006L
#define QUOTA_NUM_CANT_BE_0                         0x00002009L
#define INVALID_DOC_LIB_NAME                        0x0000200AL
#define CUSTOM_DOC_NAME_EXISTS                      0x0000200BL
#define DOC_ID_NOT_EXISTS                           0x0000200EL
#define INVALID_OAUTH_INFO                          0x00002013L
#define OAUTH_INFO_NOT_SET                          0x00002014L

#define GET_CONNECT_COUNT_ERROR                     0x00003001L
#define FAILED_TO_CREATE_ACS_LICENSE_MANAGER        0x00003002L
#define LICENSE_INFO_ERROR                          0x00003003L

// trace module
#define NC_T_EACP_SERVER_TRACE_MODULE_NAME          _T ("teacpserver")

#ifdef __WINDOWS__
#define NC_T_EACP_SERVER_TRACE(...)                 TRACE_EX2 (NC_T_EACP_SERVER_TRACE_MODULE_NAME, __VA_ARGS__)
#else
#define NC_T_EACP_SERVER_TRACE(args...)             TRACE_EX2 (NC_T_EACP_SERVER_TRACE_MODULE_NAME, args)
#endif

#endif // __T_EACP_SERVER_H
