# coding:utf-8

def get_ipaddress(request):
    '''
    获取请求中的ip，ip优先考虑HTTP_X_FORWARDED_FOR。如果没有HTTP_X_FORWARDED_FOR或者为空，则取REMOTE_ADDR
    '''

    ip = request.META.get('HTTP_X_FORWARDED_FOR', '')
    if ip == '':
        ip = request.META.get('REMOTE_ADDR', '')
    else:
        ip = ip.split(',')[0].strip()

    return ip