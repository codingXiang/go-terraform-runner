package configer

import (
	"github.com/codingXiang/go-orm/v2"
	"github.com/codingXiang/go-orm/v2/mongo"
	"github.com/codingXiang/go-terraform-runner/model"
	"github.com/jinzhu/gorm"
)

type Handler struct {
	MongoClient *mongo.Client
	DB          *gorm.DB
}

func New(mongoClient *mongo.Client, orm *orm.Orm) *Handler {
	return &Handler{
		MongoClient: mongoClient,
		DB:          orm.DB,
	}
}

func (c *Handler) GetProvider(name, alias string) (*model.TerraformProvider, error) {
	result := new(model.TerraformProvider)
	query := c.DB.Preload("Datas").Preload("Upgrade")
	query = query.Where("name = ?", name)
	if alias != "" {
		query = query.Where("alias = ?", alias)
	}
	err := query.First(&result).Error
	return result, err
}

func (c *Handler) GetModule(dataId string, provider *model.TerraformProvider, alias ...string) ([]*model.TerraformModule, error) {
	result := make([]*model.TerraformModule, 0)
	db := c.DB
	if provider != nil {
		db = db.Preload("Provider").Where("provider_id = ?", provider.ID)
	} else {
		db = db.Preload("Provider")
	}

	if dataId != "" {
		db = db.Preload("RealData", func(db *gorm.DB) *gorm.DB {
			return db.Where("terraform_module_real_data.identity_id = ?", dataId)
		})
	} else {
		db = db.Preload("Datas")
	}

	for index, n := range alias {
		if index > 0 {
			db = db.Or("alias = ?", n)
		} else {
			db = db.Where("alias = ?", n)
		}
	}
	err := db.Find(&result).Error
	return result, err
}
