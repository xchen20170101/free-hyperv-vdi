package main

import (
	"gin-vue/core"
	"gin-vue/global"
	"gin-vue/routers"
)

func main() {

	core.InitLogger()
	core.InitConf()
	global.Config.Vm.InitTemplate()


	global.DB = core.InitSqliteGorm()

	router := routers.InitRouter()
	addr := global.Config.System.Addr()
	go core.DeviceMonitor()
	core.InitTemplateInfo()



	global.Logger.Printf("gvb server运行在：%s", addr)
	router.Run(addr)
}
