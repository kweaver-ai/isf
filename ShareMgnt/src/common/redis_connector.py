import os
import time
import redis
from redis.sentinel import Sentinel
from redis.cluster import RedisCluster, ClusterNode
from src.common.sharemgnt_logger import ShareMgnt_Log


class OPRedis():
    def __init__(self):
        self.coon = OPRedis.get_redis_conn()

    @staticmethod
    def read_redis_values():
        """获取redis相关环境变量"""
        redis_configs = {
            "redis_host": os.getenv("REDIS_HOST", "proton-redis-proton-redis-sentinel.resource.svc.cluster.local"),
            "redis_port": os.getenv("REDIS_PORT", "26379"),
            "redis_user": os.getenv("REDIS_USER", ""),
            "redis_password": os.getenv("REDIS_PASSWORD", ""),
            "redis_cluster_mode": os.getenv("REDIS_CONNECTTYPE", "sentinel"),
            "redis_master_name": os.getenv("REDIS_MASTER_NAME", "mymaster"),
            "redis_sentinel_user": os.getenv("REDIS_SENTINEL_USER", ""),
            "redis_sentinel_password": os.getenv("REDIS_SENTINEL_PASSWORD", "")
        }

        return redis_configs

    @staticmethod
    def get_redis_conn():
        redis_configs = OPRedis.read_redis_values()
        # 哨兵部署
        try:
            if redis_configs["redis_cluster_mode"] == "sentinel":
                sentinel = Sentinel([(redis_configs["redis_host"], redis_configs["redis_port"])],
                                    password=redis_configs["redis_sentinel_password"],
                                    sentinel_kwargs={"password": redis_configs["redis_sentinel_password"],
                                                     "username": redis_configs["redis_sentinel_user"]})
                redis_conn = sentinel.master_for(redis_configs["redis_master_name"],
                                                 username=redis_configs["redis_user"],
                                                 password=redis_configs["redis_password"])
            
                ShareMgnt_Log("Connect sentinel redis success, the redis_host is %s, the redis_port is %s", redis_configs["redis_host"], redis_configs["redis_port"])
                return redis_conn
        except Exception as e:
            ShareMgnt_Log("Connect sentinel redis failed, the reason is %s", str(e))

        # 集群部署
        if redis_configs["redis_cluster_mode"] == "cluster":
            startup_nodes = []
            try:
                hosts = [host.strip() for host in redis_configs["redis_host"].split(",")]
                for host in hosts:
                        if ":" in host:
                            startup_nodes.append(ClusterNode(host.split(":")[0], host.split(":")[1]))
                        else:
                            startup_nodes.append(ClusterNode(host.split(":")[0], redis_configs["redis_port"]))
                            
                redis_conn = RedisCluster(startup_nodes=startup_nodes, \
                                    username=redis_configs["redis_user"], \
                                    password=redis_configs["redis_password"])
                
                ShareMgnt_Log("Connect cluster redis success, the redis_host is %s, the redis_port is %s", redis_configs["redis_host"], redis_configs["redis_port"])
                if not redis_conn.ping():
                    ShareMgnt_Log("Redis cluster connection failed, ping failed")
                    raise Exception("Redis cluster connection failed, ping failed")
                
                return redis_conn
            except Exception as e:
                ShareMgnt_Log("Connect cluster redis failed, the reason is %s", str(e))

        # 公有云的主从模式和单机模式
        if redis_configs["redis_cluster_mode"] == "master-slave" or \
                redis_configs["redis_cluster_mode"] == "standalone":
            try:
                OPRedis.pool = redis.ConnectionPool(
                    host=redis_configs["redis_host"],
                    port=redis_configs["redis_port"],
                    username=redis_configs["redis_user"],
                    password=redis_configs["redis_password"],
                    decode_responses=True)
                redis_conn = redis.Redis(connection_pool=OPRedis.pool)

                ShareMgnt_Log("Connect master-slave or standalone redis success, the redis_host is %s, the redis_port is %s", redis_configs["redis_host"], redis_configs["redis_port"])
                return redis_conn
            except Exception as e:
                ShareMgnt_Log("Connect master-slave or standalone redis failed, the reason is %s", str(e))

    def set_redis(self, key, value, timeout=None):
        """设置redis中(key: value)"""
        res = self.coon.set(key, value, ex=timeout, nx=True)
        return res

    def get_redis(self, key):
        """获取redis中(key: value)"""
        res = self.coon.get(key)
        return res

    def del_redis(self, key):
        """删除redis中(key: value)"""
        res = self.coon.delete(key)
        return res

    def set_lock(self, key, value, timeout=None):
        """向redis中写入值, 相当于上锁, 用于第三方插件下载"""
        end = time.time() + timeout
        while time.time() < end:
            if self.set_redis(key, value, timeout):
                return True
            time.sleep(0.001)

        return False

    def get_lock(self, key, timeout=None):
        """获取redis中的值, 相当于获取锁"""
        end = time.time() + timeout
        while time.time() < end:
            value = self.get_redis(key)
            # 有值, 代表有锁
            if not value:
                return False
            time.sleep(0.001)

        return True
