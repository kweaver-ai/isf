#!/usr/bin/env python
import configparser
import sys
import os
from ShareMgnt import ncTShareMgnt
from eisoo import thriftlib


def main():
    # sharemgnt 服务端口
    config_file = '/sysvol/conf/service_conf/app_default.conf'
    if os.path.isfile(config_file):
        config = configparser.ConfigParser()
        config.read(config_file)
        host = "localhost"
        port = config.getint("ShareMgnt", "port")
    else:
        host = str(sys.argv[1])
        port = int(sys.argv[2])
    client = thriftlib.BaseThriftClient(
        dest_ip=host,
        dest_port=port,
        interface=ncTShareMgnt,
        timeout_s=30
    )
    try:
        client.GetRetainOutLinkStatus()
    finally:
        client.close()


if __name__ == '__main__':
    main()
