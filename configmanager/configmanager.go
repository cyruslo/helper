package configmanager

import (
	"errors"
	"fmt"
)

type IConfig interface {
	CheckValid() error
}

type configType int
const (
	CONFIG_TYPE_UNKNOWN configType = iota
	CONFIG_TYPE_JSON
	CONFIG_TYPE_LUA
	CONFIG_TYPE_TOML
)

type innerLogger func(format string, params ...interface{})


var (
	loaderList []*fileLoader
	logger innerLogger
)


func init()  {
	logger = func(format string, params ...interface{}) {
		fmt.Printf(format, params...)
	}
}

//RegConfig 指定文件名和配置类型
//注意! 所有配置单例指针必须先实例化
func RegConfig(fileName string, config IConfig, cType configType)  {
	loader := NewLoader(fileName, config, cType)
	loaderList = append(loaderList, loader)
}

//RegAutoConfig 自动解析文件名和配置类型
//对应配置文件名和配置实例类型名去掉下划线并转成小写后必须一致
//如game_config.json或gameConfig.json都可以对应GameConfig
//注意! 所有配置单例指针必须先实例化
func RegAutoConfig(config IConfig)  {
	loader := NewAutoLoader(config)
	loaderList = append(loaderList, loader)
}

func Init() (err error) {
	for _, v := range loaderList {
		err = v.Init()
		if err != nil {
			return
		}
	}
	return
}

func ReloadAll() {
	for _, v := range loaderList {
		err := v.Reload()
		if err != nil {
			logger("reload conf:%s err:%s", v.fileName, err.Error())
			return
		}
	}
	return
}

func SetLogger(l innerLogger)  {
	logger = l
}

func makeError(format string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(format, args...))
}