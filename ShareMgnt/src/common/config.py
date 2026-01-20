import configparser

from src.common.lib import raise_exception

class Config:
    try:
        config = configparser.ConfigParser()
        config.read('/sysvol/conf/service_conf/service_access.conf')
        ossgateway_config = {
            'protocol': config.get("ossgateway", "protocol"),
            'host': config.get("ossgateway", 'privateHost'),
            'port': config.get("ossgateway", "privatePort")
        }
        authentication_config = {
            'host': config.get("authentication", 'privateHost'),
            'port': config.get("authentication", "privatePort")
        }
    except Exception:
        raise_exception(
            'read /sysvol/conf/service_conf/service_access.conf failed')