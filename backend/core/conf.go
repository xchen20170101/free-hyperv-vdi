package core

import (
	"fmt"
	"gin-vue/config"
	"gin-vue/global"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"gopkg.in/yaml.v3"
)


func InitConf() {
	var configFile string
	if runtime.GOOS == "windows" {
		configFile = "settings.yaml"
	} else {
		configFile = "/opt/pcfarm_server/settings.yaml"
	}

	c := &config.Config{}
	yamlConf, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(fmt.Errorf("get yamlConf erro:%s", err))
	}
	err = yaml.Unmarshal(yamlConf, c)
	if err != nil {
		global.Logger.Fatalf("cofig Unmarshal err:%s", err)
	}
	global.Logger.Println("config yamlFile load Init success..")
	global.Config = c
}


func InitLogger() {
	file, err := os.OpenFile("cloud_server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("can not open log file:", err)
	}
	global.Logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
}
