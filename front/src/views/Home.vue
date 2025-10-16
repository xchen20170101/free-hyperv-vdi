<template>
  <div>

    <div>
      <el-card class="back">
        <div>
          <div style="font-size: 32px;margin: 20px">欢迎来到云桌面管理平台！</div>
          <!-- <div style="font-size: 32px;margin: 20px">在这个平台上，我们将为您打造一个高效、便捷、安全的数字化工作环境。我们致力于为您提供一站式的云桌面解决方案，让您随时随地都能轻松地访问和管理您的桌面环境。</div> -->
<!--          <img src="@/assets/dashboard.png" style="width: 200px;float: right">-->
        </div>
      </el-card>
    </div>

    <el-row :gutter="12" style="margin-bottom: 30px">
      <el-col :span="12">
        <el-card style="color: #E6A23C">
          <div><i class="el-icon-user" style="margin-right: 5px"></i>用户数</div>
          <div @click="gotoUser" style="text-align: center;font-size: 30px;font-weight: bold">{{ userNum }}</div>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card style="color: #67C23A">
          <div><i class="el-icon-s-platform" style="margin-right: 5px"></i>云桌面数</div>
          <div @click="gotoDevice" style="text-align: center;font-size: 30px;font-weight: bold">{{ deviceNum }}</div>
        </el-card>
      </el-col>
    </el-row>

  </div>
</template>

<script>
import * as echarts from 'echarts'

export default {
  name: "Home",
  data(){
    return{
      userNum:0,
      deviceNum:0,
      count:10,
      article:[]
    }
  },
  created() {
    this.getcount()
  },
  methods:{
    getcount(){
      console.log("get count")
      this.request.get("/api/cloud/v1/user_count").then(res => {
        this.userNum = res.data.data.num
      })
      this.request.get("/api/cloud/v1/device_count").then(res => {
        this.deviceNum = res.data.data.num
      })
    },
    gotoUser() {
      this.$router.push("/user")
    },
    gotoDevice() {
      this.$router.push("/device")
    }
  }
}

</script>

<style scoped>
.back{
  background-image: url("../assets/dashboard.png");
  background-repeat: no-repeat;
  background-position:550px -200px;
  margin-bottom: 20px;
  height:200px;
}

</style>