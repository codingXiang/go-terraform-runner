package runner

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/codingXiang/go-logger"
	"github.com/codingXiang/go-terraform-runner/command"
	"github.com/codingXiang/go-terraform-runner/storage"
	"github.com/codingXiang/go-workflow"
	"github.com/hashicorp/go-uuid"
)

type TerraformRunner struct {
	ID     string
	worker *workflow.Workflow
	src    *storage.ConfigSource
}

func New(src *storage.ConfigSource, id string) *TerraformRunner {
	r := new(TerraformRunner)
	if id != "" {
		r.ID = id
	} else {
		r.ID = r.generateID()
	}
	return r.init(src)
}

func (r *TerraformRunner) init(src *storage.ConfigSource) *TerraformRunner {
	r.worker = r.StepWorkflow()
	r.src = src
	err := os.MkdirAll(r.ID, os.ModePerm)
	if err != nil {
		logger.Log.Error(err)
	}
	return r
}

func (r *TerraformRunner) generateID() string {
	id, _ := uuid.GenerateUUID()
	return id
}

func (r *TerraformRunner) StepWorkflow() *workflow.Workflow {
	w := workflow.New()
	w.OnFailure = workflow.RetryFailure(1)
	return w
}

func (r *TerraformRunner) AddStep(command command.Command) *TerraformRunner {
	command = setCommandPath(command, r.ID)
	step := new(workflow.Step)
	step.Label = command.GetMeta().Label
	step.Run = command.Run
	r.worker.AddStep(step)
	return r
}

func (r *TerraformRunner) Run() error {
	return r.worker.Run()
}

func (r *TerraformRunner) Clean() error {
	return os.RemoveAll(r.ID)
}

func (r *TerraformRunner) ModifyExistProject() error {
	data, err := r.src.GetConfigRecord("config", r.ID)
	if err != nil {
		return err
	}
	m, _ := json.MarshalIndent(data.Config, "", " ")
	return ioutil.WriteFile(r.ID+"/"+r.ID+".tf.json", m, 0666)
}

func setCommandPath(command command.Command, path string) command.Command {
	meta := command.GetMeta()
	meta.Path = path
	command.SetMeta(meta)
	return command
}
