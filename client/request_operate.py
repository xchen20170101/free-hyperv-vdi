#coding = 'utf-8'

import requests

class RequestOperate(object):

    def __init__(self, server_ip):
        self.host = "http://{}:8090".format(server_ip)

    def get(self, url, headers=None, param=None):
        target_url = "{}{}".format(self.host, url)
        #response = requests.get(target_url, params=param)
        response = requests.request("GET", target_url, headers=headers, data={})
        return response.json()
    
    def post(self, url, headers=None, data=None):
        target_url = "{}{}".format(self.host, url)
        #response = requests.post(target_url, data=data)
        response = requests.request("POST", target_url, headers=headers, data=data)
        return response.json()

    def put(self, url, headers=None, data=None):
        target_url = "{}{}".format(self.host, url)
        #response = requests.put(target_url, data=data)
        response = requests.request("PUT", target_url, headers=headers, data=data)
        return response.json()
    
    def delete(self, url, headers=None, data=None):
        target_url = "{}{}".format(self.host, url)
        #response = requests.delete(target_url)
        response = requests.request("DELETE", target_url, headers=headers, data=data)
        return response.json()