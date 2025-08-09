package router

import (
	"github.com/gin-gonic/gin"
	"test/web/api"
)

func InitBaseRouter(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("base")
	{
		BaseRouter.GET("captcha", api.GetCaptcha)
		BaseRouter.GET("send_sms", api.GetCaptcha)
	}
}
