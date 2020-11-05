package main

import (
	"github.com/codingXiang/configer"
	"github.com/codingXiang/go-logger"
	"github.com/codingXiang/go-orm"
	"github.com/codingXiang/go-terraform-runner/command"
	"github.com/codingXiang/go-terraform-runner/model"
	"github.com/codingXiang/go-terraform-runner/module"
	"github.com/codingXiang/go-terraform-runner/runner"
	"github.com/codingXiang/go-terraform-runner/storage"
	"log"
)

func init() {
	configer.Config = configer.NewConfiger()
	var conf = configer.NewConfigerCore("yaml", "config", "./config", ".")
	configer.Config.AddCore("config", conf)

	m := configer.NewConfigerCore("yaml", "module", "./config", ".")
	configer.Config.AddCore("module", m)

	if viper, err := conf.ReadConfig(nil); err == nil {
		logger.Log = logger.NewLoggerWithConfiger(viper)
	}
	if in, err := orm.NewOrm("database", configer.Config.GetCore("config")); err == nil {
		orm.DatabaseORM = in
		orm.DatabaseORM.CheckTable(true, &model.TerraformProvider{})
		orm.DatabaseORM.CheckTable(true, &model.TerraformProviderData{})
		orm.DatabaseORM.CheckTable(true, &model.TerraformModule{})
		orm.DatabaseORM.CheckTable(true, &model.TerraformModuleData{})
		orm.DatabaseORM.CheckTable(true, &model.TerraformProviderUpgrade{})
	} else {
		panic(err)
	}

	if conf, err := m.ReadConfig(nil); err == nil {
		module.UpgradeSource = conf
	}
}

func main() {
	if conf, err := configer.Config.GetCore("config").ReadConfig(nil); err == nil {
		s := storage.NewConfigSource(conf, orm.DatabaseORM)
		in := getConfig(s)
		var id = ""
		if output, err := s.InsertData("config", in); err == nil {
			id = output
			log.Println("上傳成功, id 為", id)
		} else {
			log.Fatal("上傳失敗，錯誤為", err)
		}
		r := runner.New(s, id)
		r.ModifyExistProject()
		r.AddStep(command.NewInitCommand(nil)).
			AddStep(command.NewPlanCommand(nil)).
			AddStep(command.NewApplyCommand(nil))
		if err := r.Run(); err != nil {
			logger.Log.Error(err)
		}

		state := module.NewTerraformState(r.ID)

		if err := s.UpdateRecordState("config", id, state); err == nil {
			log.Println("更新狀態成功，", string(state.GetRaw()))
		} else {
			logger.Log.Error("更新狀態失敗，錯誤為", err)
		}

		logger.Log.Info("output =", string(state.GetRaw()))
		r.StepWorkflow()
		r.AddStep(command.NewPlanCommand(&command.Meta{Commands: []string{"-destroy"}})).
			AddStep(command.NewApplyCommand(nil))
		if err := r.Run(); err != nil {
			logger.Log.Error(err)
		}
		state = module.NewTerraformState(r.ID)
		if err := s.UpdateRecordState("config", id, state); err == nil {
			log.Println("更新狀態成功，", string(state.GetRaw()))
		} else {
			logger.Log.Error("更新狀態失敗，錯誤為", err)
		}
		r.Clean()
	}
}

func getConfig(s *storage.ConfigSource) *module.Config {
	const IDENTITY = "gitlab"
	in := module.NewConfig()

	//透過 provider 從資料庫中取得相關的 module
	moduleVars, err := s.GetModule(nil, "group")
	if err != nil {
		panic(err)
	}
	modules := make([]*module.ModuleEntity, 0)
	for _, v := range moduleVars {
		//透過 Provider 名稱 從資料庫中取得 Provider instance
		p, _ := s.GetProvider(v.Provider.Name, "")
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
