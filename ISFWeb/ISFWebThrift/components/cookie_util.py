# coding:utf-8

def convert_cookies_to_dict(cookies_str):
    '''
    将cookie字符串转换为字典对象
    '''

    return dict([l.split('=', 1) for l in cookies_str.split('; ')])