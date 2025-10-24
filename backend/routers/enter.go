package routers

import (
	"gin-vue/global"
	"gin-vue/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(global.Config.System.Env)
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		c.Next()
	})
	
	router.Use(Cors())
	
	// 公开路由（不需要token验证）
	PublicUsersRouter(router)
	
	// 需要token验证的路由组
	authGroup := router.Group("")
	authGroup.Use(middleware.TokenAuthMiddleware())
	{
		AuthUsersRouter(authGroup)
		DevicesRouter(authGroup)
		DisksRouter(authGroup)
		GpusRouter(authGroup)
	}
	
	return router
}

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
	}
}