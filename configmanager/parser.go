package configmanager

import (
	"encoding/json"
	lua "github.com/AzureWrathCyd/gopher-lua"
	"github.com/BurntSushi/toml"
	"github.com/cyruslo/helper/luabridge"
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
	var vm *lua.LState
	if vm, err = luabridge.SafeLoad(fileContent); err != nil {
		return
	}
	if err = luabridge.SafeCall(vm, "", config); err != nil {
		return
	}
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