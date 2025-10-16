package routers

import (
	"gin-vue/api"

	"github.com/gin-gonic/gin"
)

func UsersRouter(router *gin.Engine) {
	userApi := api.ApiGroupApp.UserApi

	router.POST("/api/cloud/v1/login", userApi.UserLogin)

	router.POST("/api/cloud/v1/reset_password", userApi.AndroidResetPassword)

	router.DELETE("/api/cloud/v1/logout", userApi.UserLogout)

	router.POST("/api/cloud/v1/users", userApi.UserAdd)

	router.GET("/api/cloud/v1/users", userApi.UserGet)

	router.DELETE("/api/cloud/v1/user/:id", userApi.UserDel)

	router.PUT("/api/cloud/v1/user/:id", userApi.UserUpdate)

	router.GET("/api/cloud/v1/user_count", userApi.UserAllCountGet)

	router.POST("/api/cloud/v1/user_bind", userApi.UserBindDevices)

	router.GET("/api/cloud/v1/user_profile", userApi.UserProfileGet)

	router.PUT("/api/cloud/v1/user", userApi.UserUpdateSelfPassword)

	router.GET("/api/cloud/v1/licenses", userApi.UserLicenseGet)

	router.POST("/api/cloud/v1/licenses", userApi.LicenseActive)

	router.GET("/api/cloud/v1/licenses_all", userApi.GetLicenses)
}
