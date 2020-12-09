package runner

import (
	"encoding/json"
	"fmt"
	"github.com/codingXiang/go-orm/v2/mongo"
	configer2 "github.com/codingXiang/go-terraform-runner/configer"
	"io/ioutil"
	"os"

	"github.com/codingXiang/go-logger"
	"github.com/codingXiang/go-terraform-runner/command"
	"github.com/codingXiang/go-workflow"
)

const (
	DATA = "data"
)

type TerraformRunner struct {
	worker *workflow.Workflow
	id     string
	c      *configer2.Handler
}

func New(id string, handler *configer2.Handler) *TerraformRunner {
	r := new(TerraformRunner)
	return r.init(id, handler)
}

func (r *TerraformRunner) init(id string, handler *configer2.Handler) *TerraformRunner {
	r.worker = r.InitWorkFlow()
	r.c = handler
	r.id = id
	err := os.MkdirAll(id, os.ModePerm)
	if err != nil {
		logger.Log.Error(err)
	}
	return r
}

func (r *TerraformRunner) InitWorkFlow() *workflow.Workflow {
	w := workflow.New()
	w.OnFailure = workflow.NoRetry()
	return w
}

func (r *TerraformRunner) AddStep(command command.Command) *TerraformRunner {
	command = setCommandPath(command, r.id)
	step := new(workflow.Step)
	step.Label = command.GetMeta().Label
	step.Run = command.Run
	r.worker.AddSteps(step)
	return r
}

func (r *TerraformRunner) Run(successCallback func(objs ...interface{}) error, failCallback func(objs ...interface{}) error) error {
	return r.worker.Run(successCallback, failCallback)
}

func (r *TerraformRunner) Clean() error {
	return os.RemoveAll(r.id)
}

func (r *TerraformRunner) ModifyExistProject() error {
	data, err := r.c.MongoClient.C(DATA).First(mongo.NewSearchCondition("", r.id, nil, nil))
	if err != nil {
		return err
	}

	m, _ := json.MarshalIndent(data.Raw, "", " ")
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s.tf.json", r.id, r.id), m, 0666)
	if err != nil {
		return err
	}
	//m, _ := json.MarshalIndent(data.Config, "", " ")
	//err = ioutil.WriteFile(r.ID+"/"+r.ID+".tf.json", m, 0666)
	//if err != nil {
	//	return err
	//}
	//if config.Get("state") != nil {
	//	m, _ := json.MarshalIndent(data.State, "", " ")
	//	err = ioutil.WriteFile(r.ID+"/terraform.tfstate", m, 0666)
	//	if err != nil {
	//		return err
	//	}
	//}

	return nil
}

func setCommandPath(command command.Command, path string) command.Command {
	meta := command.GetMeta()
	meta.Path = path
	command.SetMeta(meta)
	return command
}
