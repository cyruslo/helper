package configmanager

import (
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/cyruslo/helper/luabridge"
	"reflect"
	"regexp"
	"strings"
)

type fileLoader struct {
	fileName string
	config   IConfig
	cType    configType
}

func NewAutoLoader(configPtr IConfig) (ret *fileLoader) {
	fileName, cType := judgeConfigNameAndType(configPtr)
	ret = &fileLoader{
		fileName: fileName,
		config:   configPtr,
		cType:    cType,
	}
	return
}

func NewLoader(fileName string, configPtr IConfig, cType configType) (ret *fileLoader) {
	ret = &fileLoader{
		fileName: fileName,
		config:   configPtr,
		cType:    cType,
	}
	return
}

func (self *fileLoader)Init() (err error) {
	if self.cType == CONFIG_TYPE_UNKNOWN {
		err = makeError("config type name:%s, file name:%s, unknown config type", getConfigRawTypeName(self.config), self.fileName)
		return
	}
	if self.fileName == "" {
		err = makeError("config name:%s file name empty", getConfigRawTypeName(self.config))
		return
	}
	err = paladin.Watch(self.fileName, self)
	return
}

//Reload reload config actively
func (self *fileLoader)Reload() (err error) {
	file := paladin.Get(self.fileName)
	if file == nil {
		err = makeError("reload file nil,file name:%s", self.fileName)
		return
	}
	var fileStr string
	fileStr, err = file.String()
	if err != nil {
		return
	}
	err = self.Set(fileStr)
	return
}

func (self *fileLoader) Set(fileContent string) (err error) {
	confType := reflect.TypeOf(self.config)
	if confType.Kind() == reflect.Ptr {
		confType = confType.Elem()
	}
	//reflect.New 返回指针
	config := reflect.New(confType).Interface()
	switch self.cType {
	case CONFIG_TYPE_JSON:
		err = jsonParser(fileContent, config)
	case CONFIG_TYPE_LUA:
		err = luaParser(fileContent, config)
	case CONFIG_TYPE_TOML:
		err = tomlParser(fileContent, config)
	default:
		err = makeError("file loader type err:%d", self.cType)
	}
	if err != nil {
		logger("set config file name:%s, err:%s", self.fileName, err.Error())
		return
	}
	if err = config.(IConfig).CheckValid();err != nil {
		logger("check valid config file name:%s, err:%s", self.fileName, err.Error())
		return
	}
	//Value.Set只能对assignable的值使用,具体参照Elem()和Set()的文档
	reflect.ValueOf(self.config).Elem().Set(reflect.ValueOf(config).Elem())
	return
}

func judgeConfigNameAndType(configPtr IConfig) (fileName string, cType configType) {
	typeName := getConfigTypeName(configPtr)
	for _, v := range paladin.Keys() {
		name, suffix := normalizationFileName(v)
		if name == typeName {
			fileName = v
			switch suffix {
			case "lua":
				cType = CONFIG_TYPE_LUA
			case "json":
				cType = CONFIG_TYPE_JSON
			case "toml":
				cType = CONFIG_TYPE_TOML
			default:
				fileContent, err := paladin.Get(v).String()
				if err != nil {
					cType = CONFIG_TYPE_UNKNOWN
				}else {
					cType = judgeConfigTypeByFileContent(fileContent)
				}
			}
			break
		}
	}
	return
}

func getConfigTypeName(configPtr IConfig) (typeName string) {
	t := reflect.TypeOf(configPtr)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	typeName = strings.ToLower(t.String())
	names := strings.Split(typeName, ".")
	if len(names) > 1 {
		typeName = names[len(names) - 1]
	}
	return
}

func getConfigRawTypeName(configPtr IConfig) (typeName string) {
	t := reflect.TypeOf(configPtr)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	typeName = strings.ToLower(t.String())
	return
}

func normalizationFileName(fileName string) (newName string, suffix string) {
	newName = fileName
	reg := regexp.MustCompile(`_`)
	newName = reg.ReplaceAllString(newName, ``)
	newName = strings.ToLower(newName)
	names := strings.Split(newName, ".")
	if len(names) > 1 {
		newName = names[0]
		suffix = names[1]
	}
	return
}

func judgeConfigTypeByFileContent(fileContent string) (cType configType) {
	if isJsonFile(fileContent) {
		return CONFIG_TYPE_JSON
	}else if isTomlFile(fileContent) {
		return CONFIG_TYPE_TOML
		//简单的toml文件有可能是正确的lua格式，必须先判断toml
	}else if isLuaFile(fileContent) {
		return CONFIG_TYPE_LUA
	}else {
		return CONFIG_TYPE_UNKNOWN
	}
}

func isTomlFile(fileContent string) bool {
	var tmp interface{}
	err :=  tomlParser(fileContent, &tmp)
	return err == nil
}

func isJsonFile(fileContent string) bool {
	var tmp interface{}
	return jsonParser(fileContent, &tmp) == nil
}

func isLuaFile(fileContent string) bool {
	if _, err := luabridge.SafeLoad(fileContent);err != nil {
		return false
	}
	return true
}