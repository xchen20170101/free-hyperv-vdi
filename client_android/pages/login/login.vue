<template>
    <view class="content">
        <view class="input-group">
            <view class="input-row border">
                <text class="title">账号：</text>
                <m-input class="m-input" type="text" clearable focus v-model="form.username"
                    placeholder="请输入账号"></m-input>
            </view>
            <view class="input-row border">
                <text class="title">密码：</text>
                <m-input type="password" displayable v-model="form.password" placeholder="请输入密码"></m-input>
            </view>
        </view>
        <view class="input-row border">
            <text class="title">记住密码：</text>
            <switch :checked="rememberPassword" @change="handleSwitchChange"></switch>
            <text class="reset-password" @tap="redirectToResetPassword">重置密码</text>
        </view>
        <view class="btn-row">
            <button type="primary" class="primary larger-btn" :loading="loginBtnLoading" @tap="bindLogin">登录</button>
            <button type="primary" class="primary larger-btn" :loading="loginBtnLoading"
                @tap="bindSetConfig">配置</button>
        </view>
    </view>
</template>

<script>
    import mInput from '../../components/m-input.vue'

    export default {
        components: {
            mInput
        },
        data() {
            return {
                form: {
                    username: '',
                    password: ''
                },
                positionTop: 0,
                isDevtools: false,
                loginBtnLoading: false,
                rememberPassword: true
            }
        },
        onLoad() {
            this.form.password = uni.getStorageSync("password")
            this.form.username = uni.getStorageSync("username")
            let rememberPassword = uni.getStorageSync('rememberPassword');
            console.log('Remember Password:', rememberPassword)
            if (rememberPassword === null || rememberPassword == "") {
                rememberPassword = true
                uni.setStorageSync('rememberPassword', rememberPassword);
            }
            this.rememberPassword = rememberPassword;
        },
        methods: {
            bindLogin() {
                console.log(this.rememberPassword)
                uni.setStorageSync("username", this.form.username)
                if (this.rememberPassword) {
                    uni.setStorageSync("password", this.form.password)
                }

                var serverIp = uni.getStorageSync("serverIp")
                if (serverIp == "") {
                    uni.showToast({
                        title: '未配置服务端IP',
                        icon: 'none',
                        position: 'bottom'
                    });
                }
                var targetUrl = 'http://' + serverIp + ':8090/api/cloud/v1/login'
                uni.request({
                    url: targetUrl,
                    method: 'POST',
                    header: {
                        'content-type': 'application/x-www-form-urlencoded' // 设置请求头为表单形式
                    },
                    data: {
                        username: this.form.username,
                        password: this.form.password
                    },
                    success: function(res) {
                        console.log('POST 请求成功', res)
                        // 移除了授权检查，直接跳转到设备页面
                        uni.navigateTo({
                            url: "/pages/device/device"
                        })
                    },
                    fail: function(err) {
                        console.error('POST 请求失败', err);
                        // 处理请求失败的情况
                        uni.showToast({
                            title: '网络连接错误！',
                            icon: 'none',
                            position: 'bottom'
                        });
                    }
                });
            },
            bindSetConfig() {
                uni.navigateTo({
                    url: "/pages/config/config"
                })
            },
            redirectToResetPassword() {
                uni.navigateTo({
                    url: "/pages/password/password"
                })
            },
            onload(options) {
                const serverIp = uni.getStorageSync("serverIp")
            },
            handleSwitchChange(e) {
                uni.setStorageSync('rememberPassword', e.detail.value);
            }
        },
    }
</script>

<style>
    .login-type {
        display: flex;
        justify-content: center;
    }

    .login-type-btn {
        line-height: 30px;
        margin: 0px 15px;
    }

    .login-type-btn.act {
        color: #0FAEFF;
        border-bottom: solid 1px #0FAEFF;
    }

    .send-code-btn {
        width: 120px;
        text-align: center;
        background-color: #0FAEFF;
        color: #FFFFFF;
    }

    .action-row {
        display: flex;
        flex-direction: row;
        justify-content: center;
    }

    .action-row navigator {
        color: #007aff;
        padding: 0 10px;
    }

    .oauth-row {
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: space-around;
        flex-wrap: wrap;
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
    }
</style>