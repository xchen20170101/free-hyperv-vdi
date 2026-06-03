package routers

import (
	"gin-vue/api"

	"github.com/gin-gonic/gin"
)

func DevicesRouter(router gin.IRouter) {
	deviceApi := api.ApiGroupApp.DevicesApi

	router.GET("/api/cloud/v1/devices", deviceApi.GetDevicesList)

	router.GET("/api/cloud/v1/unbind_devices", deviceApi.GetUnbindDevicesList)

	router.POST("/api/cloud/v1/vm", deviceApi.AddVM)

	router.PUT("/api/cloud/v1/vm/:id", deviceApi.UpdateVM)

	router.DELETE("/api/cloud/v1/vm/:id", deviceApi.DeleteVM)

	router.GET("/api/cloud/v1/device_count", deviceApi.DeviceAllCountGet)

	router.POST("/api/cloud/v1/vm/operate", deviceApi.OperateVM)

	router.GET("/api/cloud/v1/device_templates", deviceApi.DeviceTemplatesGet)

	router.GET("/api/cloud/v1/device_switchs", deviceApi.DeviceSwitchsGet)

	router.POST("/api/cloud/v1/unbind_user", deviceApi.DeviceUnBindUser)

	router.GET("/api/cloud/v1/templates", deviceApi.DeviceTemplateConfigGet)

	router.PUT("/api/cloud/v1/templates/:id", deviceApi.TemplateConfigUpdate)

	router.DELETE("/api/cloud/v1/template/:id", deviceApi.TemplateDelete)

	router.POST("/api/cloud/v1/reset_user_pwd", deviceApi.ResetUserPwd)

	router.POST("/api/cloud/v1/check_vm_pwd", deviceApi.CheckVmPwd)
}
