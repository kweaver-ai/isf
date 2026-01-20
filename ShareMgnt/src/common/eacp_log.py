from src.common import global_info
from src.common.business_date import BusinessDate
from src.common.sharemgnt_logger import ShareMgnt_Log
from src.driven.service_access.authentication import AuthenticaitonDriven
import uuid
LOG_MESSAGE_TYPE = {global_info.LOG_TYPE_LOGIN:  "as.audit_log.log_login", global_info.LOG_TYPE_OPERATION: "as.audit_log.log_operation", global_info.LOG_TYPE_MANAGE: "as.audit_log.log_management"}
USER_TYPE = {global_info.USER_TYPE_AUTH: "authenticated_user", global_info.USER_TYPE_ANONY: "anonymous_user", global_info.USER_TYPE_APP: "app", global_info.USER_TYPE_INTER: "internal_service"}


def eacp_log(userId=None, log_type=None, user_type=None, level=None, op_type=None, msg=None, ex_msg=None, ip=None, mac=None, raise_ex=False, log_item=None):
    """
    记录审计日志
    """
    authentication = AuthenticaitonDriven()
    if log_item:
        authentication.send_audit_log(LOG_MESSAGE_TYPE[log_type], log_item)
    else:
        date = int(BusinessDate.time()) * 1000000
        out_biz_id = str(uuid.uuid1())
        log_item = {"user_id": userId, "level": level, "op_type": op_type, "date": date, "ip": ip, "mac": mac, "msg": msg, "ex_msg": ex_msg, "out_biz_id": out_biz_id, "user_type": USER_TYPE[user_type]}

        # 内部账户需要设置用户名为用户ID
        if user_type == global_info.USER_TYPE_INTER:
            log_item["user_name"] = userId
        authentication.send_audit_log(LOG_MESSAGE_TYPE[log_type], log_item)
