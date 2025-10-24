package core

import (
	"gin-vue/api/utils"
	"gin-vue/global"
	"gin-vue/modles/models"
	"time"
)

func DeviceMonitor() {
	ticker := time.NewTicker(20 * time.Second)
	for {
		select {
		case <-ticker.C:
			CheckAllDevice()
		}
	}
}


func CheckAllDevice() {
	var allDevices []*models.Device
	global.DB.Find(&allDevices)
	for _, value := range allDevices {
		// 检查虚拟机是否已创建
		isCreated := utils.IsVMCreated(value.Name)
		
		// 如果虚拟机已创建，更新状态和配置信息
		if isCreated {
			// 如果状态是creating，更新为Off并获取配置信息
			if value.Status == "creating" {
				myNewDevice := models.Device{
					Status:     "Off",
					CpuInfo:    utils.GetCpuInfo(value.Name),
					MemoryInfo: utils.GetMemoryInfo(value.Name),
				}
				global.DB.Model(value).Updates(myNewDevice)
			} else if value.CpuInfo == "" || value.MemoryInfo == "" {
				// 如果CPU或内存信息为空，自动获取并更新这些信息
				cpuInfo := value.CpuInfo
				memoryInfo := value.MemoryInfo
				
				if cpuInfo == "" {
					cpuInfo = utils.GetCpuInfo(value.Name)
				}
				if memoryInfo == "" {
					memoryInfo = utils.GetMemoryInfo(value.Name)
				}
				
				myNewDevice := models.Device{
					CpuInfo:    cpuInfo,
					MemoryInfo: memoryInfo,
				}
				global.DB.Model(value).Updates(myNewDevice)
				global.Logger.Printf("Updated VM %s CPU/Memory info: CPU=%s, Memory=%s\n", value.Name, cpuInfo, memoryInfo)
			}
		}

		// 更新IP和运行状态
		vip := utils.GetVMIp(value.Name)
		if vip != "" {
			status := "Running"
			myNewDevice := models.Device{
				Ip:     vip,
				Status: status,
			}
			global.DB.Model(value).Updates(myNewDevice)
		}
	}

}


func TriggerLicenseDateTask() {
	currentTime := time.Now()

	nextTriggerTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 0, 0, currentTime.Location())
	if currentTime.After(nextTriggerTime) {
		nextTriggerTime = nextTriggerTime.Add(24 * time.Hour)
	}

	durationUntilNextTrigger := nextTriggerTime.Sub(currentTime)
	ticker := time.NewTicker(durationUntilNextTrigger)

	for {
		select {
		case <-ticker.C:
			global.Logger.Printf("TriggerLicenseDateTask add license date:%v", currentTime)
			AddLicenseDate()

			nextTriggerTime = nextTriggerTime.Add(24 * time.Hour)
			durationUntilNextTrigger = time.Until(nextTriggerTime)
			ticker.Reset(durationUntilNextTrigger)
		}
	}
}


func AddLicenseDate() {
	var licenseDate models.LicenseDate
	currentTime := time.Now()
	dateString := currentTime.Format("2006-01-02")
	licenseDate.ID = models.NewUUID()
	licenseDate.DayDate = dateString
	global.DB.Create(&licenseDate)
}

func InitTemplateInfo() {
	templates := global.Config.Vm.GetTemplateMap()
	for key, value := range templates {
		utils.InsertOrUpdateTemplate(key, value)
	}
}
