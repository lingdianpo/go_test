package initialize

import (
	"github.com/gin-gonic/gin"
	"test/web/router"
)

func RouterInit() *gin.Engine {
	Router := gin.Default()
	ApiGroup := Router.Group("v1")
	router.InitUserRouter(ApiGroup)
	return Router
}
