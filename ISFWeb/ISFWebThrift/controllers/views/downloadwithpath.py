# coding:utf-8
from django.http import HttpResponse
from interface.sharemgnt import ShareMgntClient
import os
from urllib.parse import quote
import re



def download_file(request):

    taskId = request.GET.get('taskId')
    sharemgnt_inst = ShareMgntClient()
    data = sharemgnt_inst.call_interface('GetSpaceReportFileInfo', taskId)
    filename = quote(data.reportName)

    response = HttpResponse(
        data.reportData, content_type="application/octet-stream")
    response['Content-Disposition'] = '''attachment; filename="%s"; filename*=utf-8''%s''' % (
        filename, filename)
    return response
