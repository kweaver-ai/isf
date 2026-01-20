#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""
MQ封装类
1.提供RPC Client
2.提供正常 Client
"""
import pika, json

class RabbitMQSend(object):
    """
    AS向cmp发送心跳
    """
    def __init__(self, host=None, port=None, virtual_host=None, username=None, password=None, exchange_name='', routing_key='', exchange_type='direct'):
        """
        :param str host: 连接RabbitMQ服务的主机名或主机IP
        :param int port: RabbitMQ服务器的端口  默认端口为5672或5671(ssl)
        :param str virtual_host: RabbitMQ使用的虚拟主机
        :param username: 用户名
        :param password: 密码
        :param exchange_name: 消息交换名称
        :param routing_key: 消息路由KEY
        :param exchange_type: 消息交换类型
        """
        self._exchange_type = exchange_type
        self._exchange_name = exchange_name
        self._routing_key = routing_key
        if username is not None and password is not None:
            credentials = pika.credentials.PlainCredentials(username, password)
        # 创建连接connection到rabbitmq server
        self._connection = pika.BlockingConnection(pika.ConnectionParameters(host, port, virtual_host, credentials))
        # 创建虚拟连接channel
        self._channel = self._connection.channel()

    def send_message(self, message, is_persist=False):
        """
        发布消息
        @param message: 要发送的消息   type : dict or str unicode
        @param is_persist: 要发送的是否要持久化   默认是False(非持久化的)
        """
        msg_props = None
        send_massage = ""
        if isinstance(message, dict):
            send_massage = json.dumps(message)
        else:
            send_massage = message
        if is_persist:
            msg_props = pika.BasicProperties(delivery_mode=2,)
        self._channel.basic_publish(self._exchange_name, self._routing_key, send_massage, properties=msg_props)

    def close(self):
        try:
            if self._connection:
                self._connection.close()
        except:
            pass

    def __del__(self):
        self.close()