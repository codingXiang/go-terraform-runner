package module

import "fmt"

type Output struct {
	Value string `json:"value"`
}

//NewOutput 建立 output 實例
func NewOutput(m *ModuleEntity) *Output {
	return &Output{
		Value: fmt.Sprintf(`${module.%s}`, m.Name),
	}
}
