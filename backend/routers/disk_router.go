package routers

import (
	"gin-vue/api"

	"github.com/gin-gonic/gin"
)

func DisksRouter(router gin.IRouter) {
	diskApi := api.ApiGroupApp.DiskApi





	router.POST("/api/cloud/v1/disks", diskApi.AddDisk)

	router.GET("/api/cloud/v1/disks", diskApi.GetDisks)

	router.DELETE("/api/cloud/v1/disk/:id", diskApi.DeleteDisk)

	router.POST("/api/cloud/v1/disk/device_bind", diskApi.AddDiskToVM)

	router.POST("/api/cloud/v1/disk/device_unbind/:id", diskApi.DiskUnbind)
}
