<template>

  <div class="login">
    <el-row :gutter="20">
      <el-col :span="6" class="login-img">
        <img src="../assets/logo.png" style="width: 250px">
      </el-col>
      <el-col :span="6">
        <el-form ref="loginRef" :model="form" class="login-form">
          <h1 style="margin-bottom: 30px;font-size: 36px">云桌面管理平台</h1>
          <el-form-item prop="username">
            <el-input
                v-model="form.username"
                type="text"
                auto-complete="off"
                placeholder="账号"
            >
              <template #prefix><i class="el-icon-user-solid"></i></template>
            </el-input>
          </el-form-item>
          <el-form-item prop="password">
            <el-input
                v-model="form.password"
                type="password"
                auto-complete="off"
                placeholder="密码"
                @keyup.enter="login"
            >
              <template #prefix><i class="el-icon-edit"></i></template>
            </el-input>
          </el-form-item>

          <el-form-item style="width:100%;">
            <el-button
                size="large"
                type="primary"
                style="width:100%;"
                @click.prevent="login"
            >
              <span>登 录</span>
            </el-button>
          </el-form-item>
        </el-form>
      </el-col>
    </el-row>
  </div>
</template>

<script>

import {Me} from '@icon-park/vue'
export default {
  name: "login",
  components:{Me},
  data() {
    return {
      form: {
        username: '',
        password: ''
      }
    }
  },
  methods: {
    login() {
      this.request.post("/api/cloud/v1/login",{username:this.form.username,password:this.form.password}).then(res=>{
        if (res.status===200){
          if (res.data.msg == "User.PasswordIsWrong") {
            this.$message.error("用户名或密码错误，请重新输入")
          } else if (res.data.msg == "User.Disable") {
            this.$message.error("用户已禁用，请联系管理员启用该账户后再登录")
          } else if (res.data.msg == "User.NoPermission") {
            this.$message.error("用户无权限，请使用管理员账户登录")
          } else {
            this.$message.success("登录成功")
            this.$router.push("/main")
          }
          
        }else {
          this.$message.error("登录失败")
        }
      })
    }
  }
}
</script>

<style  scoped >
.login {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-image: url("../assets/login_background.jpg");
  background-size: cover;
}
.login-form {
  border-radius: 6px;

  /*background: #ffffff;*/
  width: 400px;
  padding: 150px 15px 150px 150px;
  text-align: center;
}
.login-img {
  border-radius: 6px;
  /*background: #ffffff;*/
  padding: 180px 0px 100px 100px;
}
</style>