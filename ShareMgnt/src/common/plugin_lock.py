from threading import Lock

class PluginVersion():
    def __init__(self):
        if not hasattr(PluginVersion, 'AUTH_LOCAL_VERSION'):
            # 第三方认证插件在本地的版本号(用插件的object_id标识)
            PluginVersion.AUTH_LOCAL_VERSION = ""
        if not hasattr(PluginVersion, 'MSG_LOCAL_VERSION'):
            # 第三方消息插件在本地的版本号(用插件的object_id标识)
            PluginVersion.MSG_LOCAL_VERSION = ""

        self.lock1 = Lock()
        self.lock2 = Lock()

    def update_auth_local_version(self, object_id):
        """更新第三方认证插件的本地版本"""
        with self.lock1:
            PluginVersion.AUTH_LOCAL_VERSION = object_id

    def update_msg_local_version(self, object_id):
        """更新第三方消息插件的本地版本"""
        with self.lock2:
            PluginVersion.MSG_LOCAL_VERSION = object_id            
