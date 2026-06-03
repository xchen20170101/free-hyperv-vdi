package routers

import (
	"gin-vue/api"

	"github.com/gin-gonic/gin"
)

// PublicUsersRouter 公开路由（不需要token验证）
func PublicUsersRouter(router *gin.Engine) {
	userApi := api.ApiGroupApp.UserApi

	// 登录接口，不需要token
	router.POST("/api/cloud/v1/login", userApi.UserLogin)
	
	// Android端重置密码接口（可能需要根据实际需求调整是否需要认证）
	router.POST("/api/cloud/v1/reset_password", userApi.AndroidResetPassword)
}

// AuthUsersRouter 需要token验证的路由
func AuthUsersRouter(router gin.IRouter) {
	userApi := api.ApiGroupApp.UserApi

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
