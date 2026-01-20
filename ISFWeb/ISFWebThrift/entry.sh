#!/bin/bash
service nginx start

/usr/local/bin/uwsgi --ini /config/isfweb.ini;