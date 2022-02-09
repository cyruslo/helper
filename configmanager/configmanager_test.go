package configmanager

import (
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"testing"
)

type configLua struct {
	Name string
}

type configToml struct {
	Name string
}

type configJson struct {
	Name string
}

func (self *configLua)CheckValid() (err error) {
	return nil
}

func (self *configToml)CheckValid() (err error) {
	return nil
}

func (self *configJson)CheckValid() (err error) {
	return nil
}

func TestConfigManager(t *testing.T)  {
	err := paladin.Init(nil)
	if err != nil {
		t.Logf("TestConfigManager paladin init err:%s\n", err.Error())
		t.Fail()
	}
	cLua := &configLua{}
	cToml := &configToml{}
	cJson := &configJson{}
	RegAutoConfig(cLua)
	RegAutoConfig(cToml)
	RegAutoConfig(cJson)
	if err = Init();err != nil {
		t.Logf("TestConfigManager init err:%s", err.Error())
		t.Fail()
	}
	t.Logf("%v", *cLua)
	t.Logf("%v", *cToml)
	t.Logf("%v", *cJson)
	ReloadAll()
}
