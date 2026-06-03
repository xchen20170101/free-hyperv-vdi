package disk_api

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

func (DisksApi) AddDisk(c *gin.Context) {

	var form form.AddDiskForm
	if err := c.ShouldBind(&form); err != nil {
		global.Logger.Printf("add disk failed:%v\n", err.Error())
		res.FailWithMsg("Common.InvalidParam", c)
		return
	}
	name := form.Name
	capacity := form.Capacity

	iCapacity, err := strconv.Atoi(capacity)
	if err != nil {
		res.FailWithMsg("Common.InvalidParam", c)
		return
	}

	if iCapacity < 1 || iCapacity > 4*1024 {
		res.FailWithMsg("Disk.OutofLimit", c)
		return
	}

	var disk models.Disk
	global.DB.Where("name = ?", name).First(&disk)
	if disk.Name != "" {
		res.FailWithMsg("Disk.Exists", c)
		return
	}

	storagePath := global.Config.Vm.DiskPath + "\\" + name + "_disk_extra"


	err = utils.CreateDisk(name, storagePath, iCapacity)
	if err != nil {
		res.FailWithMsg("Disk.CreateDiskFailed", c)
		return
	}

	myDisk := &models.Disk{
		ID:          models.NewUUID(),
		Name:        name,
		Capacity:    capacity,
		StoragePath: storagePath,
		CreatedTime: time.Now().UnixNano() / int64(time.Millisecond),
	}
	global.DB.Create(&myDisk)
	res.OkWithData(myDisk, c)
}

func (DisksApi) GetDisks(c *gin.Context) {
	var disks []*res.GetDisksListResponse
	var totalNum int64
	var allDisks []*models.Disk

	count, _ := strconv.Atoi(c.Query("count"))
	pageNum, _ := strconv.Atoi(c.Query("index"))
	name := c.Query("name")
	keyword := "%" + name + ""
	offset := (pageNum - 1) * count

	if name == "" {
		global.DB.Model(&models.Disk{}).Count(&totalNum)
		global.DB.Limit(count).Offset(offset).Find(&allDisks)
	} else {
		global.DB.Model(&models.Disk{}).Where("name LIKE ?", keyword).Count(&totalNum)
		global.DB.Limit(count).Offset(offset).Where("name LIKE ?", keyword).Find(&allDisks)
	}

	data := make(map[string]interface{})
	for _, value := range allDisks {
		temp := &res.GetDisksListResponse{
			ID:           value.ID,
			Name:         value.Name,
			Capacity:     value.Capacity + "GB",
			StoragePath:  value.StoragePath,
			BelongDevice: utils.GetDeviceNameByDiskId(value.ID),
			CreatedTime:  utils.TransferTimeStamp(value.CreatedTime),
		}
		disks = append(disks, temp)
	}
	data["disks"] = disks
	data["totalNum"] = totalNum
	res.OkWithData(data, c)
}

func (DisksApi) DeleteDisk(c *gin.Context) {

	var disk models.Disk
	diskId := c.Param("id")
	result := global.DB.Where("id = ?", diskId).First(&disk)
	if result.Error != nil {
		res.FailWithMsg("Disk.NotExist", c)
		return
	}
	var diskBind models.DiskBind
	global.DB.Where("disk_id = ?", diskId).First(&diskBind)
	if diskBind.DeviceId != "" {

		res.FailWithMsg("Disk.HasBind", c)
		return
	}


	err := os.RemoveAll(disk.StoragePath)
	if err != nil {
		res.FailWithMsg("Disk.RemoveDiskFileFailed", c)
		return
	}

	global.DB.Where("ID = ?", diskId).Delete(&models.Disk{})
	res.OkWithData(diskId, c)
}

func (DisksApi) AddDiskToVM(c *gin.Context) {

	var form form.BindDiskForm
	if err := c.ShouldBind(&form); err != nil {
		global.Logger.Printf("add disk to vm failed:%v\n", err.Error())
		res.FailWithMsg("Common.InvalidParam", c)
		return
	}
	diskId := form.DiskId
	deviceName := form.DeviceName

	var disk models.Disk
	global.DB.Where("id = ?", diskId).First(&disk)
	if disk.Name == "" {
		res.FailWithMsg("Disk.NotExist", c)
		return
	}
	var device models.Device
	global.DB.Where("name = ?", deviceName).First(&device)
	if device.ID == "" {
		res.FailWithMsg("Disk.DeviceNotExist", c)
		return
	}
	if device.Status == "Running" {
		res.FailWithMsg("Disk.DeviceIsRunning", c)
		return
	}
	var diskBind models.DiskBind
	global.DB.Where("disk_id = ?", diskId).First(&diskBind)
	if diskBind.DeviceId != "" {
		res.FailWithMsg("Disk.BindExist", c)
		return
	}


	err := utils.BindDiskToVM(deviceName, disk.Name, disk.StoragePath)
	if err != nil {
		res.FailWithMsg("Disk.BindToVMFailed", c)
		return
	}

	myDiskBind := &models.DiskBind{
		ID:       models.NewUUID(),
		DeviceId: device.ID,
		DiskId:   diskId,
	}
	global.DB.Create(&myDiskBind)
	res.OkWithData(myDiskBind, c)
}

func (DisksApi) DiskUnbind(c *gin.Context) {

	var disk models.Disk
	diskId := c.Param("id")

	global.DB.Where("id = ?", diskId).First(&disk)
	if disk.Name == "" {
		res.FailWithMsg("Disk.NotExist", c)
		return
	}
	var diskBind models.DiskBind
	global.DB.Where("disk_id = ?", diskId).First(&diskBind)
	if diskBind.ID == "" {
		res.FailWithMsg("Disk.BindNotExist", c)
		return
	}
	var device models.Device
	global.DB.Where("id = ?", diskBind.DeviceId).First(&device)
	if device.ID == "" {
		res.FailWithMsg("Disk.DeviceNotExist", c)
		return
	}
	if device.Status == "Running" {
		res.FailWithMsg("Disk.DeviceIsRunning", c)
		return
	}

	err := utils.UnBindDiskFromVM(device.Name, disk.Name, disk.StoragePath)
	if err != nil {
		res.FailWithMsg("Disk.UnbindFailed", c)
		return
	}

	global.DB.Where("disk_id = ?", diskId).Delete(&models.DiskBind{})
	res.OkWithData(diskId, c)
}
