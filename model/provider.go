package model

type TerraformProvider struct {
	ID      int64                      `json:"id" gorm:"primary_key"`
	Name    string                     `json:"name" gorm:"comment:'名稱'"`
	Alias   string                     `json:"alias" gorm:"comment:'使用分區'"`
	Source  string                     `json:"source" gorm:"comment:'來源'"`
	Datas   []TerraformProviderData    `json:"datas"`
	Upgrade []TerraformProviderUpgrade `json:"upgrade"`
}

type TerraformProviderData struct {
	ID                  int64              `json:"id" gorm:"primary_key"`
	TerraformProviderID int64              `json:"providerId" gorm:"not null;comment:'訊息類別'"`
	TerraformProvider   *TerraformProvider `json:"provider,omitempty" gorm:"foreignkey:TerraformProviderID;association_foreignkey:ID"`
	Key                 string             `json:"key" gorm:"comment:'key'"`
	Value               string             `json:"value" gorm:"comment:'值'"`
}

type TerraformProviderUpgrade struct {
	ID                  int64              `json:"id" gorm:"primary_key"`
	TerraformProviderID int64              `json:"providerId" gorm:"not null;comment:'訊息類別'"`
	TerraformProvider   *TerraformProvider `json:"provider,omitempty" gorm:"foreignkey:TerraformProviderID;association_foreignkey:ID"`
	Content             string             `json:"content"`
	Version             string             `json:"version"`
}
