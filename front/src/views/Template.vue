<template>
  <el-card>
  <div style="margin: 10px 0">
    <el-input style="width: 200px" placeholder="请输入名称" suffix-icon="el-icon-search" v-model="name"></el-input>
<!--    <el-input style="width: 200px" placeholder="请输入邮箱" suffix-icon="el-icon-message" class="ml-5" v-model="email"></el-input>-->
<!--    <el-input style="width: 200px" placeholder="请输入地址" suffix-icon="el-icon-position" class="ml-5" v-model="address"></el-input>-->
    <el-button class="ml-5" type="primary" @click="load">搜索</el-button>
    <el-button type="warning" @click="reset">重置</el-button>
  </div>

  <el-table :data="tableData"  :header-cell-class-name="headerBg" @selection-change="handleSelectionChange">
    <!-- <el-table-column prop="id" label="ID"></el-table-column> -->
    <el-table-column prop="name" label="模板名称" ></el-table-column>
    <el-table-column prop="username" label="管理员名称" ></el-table-column>
    <el-table-column prop="userpwd" label="管理员密码" ></el-table-column>
    <el-table-column label="操作"  align="center">
      <template slot-scope="scope">
        <el-button round type="primary" @click="handleEdit(scope.row)">配置</el-button>
        <el-button round type="danger"  @click="handleDelete(scope.row)">删除</el-button>
      </template>
    </el-table-column>
  </el-table>
  <!--        分页组件-->
  <div style="padding: 10px 0">
    <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="pageNum"
        :page-sizes="[5, 10, 15, 20]"
        :page-size="pageSize"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total">
    </el-pagination>
  </div>

  <el-dialog title="模板配置" :visible.sync="dialogTemplateVisible" width="30%" >
    <el-form label-width="80px" size="small">
      <el-form-item label="管理员">
        <el-input v-model="form.username" autocomplete="off"></el-input>
      </el-form-item>
      <el-form-item label="密码">
        <el-input v-model="form.userpwd" autocomplete="off"></el-input>
      </el-form-item>
    </el-form>

    <div slot="footer" class="dialog-footer">
      <el-button @click="dialogTemplateVisible = false">取 消</el-button>
      <el-button type="primary" @click="setConfig(form.id)">配置</el-button>
    </div>
  </el-dialog>

  </el-card>
</template>

<script>
export default {
  name: "Template",
  data(){
    return{
      tableData:[],
      total:0,
      pageNum:1,
      pageSize:5,
      name:'',
      email: '',
      address: '',
      status: '',
      form: {
        username:'',
        userpwd: ''
      },
      dialogTemplateVisible: false,
      expends:[],
      checks:[],
      roleId:1,
      multipleSelection: [],
      headerBg:"headerBg",
      deviceInfos: []
    }
  },
  created() {
    this.load();
  },
  methods:{
    load() {
      this.request.get("/api/cloud/v1/templates?count="+this.pageSize+"&index="+this.pageNum+"&name="+this.name).then(res => {
        this.tableData = res.data.data.templates
        this.total = res.data.data.totalNum
      })
    },

    setConfig(id) {
      this.request.put("/api/cloud/v1/templates/"+id,this.form).then(res => {
        if (res.status == 200) {
          if (res.data.msg == "Template.NotExist") {
            this.$message.error("模板不存在");
          } else if (res.data.msg == "Common.InvalidParam") {
            this.$message.error("参数错误");
          } else {
            this.$message.success("配置成功");
          }
          this.dialogTemplateVisible = false
          this.load()
        } else {
          this.$message.error("配置失败")
        }

      })
    },

    handleEdit(row) {
      this.dialogTemplateVisible = true
      this.form = {
        id: row.id,
        username:row.username,
        userpwd: row.userpwd
      }
    },
    handleDelete(row) {
      let id = row.id
      this.request.delete("/api/cloud/v1/template/" + id).then(res => {
        if (res.status===200) {
          if (res.data.msg == "Template.NotExist") {
            this.$message.error("模板不存在")
          } else {
            this.$message.success("删除成功")
          }
          this.load()
        } else {
          this.$message.error("删除失败")
        }
      })
    },
    reset() {
      this.name = ""
      this.load()
    },

    handleSizeChange(pageSize) {
      console.log(pageSize)
      this.pageSize = pageSize
      this.load()
    },
    handleCurrentChange(pageNum) {
      console.log(pageNum)
      this.pageNum = pageNum
      this.load()
    }
  }
}
</script>

<style>
.headerBg {
  background: #eee !important;
}
.dialog-header span{
  margin-bottom: 10px; /* 设置与表单之间的垂直间距 */
  font-size: 14px; /* 设置字体大小 */
  color: #666666; /* 设置字体颜色 */
}

</style>