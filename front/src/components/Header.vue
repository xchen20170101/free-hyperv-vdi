<template>
  <div style="line-height: 60px; display: flex">

    <div style="flex: 1">
      <span :class="collapseBtnClass" style="cursor: pointer;font-size: 18px" @click="collapse"></span>
      <el-breadcrumb separator="/" style="display: inline-block; margin-left: 10px">
        <el-breadcrumb-item :to="'/home'">首页</el-breadcrumb-item>
        <el-breadcrumb-item>{{ currentPathName }}</el-breadcrumb-item>
      </el-breadcrumb>
    </div>

    <!-- <span  style="cursor: pointer;font-size: 18px;margin-right: 20px">
      <i class="el-icon-refresh-left"  @click="reload"></i>
    </span>
    <span  style="cursor: pointer;font-size: 18px;float: right">
      <i class="el-icon-full-screen"  @click="handleFullScreen"></i>
    </span> -->

    <el-dropdown style="width: 100px;cursor: pointer;text-align: right">
      <div style="display: inline-block">
        <!-- <img src="../assets/1784393-20191028150722846-959011049.jpg" alt="" style="width: 30px; border-radius: 50%; position: relative; top: 10px; right: 5px"> -->
        <span>{{ userProfile.username }}</span><i class="el-icon-arrow-down" style="margin-left: 5px"></i>
      </div>
      <el-dropdown-menu slot="dropdown">
        <!-- <el-dropdown-item><span @click="person">个人信息</span></el-dropdown-item> -->
        <el-dropdown-item>
          <span style="text-decoration: none" @click="logout">退出</span>
        </el-dropdown-item>
      </el-dropdown-menu>
    </el-dropdown>
  </div>
</template>

<script>
export default {
  name: "Header",
  data() {
    return {
      user: localStorage.getItem("user") ? JSON.parse(localStorage.getItem("user")) : {},
      fullscreen: false,  // 是否全屏
      userProfile: {}
    }
  },
  props: {
    collapseBtnClass: String,
    collapse: Function,
  },
  computed: {
    currentPathName () {
      return this.$store.state.currentPathName;　　//需要监听的数据
    }
  },
  watch: {
    currentPathName (newVal, oldVal) {
      console.log(newVal)
    }
  },
  created() {
    this.load();
  },
  methods:{
    logout() {
      this.request.delete("/api/cloud/v1/logout").then(res=>{
        if (res.status===200){
          if (res.data.msg == "User.LogoutFailed") {
            this.$message.error("退出失败")
          } else {
            this.$message.success("退出成功")
            this.$router.push("/login")
          }
          
        }else {
          this.$message.error("退出失败")
        }
      })
      //this.$router.push("/login")
      //localStorage.removeItem("user")
      //this.$message.success("退出成功")
    },
    load() {
      this.request.get("/api/cloud/v1/user_profile").then(res => {
        this.userProfile = res.data.data
      })
    },
    person(){
      this.$router.push("/person")
    },
    // 全屏事件
    handleFullScreen(){
      let element = document.documentElement;
      // 判断是否已经是全屏
      // 如果是全屏，退出
      if (this.fullscreen) {
        if (document.exitFullscreen) {
          document.exitFullscreen();
        } else if (document.webkitCancelFullScreen) {
          document.webkitCancelFullScreen();
        } else if (document.mozCancelFullScreen) {
          document.mozCancelFullScreen();
        } else if (document.msExitFullscreen) {
          document.msExitFullscreen();
        }
        console.log('已还原！');
      } else {    // 否则，进入全屏
        if (element.requestFullscreen) {
          element.requestFullscreen();
        } else if (element.webkitRequestFullScreen) {
          element.webkitRequestFullScreen();
        } else if (element.mozRequestFullScreen) {
          element.mozRequestFullScreen();
        } else if (element.msRequestFullscreen) {
          // IE11
          element.msRequestFullscreen();
        }
        console.log('已全屏！');
      }
      // 改变当前全屏状态
      this.fullscreen = !this.fullscreen;
    },
    //刷新
    reload(){
      location.reload()
    }
  }
}
</script>

<style scoped>

</style>