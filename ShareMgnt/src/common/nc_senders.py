#!/usr/bin/python3
# -*- coding:utf-8 -*-
"""一堆发送函数"""
import smtplib
import socket
from email.header import Header
from email.mime.image import MIMEImage
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText

from ShareMgnt.ttypes import ncTShareMgntError
from src.common.encrypt.simple import eisoo_rsa_decrypt
from src.common.jsonconv_Ttype import SmtpConfDec, SmtpConfEnc
from src.common.lib import check_email, check_smtp_params, raise_exception
from src.modules.smtp_manage import JsonConfManage


def email_send(conf, toList, title=None, text=None, img=None):
    """
    邮件发送函数
    """
    if not toList:
        return
    if title is None:
        title = _("test mail")
    if text is None:
        text = _("sent by server")

    check_smtp_params(conf)

    # 验证toList是电子邮件格式正确性
    for e in toList:
        if not check_email(e):
            raise_exception(exp_msg=_("email illegal"),
                            exp_num=ncTShareMgntError.
                            NCT_INVALID_EMAIL)

    # 准备邮件
    if img is None:
        msg = MIMEText(text, 'html', _charset="utf-8")
    else:
        msg = MIMEMultipart("related")
        msg_alternative = MIMEMultipart("alternative")
        msg.attach(msg_alternative)
        msg_alternative.attach(MIMEText(text, 'html', _charset="utf-8"))
        msg_img = MIMEImage(img)
        # content-Id are referenced in content
        msg_img.add_header("Content-ID", "<image>")
        msg.attach(msg_img)

    msg["from"] = conf.email
    msg["to"] = ",".join(toList)
    msg['subject'] = Header(title, 'utf-8')

    # 发送流程
    try:
        srv = None
        # 三种协议
        if conf.safeMode == 1:
            srv = smtplib.SMTP_SSL(conf.server, conf.port, timeout=30)
        else:
            srv = smtplib.SMTP(conf.server, conf.port, timeout=30)
        if conf.safeMode == 2:
            srv.starttls()
    except socket.error:
        raise_exception(exp_msg=_("NCT_SMTP_SERVER_NOT_AVAILABLE"),
                        exp_num=ncTShareMgntError.NCT_SMTP_SERVER_NOT_AVAILABLE)
    except smtplib.SMTPException:
        raise_exception(exp_msg=_("NCT_SMTP_SERVER_NOT_AVAILABLE"),
                        exp_num=ncTShareMgntError.NCT_SMTP_SERVER_NOT_AVAILABLE)


    if not conf.openRelay:      # 新增判断是否开启openRelay
        try:
            # 登录
            raw_pwd = bytes.decode(eisoo_rsa_decrypt(conf.password))
            srv.login(conf.email, raw_pwd)
        except smtplib.SMTPAuthenticationError:
            username = conf.email.split('@')[0]
            try:
                srv.login(username, raw_pwd)
            except smtplib.SMTPAuthenticationError:
                raise_exception(exp_msg=_("NCT_SMTP_LOGIN_FAILED"),
                                exp_num=ncTShareMgntError.NCT_SMTP_LOGIN_FAILED)
        except smtplib.SMTPException:
            raise_exception(exp_msg=_("No suitable authentication method was found."),
                            exp_num=ncTShareMgntError.NCT_SMTP_AUTHENTICATION_METHOD_NOT_FOUND)

    try:
        # 发送
        srv.sendmail(conf.email, toList, msg.as_string())
    except smtplib.SMTPException:
        raise_exception(exp_msg=_("NCT_SMTP_SEND_FAILED"),
                        exp_num=ncTShareMgntError.NCT_SMTP_SEND_FAILED)


def email_send_html_content(toEmailList, subject=None, content=None, img=None):
    """
    发送邮件，html格式
    """
    # 先验证mailto是电子邮件格式正确性
    if not toEmailList:
        raise_exception(exp_msg=_("NCT_SMTP_RECIPIENT_MAIL_ILLEGAL"),
                        exp_num=ncTShareMgntError.NCT_SMTP_RECIPIENT_MAIL_ILLEGAL)

    # 判断邮箱服务器是否设置
    conf = JsonConfManage("smtp_config", SmtpConfEnc, SmtpConfDec).get_config()
    if conf is None:
        raise_exception(exp_msg=_("NCT_SMTP_SERVER_NOT_SET"),
                        exp_num=ncTShareMgntError.NCT_SMTP_SERVER_NOT_SET)

    # 发送邮件
    email_send(conf, toEmailList, subject, content, img)
