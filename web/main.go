package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"test/web/global"
	"test/web/initialize"
	"test/web/utils"
	myvalidator "test/web/validator"
)

func main() {
	//1. 初始化logger
	initialize.LoggerInit()
	//2. 初始化config
	initialize.ConfigInit()
	//3. 初始化routers
	Router := initialize.RouterInit()
	//4. 初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}
	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Tarns, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} 非法的手机号码！", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	debug := global.ServerConfig.DebugInfo.Status
	zap.S().Info(debug)
	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}
	/*
		1. S()可以获取一个全局的sugar，可以让我们自己设置一个全局的logger
		2. 日志是分几倍的，debug，info，warn，error，fetal
		3.S函数和L函数很有用,提供了一个全局的安全访问logger的途径
	*/
	zap.S().Info("启动服务器，端口 ", global.ServerConfig.Port)
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Error("启动失败", zap.Any("err", err))
	}

}
