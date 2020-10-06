package module

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	TRSTATE = "terraform.tfstate"
)

type TerraformState struct {
	Version          int                              `json:"version"`
	TerraformVersion string                           `json:"terraform_version"`
	Serial           int                              `json:"serial"`
	Lineage          string                           `json:"lineage"`
	Outputs          map[string]*TerraformStateOutput `json:"outputs"`
	Resources        []interface{}                    `json:"resources"`
}

type TerraformStateOutput struct {
	Value map[string]interface{} `json:"value"`
	Type  []interface{}          `json:"type"`
}

func NewTerraformState(path string) *TerraformState {
	t := new(TerraformState)
	if output, err := t.init(path); err == nil {
		t = output
		return t
	} else {
		return nil
	}
}

func (t *TerraformState) init(path string) (*TerraformState, error) {
	file, err := os.Open(path + "/" + TRSTATE)

	if err != nil {
		return nil, err
	}
	raw, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var output *TerraformState
	err = json.Unmarshal(raw, &output)
	return output, err
}

func (t *TerraformState) GetRaw() []byte {
	raw, _ := json.Marshal(t)
	return raw
}
