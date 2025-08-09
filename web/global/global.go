package global

import (
	ut "github.com/go-playground/universal-translator"
	"test/web/config"
	"test/web/proto"
)

var (
	ServerConfig  *config.ServerConfig = &config.ServerConfig{}
	Tarns         ut.Translator
	UserSrvClient proto.UserClient
)
