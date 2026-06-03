<template>
  <el-card>
  <div style="margin: 10px 0">
    <el-input style="width: 200px" placeholder="请输入名称" suffix-icon="el-icon-search" v-model="name"></el-input>
<!--    <el-input style="width: 200px" placeholder="请输入邮箱" suffix-icon="el-icon-message" class="ml-5" v-model="email"></el-input>-->
<!--    <el-input style="width: 200px" placeholder="请输入地址" suffix-icon="el-icon-position" class="ml-5" v-model="address"></el-input>-->
    <el-button class="ml-5" type="primary" @click="load">搜索</el-button>
    <el-button type="warning" @click="reset">重置</el-button>
  </div>

  <div style="margin: 10px 0">
    <el-button type="primary" @click="handleAdd">新增 <i class="el-icon-circle-plus-outline"></i></el-button>
  </div>

  <el-table :data="tableData"  :header-cell-class-name="headerBg" @selection-change="handleSelectionChange">
    <!-- <el-table-column prop="id" label="ID"></el-table-column> -->
    <el-table-column prop="name" label="磁盘名" ></el-table-column>
    <el-table-column prop="capacity" label="总容量" ></el-table-column>
    <el-table-column prop="storagePath" label="存储路径" ></el-table-column>
    <el-table-column prop="belongDevice" label="所属云桌面" ></el-table-column>
    <el-table-column prop="createdTime" label="创建时间" ></el-table-column>
    <el-table-column label="操作"  align="center">
      <template slot-scope="scope">
        <!-- <el-button round type="success" @click="handleEdit(scope.row)" title="编辑"><i class="el-icon-edit"></i></el-button> -->
        <el-button round type="success" @click="handleBind(scope.row)">绑定</el-button>
        <el-button round type="success" @click="handleUnBind(scope.row.id)">解绑</el-button>
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

  <el-dialog title="云桌面信息" :visible.sync="dialogBindDeviceVisible" width="30%" >
    <el-form label-width="80px" size="small">
      <el-form-item label="云桌面">
        <el-select v-model="bindForm.deviceName"  placeholder="请选择" style="width: 100%">
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

  <el-dialog title="磁盘信息" :visible.sync="dialogFormVisible" width="30%" >
    <el-form label-width="80px" size="small">
      <el-form-item label="名称">
        <el-input v-model="form.name" autocomplete="off"></el-input>
      </el-form-item>
      <el-form-item label="容量">
        <el-row>
          <el-col :span="20">
            <el-input v-model="form.capacity" autocomplete="off" placeholder="最大值4TB"></el-input>
          </el-col>
          <el-col :span="4">
            <el-label>GB</el-label>
          </el-col>
        </el-row>
        
      </el-form-item>
      <!-- <el-form-item label="存储路径">
        <el-input v-model="form.storagePath" autocomplete="off"></el-input>
      </el-form-item> -->

    </el-form>
    <div slot="footer" class="dialog-footer">
      <el-button @click="dialogFormVisible = false">取 消</el-button>
      <el-button type="primary" @click="save">确 定</el-button>
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
        diskId: '',
        deviceName: ''
      },
      form: {
        name:'',
        capacity:''
        //storagePath:''
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
      this.request.get("/api/cloud/v1/disks?count="+this.pageSize+"&index="+this.pageNum+"&name="+this.name).then(res => {
        this.tableData = res.data.data.disks
        this.total = res.data.data.totalNum
      })
      this.request.get("/api/cloud/v1/devices?count=1000"+"&index=1"+"&name=").then(res => {
        this.deviceInfos = res.data.data.devices
      })
    },

    save() {
      this.request.post("/api/cloud/v1/disks", this.form).then(res => {
        if (res.status===200) {
          if (res.data.msg == "Disk.Exists") {
            this.$message.error("磁盘已存在")
          } else if (res.data.msg == "Common.InvalidParam") {
            this.$message.error("参数异常，请检查参数后再重试！")
          } else if (res.data.msg == "Disk.OutofLimit") {
            this.$message.error("磁盘容量大小超过限制范围!")
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
          } else if (res.data.msg == "Common.InvalidParam") {
            this.$message.error("参数异常，请检查参数后再重试！")
          } else if (res.data.msg == "Disk.BindToVMFailed") {
            this.$message.error("绑定失败，请检查日志查看具体原因!")
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

    del(id) {
      this.request.delete("/api/cloud/v1/disk/" + id).then(res => {
        if (res.status===200) {
          if (res.data.msg == "Disk.NotExist") {
            this.$message.error("磁盘不存在")
          } else if (res.data.msg == "Disk.HasBind") {
            this.$message.error("磁盘已绑定，请先解绑后再删除。")
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