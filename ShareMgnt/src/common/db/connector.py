import MySQLdb
import rdsdriver
from contextlib import contextmanager
from MySQLdb.cursors import DictCursor
from DBUtils.PersistentDB import PersistentDB
from src.common import global_info


@contextmanager
def safe_cursor(conn):
    """
    从连接出获取游标，执行完sql后，自动关闭游标
    conn = ConnectorManager.get_db_conn()
    with safe_cursor(conn) as cursor:
        user_id = "id1"
        cursor.execute("select * from t_user where f_user_id = %s", (user_id,))
    """
    if isinstance(conn, DBOperate):
        conn = conn.conn
    cursor = None
    try:
        cursor = conn.cursor()
        yield cursor
        conn.commit()
    except Exception:
        conn.rollback()
        raise
    finally:
        if cursor:
            cursor.close()


class DBOperate(object):
    """
    DB操作类，需要传入一个数据库连接
    封装了常用的sql操作
    """
    def __init__(self, conn):
        """
        Init
        """
        self.conn = conn

    def __del__(self):
        """
        Del
        """
        try:
            self.conn.close()
        except Exception:
            pass

    def escape(self, value):
        """
        转义非法字符
        """
        return bytes.decode(MySQLdb.escape_string(value))

    @classmethod
    def __get_columns(cls, columns):
        """
        将列名列表转换为字符串
        Args:
            columns: list，列名
        Return:
            字符串，例如：
            __get_columns(["id", "name"])
            => (`id`, `name`)
        """
        stmt = ["`%s`" % c for c in columns]
        return " (%s) " % (", ".join(stmt))

    def insert(self, table, columns):
        """
        插入一条数据
        Args:
            table: string，要插入的表名
            columns: dict or list，要插入值的列以及值
                     如果是字典，键为列名
                     如果是列表，元素为值
        Return:
            最后一个插入语句的自增ID，如果没有则为0
        Raise:
            TypeError: 参数类型错误时丢出异常
        Example:
            insert("test_table", {"id": 1, "name": "test"})
            => INSERT INTO `test_table` (`id`, `test`) VALUES ("1", "test")
            insert("test_table", [1, "test"])
            => INSERT INTO `test_table` VALUES ("1", "test")
        """
        if not isinstance(table, str):
            raise TypeError("table only use string type")

        sql = ["INSERT INTO `{0}` ".format(table)]

        if isinstance(columns, dict):
            sql.append(self.__get_columns(list(columns.keys())))
            values = list(columns.values())
        elif isinstance(columns, list):
            values = columns
        else:
            raise TypeError("columns only use list or dict type")

        sql.append(" VALUES (%s) " % (", ".join(["%s"] * len(values))))
        cursor = self.conn.cursor()
        cursor.execute("".join(sql), values)
        self.conn.commit()
        cursor.close()

    def insert_many(self, table, columns, values):
        """
        插入多条数据
        Args:
            table: 字符串，要插入的表名
            columns: 列表，要插入值的列，不需要参数使用空列表
            values: 列表元组嵌套，要插入的值
        Return:
            插入行数
        Raise:
            TypeError: 参数类型错误时丢出异常
        Example:
            insert_many("test_table", ["id", "name"], [(1, "name1"), (2, "name2")])
            => INSERT INTO `test_table` (`id`, `name`) VALUES ("1", "name1"), ("2", "name2")
            insert_many("test_table", [], [(1, "name1"), (2, "name2")])
            => INSERT INTO `test_table` VALUES ("1", "name1"), ("2", "name2")
        """
        if not isinstance(table, str):
            raise TypeError("table only use string type")

        if not isinstance(columns, list) or not isinstance(values, list):
            raise TypeError("columns or values only use list type")

        if not isinstance(values[0], tuple):
            raise TypeError("values value only use tuple type")

        sql = ["INSERT INTO `{0}`".format(table)]
        if columns:
            sql.append(self.__get_columns(columns))

        sql.append(" VALUES (%s) " % (", ".join(["%s"] * len(values[0]))))
        cursor = self.conn.cursor()
        row_affected = cursor.executemany("".join(sql), values)
        self.conn.commit()
        cursor.close()
        return row_affected

    def query(self, sql, *args):
        '''
        执行没有返回的查询
        本操作返回受影响行数与最后插入行的自增ID
        '''
        cursor = self.conn.cursor()
        cursor.execute(sql, args)
        affect_row = cursor.rowcount
        self.conn.commit()
        cursor.close()
        return affect_row

    def all(self, sql, *args):
        '''
        执行一条查询语句，并返回所有结果
        '''
        cursor = self.conn.cursor()
        cursor.execute(sql, args)
        result = cursor.fetchall()
        self.conn.commit()
        cursor.close()
        return result

    def one(self, sql, *args):
        '''
        执行一条查询语句，并返回第一个结果
        '''
        cursor = self.conn.cursor()
        cursor.execute(sql, args)
        result = cursor.fetchone()
        self.conn.commit()
        cursor.close()
        return result


class ConnectorManager:
    """
    MYSQL连接管理器，内部采用PersistentDB，一个线程分配一个连接
    """
    pool = None

    def __init__(self):
        """
        init
        """
        pass

    @staticmethod
    def get_pool():
        """
        生成连接池实例
        Args:
            db_ip: string 数据库IP
        Return:
            PooledDB
        """
        ConnectorManager.pool = PersistentDB(
            creator=rdsdriver,
            maxusage=0,      # 不限制连接使用次数
            ping=2,          # 每次使用前ping
            host=global_info.DB_WRITE_IP,  # 始终连接写库
            port=global_info.DB_PORT,
            user=global_info.DB_USER,
            password=global_info.DB_PWD,
            database=global_info.DB_NAME,
            cursorclass=rdsdriver.DictCursor,
            autocommit=True  # 自动提交，避免事务问题
        )
        return ConnectorManager.pool

    @staticmethod
    def get_db_conn():
        """
        获取数据库连接
        """
        if ConnectorManager.pool is None:
            ConnectorManager.pool = ConnectorManager.get_pool()
        return ConnectorManager.pool.connection()

    @staticmethod
    def get_db_operator(mode=2):
        """
        从连接池获取一个连接
        Args:
            mode: int 1为写，2为读（但实际都使用写连接）
        Return:
            DBOperate
        """
        return DBOperate(ConnectorManager.get_db_conn())


class DBConnector(object):
    """
    调用连接管理器，获取一个dboperator
    应用模块继承该类，就可以使用 r_db, w_db 来进行sql操作
    - r_db      读数据库连接
    - w_db      写数据库连接
    """
    __w_db = None  # 只使用写连接

    @property
    def r_db(self):
        """
        读库getter，使用写连接以确保一致性
        """
        return self.w_db  # 读操作也使用写连接

    @property
    def w_db(self):
        """
        写库getter
        """
        return ConnectorManager.get_db_operator(1)  # 始终使用写连接
