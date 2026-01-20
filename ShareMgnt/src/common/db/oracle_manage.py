# coding:utf8
from src.common.lib import exec_command
from src.common.db.connector import DBOperate



class OracleManage(object):
    """
    Oracle数据库管理类
    """

    def __init__(self, host, port, user, pwd, sid):
        """
        init
        """
        self.host = host
        self.port = port
        self.user = user
        self.pwd = pwd
        self.sid = sid
        self.conn = None

    def init_conector(self):
        """
        """
        try:
            if not self.conn:
                import cx_Oracle
                dsn = cx_Oracle.makedsn(self.host, self.port, self.sid)
                self.conn = cx_Oracle.connect(self.user, self.pwd, dsn)
        except Exception as ex:
            raise ex

    def get_conn(self):
        """
        获取一个连接
        """
        if not self.conn:
            self.init_conector()

        return self.conn

    def get_cursor(self):
        """
        获取一个游标
        """
        return self.get_conn().cursor()

    def get_db_operator(self):
        """
        """
        return DBOperate(self.get_conn())

    @property
    def operator(self):
        """
        """
        return self.get_db_operator()
