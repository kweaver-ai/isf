#!/usr/bin/python3
# -*- coding:utf-8 -*-
# pylint: disable=C0103,C0111
"""
Ldap manage
"""
import uuid
import ldap
from ldap.controls import SimplePagedResultsControl
from src.common.lib import raise_exception
from ShareMgnt.ttypes import ncTShareMgntError
from src.common.encrypt.simple import eisoo_rsa_decrypt
from EThriftException.ttypes import ncTException

WINDOWS_AD = 1
OTHER_LDAP = 2
TIME_OUT = 3.0  # 连接超时时间


class DomainOuInfo(object):
    """
    域组织信息结构
    """

    def __init__(self, ou_name=None, third_id=None,
                 dn=None, server_type=0):
        self.ou_name = ou_name
        self.third_id = third_id
        self.dn = dn
        self.server_type = server_type


class DomainUserInfo(object):
    """
    域用户信息结构
    """

    def __init__(self, login_name=None, display_name=None,
                 email=None, password=None, third_id=None,
                 dn=None, ou_dn=None, server_type=0, status=None, idcard_number=None, tel_number=None):
        self.login_name = login_name
        self.display_name = display_name
        self.email = email
        self.password = password
        self.third_id = third_id
        self.dn = dn
        self.ou_dn = ou_dn
        self.server_type = server_type
        self.status = status
        self.idcard_number = idcard_number
        self.tel_number = tel_number


def dn2name(dn):
    """
    ldap dn(Distinguished Name) 转换为域名
    DC=xxx,DC=abc,DC=com => xxx.abc.com
    不是dn则不转换
    """
    if not ldap.dn.is_dn(dn):
        return dn

    name_list = []
    result = ldap.dn.explode_dn(dn, flags=ldap.DN_FORMAT_LDAPV3)
    for rdn in result:
        key, val = rdn.split("=", 1)
        if key.upper() == "DC":
            name_list.append(val)
    return ".".join(name_list)


def name2dn(name):
    """
    域名转换为 ldap dn(Distinguished Name)
    xxx.abc.com => DC=xxx,DC=abc,DC=com
    如果是dn则不转换
    """
    if ldap.dn.is_dn(name):
        return name

    dn_list = []
    name_list = name.split(".")
    for val in name_list:
        dn_list.append("DC=%s" % val.strip())
    return ",".join(dn_list)


class LdapManage(object):
    """
    ldap manage
    """

    def __init__(self, server_ip, admin, pwd, port=389, domain_name="", use_ssl=None, pwd_enc=True):
        """
        初始化函数
        默认端口port:
            简单连接：389
            SSL连接: 636
        """
        if not (server_ip and admin and pwd):
            raise_exception(_("domain check fail"))

        self.ip = server_ip
        self.admin = admin
        if pwd_enc:
            self.pwd = bytes.decode(eisoo_rsa_decrypt(pwd))
        else:
            self.pwd = pwd
        self.conn = None
        if not use_ssl:
            self.simple_port = port
            self.simple_bind()
        else:
            self.ssl_port = port
            self.ssl_bind()
        self._server_attr = None
        self._domain_dn = None
        self._domain_name = None
        self._domain_type = None
        self._support_paged = None
        self._supply_domain_name = domain_name

    def __del__(self):
        """
        析构函数
        """
        try:
            self.conn.unbind_s()
        except Exception:
            pass

    @property
    def server_attr(self):
        if self._server_attr is not None:
            return self._server_attr

        result = self.conn.search_s("", ldap.SCOPE_BASE, 'objectClass=*',
                                    ['*', '+'])
        self._server_attr = result[0][1]
        return self._server_attr

    @property
    def domain_dn(self):
        if self._domain_dn is not None:
            if isinstance(self._domain_dn, bytes):
                self._domain_dn = bytes.decode(self._domain_dn)
            return self._domain_dn

        if "forestFunctionality" in self.server_attr:
            if "defaultNamingContext" in self.server_attr:
                self._domain_dn = self.server_attr["defaultNamingContext"][0]
            else:
                self._domain_dn = self.server_attr["rootDomainNamingContext"][0]
        elif "namingContexts" in self.server_attr:
            self._domain_dn = self.server_attr["namingContexts"][0]
        else:
            self._domain_dn = ""
        if isinstance(self._domain_dn, bytes):
            self._domain_dn = bytes.decode(self._domain_dn)
        return self._domain_dn

    @property
    def domain_name(self):
        if self._domain_name is not None:
            return self._domain_name

        domain_substrs = []
        if self.domain_dn != "":
            values = ldap.dn.explode_dn(self.domain_dn,
                                        flags=ldap.DN_FORMAT_LDAPV3)

            for value in values:
                tmpstrs = value.split("=")
                if tmpstrs[0] == 'DC' or tmpstrs[0] == 'dc':
                    domain_substrs.append(tmpstrs[1])

        # 如果通过解析域根信息后，仍然无法获取，则取界面上配置的名字
        if not domain_substrs:
            if self._supply_domain_name:
                values = self._supply_domain_name.split(",")
                for value in values:
                    domain_substrs.append(value.split("=")[-1])

        self._domain_name = ".".join(domain_substrs)
        return self._domain_name

    @property
    def domain_type(self):
        if self._domain_type is not None:
            return self._domain_type

        if "forestFunctionality" in self.server_attr:
            self._domain_type = WINDOWS_AD
        else:
            self._domain_type = OTHER_LDAP

        return self._domain_type

    @property
    def support_paged(self):
        if self._support_paged is not None:
            return self._support_paged

        if "supportedControl" in self.server_attr:
            # The OID for Simple Paged Results control
            oid1 = b"1.2.840.113556.1.4.319" in self.server_attr["supportedControl"]
            self._support_paged = oid1
        elif "IBM Lotus Software" in self.server_attr.get("vendorname", []):
            # IBM Lotus
            self._support_paged = False
        else:
            self._support_paged = False

        return self._support_paged

    def error_handle(self, ex, ip):
        """
        ldap连接异常处理,根据异常类型判断登陆错误
        """
        if isinstance(ex, ldap.SERVER_DOWN):
            raise_exception(exp_msg=_("domain server unavailable"),
                            exp_num=ncTShareMgntError.NCT_DOMAIN_SERVER_UNAVAILABLE,
                            exp_detail=ip)
        elif isinstance(ex, ldap.INVALID_CREDENTIALS):
            raise_exception(exp_msg=_("invalid account or password"),
                            exp_num=ncTShareMgntError.NCT_INVALID_ACCOUNT_OR_PASSWORD)
        else:
            raise_exception(exp_msg=str(ex),
                            exp_num=ncTShareMgntError.NCT_DOMAIN_ERROR)

    def simple_bind(self):
        """
        ldap 简单连接
        """
        try:
            host = "ldap://%s:%d" % (self.ip, self.simple_port)
            self.conn = ldap.initialize(host)
            self.conn.protocol_version = ldap.VERSION3
            self.conn.set_option(ldap.OPT_REFERRALS, 0)
            self.conn.set_option(ldap.OPT_NETWORK_TIMEOUT, TIME_OUT)
            self.conn.simple_bind_s(self.admin, self.pwd)
        except ldap.LDAPError as ex:
            self.error_handle(ex, self.ip)

    def ssl_bind(self):
        """
        获得LDAP 的ssl连接
        """
        try:
            ldap.set_option(ldap.OPT_X_TLS_REQUIRE_CERT, ldap.OPT_X_TLS_NEVER)
            host = "ldaps://%s:%d" % (self.ip, self.ssl_port)
            self.conn = ldap.initialize(host)
            self.conn.set_option(ldap.OPT_REFERRALS, 0)
            self.conn.set_option(ldap.OPT_NETWORK_TIMEOUT, TIME_OUT)
            self.conn.set_option(ldap.OPT_PROTOCOL_VERSION, 3)
            self.conn.set_option(ldap.OPT_X_TLS, ldap.OPT_X_TLS_DEMAND)
            self.conn.set_option(ldap.OPT_X_TLS_DEMAND, True)
            self.conn.set_option(ldap.OPT_DEBUG_LEVEL, 255)
            self.conn.simple_bind_s(self.admin, self.pwd)

        except ldap.LDAPError as ex:
            self.error_handle(ex, self.ip)

    def search(self, base_dn, scope, s_filter, attr):
        """
        搜索域控
        搜索结果为迭代器
        """
        if not self.support_paged:
            try:
                result = self.conn.search_s(base_dn, scope, s_filter, attr)
                return result
            except Exception as e:
                err_msg = "search ou error: base_dn=%s, scope=%s, s_filter=%s, attr=%s, ex=%s" % (
                    base_dn, scope, s_filter, attr, str(e))
                raise_exception(err_msg, exp_detail=eval(str(e))['desc'])
        try:
            COOKIE = ''
            PAGE_SIZE = 1000
            CRITICALITY = True

            result = []
            first_pass = True
            pg_ctrl = SimplePagedResultsControl(CRITICALITY, PAGE_SIZE, COOKIE)
            while first_pass or pg_ctrl.cookie:
                first_pass = False
                msgid = self.conn.search_ext(
                    base_dn, scope, s_filter, attr, serverctrls=[pg_ctrl])
                result_type, data, msgid, serverctrls = self.conn.result3(
                    msgid)
                pg_ctrl.cookie = serverctrls[0].cookie
                result += data

            return result
        except Exception as e:
            err_msg = "search ou error: base_dn=%s, scope=%s, s_filter=%s, attr=%s, ex=%s" % (
                base_dn, scope, s_filter, attr, str(e))
            raise_exception(err_msg, exp_detail=eval(str(e))['desc'])

    def search_subtree(self, base_dn, s_filter, attr=['*', '+']):
        """
        使用SCOPE_SUBTREE搜索
        """
        return self.search(base_dn, ldap.SCOPE_SUBTREE, s_filter, attr)

    def __search_ou_user(self, dn, serach_scope, serach_filter):
        """
        搜索组织用户信息
        """
        results = self.search(dn, serach_scope, serach_filter,
                              ['*', '+'])
        return results

    def __base_ou(self, dn, search_filter):
        return self.search(dn, ldap.SCOPE_BASE, search_filter, ['*', '+'])

    def base_ou(self, dn, search_filter):
        try:
            return self.search(dn, ldap.SCOPE_BASE, search_filter, ['*', '+'])
        except ncTException as search_error:
            return

    def check_base_ou(self, dn):
        return self.search(dn, ldap.SCOPE_BASE, "objectClass=*", ['*', '+'])

    def __one_level_ou(self, dn, serach_filter):
        return self.search(dn, ldap.SCOPE_ONELEVEL, serach_filter, ['*', '+'])

    def convert_domain_ous(self, results, key_config):
        """
        转换域组织信息为标准组织结构
        """
        all_ou_infos = []
        for ret in results:
            ou_info = self.__ou2obj(ret, key_config)
            if ou_info:
                all_ou_infos.append(ou_info)

        return all_ou_infos

    def __get_attr_val(self, keys, attrs, strip=True):
        for key in keys:
            if key in attrs:
                if not attrs[key][0]:
                    val = ""
                elif strip:
                    val = attrs[key][0].strip()
                else:
                    val = attrs[key][0]
                return key, val
        return "", ""

    def __ou_filter(self, ou_list, keys):
        name_filters = ["Computers", "Builtin", "Domain Controllers"]

        result = []
        for ou_info in ou_list:
            _, val = self.__get_attr_val(keys, ou_info[1])
            if isinstance(val, bytes):
                val = bytes.decode(val)
            if val in name_filters:
                continue
            result.append(ou_info)
        return result

    def __convert_third_id(self, key, third_id):
        if not third_id:
            return ""

        if key == 'objectGUID':
            return str(uuid.UUID(bytes_le=third_id))
        else:
            return third_id

    def __ou2obj(self, ou_tuple, key_config):
        """
        OU tuple convert to DomainOuInfo
        Args:
            ou_tuple (tuple): come from ldap search result
            key_config (ncTUsrmDomainKeyConfig)
        """
        dn, attr = ou_tuple
        _, val = self.__get_attr_val(key_config.departNameKeys, attr)

        # 没有ou名称，过滤
        if not val:
            return

        ou_info = DomainOuInfo()
        ou_info.dn = dn
        ou_info.ou_name = bytes.decode(val)

        key, val = self.__get_attr_val(key_config.departThirdIdKeys, attr,
                                       strip=False)

        # 没有ou thirdid，过滤
        if not val:
            return
        ou_info.third_id = self.__convert_third_id(key, val)
        ou_info.server_type = self.domain_type

        return ou_info

    def __ou_list2obj_list(self, ou_list, key_config):
        """
        OU list convert to DomainOuInfo list
        If ou cannot convert, then ignore
        Args:
            ou_list (list): come from ldap search result
            key_config (ncTUsrmDomainKeyConfig)
        """
        for ou_tuple in ou_list:
            obj = self.__ou2obj(ou_tuple, key_config)
            if not obj:
                continue

            yield obj

    def __convert_login_name(self, key, login_name):
        if not login_name:
            return ""

        if key == "userPrincipalName":
            ret_name = ""
            split_names = login_name.rsplit('@', 1)
            if len(split_names) > 1:
                ret_name = login_name
            else:
                ret_name = login_name + '@' + self.domain_name
            return ret_name
        else:
            return "%s@%s" % (login_name, self.domain_name)

    def __user2obj(self, ou_tuple, key_config):
        """
        OU tuple convert to DomainUserInfo
        Args:
            ou_tuple (tuple): come from ldap search result
            key_config (ncTUsrmDomainKeyConfig)
        """

        dn, attr = ou_tuple
        # 获取登录名
        key, val = self.__get_attr_val(key_config.loginNameKeys, attr)
        # 无登录名后续不做
        if not val:
            return

        user_info = DomainUserInfo()
        user_info.login_name = self.__convert_login_name(
            key, bytes.decode(val))
        user_info.dn = dn
        if user_info.dn:
            user_info.ou_dn = ",".join(user_info.dn.split(",")[1:]).strip()

        # 获取显示名
        _, val = self.__get_attr_val(key_config.displayNameKeys, attr)
        if val:
            user_info.display_name = bytes.decode(val)
        else:
            user_info.display_name = user_info.login_name

        # 获取邮箱
        _, val = self.__get_attr_val(key_config.emailKeys, attr)
        if isinstance(val, str):
            user_info.email = val
        else:
            user_info.email = bytes.decode(val)

        # 获取身份证号
        if key_config.idcardNumberKeys:
            _, val = self.__get_attr_val(key_config.idcardNumberKeys, attr)
            if isinstance(val, str):
                user_info.idcard_number = val
            else:
                user_info.idcard_number = bytes.decode(val)

        # 获取手机号
        if key_config.telNumberKeys:
            _, val = self.__get_attr_val(key_config.telNumberKeys, attr)
            if isinstance(val, str):
                user_info.tel_number = val
            else:
                user_info.tel_number = bytes.decode(val)

        # 获取第三方id
        key, val = self.__get_attr_val(key_config.userThirdIdKeys, attr,
                                       strip=False)
        user_info.third_id = self.__convert_third_id(key, val)

        user_info.server_type = self.domain_type

        # 获取AD域用户的状态
        if self.domain_type == WINDOWS_AD:
            key, val = self.__get_attr_val(key_config.statusKeys, attr)
            if val:
                if (int(val) & 0x02):
                    user_info.status = False
                else:
                    user_info.status = True
            else:
                # AD域中获取不到用户禁用字段时不更改AS中的用户状态
                user_info.status = None
        else:
            user_info.status = True

        return user_info

    def dn2user(self, dn, key_config):
        """
        Get OU tuple from DN first, then convert to DomainUserInfo
        """
        result = self.__base_ou(dn, "objectClass=*")
        if not result:
            return
        # 如果objectClass中存在group字段，则为安全组，不是用户
        if "group" in [x.decode() for x in result[0][1]["objectClass"]]:
            return
        return self.__user2obj(result[0], key_config)

    def convert_domain_users(self, results, key_config):
        """
        转换域用户信息为标准用户结构
        """
        all_user_infos = []
        for ret in results:
            user_info = DomainUserInfo()
            user_attributes = ret[1]

            # 获取登录名
            for key in key_config.loginNameKeys:
                if key in user_attributes:
                    if key == "userPrincipalName":
                        tmp_name = bytes.decode(
                            user_attributes['userPrincipalName'][0])
                        split_name = tmp_name.rsplit('@', 1)
                        if len(split_name) > 1:
                            user_info.login_name = tmp_name
                        else:
                            user_info.login_name = tmp_name + '@' + self.domain_name
                    else:
                        user_info.login_name = bytes.decode(
                            user_attributes[key][0]) + '@' + self.domain_name
                    break

            # 获取显示名
            for key in key_config.displayNameKeys:
                if key in user_attributes:
                    user_info.display_name = bytes.decode(
                        user_attributes[key][0])
                    break

            # 获取邮箱
            for key in key_config.emailKeys:
                if key in user_attributes:
                    user_info.email = bytes.decode(user_attributes[key][0])

                    break

            # 获取身份证号
            if key_config.idcardNumberKeys:
                for key in key_config.idcardNumberKeys:
                    if key in user_attributes:
                        user_info.idcard_number = bytes.decode(
                            user_attributes[key][0])
                        break
            else:
                user_info.idcard_number = ''

             # 获取手机号
            if key_config.telNumberKeys:
                for key in key_config.telNumberKeys:
                    if key in user_attributes:
                        user_info.tel_number = bytes.decode(
                            user_attributes[key][0])
                        break
            else:
                user_info.tel_number = ''

            # 获取第三方id
            for key in key_config.userThirdIdKeys:
                if key in user_attributes:
                    third_id = user_attributes[key][0]
                    if key == 'objectGUID':
                        third_id = str(uuid.UUID(bytes_le=third_id))
                    user_info.third_id = third_id
                    break

            # 获取启用禁用状态
            if self.domain_type == WINDOWS_AD:
                for key in key_config.statusKeys:
                    if key in user_attributes:
                        val = user_attributes[key][0]
                        if (int(val) & 0x02):
                            user_info.status = False
                        else:
                            user_info.status = True
                    else:
                        # AD域中获取不到用户禁用字段时不更改AS中的用户状态
                        user_info.status = None
            else:
                user_info.status = True

            if not user_info.email:
                user_info.email = ''

            if not user_info.idcard_number:
                user_info.idcard_number = ''

            if not user_info.tel_number:
                user_info.tel_number = ''

            if not user_info.display_name:
                user_info.display_name = user_info.login_name

            user_info.server_type = self.domain_type

            user_info.dn = ret[0]
            if user_info.dn:
                user_info.ou_dn = ",".join(user_info.dn.split(",")[1:]).strip()

            if user_info.login_name:
                all_user_infos.append(user_info)

        return all_user_infos

    def is_domain_group(self, ou_attrs, key_config):
        """
        获取安全的用户成员DN
        """
        if ou_attrs:
            # 兼容IBM Domino域， IBM Domino的类型字段为：objectclass,
            class_key = "objectClass" if "objectClass" in ou_attrs else "objectclass"
            if len(set(key_config.groupKeys) & set(ou_attrs[class_key])):
                return True
        return False

    def get_domain_ou(self, dn, key_config):
        """
        获取域组织信息
        """
        results = self.__search_ou_user(dn, ldap.SCOPE_BASE,
                                        "objectClass=*")
        all_ou_infos = self.convert_domain_ous(results, key_config)
        if len(all_ou_infos) > 0:
            return all_ou_infos[0]

    def get_domain_user(self, dn, key_config):
        """
        获取域用户信息
        """
        results = self.__search_ou_user(dn, ldap.SCOPE_BASE,
                                        "objectClass=*")
        all_ou_infos = self.convert_domain_users(results, key_config)
        if len(all_ou_infos) > 0:
            return all_ou_infos[0]

    def get_onelevel_sub_ous(self, dn, key_config):
        """
        获取域组织下一级子组织信息
        """
        # results = self.__search_ou_user(dn, ldap.SCOPE_ONELEVEL,
        #                                 key_config.subOuFilter)
        # all_ou_infos = self.convert_domain_ous(results, key_config)
        # return all_ou_infos
        return list(self.get_sub_ous_iter(dn, key_config))

    def get_search_info_ous(self, info_result, key_config):
        """
        获取搜索的组织信息
        """
        return list(self.get_search_ous_iter(info_result, key_config))

    def get_onelevel_sub_users(self, dn, key_config):
        """
        获取域组织下一级子用户信息
        """
        # all_user_infos = []

        # ou_results = self.__search_ou_user(dn, ldap.SCOPE_BASE,
        #                                    "objectClass=*")

        # if ou_results and self.is_domain_group(ou_results[0][1], key_config):
        #     ou_dn = ou_results[0][0]
        #     ou_attrs = ou_results[0][1]
        #     # 获取安全组用户的DN
        #     if 'member' in ou_attrs:
        #         member_dns = ou_attrs['member']

        #         # 获取安全组用户
        #         for user_dn in member_dns:
        #             user_info = self.get_domain_user(user_dn, key_config)
        #             if user_info:
        #                 user_info.ou_dn = ou_dn
        #                 all_user_infos.append(user_info)
        # else:
        #     results = self.__search_ou_user(dn, ldap.SCOPE_ONELEVEL,
        #                                     key_config.subUserFilter)
        #     all_user_infos = self.convert_domain_users(results, key_config)

        # return all_user_infos

        return list(self.get_sub_users_iter(dn, key_config))

    def get_search_info_users(self, info_result, key_config):
        """
        获取搜索的用户信息
        """
        return list(self.get_search_users_iter(info_result, key_config))

    def get_all_sub_ous(self, dn, key_config):
        """
        获取所有子部门信息
        本接口会判断是否为根组织
        如果是根组织，将过滤不需要的组织
        """
        # results = self.__search_ou_user(dn, ldap.SCOPE_SUBTREE, key_config.subOuFilter)
        # all_ou_infos = self.convert_domain_ous(results, key_config)
        # return all_ou_infos

        return list(self.get_all_sub_ous_iter(dn, key_config))

    def get_all_sub_ous_iter(self, dn, key_config):
        """
        获取所有子部门
        本接口会判断是否为根组织
        如果是根组织，将过滤不需要的组织
        """
        # SCOPE_ONELEVEL 不会搜索自己，所以这里需要补充
        cur_ou = self.__base_ou(dn, "objectClass=*")
        if cur_ou:
            obj = self.__ou2obj(cur_ou[0], key_config)
            if obj:
                yield obj

        # is base dn filter
        if dn == self.domain_dn:
            one_level = self.__one_level_ou(dn, key_config.subOuFilter)
            one_level = self.__ou_filter(one_level, key_config.departNameKeys)
        else:
            one_level = self.__one_level_ou(dn, key_config.subOuFilter)

        for obj in self.__recursion_ous(one_level, key_config):
            yield obj

    def get_sub_ous_iter(self, dn, key_config):
        """
        获取组织子部门信息
        """
        ou_list = self.__one_level_ou(dn, key_config.subOuFilter)
        if dn == self.domain_dn:
            ou_list = self.__ou_filter(ou_list, key_config.departNameKeys)

        for obj in self.__ou_list2obj_list(ou_list, key_config):
            yield obj

    def get_search_ous_iter(self, info_result, key_config):
        """
        获取搜索的组织子部门信息
        """
        for obj in self.__ou_list2obj_list(info_result, key_config):
            yield obj

    def __recursion_ous(self, ou_list, key_config):
        """
        递归所有组织结构
        """
        for obj in self.__ou_list2obj_list(ou_list, key_config):
            yield obj

            sub_ou_list = self.__one_level_ou(obj.dn, key_config.subOuFilter)
            for obj in self.__recursion_ous(sub_ou_list, key_config):
                yield obj

    def get_all_sub_users(self, dn, key_config):
        """
        获取域组织下所有子用户信息
        """
        all_sub_ou_dns = []
        sub_ou_infos = self.get_all_sub_ous(dn, key_config)
        tmp_ou_dns = [ou_info.dn for ou_info in sub_ou_infos]
        all_sub_ou_dns += tmp_ou_dns

        all_sub_usrs = []
        for ou_dn in all_sub_ou_dns:
            sub_users = self.get_onelevel_sub_users(ou_dn, key_config)
            if sub_users:
                all_sub_usrs += sub_users
        return all_sub_usrs

    def get_all_sub_users_iter(self, dn, key_config):
        """
        获取域组织下所有子用户信息
        """
        for ou_info in self.get_all_sub_ous_iter(dn, key_config):
            for user_info in self.get_sub_users_iter(ou_info.dn, key_config):
                yield user_info

    def __get_group(self, dn, key_config):
        """
        If this ou is group, then return
        """
        result = self.__base_ou(dn, "objectClass=*")
        if not result:
            return
        attr = result[0][1]

        # 兼容IBM Domino域， IBM Domino的类型字段为：objectclass
        class_key = "objectClass" if "objectClass" in attr else "objectclass"
        if len(set(key_config.groupKeys) & set([x.decode() for x in attr[class_key]])) > 0:
            return result[0]

    def __get_group_member(self, dn, key_config):
        group = self.__get_group(dn, key_config)
        if not group:
            return []

        users = []
        for dn in group[1].get("member", []):
            user = self.dn2user(dn.decode(), key_config)
            if not user:
                continue

            users.append(user)
        return users

    def get_sub_users_iter(self, dn, key_config):
        """
        获取指定域组织下的用户信息
        该接口只获取直接关系的用户，不会获取所有子用户
        提高性能
        """
        for user in self.__get_group_member(dn, key_config):
            yield user

        for user in self.__one_level_ou(dn, key_config.subUserFilter):
            user_info = self.__user2obj(user, key_config)
            if user_info:
                yield user_info

    def get_search_users_iter(self, info_result, key_config):
        """
        获取搜索出的用户信息
        """
        for user in info_result:
            user_info = self.__user2obj(user, key_config)
            if user_info:
                yield user_info
