# coding:utf-8

from interface.eacp import EACPClient
from EACP.ttypes import *
from components.json_wrapper import  json_response

@json_response
def batch_import(request):
    ''' 批量导入用户设备 '''
    
    filePath = request.FILES.get('file').temporary_file_path()
    fobj = open(filePath, 'r')
    content = fobj.read()
    udids=content.splitlines()
    fobj.close()
    userId = request.POST.get('userId')
    osType = int(request.POST.get('osType'))
    eacp  = EACPClient()
    return eacp.call_interface('EACP_AddDevices', userId, udids, osType)