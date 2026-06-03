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
    <el-table-column prop="name" label="GPU名称" ></el-table-column>
    <el-table-column prop="instancePath" label="实例路径" ></el-table-column>
    <el-table-column prop="bindCount" label="已绑定云桌面数量" ></el-table-column>
    <!-- <el-table-column label="操作"  align="center">
      <template slot-scope="scope">
        <el-button round type="success" @click="handleAllocation(scope.row)">分配</el-button>
      </template>
    </el-table-column> -->
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

  <el-dialog title="云桌面信息" :visible.sync="dialogBindDeviceVisible" width="30%" >
    <el-form label-width="80px" size="small">
      <el-form-item label="云桌面">
        <el-select v-model="form.vmname" multiple placeholder="请选择" style="width: 100%">
          <el-option v-for="item in deviceInfos" :key="item.name" :value="item.name">
            {{ item.name }}
          </el-option>
        </el-select>
      </el-form-item>
    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button @click="dialogBindDeviceVisible = false">取 消</el-button>
      <el-button type="primary" @click="allocation()">分配</el-button>
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
        deviceName: ''
      },
      form: {
        vmname:[],
        strVm: '',
        gpuId: ''
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
      this.request.get("/api/cloud/v1/gpus?count="+this.pageSize+"&index="+this.pageNum+"&name="+this.name).then(res => {
        this.tableData = res.data.data.gpus
        this.total = res.data.data.totalNum
      })
    },

    save() {
      this.request.post("/api/cloud/v1/disks", this.form).then(res => {
        if (res.status===200) {
          if (res.data.msg == "Disk.Exists") {
            this.$message.error("磁盘已存在")
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

    update(id) {
      this.request.put("/api/cloud/v1/user/" + id, this.form).then(res => {
        if (res.status===200) {
          if (res.data.msg == "User.NotExist") {
            this.$message.error("用户不存在")
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
      //console.log(this.bindForm)
      this.request.post("/api/cloud/v1/disk/device_bind", this.bindForm).then(res => {
        if (res.status===200) {
          if (res.data.msg == "Disk.DeviceIsRunning") {
            this.$message.error("云桌面运行中，不支持绑定，请先关闭云桌面后再进行绑定！")
          } else if (res.data.msg == "Disk.BindExist") {
            this.$message.error("该磁盘已被绑定！")
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

    handleAllocation(row) {
      this.form.gpuId = row.id
      this.dialogBindDeviceVisible = true
      this.request.get("/api/cloud/v1/devices?count=1000"+"&index=1"+"&name=").then(res => {
        this.deviceInfos = res.data.data.devices
        var totalNum = res.data.data.totalNum
        this.form.vmname = []
        for (var i = 0; i < totalNum; i++) {
          if (this.deviceInfos[i].gpuInfo === this.form.gpuId)
          {
            this.form.vmname.push(this.deviceInfos[i].name)
          }
        }
      })
      console.log(this.form)
    },
    allocation() {
      this.form.strVm = this.form.vmname.toString()
      console.log(this.form)
      this.request.post("/api/cloud/v1/gpus", this.form).then(res => {
        if (res.status===200) {
          if (res.data.msg == "Gpu.IsLimited") {
            this.$message.error("GPU最多只能绑定4个云桌面！")
          } else {
            this.$message.success("分配成功")
          }
          this.dialogBindDeviceVisible = false
          this.load()
        } else {
          this.$message.error("分配失败")
        }
      })
    },
    handleBind(row) {
      this.dialogBindDeviceVisible = true
      this.bindForm.diskId = row.id
    },
    handleUnBind(id) {
      this.request.post("/api/cloud/v1/disk/device_unbind/" + id).then(res => {
        if (res.status===200) {
          if (res.data.msg == "Disk.BindNotExist") {
            this.$message.error("绑定关系不存在，请检查是否存在绑定关系！")
          } else if (res.data.msg == "Disk.DeviceIsRunning") {
            this.$message.error("云桌面运行中，不支持解绑，请先关闭云桌面后再进行解绑！")
          } else {
            this.$message.success("解绑成功")
          }
          this.load()
        } else {
          this.$message.error("解绑失败")
        }
      })
    },

    handleSelectionChange(val) {
      console.log(val)
      this.multipleSelection = val
    },

    delBatch() {
      let ids = this.multipleSelection.map(v => v.id)  // [{}, {}, {}] => [1,2,3]
      this.request.post("/role/del/batch", ids).then(res => {
        if (res.code==='200') {
          this.$message.success("批量删除成功")
          this.load()
        } else {
          this.$message.error("批量删除失败")
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
</style>