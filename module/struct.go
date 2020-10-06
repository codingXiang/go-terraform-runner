package module

import "encoding/json"

type Config struct {
	Terraform []upgrade            `json:"terraform,omitempty"`
	Variable  map[string]*Variable  `json:"variable,omitempty"`
	Provider  map[string][]Provider `json:"provider,omitempty"`
	Module    map[string]Module     `json:"module,omitempty"`
	Output    map[string]*Output    `json:"output,omitempty"`
}

func NewConfig() *Config {
	return &Config{
		Terraform: make([]upgrade, 0),
		Variable:  make(map[string]*Variable),
		Provider:  make(map[string][]Provider),
		Module:    make(map[string]Module),
		Output:    make(map[string]*Output),
	}
}

//AddVariable 動態增加 Variable
//   - v : Variable 模組
func (m *Config) AddVariable(v *VariableEntity) *Config {
	m.Variable[v.Key] = v.New()
	return m
}
//AddMultipleVariable 批次加入 Variable
func (m *Config) AddMultipleVariable(vs... *VariableEntity) *Config {
	for _, v := range vs {
		m.AddVariable(v)
	}
	return m
}

//AddProvider 動態增加 Provider
//   - p : providerEntity 模組
func (m *Config) AddProvider(p *providerEntity) *Config {
	tmp := m.Provider[p.CloudProvider]
	if tmp == nil {
		tmp = make([]Provider, 0)
	}
	m.Provider[p.CloudProvider] = append(tmp, p.New())
	return m
}

//AddModule 動態增加 Module
//   - entity : ModuleEntity 模組
func (m *Config) AddModule(entity *ModuleEntity) *Config {
	m.Module[entity.Name] = entity.Build().New()
	m.addOutput(entity)
	return m
}

//AddUpgrade 動態增加 Upgrade
//   - entity : upgradeEntity 模組
func (m *Config) AddUpgrade(entity *upgradeEntity) *Config {
	m.Terraform = append(m.Terraform, entity.New())
	return m
}

func (m *Config) addOutput(entity *ModuleEntity) *Config {
	output := NewOutput(entity)
	m.Output[entity.Name] = output
	return m
}

func (m *Config) Build() []byte {
	result, _ := json.Marshal(m)
	return result
}
