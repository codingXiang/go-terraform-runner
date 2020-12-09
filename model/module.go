package model

import (
	"github.com/gofrs/uuid"
)

type TerraformModule struct {
	ID         string                    `json:"id" gorm:"primary_key"`
	Name       string                    `json:"name" gorm:"not_null;unique_index:idx1"`
	Alias      string                    `json:"alias" gorm:"not_null;unique_index:idx1"`
	Source     string                    `json:"source" gorm:"module的參照來源"`
	ProviderID int64                     `json:"providerId" gorm:"not null;comment:'訊息類別'"`
	Provider   *TerraformProvider        `json:"provider,omitempty" gorm:"foreignkey:ProviderID;association_foreignkey:ID"`
	Datas      []TerraformModuleData     `json:"datas"`
	RealData   []TerraformModuleRealData `json:"datas"`
}

func (g *TerraformModule) BeforeCreate() (err error) {
	u, _ := uuid.NewV4()
	g.ID = u.String()
	return
}

type TerraformModuleData struct {
	ID                int64            `json:"id" gorm:"primary_key;comment:'id'"`
	TerraformModuleID string           `json:"moduleId" gorm:"Column:module_id;not null;comment:'模組id'"`
	TerraformModule   *TerraformModule `json:"module,omitempty" gorm:"foreignkey:TerraformModuleID;association_foreignkey:ID"`
	Key               string           `json:"key"`
	IsRequired        bool             `json:"isRequired" gorm:"default:1;comment:'是否為必要欄位'"`
	Default           string           `json:"default,omitempty" gorm:"comment:'預設資料'"`
	IsModuleLink      bool             `json:"isModuleLink" gorm:"default:0;是否為動態的模組連結"`
	Value             string           `json:"value" gorm:"comment:'實際參數值'"`
	Description       string           `json:"description" gorm:"comment:'參數描述'"`
}

type TerraformModuleRealData struct {
	ID                int64            `json:"id" gorm:"primary_key;comment:'id'"`
	IdentityID        string           `json:"identityId" gorm:"unique_index:idx1"`
	TerraformModuleID string           `json:"moduleId" gorm:"Column:module_id;unique_index:idx1;not null;comment:'模組id'"`
	TerraformModule   *TerraformModule `json:"module,omitempty" gorm:"foreignkey:TerraformModuleID;association_foreignkey:ID"`
	Key               string           `json:"key" gorm:"unique_index:idx1"`
	Value             string           `json:"value" gorm:"comment:'實際參數值'"`
}
