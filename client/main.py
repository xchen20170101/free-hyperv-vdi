#coding = 'utf-8'

import sys
import typing
from PyQt5 import QtCore, QtGui, QtWidgets
from PyQt5.QtWidgets import QApplication, QDialog, QWidget, QTableWidgetItem, QPushButton, QVBoxLayout, QMessageBox, QGridLayout, QLabel
import Ui_login
import Ui_main
import Ui_config
import Ui_password
import subprocess
from request_operate import RequestOperate
import requests
from config_util import ConfigUtil
from qt_material import apply_stylesheet

config = ConfigUtil()

def show_tips(constant, tip=0):
    msg_box = QMessageBox()
    if tip == 1:
        msg_box.setIcon(QMessageBox.Warning)
        msg_box.setWindowTitle("警告")
    else:
        msg_box.setIcon(QMessageBox.Information)
        msg_box.setWindowTitle("信息")
    msg_box.setText(constant)
    
    msg_box.exec_()


class ConfigDialog(QtWidgets.QDialog, Ui_config.Ui_Dialog_setConfig):

    def __init__(self, parent=None):
        super(ConfigDialog, self).__init__(parent)
        self.setupUi(self)
        self.pushButton_setconfig.clicked.connect(self.set_config)
        self.lineEdit_serverip.setText(config.read_ini("server", "ip"))

    def set_config(self):
        server_ip = self.lineEdit_serverip.text()
        config.write_ini("server", "ip", server_ip)
        self.close()    
    

class PasswordDialog(QtWidgets.QDialog, Ui_password.Ui_DialogPassword):
    def __init__(self, parent=None):
        super(PasswordDialog, self).__init__(parent)
        self.setupUi(self)
        self.pushButton_modify.clicked.connect(self.set_password)
        self.pushButton_cancel.clicked.connect(self.cancel)

    def cancel(self):
        self.close()
    
    def set_password(self):
        password = self.lineEdit_password.text()
        repeat = self.lineEdit_repeat.text()
        if password != repeat:
            show_tips("两次输入的密码不一致，请重新输入！", 1)
            return
        # 移除了密码修改的API调用，直接保存密码
        config.write_ini("user", "password", password)
        self.close()
        
    def do_login_cookie(self, data):
        # 移除了cookie处理
        pass


class LoginDialog(QtWidgets.QDialog, Ui_login.Ui_Dialog):

    # 登录窗口信号
    submitSingal = QtCore.pyqtSignal(str)

    def __init__(self, parent=None):
        super(LoginDialog, self).__init__(parent)
        self.setupUi(self)
        self.mainWindow = MainWindow()
        self.submitSingal.connect(self.mainWindow.do_login_cookie)
        self.pushButton_login.clicked.connect(self.do_login)
        self.pushButton_config.clicked.connect(self.do_config)
        self.lineEdit_username.setText(config.read_ini("user", "name"))
        self.lineEdit_password.setText(config.read_ini("user", "password"))

    def do_login(self):
        server_ip = config.read_ini("server", "ip")
        if server_ip == "":
            show_tips("请先配置服务器地址", 1)
            return
            
        # 移除了登录API调用，直接跳转到主窗口
        self.hide()
        self.mainWindow.show()
        self.submitSingal.emit("")

    def do_config(self):
        self.config_dlg = ConfigDialog()
        self.config_dlg.show()


class MainWindow(QtWidgets.QMainWindow, Ui_main.Ui_MainWindow):

    submitSingal1 = QtCore.pyqtSignal(str)

    def __init__(self, parent=None):
        super().__init__(parent)
        self.setupUi(self)
        self.cookie = {}
        self.devices = []
        self.pos_dict = {}
        

    # def initTableRow(self):
    #     self.tableWidget.setRowCount(len(self.devices))
    #     self.tableWidget.setColumnCount(7)
    #     self.tableWidget.setHorizontalHeaderLabels(["云桌面", "状态", "IP", "", "", "", ""])

    #     for row, item in enumerate(self.devices):
    #         name_item = QTableWidgetItem(item["name"])
    #         status_item = QTableWidgetItem(item["status"])
    #         ip_item = QTableWidgetItem(item["virtualIp"])
    #         self.tableWidget.setItem(row, 0, name_item)
    #         self.tableWidget.setItem(row, 1, status_item)
    #         self.tableWidget.setItem(row, 2, ip_item)

    #         btn_openvm = QPushButton("开机")
    #         btn_openvm.clicked.connect(self.open_vm)
    #         btn_closevm = QPushButton("关机")
    #         btn_closevm.clicked.connect(self.close_vm)
    #         btn_reset = QPushButton("重启")
    #         btn_reset.clicked.connect(self.reset_vm)
    #         btn_connect = QPushButton("连接")
    #         btn_connect.clicked.connect(self.connect_vm)
    #         self.tableWidget.setCellWidget(row, 3, btn_openvm)
    #         self.tableWidget.setCellWidget(row, 4, btn_closevm)
    #         self.tableWidget.setCellWidget(row, 5, btn_reset)
    #         self.tableWidget.setCellWidget(row, 6, btn_connect)
            
    #     self.setCentralWidget(self.tableWidget)

    def initTableWidget(self):
        # 创建一个QWidget作为容器
        central_widget = QWidget(self)
        central_widget.setFixedSize(800, 800)
        self.setCentralWidget(central_widget)

        # 使用QGridLayout布局管理器
        grid_layout = QGridLayout(central_widget)


        # 创建多个QWidget并添加到布局管理器中
        for i in range(len(self.devices)):
            widget = QWidget()
            widget.setFixedSize(250, 200)
            widget.setObjectName("myWidget")
            widget.setStyleSheet("QWidget#myWidget { border: 2px solid green; }")
            vm_name = self.devices[i]["name"]
            status = self.devices[i]["status"]
            ip = self.devices[i]["virtualIp"]
            label1 = QLabel(f'云桌面: {vm_name}')
            label2 = QLabel(f'状态: {status}')
            label3 = QLabel(f'IP: {ip}')
            button1 = QPushButton(f'开机')
            button2 = QPushButton(f'关机')
            button3 = QPushButton(f'重启')
            button4 = QPushButton(f'连接')
            layout = QVBoxLayout()
            layout.addWidget(label1)
            layout.addWidget(label2)
            layout.addWidget(label3)
            layout.addWidget(button1)
            layout.addWidget(button2)
            layout.addWidget(button3)
            layout.addWidget(button4)
            widget.setLayout(layout)

            id = self.devices[i]["id"]
            #self.pos_dict[""]
            
            #button = QPushButton(f'Button {i+1}')
            row = i // 3  # 计算行数
            col = i % 3   # 计算列数
            button1.clicked.connect(self.open_vm1(row, col))
            button2.clicked.connect(self.close_vm1(row, col))
            button3.clicked.connect(self.reset_vm1(row, col))
            button4.clicked.connect(self.connect_vm1(row, col))

            grid_layout.addWidget(widget, row, col)

    def operate_vm1(self, operate, index):
        device = self.devices[index]
        print(device)
        url = "/api/cloud/v1/vm/operate"
        data = {
            "vm_id": device.get("id"),
            "action": str(operate)
        }
        try:
            resp = RequestOperate(server_ip=self.server_ip).post(url, data=data)
            return resp
        except Exception as e:
            print(str(e))
        return None

    def open_vm1(self, row, col):
        def slot():
            index = row * 3 + col
            device = self.devices[index]
            if device.get("status") == "running":
                show_tips("已经是开机状态！")
                return
            resp = self.operate_vm1(1, index)
            if not resp:
                show_tips("开机失败，请查看云桌面状态或联系管理员检查原因", 1)
            if resp.get('code') == 0:
                show_tips("开机中, 请等待状态为running后再进行连接操作!")  
        return slot
    
    def close_vm1(self, row, col):
        def slot():
            index = row * 3 + col
            device = self.devices[index]
            if device.get("status") == "Off":
                show_tips("已经是关机状态！")
                return
            resp = self.operate_vm1(2, index)
            if not resp:
                show_tips("关机失败，请查看云桌面状态或联系管理员检查原因", 1)
        return slot
    
    def reset_vm1(self, row, col):
        def slot():
            index = row * 3 + col
            device = self.devices[index]
            if device.get("status") != "running":
                show_tips("无法操作，请检查云桌面状态或联系管理员检查原因", 1)
                return
            resp = self.operate_vm1(3, index)
            if not resp:
                show_tips("重启失败，请查看云桌面状态或联系管理员检查原因", 1)
        return slot

    def do_modify_pwd(self):
        self.pwdDlg = PasswordDialog()
        self.pwdDlg.show()

    def do_login_cookie(self, data):
        # 移除了cookie处理，直接设置服务器IP并加载虚拟机
        self.server_ip = config.read_ini("server", "ip")
        self.get_user_profile()
        self.load_vms()
        self.timer = QtCore.QTimer(self)
        self.timer.timeout.connect(self.load_vms)
        # 每隔10s触发一次
        self.timer.start(10000)
        self.actionpwd.triggered.connect(self.do_modify_pwd)

    def load_vms(self):
        # 获取云桌面列表
        self.devices = []
        url = "/api/cloud/v1/devices?count=100&index=1&name="
        try:
            resp = RequestOperate(server_ip=self.server_ip).get(url)
            temps = resp['data']['devices']
            user_name = config.read_ini("user", "name")
            for temp in temps:
                if temp["username"] == user_name:
                    self.devices.append(temp)
            self.initTableWidget()
        except Exception as e:
            print(str(e))

    def get_user_profile(self):
        url = "/api/cloud/v1/user_profile"
        try:
            resp = RequestOperate(server_ip=self.server_ip).get(url)
            if resp.get('code') == 0:
                user_info = resp.get('data')
                config.write_ini("user", "name", user_info.get('username'))
                config.write_ini("user", "password",  user_info.get('password'))
                return user_info
        except Exception as e:
            print(str(e))
        return None

    def connect_vm1(self, row, col):
        def slot():
            index = row * 3 + col
            try:
                from pywinauto import Application
                # 通过mstsc连接云桌面
                user_info = self.get_user_profile()
                if not user_info:
                    show_tips("获取用户信息失败", 1)
                    return
                device = self.devices[index]
                mstsc_path = 'C:\\Windows\\System32\\mstsc.exe'
                remote_vm = device.get('virtualIp')
                if not remote_vm:
                    show_tips("IP非法，请重新连接", 1)
                    return
                username = user_info.get('username')
                password = user_info.get('password')

                app = Application().start(mstsc_path)
                #print(app.windows())
                dlg = app.window(title='远程桌面连接')
                #dlg.print_control_identifiers()
                dlg['Edit'].type_keys(remote_vm)
                dlg['ConnectButton'].click()
                cred_dlg = app.window(title='Windows安全中心')
                #cred_dlg.print_control_identifiers()
                cred_dlg['Edit'].type_keys(username)
                cred_dlg['Edit2'].type_keys(password)
                cred_dlg['Button'].click()
                
            except Exception as e:
                pass
        return slot

    # def connect_vm(self):
    #     try:
    #         from pywinauto import Application
    #         # 通过mstsc连接云桌面
    #         user_info = self.get_user_profile()
    #         if not user_info:
    #             show_tips("获取用户信息失败", 1)
    #             return
    #         button = self.sender()
    #         index = self.tableWidget.indexAt(button.pos())
    #         if not index.isValid():
    #             show_tips("云桌面不存在", 1)
    #             return
    #         device = self.devices[index.row()]
    #         mstsc_path = 'C:\\Windows\\System32\\mstsc.exe'
    #         remote_vm = device.get('virtualIp')
    #         if not remote_vm:
    #             show_tips("IP非法，请重新连接", 1)
    #             return
    #         username = user_info.get('username')
    #         password = user_info.get('password')

    #         app = Application().start(mstsc_path)
    #         #print(app.windows())
    #         dlg = app.window(title='远程桌面连接')
    #         #dlg.print_control_identifiers()
    #         dlg['Edit'].type_keys(remote_vm)
    #         dlg['ConnectButton'].click()
    #         cred_dlg = app.window(title='Windows安全中心')
    #         #cred_dlg.print_control_identifiers()
    #         cred_dlg['Edit'].type_keys(username)
    #         cred_dlg['Edit2'].type_keys(password)
    #         cred_dlg['Button'].click()
            
    #     except Exception as e:
    #         pass


if __name__ == '__main__':
    app = QApplication(sys.argv)

    apply_stylesheet(app, theme='dark_teal.xml')
    dialog = LoginDialog()
    dialog.show()
    sys.exit(app.exec_())