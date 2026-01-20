#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""This is login access control manage class"""
from netaddr import IPNetwork, IPAddress
from src.common.db.connector import DBConnector

class LoginAccessControlManage(DBConnector):

    """
    Login access control manage
    """

    def __init__(self):
        """
        init
        """
        pass

    def check_login_ip_in_net_config(self, login_ip, ip, sub_net_mask):
        """
        检查登录ip是否在指定网段配置访问内
        """
        # 获取子网掩码长度
        o = list(map(int, sub_net_mask.split('.')))
        res = (16777216 * o[0]) + (65536 * o[1]) + (256 * o[2]) + o[3]
        net_mask_length = bin(res).count('1')
        net = ip + '/' + str(net_mask_length)
        return True if IPAddress(login_ip) in IPNetwork(net) else False
