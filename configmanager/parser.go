package configmanager

import (
	"encoding/json"
	//"git.huoys.com/qp/luabridge"
	"github.com/BurntSushi/toml"
)

type parser func(fileContent string, config interface{}) (err error)

func luaParser(fileContent string, config interface{}) (err error) {
	//var vm *luabridge.LuaVM
	//if vm ,err = luabridge.NewLuaVMByString(fileContent);err != nil {
	//	return
	//}
	//if err = vm.Get(-1, config);err != nil {
	//	return
	//}
	return
}

func jsonParser(fileContent string, config interface{}) (err error) {
	if err = json.Unmarshal([]byte(fileContent), config);err != nil {
		return
	}
	return
}

func tomlParser(fileContent string, config interface{}) (err error) {
	if err = toml.Unmarshal([]byte(fileContent), config);err != nil {
		return
	}
	return
}