import os
import json
import requests


requests.packages.urllib3.disable_warnings()


class Client(object):

    def __init__(self, addr, token, **kws):
        self.addr = addr
        self.token = token
        self.username = token.split(":")[0]
        if token.find(':') >= 0:
            self.password = token.split(':')[1]
        else:
            self.password = ''
        self.debug = kws.get('debug', False)
        self.default()

    def default(self):
        if self.addr.find(':') < 0:
            self.addr = '{}:{}'.format(self.addr, 10000) 
  
    def request(self, url, method, data=""):
        url = "{}/{}".format(self.addr, url)
        if self.debug:
            print("Client.request {} {} {}".format(method, url, data))

        resp = requests.request(method, url,
                                json=data, verify=False,
                                auth=(self.username, self.password))
        if self.debug:
            print("Client.request RESPONSE {}".format(resp.text))
        if not resp.ok:
            resp.raise_for_status()
        return resp


# noinspection PyInterpreter
if __name__ == '__main__':
    addr = os.environ.get("LS_SERVER", "https://localhost:10080")
    token = os.environ.get("LS_AUTH", "")

    c = Client(addr, token)
    resp = c.request("user", "GET")
    print(json.dumps(resp.json(), indent=4))

    resp = c.request("link", "GET")
    print(json.dumps(resp.json(), indent=4))

    resp = c.request("neighbor", "GET")
    print(json.dumps(resp.json(), indent=4))
