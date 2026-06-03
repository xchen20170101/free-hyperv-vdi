<template>
  <el-card>
  <div style="margin: 10px 0">
    <el-input style="width: 200px" placeholder="请输入名称" suffix-icon="el-icon-search" v-model="name"></el-input>
    <el-button class="ml-5" type="primary" @click="load">搜索</el-button>
    <el-button type="warning" @click="reset">重置</el-button>
  </div>

  <div style="margin: 10px 0">
    <el-button type="primary" @click="handleAdd">新增 <i class="el-icon-circle-plus-outline"></i></el-button>
  </div>

  <el-table :data="tableData"  :header-cell-class-name="headerBg" @selection-change="handleSelectionChange">
    <!-- <el-table-column prop="id" label="ID"></el-table-column> -->
    <el-table-column prop="name" label="用户名" ></el-table-column>
    <el-table-column prop="password" label="密码" ></el-table-column>
    <el-table-column prop="role" label="角色" ></el-table-column>
    <el-table-column prop="status" label="状态" ></el-table-column>
    <el-table-column label="操作"  align="center">
      <template slot-scope="scope">
        <el-button round type="success" @click="handleEdit(scope.row)">编辑</el-button>
        <el-button v-if="scope.row.name !== 'vmadmin'" round type="success" @click="handleBind(scope.row)">绑定</el-button>
        <el-popconfirm
            v-if="scope.row.name !== 'vmadmin'"
            class="ml-5"
            confirm-button-text='确定'
            cancel-button-text='我再想想'
            icon="el-icon-info"
            icon-color="red"
            title="您确定删除吗？删除后不可恢复，请谨慎操作！"
            @confirm="del(scope.row.id)"
        >
          <el-button round type="danger" slot="reference">删除</el-button>
        </el-popconfirm>
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

  <el-dialog title="设备信息" :visible.sync="dialogBindDeviceVisible" width="30%" >
    <el-form label-width="80px" size="small">
      <el-form-item label="云桌面">
        <el-select v-model="bindForm.deviceNames" multiple placeholder="请选择" style="width: 100%">
          <el-option v-for="item in deviceInfos" :key="item.name" :value="item.name">
            {{ item.name }}
          </el-option>
        </el-select>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button @click="dialogBindDeviceVisible = false">取 消</el-button>
      <el-button type="primary" @click="bind()">绑定</el-button>
    </div>
  </el-dialog>

  <el-dialog title="用户信息" :visible.sync="dialogFormVisible" width="30%" >
    <el-form label-width="80px" size="small">
      <el-form-item label="名称">
        <el-input v-model="form.name" autocomplete="off"></el-input>
      </el-form-item>
      <el-form-item label="密码">
        <el-input v-model="form.password" type="password" autocomplete="off" @input="validatePassword"></el-input>
      </el-form-item>
      <el-form-item label="角色">
        <el-select v-model="form.role" placeholder="请选择" style="width: 100%">
          <el-option value="管理员">管理员</el-option>
          <el-option value="云桌面用户">云桌面用户</el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="状态">
        <el-select v-model="form.status" placeholder="请选择" style="width: 100%">
          <el-option value="启用">启用</el-option>
          <el-option value="禁用">禁用</el-option>
        </el-select>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button @click="dialogFormVisible = false">取 消</el-button>
      <el-button type="primary" @click="save">确 定</el-button>
    </div>
  </el-dialog>

  <el-dialog title="用户信息" :visible.sync="dialogUpdateFormVisible" width="30%" >
    <el-form label-width="80px" size="small">
      <el-form-item label="密码">
        <el-input v-model="form.password" type="password" autocomplete="off" @input="validatePassword"></el-input>
      </el-form-item>
      <el-form-item label="角色">
        <!-- <el-input v-model="form.role" autocomplete="off"></el-input> -->
        <el-select v-model="form.role" placeholder="请选择" style="width: 100%">
          <el-option value="管理员">管理员</el-option>
          <el-option value="云桌面用户">云桌面用户</el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="状态">
        <el-select v-model="form.status" placeholder="请选择" style="width: 100%">
          <el-option value="启用">启用</el-option>
          <el-option value="禁用">禁用</el-option>
        </el-select>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button @click="dialogUpdateFormVisible = false">取 消</el-button>
      <el-button type="primary" @click="update(form.id)">修改</el-button>
    </div>
  </el-dialog>

  </el-card>
</template>

<script>
export default {
  name: "Role",
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
      bindForm: {
        userId: '',
        deviceNames: [],
        strDeviceName: ''
      },
      form: {
        name:'',
        password:'',
        role:'',
        status: ''
      },
      dialogFormVisible: false,
      dialogUpdateFormVisible: false,
      dialogBindDeviceVisible: false,
      menuDialogVis:false,
      menuData:[],
      props:{
        label: 'name',
      },
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
      this.request.get("/api/cloud/v1/users?count="+this.pageSize+"&index="+this.pageNum+"&name="+this.name).then(res => {
        this.tableData = res.data.data.users
        this.total = res.data.data.totalNum
      })
      this.request.get("/api/cloud/v1/unbind_devices").then(res => {
        this.deviceInfos = res.data.data.devices
      })
    },

    save() {
      this.request.post("/api/cloud/v1/users", this.form).then(res => {
        if (res.status===200) {
          if (res.data.msg == "User.Exists") {
            this.$message.error("用户已存在")
          } else if (res.data.msg == "Common.InvalidParam") {
            this.$message.error("参数异常，请检查参数后再重试！")
          } else {
            this.$message.success("保存成功")
          }
          
          this.dialogFormVisible = false
          this.load()
        } else {
          this.$message.error("保存失败")
        }
      })
    },
    validatePassword(event) {
      // 正则表达式匹配字母和数字
      const regex = /^[a-zA-Z0-9]*$/;
      if (!regex.test(event)) {
        // 如果输入包含非字母和数字的字符，则清除该输入
        this.form.password = this.form.password.replace(/[^a-zA-Z0-9]/g, '');
      }
    },

    update(id) {
      this.request.put("/api/cloud/v1/user/" + id, this.form).then(res => {
        if (res.status===200) {
          if (res.data.msg == "User.NotExist") {
            this.$message.error("用户不存在")
          } else if (res.data.msg == "Common.InvalidParam") {
            this.$message.error("参数异常，请检查参数后再重试！")
          } else {
            this.$message.success("更新成功")
          }
          this.dialogUpdateFormVisible = false
          this.load()
        } else {
          this.$message.error("保存失败")
        }
      })
    },

    bind() {
      console.log(this.bindForm)
      this.bindForm.strDeviceName = this.bindForm.deviceNames.join(',')
      this.request.post("/api/cloud/v1/user_bind", this.bindForm).then(res => {
        if (res.status===200) {
          if (res.data.msg == "Common.InvalidParam") {
            this.$message.error("参数异常，请检查参数后再重试！")
          } else {
            this.$message.success("绑定成功")
          }
          this.dialogBindDeviceVisible = false
          this.load()
        } else {
          this.$message.error("绑定失败")
        }
        this.bindForm = {}
      })
    },

    handleAdd() {
      this.dialogFormVisible = true
      this.form = {}
    },

    handleEdit(row) {
      this.form = row
      this.dialogUpdateFormVisible = true
    },
    handleBind(row) {
      this.dialogBindDeviceVisible = true
      this.bindForm.userId = row.id
    },

    del(id) {
      this.request.delete("/api/cloud/v1/user/" + id).then(res => {
        if (res.status===200) {
          if (res.data.msg == "User.NotExist") {
            this.$message.error("用户不存在")
          } else if (res.data.msg == "User.HasBind") {
            this.$message.error("该用户已绑定设备,请解除绑定后再删除")
          } else {
            this.$message.success("删除成功")
          }
          this.load()
        } else {
          this.$message.error("删除失败")
        }
      })
    },

    handleSelectionChange(val) {
      console.log(val)
      this.multipleSelection = val
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
</style>