from django.http import HttpResponse
from django.utils.deprecation import MiddlewareMixin

from thrift.Thrift import TException
from EThriftException.ttypes import ncTException
from components.json_wrapper import json_encode

class ExceptMiddleware(MiddlewareMixin):
    def process_exception(self, request, ex):
        if isinstance(ex, ncTException):
            exp = vars(ex)
            exp.update({'errMsg': ex.expMsg})
            exp.pop('expMsg')     
            result = {
                'error': exp
            }
        elif isinstance(ex, TException):
            result = {
                'error': {
                    'errMsg': ex.message,
                }
            }
        else:
            result = {
                'error': {
                    'errMsg': ex.__str__(),
                }
            }

        return HttpResponse(json_encode(result), content_type='application/json', status=501)
