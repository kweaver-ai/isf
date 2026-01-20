# coding:utf-8

from django.http import HttpResponse
from interface.sharemgnt import ShareMgntClient

from urllib.parse import quote

def download_activity(request):
    taskId = request.GET.get('taskId')
    sharemgnt_inst = ShareMgntClient()
    data = sharemgnt_inst.call_interface('GetActiveReportFileInfo', taskId)
    filename = quote(data.reportName)

    response = HttpResponse(
        data.reportData, content_type="application/octet-stream")
    response['Content-Disposition'] = '''attachment; filename="%s"; filename*=utf-8''%s''' % (
        filename, filename)
    return response
