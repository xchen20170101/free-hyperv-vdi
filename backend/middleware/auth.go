package middleware

import (
	"gin-vue/global"
	"gin-vue/modles/models"
	"gin-vue/modles/res"
	"time"

	"github.com/gin-gonic/gin"
)

// TokenAuthMiddleware Token验证中间件
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取token（优先从Cookie，其次从Header）
		token, err := c.Cookie("accessToken")
		if err != nil || token == "" {
			// 从Header中获取（支持 Bearer token 格式）
			authHeader := c.GetHeader("Authorization")
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				token = authHeader[7:]
			} else {
				token = authHeader
			}
		}

		// 如果没有token，返回未授权
		if token == "" {
			res.FailWithMsg("Token.NotExist", c)
			c.Abort()
			return
		}

		// 查询数据库中的token
		var dbToken models.Token
		result := global.DB.Where("value = ?", token).First(&dbToken)
		if result.Error != nil {
			// token不存在或已被其他登录覆盖
			res.FailWithMsg("Token.Invalid", c)
			c.Abort()
			return
		}

	// 验证token是否过期
	now := time.Now().Unix()
	if now-dbToken.CreatedTime > dbToken.Cryptoperiod {
		// token已过期，删除数据库记录
		global.DB.Delete(&dbToken)
		res.FailWithMsg("Token.Expired", c)
		c.Abort()
		return
	}

	// 滑动过期时间：更新token的创建时间，使其保持活跃
	// 只要用户在使用系统，token就不会过期
	dbToken.CreatedTime = now
	global.DB.Save(&dbToken)

	// 计算剩余有效期（用于Cookie的MaxAge）
	remainingTime := int(dbToken.Cryptoperiod)

	// token有效，将userId存入上下文供后续使用
	c.Set("userId", dbToken.UserId)
	// 同时更新Cookie，保持会话（使用完整的有效期）
	c.SetCookie("accessToken", dbToken.Value, remainingTime, "/", "", false, true)
	c.SetCookie("userId", dbToken.UserId, remainingTime, "/", "", false, true)
	
	c.Next()
	}
}

