#coding = 'utf-8'

import configparser

"""
1. 保存服务端IP，用户名
2. 首次登录，需要配置服务端IP，若不配置，则无法登录
3. 首次登录成功后，保存用户名到配置文件
4. 打开登录页面，从配置文件中加载用户名和服务端IP
"""

class ConfigUtil(object):

    def __init__(self, file_path="conf.ini"):
        self.file_path = file_path
        self.config = configparser.ConfigParser()
        
    def read_ini(self, section, key):
        self.config.read(self.file_path)
        try:
            return self.config.get(section, key)
        except:
            return ""

    def write_ini(self, section, key, value):
        if not self.config.has_section(section):
            self.config.add_section(section)
        self.config.set(section, key, value)
        with open(self.file_path, 'w') as config_file:
            self.config.write(config_file)
