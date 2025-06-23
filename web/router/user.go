package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"test/web/api"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user")
	zap.S().Info("配置用户相关url")
	{
		UserRouter.GET("list", api.GetUserList)
	}

}
