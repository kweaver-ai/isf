import json

from src.common.http import send_request
from src.common.config import Config
from src.common.lib import raise_exception
from src.common.sharemgnt_logger import ShareMgnt_Log
import urllib.parse


class OssgatewayDriven:
    def get_local_storages(self):
        """获取所有存储信息"""
        try:
            url = f'{Config.ossgateway_config["protocol"]}://{Config.ossgateway_config["host"]}:{Config.ossgateway_config["port"]}/api/ossgateway/v1/local-storages?enabled=true'

            code, content = send_request(url)
            return json.loads(content)

        except Exception as ex:
            ShareMgnt_Log("ossgateway send request error: %s", str(ex))
            raise_exception("ossgateway send request error")

    def get_upload_info(self, oss_id, key, storage_prefix=False):
        """获取上传url信息"""
        try:
            url = f'{Config.ossgateway_config["protocol"]}://{Config.ossgateway_config["host"]}:{Config.ossgateway_config["port"]}/api/ossgateway/v1/upload/{oss_id}/{key}'
            reqArgs = {}
            reqArgs["request_method"] = "PUT"
            reqArgs["internal_request"] = "true"
            if storage_prefix:
                reqArgs["storage_prefix"] = "true"

            url = url + "?" + urllib.parse.urlencode(reqArgs)
            code, content = send_request(url)
            return code, json.loads(content)

        except Exception as ex:
            ShareMgnt_Log("ossgateway send request error: %s", str(ex))
            raise_exception("ossgateway send request error")

    def get_download_info(self, oss_id, key, expires_time=None, save_name=None, storage_prefix=False):
        """获取下载url信息"""
        try:
            url = f'{Config.ossgateway_config["protocol"]}://{Config.ossgateway_config["host"]}:{Config.ossgateway_config["port"]}/api/ossgateway/v1/download/{oss_id}/{key}'
            reqArgs = {}
            reqArgs["type"] = "query_string"
            reqArgs["internal_request"] = "true"
            if save_name is not None:
                reqArgs["save_name"] = save_name
            if expires_time is not None:
                reqArgs["expires"] = expires_time
            if storage_prefix:
                reqArgs["storage_prefix"] = "true"

            url = url + "?" + urllib.parse.urlencode(reqArgs)
            code, content = send_request(url)
            return code, json.loads(content)

        except Exception as ex:
            ShareMgnt_Log("ossgateway send request error: %s", str(ex))
            raise_exception("ossgateway send request error")

    def get_delete_info(self, oss_id, key, storage_prefix=False):
        """获取删除url信息"""
        try:
            url = f'{Config.ossgateway_config["protocol"]}://{Config.ossgateway_config["host"]}:{Config.ossgateway_config["port"]}/api/ossgateway/v1/delete/{oss_id}/{key}'
            reqArgs = {}
            reqArgs["internal_request"] = "true"
            if storage_prefix:
                reqArgs["storage_prefix"] = "true"

            url = url + "?" + urllib.parse.urlencode(reqArgs)
            code, content = send_request(url)
            return code, json.loads(content)

        except Exception as ex:
            ShareMgnt_Log("ossgateway send request error: %s", str(ex))
            raise_exception("ossgateway send request error")

    def get_as_storage_info(self):
        """"获取as存储"""
        try:
            url = f'{Config.ossgateway_config["protocol"]}://{Config.ossgateway_config["host"]}:{Config.ossgateway_config["port"]}/api/ossgateway/v1/objectstorageinfo?app=as'

            code, content = send_request(url)
            return code, json.loads(content)

        except Exception as ex:
            ShareMgnt_Log("ossgateway send request error: %s", str(ex))
            raise_exception("ossgateway send request error")

    def get_storage_info_by_id(self, oss_id):
        """"通过id获取存储信息"""
        try:
            url = f'{Config.ossgateway_config["protocol"]}://{Config.ossgateway_config["host"]}:{Config.ossgateway_config["port"]}/api/ossgateway/v1/objectstorageinfo/{oss_id}'

            code, content = send_request(url)
            return code, json.loads(content)

        except Exception as ex:
            ShareMgnt_Log("ossgateway send request error: %s", str(ex))
            raise_exception("ossgateway send request error")
