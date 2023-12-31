package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/user-web/api"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	zap.S().Infof("配置用户相关url")
	{
		UserRouter.GET("list", api.GetUserList)
		UserRouter.POST("pwd_login", api.PassWordLogin)
	}
}
