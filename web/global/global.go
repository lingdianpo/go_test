package global

import (
	ut "github.com/go-playground/universal-translator"
	"go_test/web/config"
	"go_test/web/proto"
)

var (
	ServerConfig  *config.ServerConfig = &config.ServerConfig{}
	Tarns         ut.Translator
	UserSrvClient proto.UserClient
)
