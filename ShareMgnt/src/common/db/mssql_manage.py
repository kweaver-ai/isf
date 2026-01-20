# coding: utf-8
import imp
from src.common.db.connector import DBOperate
from src.common.sharemgnt_logger import ShareMgnt_Log

ODBC_NAME = "pymssql"

class MsSQLManage(object):
    """
    Windows SqlServer管理类
    """
    def __init__(self, host, port, user, password, database, charset="utf8", as_dict=True):
        """
        """
        self.host = host
        self.user = user
        self.password = password
        self.port = port
        self.database = database
        self.charset = charset
        self.as_dict = as_dict
        self.conn = None
        self.cursor = None

    def __del__(self):
        """
        """
        try:
            if self.cursor:
                self.cursor.close()
            if self.conn:
                self.conn.close()
        except Exception:
            pass

    def init_conector(self):
        """
        """
        try:
            mod_info = imp.find_module(ODBC_NAME)
        except ImportError:
            ShareMgnt_Log(_("MSSQL_ODBC_NOT_FOUND"))
            return

        pymssql = imp.load_module(ODBC_NAME, *mod_info)
        try:
            if not self.conn:
                self.conn = pymssql.connect(host=self.host,
                                            port=self.port,
                                            user=self.user,
                                            password=self.password,
                                            database=self.database,
                                            charset=self.charset,
                                            as_dict=True)
        except Exception as ex:
            raise ex

    def get_conn(self):
        """
        """
        if not self.conn:
            self.init_conector()

        return self.conn

    def get_cursor(self):
        """
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

