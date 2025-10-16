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

		isCreated := utils.IsVMCreated(value.Name)
		if isCreated && value.Status == "creating" {
			myNewDevice := models.Device{
				Status:     "Off",
				CpuInfo:    utils.GetCpuInfo(value.Name),
				MemoryInfo: utils.GetMemoryInfo(value.Name),
			}
			global.DB.Model(value).Updates(myNewDevice)
		}


		vip := utils.GetVMIp(value.Name)
		if vip != "" {
			status := "running"
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
