package devices_api

import (
	"gin-vue/api/utils"
	"gin-vue/global"
	"gin-vue/modles/form"
	"gin-vue/modles/models"
	"gin-vue/modles/res"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)


func (DevicesApi) GetDevicesList(c *gin.Context) {
	var devices []*res.GetDevicesListResponse
	var totalNum int64
	count, _ := strconv.Atoi(c.Query("count"))
	pageNum, _ := strconv.Atoi(c.Query("index"))
	name := c.Query("name")
	var allDevices []*models.Device
	keyword := "%" + name + "%"
	offset := (pageNum - 1) * count
	if name == "" {
		global.DB.Model(&models.Device{}).Count(&totalNum)
		global.DB.Limit(count).Offset(offset).Find(&allDevices)
	} else {
		global.DB.Model(&models.Device{}).Where("name LIKE ?", keyword).Count(&totalNum)
		global.DB.Limit(count).Offset(offset).Where("name LIKE ?", keyword).Find(&allDevices)
	}

	vmNames := make([]string, len(allDevices))
	for i, device := range allDevices {
		vmNames[i] = device.Name
	}
	vmStates := utils.GetAllVMStates(vmNames)
	vmIps := utils.GetAllVMIps(vmNames)

	data := make(map[string]interface{})
	for _, value := range allDevices {

		realTimeStatus := vmStates[value.Name]
		
		if realTimeStatus != "Unknown" && realTimeStatus != value.Status {
			myNewDevice := models.Device{
				Status: realTimeStatus,
			}
			global.DB.Model(value).Updates(myNewDevice)
			global.Logger.Printf("Updated VM %s status from %s to %s\n", value.Name, value.Status, realTimeStatus)
		}


		currentStatus := realTimeStatus
		if currentStatus == "Unknown" {
			currentStatus = value.Status
		}


		realTimeIp := vmIps[value.Name]

		currentIp := realTimeIp

		temp := &res.GetDevicesListResponse{
			ID:           value.ID,
			Name:         value.Name,
			UserName:     utils.GetUserInfoByDevice(value.ID),
			TemplateInfo: utils.GetTemplateNameById(value.TemplateId),
			VirtualIp:    currentIp,
			Status:       currentStatus,
			CreatedTime:  utils.TransferTimeStamp(value.CreatedTime),
			MemoryInfo:   value.MemoryInfo,
			CpuInfo:      value.CpuInfo,
			GpuInfo:      value.GpuInfo,
		}
		devices = append(devices, temp)
	}

	data["devices"] = devices
	data["totalNum"] = totalNum

	res.OkWithData(data, c)
}


func (DevicesApi) GetUnbindDevicesList(c *gin.Context) {
	var respDevices []*res.GetDevicesListResponse
	var totalNum int64
	var devices []*models.Device

	bindDeviceIds := utils.GetBindDeviceIds()
	global.DB.Find(&devices)
	data := make(map[string]interface{})
	for _, value := range devices {
		isBind := utils.IsInArray(value.ID, bindDeviceIds)
		if isBind {
			continue
		}
		if value.Status != "running" {
			continue
		}
		temp := &res.GetDevicesListResponse{
			ID:          value.ID,
			Name:        value.Name,
			UserName:    utils.GetUserInfoByDevice(value.ID),
			VirtualIp:   "",
			Status:      value.Status,
			CreatedTime: utils.TransferTimeStamp(value.CreatedTime),
		}
		respDevices = append(respDevices, temp)
	}
	data["devices"] = respDevices
	data["totalNum"] = totalNum

	res.OkWithData(data, c)
}


func (DevicesApi) UpdateVM(c *gin.Context) {
	var device models.Device



	id := c.Param("id")


	result := global.DB.Where("id = ?", id).First(&device)
	if result.Error != nil {
		res.FailWithMsg("Device.NotExist", c)
		return
	}
	if device.Name == "" {
		res.FailWithMsg("Device.NotExist", c)
		return
	}
	if device.Status == "running" {
		res.FailWithMsg("Device.IsRunning", c)
		return
	}
	var form form.UpdateVmForm
	if err := c.ShouldBind(&form); err != nil {
		global.Logger.Printf("update vm failed:%v\n", err.Error())
		res.FailWithMsg("Common.InvalidParam", c)
		return
	}


	flag := false


	cpuInfo := form.CpuInfo
	memoryInfo := form.MemoryInfo

	if cpuInfo != device.CpuInfo {
		err := utils.SetCpuInfo(device.Name, cpuInfo)
		if err != nil {
			res.FailWithMsg("Device.UpdateFailed", c)
			return
		}
		flag = true
	}

	if memoryInfo != device.MemoryInfo {
		err := utils.SetMemoryInfo(device.Name, memoryInfo)
		if err != nil {
			res.FailWithMsg("Device.UpdateFailed", c)
			return
		}
		flag = true
	}

	if flag {
		myNewDevice := models.Device{
			CpuInfo:    cpuInfo,
			MemoryInfo: memoryInfo,
		}
		global.DB.Model(device).Updates(myNewDevice)
	}

	res.OkWithData(device, c)
}


func (DevicesApi) DeleteVM(c *gin.Context) {

	var device models.Device
	var bind models.Bind
	id := c.Param("id")
	result := global.DB.Where("id = ?", id).First(&device)
	if result.Error != nil {
		res.FailWithMsg("Device.NotExist", c)
		return
	}
	if device.Status == "running" {
		res.FailWithMsg("Device.IsRunning", c)
		return
	}
	global.DB.Where("device_id = ?", id).First(&bind)
	go utils.DeleteVMInForce(device.Name)

	result = global.DB.Delete(&device)
	if result.Error != nil {
		res.FailWithMsg("Device.DeleteVMFailed", c)
		return
	}
	if bind.UserId != "" {
		result = global.DB.Delete(&bind)
		if result.Error != nil {
			res.FailWithMsg("Device.DeleteBindFailed", c)
			return
		}
	}
	result = global.DB.Where("device_id = ?", id).Delete(&models.DiskBind{})
	if result.Error != nil {
		res.FailWithMsg("Device.DeleteDiskBindFailed", c)
		return
	}
	res.OkWithData(device, c)
}


func (DevicesApi) AddVM(c *gin.Context) {

	var device models.Device
	var form form.AddVmForm
	if err := c.ShouldBind(&form); err != nil {
		global.Logger.Printf("add vm failed:%v\n", err.Error())
		res.FailWithMsg("Common.InvalidParam", c)
		return
	}

	name := form.Name
	
	if !utils.ValidateVMName(name) {
		global.Logger.Printf("invalid vm name:%s\n", name)
		res.FailWithMsg("Device.InvalidName", c)
		return
	}
	
	vmSwitch := form.VmSwitch
	template := form.SrcVmPath
	templateName := global.Config.Vm.GetTemplateFileName(template)
	global.Logger.Printf("add vm template:%s, templateName:%s\n", template, templateName)
	_ = global.DB.Where("name = ?", name).First(&device)
	if device.Name != "" {
		res.FailWithMsg("Device.Exist", c)
		return
	}
	srcVmPath := global.Config.Vm.Template + "\\" + template + "\\Virtual Machines\\" + templateName
	_, errFile := os.Stat(srcVmPath)
	if errFile != nil {
		res.FailWithMsg("Device.TemplateNotExist", c)
		return
	}
	targetVmPath := global.Config.Vm.Path + "\\" + name
	targetDiskPath := global.Config.Vm.DiskPath + "\\" + name + "_disk"
	global.Logger.Printf("srcVmPath:%s, targetVmPath:%s, targetDiskPath:%s\n", srcVmPath, targetVmPath, targetDiskPath)


	go utils.ImportVM(name, template, srcVmPath, targetVmPath, targetDiskPath, vmSwitch)

	myDevice := &models.Device{
		ID:          models.NewUUID(),
		Name:        name,
		TemplateId:  utils.GetTemplateIdByName(template),
		Ip:          "",
		Status:      "creating",
		CreatedTime: time.Now().UnixNano() / int64(time.Millisecond),
	}
	global.DB.Create(&myDevice)

	res.OkWithData(myDevice, c)
}


func (DevicesApi) DeviceAllCountGet(c *gin.Context) {
	var countNum int64
	global.DB.Model(&models.Device{}).Count(&countNum)
	data := make(map[string]interface{})
	data["num"] = countNum
	res.OkWithData(data, c)

}


func (DevicesApi) DeviceTemplatesGet(c *gin.Context) {
	data := make(map[string]interface{})
	data["templates"] = global.Config.Vm.GetTemplates()
	res.OkWithData(data, c)
}


func (DevicesApi) DeviceTemplateConfigGet(c *gin.Context) {
	var responses []*res.GetTemplatesListResponse
	var templates []*models.Template
	var totalNum int64
	count, _ := strconv.Atoi(c.Query("count"))
	pageNum, _ := strconv.Atoi(c.Query("index"))
	name := c.Query("name")
	keyword := "%" + name + "%"
	offset := (pageNum - 1) * count

	if name == "" {
		global.DB.Model(&models.Template{}).Count(&totalNum)
		global.DB.Limit(count).Offset(offset).Find(&templates)
	} else {
		global.DB.Model(&models.Device{}).Where("name LIKE ?", keyword).Count(&totalNum)
		global.DB.Limit(count).Offset(offset).Where("name LIKE ?", keyword).Find(&templates)
	}

	data := make(map[string]interface{})
	for _, value := range templates {
		temp := &res.GetTemplatesListResponse{
			ID:       value.ID,
			Name:     value.Name,
			UserName: value.UserName,
			UserPwd:  value.UserPwd,
		}
		responses = append(responses, temp)
	}
	data["templates"] = responses
	data["totalNum"] = totalNum
	res.OkWithData(data, c)
}


func (DevicesApi) TemplateConfigUpdate(c *gin.Context) {
	var template models.Template
	id := c.Param("id")

	global.DB.Where("id = ?", id).First(&template)
	if template.Name == "" {
		res.FailWithMsg("Template.NotExist", c)
		return
	}
	var form form.TemplateConfigUpdateForm
	if err := c.ShouldBind(&form); err != nil {
		global.Logger.Printf("update template failed:%v\n", err.Error())
		res.FailWithMsg("Common.InvalidParam", c)
		return
	}
	myNewTemplate := models.Template{
		UserName: form.UserName,
		UserPwd:  form.UserPwd,
	}
	global.DB.Model(template).Updates(myNewTemplate)
	res.OkWithData(template.ID, c)
}


func (DevicesApi) TemplateDelete(c *gin.Context) {
	var template models.Template
	id := c.Param("id")

	global.DB.Where("id = ?", id).First(&template)
	if template.Name == "" {
		res.FailWithMsg("Template.NotExist", c)
		return
	}
	go utils.DeleteTemplateFile(template.Name)
	global.DB.Where("id = ?", id).Delete(&models.Template{})

}


func (DevicesApi) DeviceSwitchsGet(c *gin.Context) {
	data := make(map[string]interface{})
	data["switchs"] = utils.GetSwitchs()
	res.OkWithData(data, c)
}


func (DevicesApi) DeviceUnBindUser(c *gin.Context) {
	var form form.DeviceUnbindUserForm
	if err := c.ShouldBind(&form); err != nil {
		global.Logger.Printf("device unbind user failed:%v\n", err.Error())
		res.FailWithMsg("Common.InvalidParam", c)
		return
	}
	var bind models.Bind
	var device models.Device
	var user models.User
	var template models.Template


	global.DB.Where("device_id =?", form.DeviceId).First(&bind)
	if bind.UserId == "" {
		res.FailWithMsg("Device.UnBindUserFailed", c)
		return
	}
	global.DB.Where("id =?", bind.UserId).First(&user)
	if user.Name == "" {
		res.FailWithMsg("Device.UnBindUserFailed", c)
		return
	}
	global.DB.Where("id = ?", bind.DeviceId).First(&device)
	if device.Name == "" {
		res.FailWithMsg("Device.UnBindUserFailed", c)
		return
	}
	if device.Status != "running" {
		res.FailWithMsg("Device.MustBeRunning", c)
		return
	}
	global.DB.Where("id =?", device.TemplateId).First(&template)
	if template.Name == "" {
		res.FailWithMsg("Device.TemplateNotExist", c)
		return
	}
	err := utils.DeleteVMLocalUser(user.Name, device.Name, template.UserName, template.UserPwd)
	if err != nil {
		res.FailWithMsg("Device.UnBindUserFailed", c)
		return
	}
	global.DB.Where("device_id = ?", form.DeviceId).Delete(&models.Bind{})
	res.OkWithData(form.DeviceId, c)
}


func (DevicesApi) OperateVM(c *gin.Context) {
	var form form.OperateVmForm
	if err := c.ShouldBind(&form); err != nil {
		global.Logger.Printf("operate vm failed:%v\n", err.Error())
		res.FailWithMsg("Common.InvalidParam", c)
		return
	}

	vmId := form.VmId
	action := form.Action

	var device models.Device
	result := global.DB.Where("id = ?", vmId).First(&device)
	if result.Error != nil {
		res.FailWithMsg("Device.NotExist", c)
		return
	}
	if device.Name == "" {
		res.FailWithMsg("Device.NotExist", c)
		return
	}




	if action == "1" {
		if device.Status == "running" {
			res.FailWithMsg("Device.StatusInvalid", c)
			return
		}
		device.Status = "booting"
		global.DB.Save(&device)
		go utils.OpenVM(device.Name, device.Status)
	} else if action == "2" {
		if device.Status == "Off" {
			res.FailWithMsg("Device.StatusInvalid", c)
			return
		}
		status := utils.CloseVM(device.Name, device.Status)
		device.Status = status
		device.Ip = ""
		global.DB.Save(&device)
	} else if action == "3" {
		if device.Status == "Off" {
			res.FailWithMsg("Device.StatusInvalid", c)
			return
		}
		_ = utils.CloseVM(device.Name, device.Status)
		status := utils.OpenVM(device.Name, device.Status)
		device.Status = status
		global.DB.Save(&device)
	}

	res.OkWithData(device, c)
}


func (DevicesApi) ResetUserPwd(c *gin.Context) {
	var form form.ResetUserPwdForm
	if err := c.ShouldBind(&form); err != nil {
		global.Logger.Printf("reset user password failed:%v\n", err.Error())
		res.FailWithMsg("Common.InvalidParam", c)
		return
	}
	
	var device models.Device
	var template models.Template

	global.DB.Where("id = ?", form.DeviceId).First(&device)
	if device.Name == "" {
		res.FailWithMsg("Device.NotExist", c)
		return
	}
	

	if device.Status != "running" {
		res.FailWithMsg("Device.MustBeRunning", c)
		return
	}
	

	userName := utils.GetUserInfoByDevice(form.DeviceId)
	if userName == "" {
		res.FailWithMsg("Device.UserNotBind", c)
		return
	}
	

	global.DB.Where("id = ?", device.TemplateId).First(&template)
	if template.Name == "" {
		res.FailWithMsg("Device.TemplateNotExist", c)
		return
	}
	

	err := utils.ResetVMUserPwd(userName, form.NewPassword, device.Name, template.UserName, template.UserPwd)
	if err != nil {
		global.Logger.Printf("reset user password failed:%v\n", err.Error())
		res.FailWithMsg("Device.ResetUserPwdFailed", c)
		return
	}
	
	res.OkWithData(form.DeviceId, c)
}


func (DevicesApi) CheckVmPwd(c *gin.Context) {
	var form form.CheckVmPwdForm
	if err := c.ShouldBind(&form); err != nil {
		global.Logger.Printf("check vm password failed:%v\n", err.Error())
		res.FailWithMsg("Common.InvalidParam", c)
		return
	}
	

	isValid, err := utils.ValidVMUserPwd(form.VmName, form.Username, form.UserPwd)
	if err != nil {
		global.Logger.Printf("check vm password error:%v\n", err.Error())
		res.FailWithMsg("Device.CheckPasswordFailed", c)
		return
	}
	

	data := make(map[string]interface{})
	data["valid"] = isValid
	res.OkWithData(data, c)
}
