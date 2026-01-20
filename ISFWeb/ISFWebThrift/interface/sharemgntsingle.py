from interface.baseclient import BaseClient
from ShareMgnt import ncTShareMgnt
from ShareMgnt.constants import NCT_SHAREMGNT_PORT
from components.getconfigfromfile import getconfigfromfile

class ShareMgntSingleClient(BaseClient):
    def __init__ (self):
        '''
        初始化Thrift连接
        '''
        info = getconfigfromfile()
        self.host = info['sharemgnt-single']['host']
        self.port = NCT_SHAREMGNT_PORT
        self.client = ncTShareMgnt
        super(ShareMgntSingleClient, self).__init__()
