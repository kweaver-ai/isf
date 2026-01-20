from interface.baseclient import BaseClient
from EACP import ncTEACP
from EACP.constants import NC_T_EACP_PORT
from components.getconfigfromfile import getconfigfromfile

class EACPClient(BaseClient):
    def __init__ (self):
        '''
        初始化Thrift连接
        '''
        info = getconfigfromfile()
        eacp = info['eacp']
        self.host = eacp['thriftHost'] 
        self.port = NC_T_EACP_PORT
        self.client = ncTEACP
        
        super(EACPClient, self).__init__()
