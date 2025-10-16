package routers

import (
	"gin-vue/api"

	"github.com/gin-gonic/gin"
)

func GpusRouter(router *gin.Engine) {
	gpuApi := api.ApiGroupApp.GpuApi

	router.GET("/api/cloud/v1/gpus", gpuApi.GetGpus)
	router.POST("/api/cloud/v1/gpus", gpuApi.AllocationGpu)
	router.POST("/api/cloud/v1/bind_gpu", gpuApi.BindGpu)
	router.POST("/api/cloud/v1/unbind_gpu", gpuApi.UnBindGpu)
}
