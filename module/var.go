package module

import (
	"encoding/json"
	"github.com/codingXiang/go-logger"
	"strings"
)

const (
	REF         = `${var.@key@}`
	TYPE_STRING = "string"
	TYPE_NUMBER = "number"
	TYPE_MAP    = "map"
	TYPE_ARRAY  = "array"
)

type Variable struct {
	Description string      `json:"description,omitempty"`
	Default     interface{} `json:"default,omitempty"`
}
type VariableEntity struct {
	Dependency string    `json:"dependency"`
	Key        string    `json:"key"`
	Type       string    `json:"type"`
	Data       *Variable `json:"data"`
}

func NewVariableEntity(key string, def interface{}, description string) *VariableEntity {
	return &VariableEntity{
		Key: key,
		Data: &Variable{
			Description: description,
			Default:     def,
		},
	}
}

func (v *VariableEntity) setType(typ string) *VariableEntity {
	v.Type = typ
	return v
}

func (v *VariableEntity) GetRef() string {
	ref := strings.ReplaceAll(REF, "@key@", v.Key)
	return ref
}
func (v *VariableEntity) GetRefWithKey(key string) string {
	if v.Type == TYPE_MAP {
		return strings.ReplaceAll(REF, "@key@", v.Key+"[\""+key+"\"]")
	}
	return v.GetRef()
}

func (v *VariableEntity) New() *Variable {
	return v.Data
}

func (v *VariableEntity) GetMap() map[string]interface{} {
	if v.Type == TYPE_MAP {
		in := map[string]interface{}{}
		tmp, _ := json.Marshal(v.Data.Default)
		err := json.Unmarshal(tmp, &in)
		if err != nil {
			logger.Log.Error(err)
		}
		return in
	}
	return nil
}

func NewStringVar(key, value, description string) *VariableEntity {
	return NewVariableEntity(key, value, description).setType(TYPE_STRING)
}

func NewNumberVar(key string, value interface{}, description string) *VariableEntity {
	return NewVariableEntity(key, value, description).setType(TYPE_NUMBER)
}

func NewArrayVar(key string, value []interface{}, description string) *VariableEntity {
	return NewVariableEntity(key, value, description).setType(TYPE_ARRAY)
}

func NewMapVar(key string, value map[string]interface{}, description string) *VariableEntity {
	return NewVariableEntity(key, value, description).setType(TYPE_MAP)
}
