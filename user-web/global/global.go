package global

import (
	ut "github.com/go-playground/universal-translator"
	"mxshop-api/user-web/config"
)

var (
	Translator   ut.Translator
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
)
