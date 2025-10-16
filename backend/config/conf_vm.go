package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type VMConfig struct {
	Template string `yaml:"template"`
	Path     string `yaml:"path"`
	DiskPath string `yaml:"disk_path"`
}

var VmMap = map[string]string{}


















func (vmConfig VMConfig) GetTemplates() []string {
	var result []string
	for key := range VmMap {
		result = append(result, key)
	}
	return result
}

func (vmConfig VMConfig) GetTemplateMap() map[string]string {
	return VmMap
}

func (vmConfig VMConfig) DeleteTemplateMap(key string) {
	delete(VmMap, key)
}

func (vmConfig VMConfig) GetTemplateFileName(template string) string {

	return VmMap[template]
}

func (vmConfig VMConfig) InitTemplate() {

	files, err := ioutil.ReadDir(vmConfig.Template)
	if err != nil {
		return
	}
	for _, file := range files {
		if file.IsDir() {
			dir := fmt.Sprintf("%s\\%s\\Virtual Machines", vmConfig.Template, file.Name())
			value := vmConfig.GetTemplateFile(dir)
			VmMap[file.Name()] = value

		}
	}
}

func (vmConfig VMConfig) GetTemplateFile(dir string) string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".vmcx" {
			return file.Name()
		}
	}
	return ""
}
