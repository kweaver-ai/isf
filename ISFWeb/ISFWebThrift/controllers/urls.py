from django.urls import path, re_path
from controllers.views.user import download_user, import_user
from controllers.views.device import batch_import
from controllers.views.onlinestatistics import download_activity
from controllers.views.update_thirdparty import update_thirdparty_package
from controllers.views.log import download_log
from controllers.views.downloadwithpath import download_file
from controllers.views.thrift import controller
import controllers.privateAPI.view as privateapi_view

urlpatterns = [
	path('user/downloaduser/', download_user),
	path('user/importuser/', import_user),
    path('device/batchimport/', batch_import),
    path('onlinestatistics/downloadActivity/', download_activity),
    path('update_thirdparty/update_thirdparty/', update_thirdparty_package),
    path('log/downloadLog/', download_log),
    path('downloadwithpath/downloadfile/', download_file),
    re_path(r'(ossgateway)/(\w+)', privateapi_view.controller),
    re_path(r'(audit-log)/(log/management)$', privateapi_view.controller),
    re_path(r'(\w+)/(\w+)$', controller),
]
