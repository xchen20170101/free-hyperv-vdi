<template>
    <view class="content">
        <view class="input-group">
            <view class="input-row border">
                <text class="title">用户名：</text>
                <m-input class="m-input" type="text" clearable focus v-model="userName" placeholder="请输入用户名"></m-input>
            </view>
            <view class="input-row border">
                <text class="title">原密码：</text>
                <m-input class="m-input" type="password" clearable focus v-model="oldPwd"
                    placeholder="请输入原密码"></m-input>
            </view>
            <view class="input-row border">
                <text class="title">新密码：</text>
                <m-input class="m-input" type="password" clearable focus v-model="newPwd"
                    placeholder="请输入新密码"></m-input>
            </view>
        </view>
        <view class="btn-row">
            <button type="primary" class="primary larger-btn" :loading="loginBtnLoading"
                @tap="bindSetServer">保存</button>
            <button type="primary" class="primary larger-btn" :loading="loginBtnLoading"
                @tap="bindCancelServer">取消</button>
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
                userName: '',
                oldPwd: '',
                newPwd: '',
                positionTop: 0,
                isDevtools: false,
                loginBtnLoading: false,
            }
        },
        methods: {
            initPosition() {
                /**
                 * 使用 absolute 定位，并且设置 bottom 值进行定位。软键盘弹出时，底部会因为窗口变化而被顶上来。
                 * 反向使用 top 进行定位，可以避免此问题。
                 */
                this.positionTop = uni.getSystemInfoSync().windowHeight - 100;
            },
            bindSetServer() {
                var serverIp = uni.getStorageSync("serverIp")
                if (serverIp == "") {
                    uni.showToast({
                        title: '未配置服务端IP',
                        icon: 'none',
                        position: 'bottom'
                    });
                    uni.navigateBack()
                    return
                }
                //console.log("serverIp:", serverIp)
                var targetUrl = 'http://' + serverIp + ':8090/api/cloud/v1/reset_password'
                //console.log("targetUrl:", targetUrl)
                uni.request({
                    url: targetUrl,
                    method: 'POST',
                    header: {
                        'content-type': 'application/x-www-form-urlencoded' // 设置请求头为表单形式
                    },
                    data: {
                        username: this.userName,
                        oldpassword: this.oldPwd,
                        newpassword: this.newPwd
                    },
                    success: function(res) {
                        if (res.data.code == 0) {
                            uni.setStorageSync("password", this.newPwd)
                            uni.navigateTo({
                                url: "/pages/login/login"
                            })
                        } else {
                            if (res.data.msg == "User.PasswordIsWrong") {
                                uni.showToast({
                                    title: '用户名或密码错误!',
                                    icon: 'none',
                                    position: 'bottom'
                                });
                            } else if (res.data.msg == "User.Disable") {
                                uni.showToast({
                                    title: '用户已禁用，请联系管理员启用后再登录!',
                                    icon: 'none',
                                    position: 'bottom'
                                })
                            } else if (res.data.msg == "User.NotExist") {
                                uni.showToast({
                                    title: '用户不存在，请检查用户是否正确后再重试!',
                                    icon: 'none',
                                    position: 'bottom'
                                })
                            }
                        }
                        // 在这里处理获取到的数据
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
                uni.navigateBack()
            },
            bindCancelServer() {
                console.log("cancel server")
                // uni.redirectTo({
                //     url: "/pages/login/login"
                // })
                uni.navigateBack()
            },
            onReady() {
                this.initPosition();
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

    .oauth-image {
        position: relative;
        width: 50px;
        height: 50px;
        border: 1px solid #dddddd;
        border-radius: 50px;
        background-color: #ffffff;
    }

    .oauth-image image {
        width: 30px;
        height: 30px;
        margin: 10px;
    }

    .oauth-image button {
        position: absolute;
        left: 0;
        top: 0;
        width: 100%;
        height: 100%;
        opacity: 0;
    }

    .captcha-view {
        line-height: 0;
        justify-content: center;
        align-items: center;
        display: flex;
        position: relative;
        background-color: #f3f3f3;
    }

    .btn-row {
        display: flex;
        justify-content: space-between;
        /* Adjust this according to your layout needs */
        /* Other styles as needed */
    }

    /* Additional styles for buttons (modify as needed) */
    .btn-row button {
        /* Add specific button styles here */
        margin: 5px;
        /* Adjust spacing between buttons */
    }

    .larger-btn {
        padding: 5px 40px;
        /* Adjust padding to increase button size */
        font-size: 16px;
        /* Adjust font size */
        /* Other styles for larger buttons */
    }
</style>