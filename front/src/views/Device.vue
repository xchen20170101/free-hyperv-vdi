<template>
  <el-card>
  <div style="margin: 10px 0">
    <el-input style="width: 200px" placeholder="请输入名称" suffix-icon="el-icon-search" v-model="name"></el-input>
<!--    <el-input style="width: 200px" placeholder="请输入邮箱" suffix-icon="el-icon-message" class="ml-5" v-model="email"></el-input>-->
<!--    <el-input style="width: 200px" placeholder="请输入地址" suffix-icon="el-icon-position" class="ml-5" v-model="address"></el-input>-->
    <el-button class="ml-5" type="primary" @click="load">搜索</el-button>
    <el-button type="warning" @click="reset">刷新</el-button>
  </div>

  <div style="margin: 10px 0">
    <el-button type="primary" @click="handleAdd" title="基于模板创建">快速新建 <i class="el-icon-circle-plus-outline"></i></el-button>
  </div>

  <el-table :data="tableData"  :header-cell-class-name="headerBg" @selection-change="handleSelectionChange">
    <!-- <el-table-column prop="id" label="ID"></el-table-column> -->
    <el-table-column prop="name" label="云桌面名称" ></el-table-column>
    <el-table-column prop="username" label="所属用户" ></el-table-column>
    <el-table-column prop="templateInfo" label="所属模板" ></el-table-column>
    <el-table-column prop="virtualIp" label="IP" ></el-table-column>
    <el-table-column prop="memoryInfo" label="内存" ></el-table-column>
    <el-table-column prop="cpuInfo" label="CPU" ></el-table-column>
    <el-table-column prop="gpuInfo" label="GPU" ></el-table-column>
    <el-table-column prop="status" label="状态" ></el-table-column>
    <el-table-column prop="createdTime" label="创建时间" ></el-table-column>
    <el-table-column label="操作">
      <template slot-scope="scope">
        <div class="button-container">
          <el-button round type="success" @click="openVM(scope.row.id)">开机</el-button>
          <el-popconfirm
              class="ml-5"
              confirm-button-text='确定'
              cancel-button-text='我再想想'
              icon="el-icon-info"
              icon-color="red"
              title="您确定关机吗？"
              @confirm="closeVM(scope.row.id)"
          >
            <el-button round type="danger" slot="reference">关机</el-button>
          </el-popconfirm>
        </div>
      </template>
    </el-table-column>
    <el-table-column label="" >
      <template slot-scope="scope">
        <div class="button-container">
        <el-popconfirm
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
          <el-button round type="success" @click="handleEdit(scope.row)">编辑</el-button>
        </div>
      </template>
    </el-table-column>
    <el-table-column label="" >
      <template slot-scope="scope">
        <div class="button-container">
          <el-popconfirm
              class="ml-5"
              confirm-button-text='确定'
              cancel-button-text='我再想想'
              icon="el-icon-info"
              icon-color="red"
              title="您确定重启吗？"
              @confirm="resetVM(scope.row.id)"
          >
            <el-button round type="danger" slot="reference">重启</el-button>
          </el-popconfirm>
          <el-button round type="danger" @click="unBindUser(scope.row.id)">用户解绑</el-button>
        </div>
      </template>
    </el-table-column>
    <el-table-column label=""  align="center">
      <template slot-scope="scope">
        <div class="button-container">
          <el-button round type="success" @click="bindGpu(scope.row)">GPU绑定</el-button>
          <el-button round type="danger" @click="unBindGPU(scope.row)">GPU解绑</el-button>
        </div>
      </template>
    </el-table-column>
    <el-table-column label=""  align="center">
      <template slot-scope="scope">
        <div class="button-container">
          <el-button round type="warning" @click="resetUserPwd(scope.row)" class="reset-pwd-btn">重置密码</el-button>
        </div>
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



  <el-dialog title="云桌面信息" :visible.sync="dialogFormVisible" width="50%" >
    <el-form label-width="80px" size="small">
      <el-form-item label="名称">
        <el-input v-model="addForm.name" autocomplete="off" @input="validateVMName" placeholder="只能包含字母、数字和下划线，不能以下划线开头"></el-input>
        <div v-if="nameError" style="color: #F56C6C; font-size: 12px; margin-top: 5px;">{{ nameError }}</div>
      </el-form-item>

    <el-form-item label="模板" prop="srcVmPath">
      <!-- <el-input v-model="addForm.srcVmPath" autocomplete="off" placeholder="请输入模板云桌面文件的路径，格式：D:\VM\Virtual Machines\example.vmcx"></el-input> -->
      <el-select v-model="addForm.srcVmPath" placeholder="请选择" style="width: 100%">
          <el-option v-for="item in templates" :key="item" :value="item">
            {{ item }}
          </el-option>
        </el-select>
    </el-form-item>

    <el-form-item label="网络" prop="vmSwitch">
      <el-select v-model="addForm.vmSwitch" placeholder="请选择" style="width: 100%">
          <el-option v-for="item in switchs" :key="item" :value="item">
            {{ item }}
          </el-option>
        </el-select>
    </el-form-item>

    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button @click="dialogFormVisible = false">取 消</el-button>
      <el-button type="primary" @click="save">确 定</el-button>
    </div>
  </el-dialog>

  <el-dialog title="云桌面信息" :visible.sync="dialogUpdateFormVisible" width="30%" >
    <el-form label-width="80px" size="small">
      <!-- <el-form-item label="名称">
        <el-input v-model="editForm.name" autocomplete="off"></el-input>
      </el-form-item> -->

      <el-form-item label="内存">
        <el-select v-model="editForm.memoryInfo" placeholder="请选择内存大小" style="width: 100%">
          <el-option v-for="item in memoryOptions" :key="item" :label="item" :value="item">
            {{ item }}
          </el-option>
        </el-select>
      </el-form-item>

      <el-form-item label="CPU">
        <el-select v-model="editForm.cpuInfo" placeholder="请选择CPU数量" style="width: 100%">
        <el-option v-for="item in cpuOptions" :key="item.value" :label="item.label" :value="item.value">
          {{ item.label }}
        </el-option>
      </el-select>
      </el-form-item>

      <!-- <el-form-item label="所属用户">
        <el-select v-model="editForm.username" placeholder="请选择" style="width: 100%">
          <el-option v-for="item in userInfos" :key="item.name" :value="item.name">
            {{ item.name }}
          </el-option>
        </el-select>
      </el-form-item> -->

    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button @click="dialogUpdateFormVisible = false">取 消</el-button>
      <el-button type="primary" @click="update(editForm.id)">修改</el-button>
    </div>
  </el-dialog>


  <el-dialog title="GPU信息" :visible.sync="dialogGpuVisible" width="30%" >
    <el-form label-width="80px" size="small">
      <el-form-item label="GPU">
        <el-select v-model="gpuForm.gpuId" placeholder="请选择GPU" style="width: 100%">
          <el-option v-for="item in gpus" :key="item.name" :label="item.name" :value="item.id">
            {{ item.name }}
          </el-option>
        </el-select>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button @click="dialogGpuVisible = false">取 消</el-button>
      <el-button type="primary" @click="gpuUpdate()">绑定</el-button>
    </div>
  </el-dialog>

  <el-dialog title="重置Windows账户密码" :visible.sync="dialogPwdVisible" width="30%" >
    <el-form label-width="80px" size="small">
      <el-form-item label="云桌面">
        <el-input v-model="pwdForm.deviceName" disabled></el-input>
      </el-form-item>
      <el-form-item label="新密码">
        <el-input v-model="pwdForm.newPassword" show-password placeholder="请输入新密码"></el-input>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button @click="dialogPwdVisible = false">取 消</el-button>
      <el-button type="primary" @click="submitResetPwd()">确 定</el-button>
    </div>
  </el-dialog>

  </el-card>
</template>

<script>

export default {
  name: "Device",
  data(){
    return{
      tableData:[],
      total:0,
      pageNum:1,
      pageSize:5,
      name:'',
      username: '',
      virtualIp: '',
      status: '',
      createdTime: '',
      nameError: '',
      gpuForm: {
        gpuId: '',
        deviceName: ''
      },
      addForm: {
        name:'',
        srcVmPath: '',
        vmSwitch: ''
      },
      editForm: {
        name:'',
        memoryInfo: '',
        cpuInfo: ''
      },
      configForm: {
        memoryInfo: '',
        cpuInfo: ''
      },
      dialogFormVisible: false,
      dialogUpdateFormVisible: false,
      dialogConfigFormVisible: false,
      dialogGpuVisible: false,
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
      userInfos: [],
      memoryOptions: ['2GB', '4GB', '8GB', '16GB'],
      cpuOptions: [
        {value: '1', label: '1'},
        {value: '2', label: '2'},
        {value: '4', label: '4'},
        {value: '8', label: '8'}
      ],
      timer: null,
      templates: [],
      gpus: [],
      switchs: [],
      pwdForm: {
        deviceId: '',
        deviceName: '',
        newPassword: ''
      },
      dialogPwdVisible: false
    }
  },
  created() {
    this.load();
  },
  mounted() {
    this.timer = setInterval(this.load, 10000)
  },
  beforeDestroy() {
    clearInterval(this.timer)
  },
  methods:{
    load() {
      this.request.get("/api/cloud/v1/devices?count="+this.pageSize+"&index="+this.pageNum+"&name="+this.name).then(res => {
        this.tableData = res.data.data.devices
        this.total = res.data.data.totalNum
      })
    },
    validateVMName() {
      const name = this.addForm.name;
      if (!name) {
        this.nameError = '';
        return false;
      }
      
      // 检查长度
      if (name.length > 64) {
        this.nameError = '云桌面名称长度不能超过64个字符';
        return false;
      }
      
      // 检查格式：只能包含字母、数字和下划线，不能以下划线开头
      const namePattern = /^[a-zA-Z0-9][a-zA-Z0-9_]{0,63}$/;
      if (!namePattern.test(name)) {
        this.nameError = '云桌面名称只能包含字母、数字和下划线，不能以下划线开头，最大长度64个字符';
        return false;
      }
      
      this.nameError = '';
      return true;
    },
    save() {
      console.log(this.addForm)
      
      // 验证云桌面名称
      if (!this.validateVMName()) {
        return;
      }
      
      this.request.post("/api/cloud/v1/vm", this.addForm).then(res => {
        if (res.status===200) {
          if (res.data.msg == "Device.Exist") {
            this.$message.error("云桌面已存在")
          } else if (res.data.msg == "Common.InvalidParam") {
            this.$message.error("参数异常，请检查参数后再重试！")
          } else if (res.data.msg == "Device.TemplateNotExist") {
            this.$message.error("模板文件不存在，请检查后再重试!")
          } else if (res.data.msg == "Device.InvalidName") {
            this.$message.error("云桌面名称只能包含字母、数字和下划线，不能以下划线开头，最大长度64个字符")
          } else {
            this.$message.success("云桌面创建中...")
          }
          this.dialogFormVisible = false
          this.load()
        } else {
          this.$message.error("保存失败")
        }
      })
    },

    update(id) {
      console.log("update vm")
      this.request.put("/api/cloud/v1/vm/" + id, this.editForm).then(res => {
        if (res.status===200) {
          if (res.data.msg == "Common.InvalidParam") {
            this.$message.error("参数异常，请检查参数后再重试！")
          } else if (res.data.msg == "Device.IsRunning") {
            this.$message.error("云桌面处于运行状态，无法删除，请先关闭云桌面后再操作！")
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

    unBindUser(id) {
      let unBindForm = {
        deviceId: id
      }
      this.request.post("/api/cloud/v1/unbind_user", unBindForm).then(res => {
        if (res.status === 200) {
          if (res.data.msg == "Common.InvalidParam") {
            this.$message.error("参数异常，请检查参数后再重试！")
          } else if (res.data.msg == "Device.UnBindUserFailed") {
            this.$message.error("云桌面用户解绑失败，请检查后再重试!")
          } else if (res.data.msg == "Device.MustBeRunning") {
            this.$message.error("云桌面必须是运行状态才能解绑!")
          } else if (res.data.msg == "Device.TemplateNotExist") {
            this.$message.error("模板文件不存在，请检查后再重试!")
          } else {
            this.$message.success("解绑成功")
          }
          this.load()
        } else {
          this.$message.error("解绑失败")
        }
      })
    },
    handleAdd() {
      this.request.get("/api/cloud/v1/device_templates").then(res => {
        this.templates = res.data.data.templates
      })
      this.request.get("/api/cloud/v1/device_switchs").then(res=> {
        this.switchs = res.data.data.switchs
      })
      this.dialogFormVisible = true
      this.addForm = {}
      this.nameError = ''
    },

    handleEdit(row) {
      this.editForm = row
      this.dialogUpdateFormVisible = true
    },
    bindGpu(row) {
      this.gpuForm.deviceName = row.name
      this.dialogGpuVisible = true
      this.request.get("/api/cloud/v1/gpus?count=1000&index=1").then(res => {
        this.gpus = res.data.data.gpus
      })
    },
    unBindGPU(row) {
      console.log(row)
      let unbind_form = {
        deviceId: row.id
      }
      this.request.post("/api/cloud/v1/unbind_gpu", unbind_form).then(res => {
        if (res.status === 200) {
          if (res.data.msg == "Common.InvalidParam") {
            this.$message.error("参数异常，请检查参数后再重试！")
          } else {
            this.$message.success("解绑成功")
            this.gpuForm = {
              gpuId: '',
              deviceName: ''
            }
            this.dialogGpuVisible = false
          }
          this.load()
        } else {
          this.$message.error("解绑失败")
        }
      })
    },
    gpuUpdate() {
      console.log(this.gpuForm)
      this.request.post("/api/cloud/v1/bind_gpu", this.gpuForm).then(res => {
        if (res.status === 200) {
          if (res.data.msg == "Common.InvalidParam") {
            this.$message.error("参数异常，请检查参数后再重试！")
          } else if (res.data.msg == "Gpu.NotExist") {
            this.$message.error("GPU不存在")
          } else if (res.data.msg == "Gpu.IsLimited") {
            this.$message.error("一个GPU最多只能被4个云桌面使用")
          } else if (res.data.msg == "Gpu.BindingExist") {
            this.$message.error("已绑定GPU,请先解绑再进行绑定")
          } else if (res.data.msg == "Gpu.BindFailed") {
            this.$message.error("绑定失败,请在设备管理器上检查GPU状态后再重试")
          } else {
            this.$message.success("绑定成功")
            this.gpuForm = {
              gpuId: '',
              deviceName: ''
            }
            this.dialogGpuVisible = false
          }
          this.load()
        } else {
          this.$message.error("绑定失败")
        }
      })
    },
    // handleFileChange(field, file) {
    //   console.log(field)
    //   console.log(file)
    // },
    // handleSuccess(file) {
    //   this.addForm.srcVmPath = file.raw.name
    //   console.log(file)
    //   console.log(URL.createObjectURL(file.raw))
    //   console.log(document.getElementsByClassName("el-upload__input")[0].value)
    // },
    handleVMConfig(row) {
      // console.log(row)
      this.configForm = row
      this.dialogConfigFormVisible = true
    },
    setVMConfig() {
      console.log(this.configForm)
    },
    openVM(id) {
      let obj = {
        "vm_id": id,
        "action": 1
      }
      this.request.post("/api/cloud/v1/vm/operate", obj).then(res => {
        if (res.status===200) {
          if (res.data.code == 2) {
            this.$message.error("无法操作，请检查云桌面的状态或联系管理员！")
          } else {
            this.$message.success("操作成功")
          }
          this.load()
        } else {
          this.$message.error("操作失败")
        }
      })
    },
    closeVM(id) {
      let obj = {
        "vm_id": id,
        "action": 2
      }
      this.request.post("/api/cloud/v1/vm/operate", obj).then(res => {
        if (res.status===200) {
          if (res.data.code == 2) {
            this.$message.error("无法操作，请检查云桌面的状态或联系管理员！")
          } else {
            this.$message.success("操作成功")
          }
          this.load()
        } else {
          this.$message.error("操作失败")
        }
      })
    },
    resetVM(id) {
      let obj = {
        "vm_id": id,
        "action": 3
      }
      this.request.post("/api/cloud/v1/vm/operate", obj).then(res => {
        if (res.status===200) {
          if (res.data.code == 2) {
            this.$message.error("无法操作，请检查云桌面的状态或联系管理员！")
          } else {
            this.$message.success("操作成功")
          }
          this.load()
        } else {
          this.$message.error("操作失败")
        }
      })
    },

    del(id) {
      this.request.delete("/api/cloud/v1/vm/" + id).then(res => {
        if (res.status===200) {
          if (res.data.msg == "Device.IsRunning") {
            this.$message.error("云桌面处于运行状态，无法删除，请先关闭云桌面后再操作！")
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
    },
    resetUserPwd(row) {
      this.pwdForm.deviceId = row.id;
      this.pwdForm.deviceName = row.name;
      this.dialogPwdVisible = true;
    },
    
    submitResetPwd() {
      if (!this.pwdForm.newPassword) {
        this.$message.error("请输入新密码");
        return;
      }
      this.request.post("/api/cloud/v1/reset_user_pwd", this.pwdForm).then(res => {
        if (res.status === 200) {
          if (res.data.msg == "Common.InvalidParam") {
            this.$message.error("参数异常，请检查参数后再重试！")
          } else if (res.data.msg == "Device.ResetUserPwdFailed") {
            this.$message.error("重置Windows账户密码失败，请检查后再重试!")
          } else if (res.data.msg == "Device.MustBeRunning") {
            this.$message.error("云桌面必须是运行状态才能重置密码!")
          } else if (res.data.msg == "Device.TemplateNotExist") {
            this.$message.error("模板文件不存在，请检查后再重试!")
          } else {
            this.$message.success("密码重置成功")
            this.dialogPwdVisible = false
            this.pwdForm = {
              deviceId: '',
              deviceName: '',
              newPassword: ''
            }
          }
          this.load()
        } else {
          this.$message.error("重置失败")
        }
      })
    }

  }
}
</script>

<style>
.headerBg {
  background: #eee !important;
}

.button-container {
  display: flex;
  flex-direction: column;
  align-items: center; /* 将按钮容器内元素垂直对齐方式设置为顶部对齐 */
}

.button-container > * {
  margin-bottom: 5px; /* 调整按钮之间的垂直间距 */
}

.button-container .el-button {
  width: 100px; /* 设置按钮宽度 */
  height: 30px; /* 设置按钮高度 */
}

.reset-pwd-btn {
  width: 100px; /* Adjust the width as needed */
  text-align: center; /* Center the text */
}

</style>