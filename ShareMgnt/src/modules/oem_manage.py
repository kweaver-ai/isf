from src.common.db.connector import DBConnector
from src.common.lib import raise_exception
from ShareMgnt.ttypes import ncTOEMInfo, ncTShareMgntError


class OEMManage(DBConnector):
    """
    OEM manage
    """
    def __init__(self):
        """
        init
        """
        pass

    def _check_section(self, section):
        if len(section) <= 0 or len(section) > 32:
            raise_exception(exp_msg=_("IDS_INVALID_SECTION"),
                            exp_num=ncTShareMgntError.NCT_INVALID_SECTION)

    def _check_option(self, option):
        if len(option) <= 0 or len(option) > 32:
            raise_exception(exp_msg=_("IDS_INVALID_OPTION"),
                            exp_num=ncTShareMgntError.NCT_INVALID_OPTION)

    def _check_value(self, value):
        if len(value) < 0 or len(value) > 16777215:
            raise_exception(exp_msg=_("IDS_INVALID_VALUE"),
                            exp_num=ncTShareMgntError.NCT_INVALID_VALUE)

    def set_config(self, oemInfo):
        """
        设置config
        """
        self._check_section(oemInfo.section)
        self._check_option(oemInfo.option)
        self._check_value(oemInfo.value)

        str_sql = """
        select f_section from t_oem_config
        where f_section = %s and f_option = %s
        """
        if self.w_db.one(str_sql, oemInfo.section, oemInfo.option):
            str_sql = """
            update t_oem_config set f_value = %s
            where f_section = %s and f_option = %s
            """
            self.w_db.query(str_sql, oemInfo.value, oemInfo.section, oemInfo.option)
        else:
            str_sql = """
            insert into t_oem_config(f_section, f_option, f_value)
            values(%s, %s, %s)
            """
            self.w_db.query(str_sql, oemInfo.section, oemInfo.option, oemInfo.value)

        # 只要开启用户弹出协议，需要重新清空所有用户已经同意的协议状态
        if oemInfo.option == "autoPopUserAgreement" and oemInfo.value == b"true":
            str_sql = """
            update t_user set f_agreed_to_terms_of_use = 0
            """
            self.w_db.query(str_sql)

    def get_config_by_section(self, section):
        """
        根据section获取批量配置
        """
        self._check_section(section)

        str_sql = """
        select f_option,f_value from t_oem_config
        where f_section = %s
        """
        results = self.w_db.all(str_sql, section)
        if len(results) == 0:
            raise_exception(exp_msg=_("IDS_SECTION_NOT_EXISTS"),
                            exp_num=ncTShareMgntError.NCT_SECTION_NOT_EXISTS)

        retInfos = []
        for result in results:
            tmpInfo = ncTOEMInfo()
            tmpInfo.section = section
            tmpInfo.option = result["f_option"]
            if not isinstance(result["f_value"], bytes):
                result["f_value"] = bytes(result["f_value"], encoding="utf8")
            tmpInfo.value = result["f_value"]
            retInfos.append(tmpInfo)

        return retInfos

    def get_config_by_option(self, section, option):
        """
        根据section和option获取配置
        """
        self._check_section(section)
        self._check_option(option)

        str_sql = """
        select f_value from t_oem_config
        where f_section = %s and f_option = %s
        """
        result = self.w_db.one(str_sql, section, option)
        if result:
            if not isinstance(result["f_value"], bytes):
                result["f_value"] = bytes(result["f_value"], encoding="utf8")
            return result["f_value"]
        else:
            raise_exception(exp_msg=_("IDS_OPTION_NOT_EXISTS"),
                            exp_num=ncTShareMgntError.NCT_OPTION_NOT_EXISTS)
