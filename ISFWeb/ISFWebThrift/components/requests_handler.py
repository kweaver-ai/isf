import requests as req
from requests.adapters import HTTPAdapter, Retry

retry_times = 4
retry_backoff_factor = 0.1

session = req.Session()
retry = Retry(total=retry_times, backoff_factor=retry_backoff_factor)
adapter = HTTPAdapter(max_retries=retry)
session.mount('http://', adapter)
session.mount('https://', adapter)

class RequestHandler:
    def __init__(self):
        self.session = session
    def request(self, method=None, url=None, params=None, data=None, json=None, timeout=(0.5, 30), **kwargs):
        return self.session.request(method=method, url=url, params=params, data=data, json=json, timeout=timeout, **kwargs)

requests = RequestHandler()