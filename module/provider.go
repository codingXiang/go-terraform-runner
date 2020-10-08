package module

import "github.com/codingXiang/go-terraform-runner/model"

type Provider map[string]interface{}

const (
	ALIAS = "alias"
)

type providerEntity struct {
	CloudProvider string                 `json:"-"`
	Data          map[string]interface{} `json:"-"`
}

func NewProvider(cloudProvider string) *providerEntity {
	return &providerEntity{
		CloudProvider: cloudProvider,
		Data:          make(map[string]interface{}),
	}
}

//func NewProviderWithVar(cloudProvider string, v *VariableEntity) *providerEntity {
//	if v.Dependency == PROVIDER {
//		p := NewProvider(cloudProvider).init(v)
//		return p
//	}
//	return nil
//}

func NewProviderWithVar(v *model.TerraformProvider) *providerEntity {
	p := NewProvider(v.Name)
	p.AddConfig(ALIAS, v.Alias)
	for _, data := range v.Datas {
		p.AddConfig(data.Key, data.Value)
	}
	return p
}

//func (p *providerEntity) init(v *VariableEntity) *providerEntity {
//	in := v.GetMap()
//	p.CloudProvider = v.Key
//	for key, value := range in {
//		if key == ALIAS {
//			p.AddConfig(key, value)
//		} else {
//			p.AddConfig(key, v.GetRefWithKey(key))
//		}
//	}
//	return p
//}

//AddConfig 加入各別的設定值
//	@param
// 		- key : 設定的 key
//		- value : 設定的值
//	@return
//		- return providerEntity 自己
func (p *providerEntity) AddConfig(key string, value interface{}) *providerEntity {
	p.Data[key] = value
	return p
}

func (p *providerEntity) GetProviderLink() string {
	result := p.CloudProvider
	if in := p.Data[ALIAS]; in != "" {
		result += "." + in.(string)
	}
	return result
}

//New 產生 Provider 的 Map
//	@return
//		- map[string]interface{}
func (p *providerEntity) New() Provider {
	return p.Data
}
