package storage

import (
	"github.com/hashicorp/go-uuid"
	"gopkg.in/mgo.v2/bson"
	"io"
	"os"
	"terraformRunner/model"
	"terraformRunner/module"
)

const (
	VAR = "var"
)

type ConfigRecord struct {
	Identity string                 `json:"identity"`
	Config   interface{}            `json:"config"`
	State    *module.TerraformState `json:"state,omitempty"`
}

func (c *ConfigSource) InsertData(collection string, data interface{}) (string, error) {
	id, _ := uuid.GenerateUUID()
	config := &ConfigRecord{
		Identity: id,
		Config:   data,
	}
	return id, c.mongo.db.C(collection).Insert(config)
}

func (c *ConfigSource) UpdateRecordState(collection, id string, state *module.TerraformState) error {
	update := bson.M{
		"$set": bson.M{
			"state": state,
		},
	}
	selector := bson.M{"identity": id}
	_, err := c.mongo.db.C(collection).Upsert(selector, update)
	return err
}

func (c *ConfigSource) GetConfigRecord(collection string, id string) (*ConfigRecord, error) {
	target := new(ConfigRecord)
	err := c.mongo.db.C(collection).Find(bson.M{"identity": id}).One(&target)
	return target, err
}

func (c *ConfigSource) GetProvider(name, alias string) (*model.TerraformProvider, error) {
	result := new(model.TerraformProvider)
	query := c.orm.Preload("Datas").Preload("Upgrade")
	query = query.Where("name = ?", name)
	if alias != "" {
		query = query.Where("alias = ?", alias)
	}
	err := query.First(&result).Error
	return result, err
}

func (c *ConfigSource) GetModule(provider *model.TerraformProvider, moduleName ...string) ([]*model.TerraformModule, error) {
	result := make([]*model.TerraformModule, 0)
	db := c.orm.Preload("Provider").Preload("Datas").Where("provider_id = ?", provider.ID)
	condition := ""
	for i := 0; i < len(moduleName); i += 1 {
		if i > len(moduleName)-1 {
			condition += "id = ? OR "
		} else {
			condition += "id = ?"
		}
	}
	err := db.Where(condition, moduleName).Find(&result).Error
	return result, err
}

func (c *ConfigSource) GetVar(identity, dependency string) ([]*module.VariableEntity, error) {
	target := make([]*module.VariableEntity, 0)
	err := c.mongo.db.C(VAR).Find(bson.M{"identity": identity, "dependency": dependency}).All(&target)
	return target, err
}

func (c *ConfigSource) UploadFile(filename, sourceFilePath string) error {
	file, err := c.mongo.db.GridFS("terraform").Create(filename)
	check(err)
	source, err := os.Open(sourceFilePath)
	check(err)
	defer source.Close()
	_, err = io.Copy(file, source)
	check(err)
	err = file.Close()
	check(err)
	return err
}

func (c *ConfigSource) DownloadFile(filename, targetFilePath string) error {
	file, err := c.mongo.db.GridFS("terraform").Open(filename)
	check(err)
	out, _ := os.OpenFile(targetFilePath, os.O_CREATE|os.O_RDWR, 0666)
	defer out.Close()
	_, err = io.Copy(out, file)
	check(err)
	err = file.Close()
	check(err)
	return err
}

func check(err error) {
	if err != nil {
		panic(err.Error())
	}
}
