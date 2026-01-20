# coding:utf-8

from django.http import HttpResponse

from interface.sharemgnt import ShareMgntClient
from ShareMgnt.ttypes import *

from EACP.ttypes import *

from components.json_wrapper import json_encode

from urllib.parse import quote

def download_user(request):
    '''
    下载导出用户组织信息/下载导入失败记录
    '''
    
    taskId = request.GET.get('taskId')
    sharemgnt_inst = ShareMgntClient()
    if taskId != "0":
        data = sharemgnt_inst.call_interface('Usrm_DownloadBatchUsers', taskId)
    else:
        data = sharemgnt_inst.call_interface('Usrm_DownloadImportFailedUsers')
    filename = quote(data.reportName)

    response = HttpResponse(
        data.reportData, content_type="application/octet-stream")
    response['Content-Disposition'] = '''attachment; filename="%s"; filename*=utf-8''%s''' % (
        filename, filename)
    return response

def import_user(request):
    ''' 
    导入用户组织信息 
    '''

    filePath = request.FILES.get('file').temporary_file_path()
    fobj = open(filePath, 'rb')
    content = fobj.read()
    fobj.close()
    fileName = request.POST.get('fileName')
    isTrue = request.POST.get('userCover')
    userCover = isTrue == str('true')
    responsiblePersonId = request.POST.get('responsiblePersonId')
    sharemgnt_client = ShareMgntClient()
    binaryData = ncTBatchUsersFile()
    binaryData.fileName = fileName
    binaryData.data = content
    result = sharemgnt_client.call_interface('Usrm_ImportBatchUsers', binaryData, userCover, responsiblePersonId)
    response = json_encode(result)
    return HttpResponse(response, content_type="application/json")


