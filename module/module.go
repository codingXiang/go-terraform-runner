package module

import (
	"fmt"
	"github.com/codingXiang/go-terraform-runner/model"
)

const (
	PROVIDER = "providers"
	SOURCE   = "source"
	MODULE   = "module"
)

type Module map[string]interface{}

type ModuleEntity struct {
	Name      string                 `json:"-"`
	Providers map[string]interface{} `json:"-"`
	Source    string                 `json:"-"`
	Data      map[string]interface{} `json:"-"`
	Vars      []*VariableEntity      `json:"-"`
}

//NewModule 建立 Module 實例
func NewModule(name string) *ModuleEntity {
	return &ModuleEntity{
		Name:      name,
		Providers: make(map[string]interface{}),
		Data:      make(map[string]interface{}),
		Vars:      make([]*VariableEntity, 0),
	}
}

//func NewModuleWithVars(v *VariableEntity) *ModuleEntity {
//	return NewModule(v.Key).init(v)
//}

func NewModuleWithVar(v *model.TerraformModule) *ModuleEntity {
	m := NewModule(v.Name).init(v)
	if v.Source == "" {
		m.SetSource(v.Provider.Source)
	} else {
		m.SetSource(v.Source)
	}
	return m
}

func (e *ModuleEntity) init(v *model.TerraformModule) *ModuleEntity {
	if v.RealData != nil {
		for _, data := range v.RealData {
			varID := fmt.Sprintf("%s-%s-%s-%s", v.Provider.Name, v.Provider.Alias, data.TerraformModuleID, data.Key)
			va := NewVariableEntity(varID, data.Value, "")
			e.Vars = append(e.Vars, va)
			e.AddData(data.Key, va.GetRef())
		}
		return e
	}
	for _, data := range v.Datas {
		if data.IsModuleLink {
			e.AddData(data.Key, data.Default)
		} else {
			varID := fmt.Sprintf("%s-%s-%s-%s", v.Provider.Name, v.Provider.Alias, data.TerraformModuleID, data.Key)
			va := NewVariableEntity(varID, data.Default, data.Description)
			e.Vars = append(e.Vars, va)
			e.AddData(data.Key, va.GetRef())
		}
	}
	return e
}

//AddProvider	加入 provider 設定
//	@param
// 		- p : providerEntity 實例
//	@return
//		- return Module 自己
func (e *ModuleEntity) AddProvider(p *providerEntity) *ModuleEntity {
	e.Providers[p.CloudProvider] = p.GetProviderLink()
	return e
}

//SetSource	設定 Module Source
//	@param
// 		- source : source 的位置
//	@return
//		- return Module 自己
func (e *ModuleEntity) SetSource(source string) *ModuleEntity {
	e.Source = source
	return e
}

//AddData 加入各別的設定值
//	@param
// 		- key : 設定的 key
//		- value : 設定的值
//	@return
//		- return Module 自己
func (e *ModuleEntity) AddData(key string, value interface{}) *ModuleEntity {
	e.Data[key] = value
	return e
}

//Build 建立一個完整的 Module
//	@return
//		- return Module 自己
func (e *ModuleEntity) Build() *ModuleEntity {
	e.Data[PROVIDER] = e.Providers
	e.Data[SOURCE] = e.Source
	return e
}

//New 產生 Module 的 Map
//	@return
//		- map[string]interface{}
func (e *ModuleEntity) New() Module {
	return e.Data
}
