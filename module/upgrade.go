package module

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/codingXiang/go-terraform-runner    /model"
)

const (
	required_version   = "required_version"
	required_providers = "required_providers"
)

var UpgradeSource *viper.Viper

type upgrade map[string]interface{}

type sourceEntity struct {
	Source string `json:"source"`
}

func newSourceEntity(provider *providerEntity, version string) *sourceEntity {
	if in := UpgradeSource.GetString(SOURCE + "." + version + "." + provider.CloudProvider); in != "" {
		up := new(sourceEntity)
		up.Source = in
		return up
	} else {
		return nil
	}
}

type upgradeEntity struct {
	RequiredProvider map[string]*sourceEntity `json:"required_providers"`
	RequiredVersion  string                   `json:"required_version"`
	Data             map[string]interface{}   `json:"-"`
}

func NewUpgrade(requiredVersion string) *upgradeEntity {
	u := new(upgradeEntity)
	u.Data = make(map[string]interface{})
	u.RequiredVersion = requiredVersion
	u.RequiredProvider = make(map[string]*sourceEntity)
	return u
}

func NewUpgradeWithProvider(p *model.TerraformProvider) *upgradeEntity {
	u := NewUpgrade(p.Upgrade[0].Version)
	u.RequiredProvider[p.Name] = &sourceEntity{Source: p.Upgrade[0].Content}
	return u
}

func (u *upgradeEntity) AddProvider(provider *providerEntity) *upgradeEntity {
	up := newSourceEntity(provider, u.RequiredVersion)
	u.RequiredProvider[provider.CloudProvider] = up
	return u
}

func (u *upgradeEntity) New() upgrade {
	u.Data[required_providers] = u.RequiredProvider
	u.Data[required_version] = fmt.Sprintf(">= %s", u.RequiredVersion)
	return u.Data
}
