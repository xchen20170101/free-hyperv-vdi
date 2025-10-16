package form


type AddUserForm struct {
	Name     string `form:"name" binding:"required"`
	Password string `form:"password" binding:"required"`
	Role     string `form:"role" binding:"required"`
	Status   string `form:"status" binding:"required"`
}


type UpdateUserForm struct {
	Password string `form:"password" binding:"required"`
	Role     string `form:"role" binding:"required"`
	Status   string `form:"status" binding:"required"`
}


type BindUserForm struct {
	UserId        string `form:"userId" binding:"required"`
	StrDeviceName string `form:"strDeviceName" binding:"required"`
}


type ActiveLicenseForm struct {
	MachineCode string `form:"machineCode" binding:"required"`
	LicenseCode string `form:"licenseCode" binding:"required"`
}


type UpdateVmForm struct {
	CpuInfo    string `form:"cpuInfo" binding:"required"`
	MemoryInfo string `form:"memoryInfo" binding:"required"`
}


type AddVmForm struct {
	Name      string `form:"name" binding:"required"`
	VmSwitch  string `form:"vmSwitch" binding:"required"`
	SrcVmPath string `form:"srcVmPath" binding:"required"`
}


type OperateVmForm struct {
	VmId   string `form:"vm_id" binding:"required"`
	Action string `form:"action" binding:"required"`
}


type AddDiskForm struct {
	Name     string `form:"name" binding:"required"`
	Capacity string `form:"capacity" binding:"required"`
}


type BindDiskForm struct {
	DiskId     string `form:"diskId" binding:"required"`
	DeviceName string `form:"deviceName" binding:"required"`
}


type DeviceUnbindUserForm struct {
	DeviceId string `form:"deviceId" binding:"required"`
}


type DeviceBindGpuForm struct {
	GpuId      string `form:"gpuId" binding:"required"`
	DeviceName string `form:"deviceName" binding:"required"`
}


type DeviceUnBindGpuForm struct {
	DeviceId string `form:"deviceId" binding:"required"`
}


type TemplateConfigUpdateForm struct {
	UserName string `form:"username" binding:"required"`
	UserPwd  string `form:"userpwd" binding:"required"`
}


type ResetUserPwdForm struct {
	DeviceId    string `form:"deviceId" binding:"required"`
	DeviceName  string `form:"deviceName" binding:"required"` 
	NewPassword string `form:"newPassword" binding:"required"`
}


type CheckVmPwdForm struct {
	VmName   string `form:"vmname" binding:"required"`
	Username string `form:"username" binding:"required"`
	UserPwd  string `form:"userpwd" binding:"required"`
}
