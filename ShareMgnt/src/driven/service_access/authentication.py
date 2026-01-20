import json

from src.common.http import send_request
from src.common.config import Config
from src.common.lib import raise_exception
from src.common.sharemgnt_logger import ShareMgnt_Log
import urllib.parse


class AuthenticaitonDriven:
    def send_audit_log(self, topic, log_item):
        """发送审计日志"""
        try:
            url = f'http://{Config.authentication_config["host"]}:{Config.authentication_config["port"]}/api/authentication/v1/audit-log'

            # 定义请求体
            request_body = {
                "topic": topic,
                "message": log_item
            }

            code, content = send_request(url, data=request_body, method="POST")
            return

        except Exception as ex:
            ShareMgnt_Log("Authenticaiton send request error: %s", str(ex))
            raise_exception("Authenticaiton send request error")