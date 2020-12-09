package main

import (
	"github.com/codingXiang/configer/v2"
	"github.com/codingXiang/go-logger/v2"
	"github.com/codingXiang/go-orm/v2"
	"github.com/codingXiang/go-orm/v2/mongo"
	"github.com/codingXiang/go-terraform-runner/command"
	configer2 "github.com/codingXiang/go-terraform-runner/configer"
	"github.com/codingXiang/go-terraform-runner/model"
	"github.com/codingXiang/go-terraform-runner/module"
	"github.com/codingXiang/go-terraform-runner/runner"
	uuid "github.com/satori/go.uuid"
	"log"
)

func init() {
	configer.Config = configer.NewConfiger()
	var conf = configer.NewCore(configer.YAML, "config", "./config", ".")
	configer.Config.AddCore("config", conf)

	if c, err := conf.ReadConfig(); err == nil {
		logger.Log = logger.New(c)
		if o, err := orm.New(c); err == nil {
			orm.DatabaseORM = o
			orm.DatabaseORM.CheckTable(true, &model.TerraformProvider{})
			orm.DatabaseORM.CheckTable(true, &model.TerraformProviderData{})
			orm.DatabaseORM.CheckTable(true, &model.TerraformModule{})
			orm.DatabaseORM.CheckTable(true, &model.TerraformModuleData{})
			orm.DatabaseORM.CheckTable(true, &model.TerraformProviderUpgrade{})
			orm.DatabaseORM.CheckTable(true, &model.TerraformModuleRealData{})
		}
		mongo.MongoClient = mongo.New(c)
	}

	//if conf, err := m.ReadConfig(nil); err == nil {
	//	module.UpgradeSource = conf
	//}
}

func main() {
	if _, err := configer.Config.GetCore("config").ReadConfig(); err == nil {
		handler := configer2.New(mongo.MongoClient, orm.DatabaseORM)
		in := getConfig(handler)
		var id = uuid.NewV4().String()

		if err := handler.MongoClient.C(runner.DATA).Insert(mongo.NewRawData(id, nil, in)); err == nil {
			log.Println("上傳成功, id 為", id)
		} else {
			log.Fatal("上傳失敗，錯誤為", err)
		}

		r := runner.New(id, handler)
		r.ModifyExistProject()

		r.AddStep(command.NewInitCommand(nil)).
			AddStep(command.NewPlanCommand(nil)).
			AddStep(command.NewApplyCommand(nil))
		if err := r.Run(func(objs ...interface{}) error {
			logger.Log.Info("執行成功")
			return nil
		}, func(objs ...interface{}) error {
			logger.Log.Error("執行失敗")
			return nil
		}); err != nil {
			logger.Log.Error(err)
		}

		state := module.NewTerraformState(id)
		if err := handler.MongoClient.C("state").Insert(mongo.NewRawData(id, nil, state)); err == nil {
			log.Println("更新狀態成功，", string(state.GetRaw()))
		} else {
			logger.Log.Error("更新狀態失敗，錯誤為", err)
		}

		logger.Log.Info("output =", string(state.GetRaw()))
		r.InitWorkFlow()
		r.AddStep(command.NewPlanCommand(&command.Meta{Commands: []string{"-destroy"}})).
			AddStep(command.NewApplyCommand(nil))
		if err := r.Run(func(objs ...interface{}) error {
			logger.Log.Info("刪除成功")
			return nil
		}, func(objs ...interface{}) error {
			logger.Log.Error("刪除失敗")
			return nil
		}); err != nil {
			logger.Log.Error(err)
		}
		state = module.NewTerraformState(id)
		if data, err := handler.MongoClient.C("state").First(mongo.NewSearchCondition("", id, nil, nil)); err == nil {
			log.Println("更新狀態成功，", data.GetRaw())
		} else {
			logger.Log.Error("更新狀態失敗，錯誤為", err)
		}
		r.Clean()
	}
}

func getConfig(handler *configer2.Handler) *module.Config {
	in := module.NewConfig()

	//透過 provider 從資料庫中取得相關的 module
	moduleVars, err := handler.GetModule("test1234", nil, "Aliyun.paas.group")
	if err != nil {
		panic(err)
	}
	modules := make([]*module.ModuleEntity, 0)
	for _, v := range moduleVars {
		//透過 Provider 名稱 從資料庫中取得 Provider instance
		p, _ := handler.GetProvider(v.Provider.Name, "")
		//轉換成 provider entity
		provider := module.NewProviderWithVar(p)
		//透過 provider 取得 upgrade
		upgrade := module.NewUpgradeWithProvider(p)
		//加入 provider
		in.AddProvider(provider)
		in.AddUpgrade(upgrade)
		m := module.NewModuleWithVar(v).AddProvider(provider)
		modules = append(modules, m)
		in.AddMultipleVariable(m.Vars...)
	}

	for _, m := range modules {
		in.AddModule(m)
	}
	return in
}
