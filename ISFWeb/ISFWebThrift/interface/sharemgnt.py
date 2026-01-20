from interface.baseclient import BaseClient
from ShareMgnt import ncTShareMgnt
from ShareMgnt.constants import NCT_SHAREMGNT_PORT
from components.getconfigfromfile import getconfigfromfile

class ShareMgntClient(BaseClient):
    def __init__ (self):
        '''
        初始化Thrift连接
        '''
        info = getconfigfromfile()
        sharemgnt = info['sharemgnt']
        self.host = sharemgnt['host']
        self.port = NCT_SHAREMGNT_PORT
        self.client = ncTShareMgnt
        super(ShareMgntClient, self).__init__()
