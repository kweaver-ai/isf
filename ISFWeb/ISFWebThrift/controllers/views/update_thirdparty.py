# coding:utf-8

from django.http import HttpResponse
from components.json_wrapper import json_encode
from ShareMgnt.ttypes import ncTThirdPartyPluginInfo
from interface.sharemgntsingle import ShareMgntSingleClient

def update_thirdparty_package(request):
    '''
    上传第三方认证插件
    '''
    package = request.FILES.get('package')

    if package is None:
        raise Exception('no package uploaded')

    filename = request.POST.get('filename').strip()
    indexId = int(request.POST.get('indexId').encode('utf-8').strip())
    thirdPartyId = request.POST.get('thirdPartyId').strip()
    pluginType = int(request.POST.get('type').encode('utf-8').strip())

    info = ncTThirdPartyPluginInfo()
    info.filename = filename
    info.data = package.read()
    info.indexId = indexId
    info.thirdPartyId = thirdPartyId
    info.type = pluginType

    shargmgnt_single_inst = ShareMgntSingleClient()
    shargmgnt_single_inst.call_interface('AddGlobalThirdPartyPlugin', info)

    return HttpResponse(json_encode(None), content_type="application/json")
