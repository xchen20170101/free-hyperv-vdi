<template>

    <el-container style="height: 100vh">
      <el-aside :width="sideWidth+'px'" style="background-color: rgb(238, 241, 246);box-shadow: 2px 0 6px rgb(0 21 41/35%)">
        <Aside :isCollapse="isCollapse" :menus="user.menus"/>
      </el-aside>

      <el-container>
        <el-header style="border-bottom: #ababab 1px solid">
          <Header :collapseBtnClass="collapseBtnClass" :collapse="collapse"/>
        </el-header>

        <el-main>
          <router-view/>
        </el-main>
      </el-container>
    </el-container>

</template>

<script>

import Aside from "@/components/Aside";
import Header from "@/components/Header";

export default {
  name: 'Manage',
  components: { Header, Aside},
  data() {
    return {
      isCollapse: false,
      collapseBtnClass:'el-icon-s-fold',
      sideWidth:200,
      user: localStorage.getItem("user") ? JSON.parse(localStorage.getItem("user")) : {}
    }
  },
  methods:{
    //折叠菜单栏
    collapse(){
      this.isCollapse=!this.isCollapse;
      if (this.isCollapse){
        this.sideWidth=64
        this.collapseBtnClass='el-icon-s-unfold'
      }else {
        this.sideWidth=200
        this.collapseBtnClass='el-icon-s-fold'
      }
    }

  }
}
</script>


<style>
.el-header {

  color: #333;
  line-height: 60px;
}

.el-aside {
  color: #333;
}

</style>