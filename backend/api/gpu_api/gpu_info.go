package gpu_api

import (
	"gin-vue/api/utils"
	"gin-vue/global"
	"gin-vue/modles/form"
	"gin-vue/modles/models"
	"gin-vue/modles/res"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)


func (GpusApi) GetGpus(c *gin.Context) {
	var responses []*res.GetGpusListResponse
	var gpus []*models.Gpu
	var totalNum int64
	count, _ := strconv.Atoi(c.Query("count"))
	pageNum, _ := strconv.Atoi(c.Query("index"))
	name := c.Query("name")
	keyword := "%" + name + "%"
	offset := (pageNum - 1) * count

	if name == "" {
		global.DB.Model(&models.Gpu{}).Count(&totalNum)
		global.DB.Limit(count).Offset(offset).Find(&gpus)
	} else {
		global.DB.Model(&models.Device{}).Where("name LIKE ?", keyword).Count(&totalNum)
		global.DB.Limit(count).Offset(offset).Where("name LIKE ?", keyword).Find(&gpus)
	}

	data := make(map[string]interface{})
	for _, value := range gpus {
		var bindCount int64
		global.DB.Model(&models.Device{}).Where("gpu_info = ?", value.Name).Count(&bindCount)
		temp := &res.GetGpusListResponse{
			ID:           value.ID,
			Name:         value.Name,
			InstancePath: value.InstancePath,
			BindCount:    int(bindCount),
		}
		responses = append(responses, temp)
	}
	data["gpus"] = responses
	data["totalNum"] = totalNum

	res.OkWithData(data, c)
}


func (GpusApi) AllocationGpu(c *gin.Context) {
	var oldDevices []*models.Device
	var devices []*models.Device
	var gpu models.Gpu
	gpuId := c.PostForm("gpuId")
	global.DB.Where("id = ?", gpuId).First(&gpu)
	if gpu.Name == "" {
		res.FailWithMsg("Gpu.NotExist", c)
		return
	}
	strVm := c.PostForm("strVm")
	vms := strings.Split(strVm, ",")

	if len(vms) > 4 {
		res.FailWithMsg("Gpu.IsLimited", c)
		return
	}

	target := utils.GetInstancePathFromDB(gpu.Desc)
	global.Logger.Println("target:", target)

	global.DB.Where("gpu_info = ?", gpuId).Find(&oldDevices)
	for _, oldValue := range oldDevices {
		go utils.RemoveGpuAdpater(oldValue.Name, target)
	}

	global.DB.Model(&models.Device{}).Where("gpu_info = ?", gpuId).Update("gpu_info", "")

	global.DB.Where("name in ?", vms).Find(&devices)

	for _, value := range devices {

		if value.GpuInfo == gpu.ID {
			continue
		}
		myDevice := models.Device{
			GpuInfo: gpu.ID,
		}
		global.DB.Model(value).Updates(myDevice)
		go utils.AddGpuToVm(value.Name, target)

	}
	res.OkWithData(gpu, c)

}


func (GpusApi) BindGpu(c *gin.Context) {
	var form form.DeviceBindGpuForm
	var devices []*models.Device
	var device models.Device
	var gpu models.Gpu
	if err := c.ShouldBind(&form); err != nil {
		global.Logger.Printf("BindGpu failed:%v\n", err.Error())
		res.FailWithMsg("Common.InvalidParam", c)
		return
	}
	gpuId := form.GpuId
	deviceName := form.DeviceName
	global.DB.Where("name = ?", deviceName).First(&device)
	if device.GpuInfo != "" {
		res.FailWithMsg("Gpu.BindingExist", c)
		return
	}
	global.DB.Where("id = ?", gpuId).First(&gpu)
	if gpu.Name == "" {
		res.FailWithMsg("Gpu.NotExist", c)
		return
	}
	global.DB.Where("gpu_info = ?", gpu.Name).Find(&devices)

	if len(devices) >= 8 {
		res.FailWithMsg("Gpu.IsLimited", c)
		return
	}
	target := utils.GetInstancePathFromDB(gpu.Desc)
	err := utils.AddGpuToVm(deviceName, target)
	if err != nil {
		go utils.RemoveGpuAdpater(deviceName, target)
		res.FailWithMsg("Gpu.BindFailed", c)
		return
	}
	global.DB.Model(&models.Device{}).Where("name = ?", deviceName).Update("gpu_info", gpu.Name)
	res.OkWithData(gpu, c)
}


func (GpusApi) UnBindGpu(c *gin.Context) {
	var form form.DeviceUnBindGpuForm
	var device models.Device
	var gpu models.Gpu
	if err := c.ShouldBind(&form); err != nil {
		global.Logger.Printf("UnBindGpu failed:%v\n", err.Error())
		res.FailWithMsg("Common.InvalidParam", c)
		return
	}
	deviceId := form.DeviceId
	global.DB.Where("id = ?", deviceId).First(&device)
	if device.Name == "" {
		res.FailWithMsg("Device.NotExist", c)
		return
	}
	global.DB.Where("name = ?", device.GpuInfo).First(&gpu)
	if gpu.InstancePath == "" {
		res.FailWithMsg("Gpu.NotExist", c)
		return
	}
	target := utils.GetInstancePathFromDB(gpu.Desc)
	global.DB.Model(&models.Device{}).Where("id = ?", deviceId).Update("gpu_info", "")
	go utils.RemoveGpuAdpater(device.Name, target)
	res.OkWithData(gpu, c)

}
